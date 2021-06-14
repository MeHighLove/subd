package usecase

import (
	"github.com/go-openapi/strfmt"
	"net/http"
	"strconv"
	"strings"
	smth "subd"
	"subd/constants"
	"subd/models"
	"time"
)

type Smth struct {
	repo    smth.Repository
}

func NewSmth(e smth.Repository) smth.UseCase {
	return &Smth{repo: e}
}

func (s Smth) GetThreads(slug string, limit int, since string, desc bool) (models.Threads, int) {
	isExisted, err := s.repo.CheckForum(slug)
	if err != nil {
		return models.Threads{}, http.StatusInternalServerError
	}
	if !isExisted {
		return models.Threads{}, constants.NotFound
	}

	threads, err := s.repo.GetForumThreads(slug, limit, since, desc)
	if err != nil {
		return models.Threads{}, http.StatusInternalServerError
	}

	return threads, http.StatusOK
}

func (s Smth) GetUser(nickname string) (models.User, int) {
	user, status := s.repo.GetUser(nickname)

	if status == http.StatusNotFound {
		return models.User{}, constants.NotFound
	}

	return user, http.StatusOK
}

func (s Smth) GetForumUsers(slug string, limit int, since string, desc bool) (models.Users, int) {
	isExisted, err := s.repo.CheckForum(slug)
	if err != nil {
		return models.Users{}, http.StatusInternalServerError
	}
	if !isExisted {
		return models.Users{}, constants.NotFound
	}

	users, err := s.repo.GetForumUsers(slug, limit, since, desc)
	if err != nil {
		return models.Users{}, http.StatusInternalServerError
	}

	return users, http.StatusOK
}

func (s Smth) CreateNewThread(newThread *models.Thread) (models.Thread, int) {
	user, status := s.repo.GetUser(newThread.Author)
	if status == constants.NotFound {
		return models.Thread{}, constants.NotFound
	}
	newThread.Author = user.Nickname

	forum, status := s.repo.GetForum(newThread.Forum)
	if status == constants.NotFound {
		return models.Thread{}, constants.NotFound
	}
	newThread.Forum = forum.Slug

	var err error
	newThread.Id, err = s.repo.AddNewThread(*newThread)
	if err != nil {
		thread, _ := s.repo.GetThread(newThread.Slug)
		return thread, http.StatusConflict
	}
	err = s.repo.IncrementThreads(newThread.Forum)

	s.repo.AddForumUsers(newThread.Forum, newThread.Author)

	return *newThread, http.StatusCreated
}

func (s Smth) GetThread(slugOrId string) (models.Thread, int) {
	var thread models.Thread
	var status int
	if id, err := strconv.Atoi(slugOrId); err != nil {
		isExist, err := s.repo.CheckThread(slugOrId)
		if err != nil {
			return models.Thread{}, http.StatusInternalServerError
		}
		if !isExist {
			return models.Thread{}, http.StatusNotFound
		}
		thread, err = s.repo.GetThread(slugOrId)
		if err != nil {
			return models.Thread{}, http.StatusInternalServerError
		}
	} else {
		thread, status = s.repo.GetThreadById(id)
		if status == http.StatusNotFound {
			return models.Thread{}, status
		}
		if status == http.StatusInternalServerError {
			return models.Thread{}, http.StatusInternalServerError
		}
	}
	return thread, http.StatusOK
}

func (s Smth) CreateNewPosts(newPosts []*models.Post, slugOrId string)  int {
	var thread models.Thread
	var status int
	if id, err := strconv.Atoi(slugOrId); err != nil {
		thread, status = s.repo.GetThreadStatus(slugOrId)
		if status == constants.NotFound {
			return http.StatusNotFound
		}
	} else {
		thread, status = s.repo.GetThreadById(id)
		if status == http.StatusNotFound {
			return status
		}
		if status == http.StatusInternalServerError {
			return http.StatusInternalServerError
		}
	}

	if len(newPosts) == 0 {
		return http.StatusCreated
	}

	now := time.Now()
	var err error
	var isExist bool
	//Здесь тоже можно будет делать добавление в форум-юзер функцией!!!
	for i := range newPosts {
		isExist, err = s.repo.CheckUser(newPosts[i].Author)
		if err != nil {
			return http.StatusConflict
		}
		if !isExist {
			return http.StatusNotFound
		}
		newPosts[i].Thread = int(thread.Id)
		newPosts[i].Forum = thread.Forum
		newPosts[i].Created = strfmt.DateTime(now)
		err = s.repo.AddPost(newPosts[i])
		if err != nil {
			return http.StatusConflict
		}
		err = s.repo.IncrementPosts(newPosts[i].Forum)
		s.repo.AddForumUsers(newPosts[i].Forum, newPosts[i].Author)
	}

	return http.StatusCreated
}

