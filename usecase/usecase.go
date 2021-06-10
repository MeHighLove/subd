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

	if since == "" {
		since = "0001-01-01 00:00:00.000000"
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
	isExist, err := s.repo.CheckUser(newThread.Author)
	if err != nil {
		return models.Thread{}, http.StatusInternalServerError
	}
	if !isExist {
		return models.Thread{}, constants.NotFound
	}

	isExisted, err := s.repo.CheckForum(newThread.Forum)
	if err != nil {
		return models.Thread{}, http.StatusInternalServerError
	}
	if !isExisted {
		return models.Thread{}, constants.NotFound
	}

	isExisted, err = s.repo.CheckThread(newThread.Slug)
	if err != nil {
		return models.Thread{}, http.StatusInternalServerError
	}
	if isExisted {
		newThread2, err := s.repo.GetThread(newThread.Slug)
		if err != nil {
			return models.Thread{}, http.StatusInternalServerError
		}
		return newThread2, http.StatusConflict
	}

	newThread.Id, err = s.repo.AddNewThread(*newThread)
	if err != nil {
		return models.Thread{}, http.StatusInternalServerError
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

func (s Smth) CreateNewPosts(newPosts models.Posts, slugOrId string) (models.Posts, int) {
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
		if status == http.StatusNotFound {
			return models.Posts{}, status
		}
		if status == http.StatusInternalServerError {
			return models.Posts{}, http.StatusInternalServerError
		}
	}

	if len(newPosts) == 0 {
		return models.Posts{}, http.StatusCreated
	}

	now := time.Now()

	var err error
	for i := range newPosts {
		newPosts[i].Thread = int(thread.Id)
		newPosts[i].Forum = thread.Forum
		newPosts[i].Created = strfmt.DateTime(now)
		newPosts[i], err = s.repo.AddPost(newPosts[i])
		if err != nil {
			return models.Posts{}, http.StatusConflict
		}
		err = s.repo.IncrementPosts(newPosts[i].Forum)
	}

	return newPosts, http.StatusCreated
}

//Здесь может быть ошибка, потому что отдаем данные не ранее созданного форума, а те, которые пришли на вход(пространство для + рпс смотри выше)
func (s Smth) CreateNewForum(newForum *models.Forum) (models.Forum, int) {
	isExist, err := s.repo.CheckUser(newForum.Owner)
	if err != nil {
		return models.Forum{}, http.StatusInternalServerError
	}
	if !isExist {
		return models.Forum{}, constants.NotFound
	}

	isExisted, err := s.repo.CheckForum(newForum.Slug)
	if err != nil {
		return models.Forum{}, http.StatusInternalServerError
	}
	if isExisted {
		newForum.Posts, newForum.Threads, err = s.repo.GetForumCounts(newForum.Slug)
		if err != nil {
			return models.Forum{}, http.StatusInternalServerError
		}
		return *newForum, http.StatusConflict
	}
	err, _ = s.repo.AddNewForum(newForum)
	if err != nil {
		return models.Forum{}, http.StatusInternalServerError
	}
	s.repo.AddForumUsers(newForum.Slug, newForum.Owner)

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

func (s Smth) EditMessage(id int, message string) (models.Post, int) {
	isExisted, err := s.repo.CheckPost(id)
	if err != nil {
		return models.Post{}, http.StatusInternalServerError
	}
	if !isExisted {
		return models.Post{}, constants.NotFound
	}

	err = s.repo.EditMessage(id, message)
	if err != nil {
		return models.Post{}, http.StatusInternalServerError
	}

	post, _ := s.repo.GetPost(id)

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
	isExist, err := s.repo.CheckUser(nickname)
	if err != nil {
		return models.User{}, http.StatusInternalServerError
	}
	if !isExist {
		return models.User{}, http.StatusNotFound
	}

	isExist, err = s.repo.CheckUserByEmail(user.Email.String())
	if err != nil {
		return models.User{}, http.StatusInternalServerError
	}
	if isExist {
		return models.User{}, http.StatusConflict
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
		isExist, err := s.repo.CheckThread(slugOrId)
		if err != nil {
			return models.Thread{}, http.StatusInternalServerError
		}
		if !isExist {
			return models.Thread{}, http.StatusNotFound
		}
		thread, err = s.repo.UpdateThread(slugOrId, newThread)
		if err != nil {
			return models.Thread{}, http.StatusInternalServerError
		}
	} else {
		isExist, err := s.repo.CheckThreadById(id)
		if err != nil {
			return models.Thread{}, http.StatusInternalServerError
		}
		if !isExist {
			return models.Thread{}, http.StatusNotFound
		}
		thread, err = s.repo.UpdateThreadById(id, newThread)
	}

	return thread, http.StatusOK
}

func (s Smth) Vote(slugOrId string, vote models.Vote) (models.Thread, int) {
	var thread models.Thread
	isExist, err := s.repo.CheckUser(vote.Nickname)
	if err != nil {
		return models.Thread{}, http.StatusInternalServerError
	}
	if !isExist {
		return models.Thread{}, http.StatusNotFound
	}
	if id, err := strconv.Atoi(slugOrId); err != nil {
		isExist, err := s.repo.CheckThread(slugOrId)
		if err != nil {
			return models.Thread{}, http.StatusInternalServerError
		}
		if !isExist {
			return models.Thread{}, http.StatusNotFound
		}
		isExist, err = s.repo.CheckVote(slugOrId, vote.Nickname)
		if err != nil {
			return models.Thread{}, http.StatusInternalServerError
		}
		if !isExist {
			err = s.repo.AddVote(slugOrId, vote)
			if err != nil {
				return models.Thread{}, http.StatusInternalServerError
			}
		} else {
			num, err := s.repo.GetValueVote(slugOrId, vote.Nickname)
			if err != nil {
				return models.Thread{}, http.StatusInternalServerError
			}
			if num != vote.Voice {
				err = s.repo.UpdateVote(slugOrId, vote)
				if err != nil {
					return models.Thread{}, http.StatusInternalServerError
				}
			}
			thread, _ = s.repo.GetThread(slugOrId)
		}
	} else {
		thread, status := s.repo.GetThreadById(id)
		if status == constants.NotFound {
			return models.Thread{}, http.StatusNotFound
		}
		isExist, err = s.repo.CheckVote(thread.Slug, vote.Nickname)
		if err != nil {
			return models.Thread{}, http.StatusInternalServerError
		}
		if !isExist {
			err = s.repo.AddVote(thread.Slug, vote)
			if err != nil {
				return models.Thread{}, http.StatusInternalServerError
			}
		} else {
			num, err := s.repo.GetValueVote(thread.Slug, vote.Nickname)
			if err != nil {
				return models.Thread{}, http.StatusInternalServerError
			}
			if num != vote.Voice {
				err = s.repo.UpdateVote(slugOrId, vote)
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

/*func (e Event) GetNear(coord models.Coordinates, page int) (models.EventCardsWithCoords, error) {
	sqlEvents, err := e.repo.GetNearEvents(time.Now(), coord, page)
	if err != nil {
		e.logger.Warn(err)
		return models.EventCardsWithCoords{}, err
	}

	var pageEvents models.EventCardsWithCoords

	for i := range sqlEvents {
		pageEvents = append(pageEvents, models.ConvertCoordsCard(sqlEvents[i], coord))
	}
	if len(pageEvents) == 0 {
		e.logger.Debug("page" + fmt.Sprint(page) + "is empty")
		return models.EventCardsWithCoords{}, nil
	}

	return pageEvents, nil
}

func (e Event) GetAllEvents(page int) (models.EventCards, error) {
	sqlEvents, err := e.repo.GetAllEvents(time.Now(), page)
	if err != nil {
		e.logger.Warn(err)
		return models.EventCards{}, err
	}

	var pageEvents models.EventCards

	for i := range sqlEvents {
		pageEvents = append(pageEvents, models.ConvertDateCard(sqlEvents[i]))
	}
	if len(pageEvents) == 0 {
		e.logger.Debug("page" + fmt.Sprint(page) + "is empty")
		return models.EventCards{}, nil
	}

	return pageEvents, nil
}

func (e Event) GetOneEvent(eventId uint64) (models.Event, error) {
	ev, err := e.repo.GetOneEventByID(eventId)
	if err != nil {
		e.logger.Warn(err)
		return models.Event{}, err
	}

	jsonEvent := models.ConvertEvent(ev)

	tags, err := e.repo.GetTags(eventId)
	if err != nil {
		e.logger.Warn(err)
		return jsonEvent, err
	}

	jsonEvent.Tags = tags

	followers, err := e.repoSub.GetEventFollowers(eventId)
	if err != nil {
		e.logger.Warn(err)
		return jsonEvent, err
	}

	jsonEvent.Followers = followers

	return jsonEvent, nil
}

func (e Event) GetOneEventName(eventId uint64) (string, error) {
	name, err := e.repo.GetOneEventNameByID(eventId)
	if err != nil {
		e.logger.Warn(err)
		return "", err
	}

	return name, nil
}

func (e Event) Delete(eventId uint64) error {
	return e.repo.DeleteById(eventId)
}

func (e Event) CreateNewEvent(newEvent *models.Event) error {
	// TODO где-то здесь должна быть проверка на поля
	return e.repo.AddEvent(newEvent)
}

func (e Event) SaveImage(eventId uint64, img *multipart.FileHeader) error {
	src, err := img.Open()
	if err != nil {
		e.logger.Warn(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	fileName := constants.EventsPicDir + fmt.Sprint(eventId) + generator.RandStringRunes(6) + img.Filename

	dst, err := os.Create(fileName)
	if err != nil {
		e.logger.Warn(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		e.logger.Warn(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return e.repo.UpdateEventAvatar(eventId, fileName)
}

func (e Event) GetEventsByCategory(typeEvent string, page int) (models.EventCards, error) {
	var sqlEvents []models.EventCardWithDateSQL
	var err error
	if typeEvent == "" {
		sqlEvents, err = e.repo.GetAllEvents(time.Now(), page)
	} else {
		sqlEvents, err = e.repo.GetEventsByCategory(typeEvent, time.Now(), page)
	}
	if err != nil {
		e.logger.Warn(err)
		return models.EventCards{}, err
	}

	var pageEvents models.EventCards

	for i := range sqlEvents {
		pageEvents = append(pageEvents, models.ConvertDateCard(sqlEvents[i]))
	}
	if len(pageEvents) == 0 {
		e.logger.Debug("page" + fmt.Sprint(page) + "in category" + typeEvent + "empty")
		return models.EventCards{}, nil
	}

	return pageEvents, nil
}

func (e Event) GetImage(eventId uint64) ([]byte, error) {
	ev, err := e.repo.GetOneEventByID(eventId)
	if err != nil {
		e.logger.Warn(err)
		return []byte{}, err
	}

	file, err := ioutil.ReadFile(ev.Image.String)
	if err != nil {
		e.logger.Warn(errors.New("Cannot open file: " + ev.Image.String))
		return []byte{}, err
	}

	return file, nil
}

func (e Event) FindEvents(str string, category string, page int) (models.EventCards, error) {
	str = strings.ToLower(str)

	var sqlEvents []models.EventCardWithDateSQL
	var err error
	if category == "" {
		sqlEvents, err = e.repo.FindEvents(str, time.Now(), page)
	} else {
		sqlEvents, err = e.repo.CategorySearch(str, category, time.Now(), page)
	}
	if err != nil {
		e.logger.Warn(err)
		return models.EventCards{}, err
	}

	var pageEvents models.EventCards

	for i := range sqlEvents {
		pageEvents = append(pageEvents, models.ConvertDateCard(sqlEvents[i]))
	}

	if len(pageEvents) == 0 {
		e.logger.Debug("empty result for method FindEvents")
		return models.EventCards{}, nil
	}

	return pageEvents, nil
}

func (e Event) RecomendSystem(uid uint64, category string) error {
	if err := e.repo.RecomendSystem(uid, category); err != nil {
		time.Sleep(1 * time.Second)
		if err := e.repo.RecomendSystem(uid, category); err != nil {
			e.logger.Warn(err)
			return errors.New("cannot add record in user_prefer")
		}
	}
	return nil
}

func (e Event) GetRecommended(uid uint64, page int) (models.EventCards, error) {
	sqlEvents, err := e.repo.GetRecommended(uid, time.Now(), page)
	if err != nil {
		e.logger.Warn(err)
		return models.EventCards{}, err
	}

	var pageEvents models.EventCards

	for i := range sqlEvents {
		pageEvents = append(pageEvents, models.ConvertDateCard(sqlEvents[i]))
	}
	if len(pageEvents) == 0 {
		e.logger.Debug("empty result for method GetRecomended")
		return models.EventCards{}, nil
	}

	return pageEvents, nil
}*/