func (s Smth) CreateNewForum(newForum *models.Forum) (models.Forum, int) {
	user, status := s.repo.GetUser(newForum.Owner)
	if status == constants.NotFound {
		return models.Forum{}, constants.NotFound
	}
	newForum.Owner = user.Nickname

	oldForum, status := s.repo.GetForum(newForum.Slug)
	if status != constants.NotFound {
		return oldForum, http.StatusConflict
	}
	err, _ := s.repo.AddNewForum(newForum)
	if err != nil {
		return models.Forum{}, http.StatusInternalServerError
	}
	//s.repo.AddForumUsers(newForum.Slug, newForum.Owner)

	return *newForum, http.StatusCreated
}

func (s Smth) GetForum(slug string) (models.Forum, int) {
	forum, status := s.repo.GetForum(slug)

	return forum, status
}

func (s Smth) GetPost(id int, related string) (models.FullPost, int) {
	fullPost := models.FullPost{}
	var err error

	isExisted, err := s.repo.CheckPost(id)
	if err != nil {
		return models.FullPost{}, http.StatusInternalServerError
	}
	if !isExisted {
		return models.FullPost{}, constants.NotFound
	}

	post, _ := s.repo.GetPost(id)
	fullPost.Post = &post

	if related != "" {
		split := strings.Split(related, ",")

		for _, elem := range split {
			switch elem {
			case "user":
				user, _ := s.repo.GetUser(fullPost.Post.Author)
				fullPost.Author = &user
			case "thread":
				thread, _ := s.repo.GetThreadById(fullPost.Post.Thread)
				fullPost.Thread = &thread
			case "forum":
				forum, _ := s.repo.GetForum(fullPost.Post.Forum)
				fullPost.Forum = &forum
			}
		}
	}

	return fullPost, http.StatusOK

}

func (s Smth) EditMessageNull(id int) (models.PostNullMessage, int) {
	post, status := s.repo.GetPostNull(id)
	if status == constants.NotFound {
		return models.PostNullMessage{}, constants.NotFound
	}


	return post, http.StatusOK
}

func (s Smth) EditMessage(id int, message string) (models.Post, int) {
	post, status := s.repo.GetPost(id)
	if status == constants.NotFound {
		return models.Post{}, constants.NotFound
	}
	if post.Message == message {
		return post, http.StatusConflict
	}
	post.IsEdited = true
	post.Message = message

	err := s.repo.EditMessage(id, message)
	if err != nil {
		return models.Post{}, http.StatusInternalServerError
	}


	return post, http.StatusOK
}

func (s Smth) Clear() error {
	err := s.repo.Clear()
	if err != nil {
		return err
	}

	return nil
}

func (s Smth) Status() (models.Status, error) {
	status, err := s.repo.Status()
	if err != nil {
		return models.Status{}, err
	}

	return status, nil
}

func (s Smth) CreateUser(nickname string, user models.User) (models.Users, int) {
	isExist, err := s.repo.CheckUserByNicknameOrEmail(nickname, user.Email.String())
	if err != nil {
		return models.Users{}, http.StatusInternalServerError
	}

	var users models.Users

	if isExist {
		users, err = s.repo.GetUserByNicknameOrEmail(nickname, user.Email.String())
		if err != nil {
			return models.Users{}, http.StatusInternalServerError
		}
		return users, http.StatusConflict
	}


	err = s.repo.CreateUser(nickname, user)
	if err != nil {
		return models.Users{}, http.StatusInternalServerError
	}

	newUser, _ := s.repo.GetUser(nickname)
	users = append(users, newUser)

	return users, http.StatusCreated
}

func (s Smth) UpdateUser(nickname string, user models.User) (models.User, int) {
	oldUser, status := s.repo.GetUser(nickname)
	if status == constants.NotFound {
		return models.User{}, http.StatusNotFound
	}

	isExist, err := s.repo.CheckUserByEmail(user.Email.String())
	if err != nil {
		return models.User{}, http.StatusInternalServerError
	}
	if isExist {
		return models.User{}, http.StatusConflict
	}

	if user.Email == "" {
		user.Email = oldUser.Email
	}
	if user.About == "" {
		user.About = oldUser.About
	}
	if user.Fullname == "" {
		user.Fullname = oldUser.Fullname
	}

	err = s.repo.UpdateUser(nickname, user)
	if err != nil {
		return models.User{}, http.StatusInternalServerError
	}

	newUser, _ := s.repo.GetUser(nickname)

	return newUser, http.StatusOK
}

func (s Smth) UpdateThread(slugOrId string, newThread models.Thread) (models.Thread, int) {
	var thread models.Thread
	if id, err := strconv.Atoi(slugOrId); err != nil {
		oldThread, status := s.repo.GetThreadStatus(slugOrId)
		if status == constants.NotFound {
			return models.Thread{}, http.StatusNotFound
		}
		if newThread.Message == "" {
			newThread.Message = oldThread.Message
		}
		if newThread.Title == "" {
			newThread.Title = oldThread.Title
		}
		thread, err = s.repo.UpdateThread(slugOrId, newThread)
		if err != nil {
			return models.Thread{}, http.StatusInternalServerError
		}
	} else {
		oldThread, status := s.repo.GetThreadById(id)
		if status == constants.NotFound {
			return models.Thread{}, http.StatusNotFound
		}
		if newThread.Message == "" {
			newThread.Message = oldThread.Message
		}
		if newThread.Title == "" {
			newThread.Title = oldThread.Title
		}
		thread, err = s.repo.UpdateThreadById(id, newThread)
	}

	return thread, http.StatusOK
}

func (s Smth) Vote(slugOrId string, vote models.Vote) (models.Thread, int) {
	var thread models.Thread
	var status int
	isExist, err := s.repo.CheckUser(vote.Nickname)
	if err != nil {
		return models.Thread{}, http.StatusInternalServerError
	}
	if !isExist {
		return models.Thread{}, http.StatusNotFound
	}
	if id, err := strconv.Atoi(slugOrId); err != nil {
		thread, status = s.repo.GetThreadStatus(slugOrId)
		if status == constants.NotFound {
			return models.Thread{}, http.StatusNotFound
		}
		isExist, err = s.repo.CheckVote(int(thread.Id), vote.Nickname)
		if err != nil {
			return models.Thread{}, http.StatusInternalServerError
		}
		if !isExist {
			err = s.repo.AddVote(int(thread.Id), vote)
			if err != nil {
				return models.Thread{}, http.StatusInternalServerError
			}
			thread.Votes += vote.Voice
		} else {
			num, err := s.repo.GetValueVote(int(thread.Id), vote.Nickname)
			if err != nil {
				return models.Thread{}, http.StatusInternalServerError
			}
			if num != vote.Voice {
				err = s.repo.UpdateVote(int(thread.Id), vote)
				if err != nil {
					return models.Thread{}, http.StatusInternalServerError
				}
				thread.Votes += 2 * vote.Voice
			}
		}
	} else {
		thread, status = s.repo.GetThreadById(id)
		if status == constants.NotFound {
			return models.Thread{}, http.StatusNotFound
		}
		isExist, err = s.repo.CheckVote(id, vote.Nickname)
		if err != nil {
			return models.Thread{}, http.StatusInternalServerError
		}
		if !isExist {
			err = s.repo.AddVote(int(thread.Id), vote)
			if err != nil {
				return models.Thread{}, http.StatusInternalServerError
			}
			thread.Votes += vote.Voice
		} else {
			num, err := s.repo.GetValueVote(int(thread.Id), vote.Nickname)
			if err != nil {
				return models.Thread{}, http.StatusInternalServerError
			}
			if num != vote.Voice {
				err = s.repo.UpdateVote(id, vote)
				if err != nil {
					return models.Thread{}, http.StatusInternalServerError
				}
				thread.Votes += 2 * vote.Voice
			}
		}
	}

	return thread, http.StatusOK
}

func (s Smth) GetThreadSortFlat(slugOrId string, limit int, since int, desc bool) (models.Posts, int) {
	var thread models.Thread
	var status int
	if id, err := strconv.Atoi(slugOrId); err != nil {
		isExist, err := s.repo.CheckThread(slugOrId)
		if err != nil {
			return models.Posts{}, http.StatusInternalServerError
		}
		if !isExist {
			return models.Posts{}, http.StatusNotFound
		}
		thread, err = s.repo.GetThread(slugOrId)
		if err != nil {
			return models.Posts{}, http.StatusInternalServerError
		}
	} else {
		thread, status = s.repo.GetThreadById(id)
		if status == constants.NotFound {
			return models.Posts{}, http.StatusNotFound
		}
	}
	if desc == true {
		posts, err := s.repo.GetPostsFlatDesc(int(thread.Id) ,limit, since)
		if err != nil {
			return models.Posts{}, http.StatusInternalServerError
		}
		return posts, http.StatusOK
	} else {
		posts, err := s.repo.GetPostsFlat(int(thread.Id) ,limit, since)
		if err != nil {
			return models.Posts{}, http.StatusInternalServerError
		}
		return posts, http.StatusOK
	}
}

func (s Smth) GetThreadSortTree(slugOrId string, limit int, since int, desc bool) (models.Posts, int) {
	var thread models.Thread
	var status int
	if id, err := strconv.Atoi(slugOrId); err != nil {
		isExist, err := s.repo.CheckThread(slugOrId)
		if err != nil {
			return models.Posts{}, http.StatusInternalServerError
		}
		if !isExist {
			return models.Posts{}, http.StatusNotFound
		}
		thread, err = s.repo.GetThread(slugOrId)
		if err != nil {
			return models.Posts{}, http.StatusInternalServerError
		}
	} else {
		thread, status = s.repo.GetThreadById(id)
		if status == constants.NotFound {
			return models.Posts{}, http.StatusNotFound
		}
	}
	if since != 0 {
		if desc == true {
			posts, err := s.repo.GetPostsTreeSinceDesc(int(thread.Id) ,limit, since)
			if err != nil {
				return models.Posts{}, http.StatusInternalServerError
			}
			return posts, http.StatusOK
		} else {
			posts, err := s.repo.GetPostsTreeSince(int(thread.Id) ,limit, since)
			if err != nil {
				return models.Posts{}, http.StatusInternalServerError
			}
			return posts, http.StatusOK
		}
	} else {
		if desc == true {
			posts, err := s.repo.GetPostsTreeDesc(int(thread.Id) ,limit)
			if err != nil {
				return models.Posts{}, http.StatusInternalServerError
			}
			return posts, http.StatusOK
		} else {
			posts, err := s.repo.GetPostsTree(int(thread.Id) ,limit)
			if err != nil {
				return models.Posts{}, http.StatusInternalServerError
			}
			return posts, http.StatusOK
		}
	}
}

func (s Smth) GetThreadSortParentTree(slugOrId string, limit int, since int, desc bool) (models.Posts, int) {
	var thread models.Thread
	var status int
	if id, err := strconv.Atoi(slugOrId); err != nil {
		isExist, err := s.repo.CheckThread(slugOrId)
		if err != nil {
			return models.Posts{}, http.StatusInternalServerError
		}
		if !isExist {
			return models.Posts{}, http.StatusNotFound
		}
		thread, err = s.repo.GetThread(slugOrId)
		if err != nil {
			return models.Posts{}, http.StatusInternalServerError
		}
	} else {
		thread, status = s.repo.GetThreadById(id)
		if status == constants.NotFound {
			return models.Posts{}, http.StatusNotFound
		}
	}
	if since != 0 {
		if desc == true {
			posts, err := s.repo.GetPostsParentTreeSinceDesc(int(thread.Id) ,limit, since)
			if err != nil {
				return models.Posts{}, http.StatusInternalServerError
			}
			return posts, http.StatusOK
		} else {
			posts, err := s.repo.GetPostsParentTreeSince(int(thread.Id) ,limit, since)
			if err != nil {
				return models.Posts{}, http.StatusInternalServerError
			}
			return posts, http.StatusOK
		}
	} else {
		if desc == true {
			posts, err := s.repo.GetPostsParentTreeDesc(int(thread.Id) ,limit)
			if err != nil {
				return models.Posts{}, http.StatusInternalServerError
			}
			return posts, http.StatusOK
		} else {
			posts, err := s.repo.GetPostsParentTree(int(thread.Id) ,limit)
			if err != nil {
				return models.Posts{}, http.StatusInternalServerError
			}
			return posts, http.StatusOK
		}
	}
}
