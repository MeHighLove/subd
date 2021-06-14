package event

import "subd/models"

type Repository interface {
	CheckUser(user string) (bool, error)
	CheckUserByEmail(email string) (bool, error)
	CheckUserByNicknameOrEmail(nickname string, email string) (bool, error)
	AddNewForum(newForum *models.Forum) (error, bool)
	GetForumCounts(slug string) (uint64, uint64, error)
	GetForum(slug string) (models.Forum, int)
	CheckForum(slug string) (bool, error)
	CheckThread(slug string) (bool, error)
	CheckThreadById(id int) (bool, error)
	CheckPost(id int) (bool, error)
	GetThread(slug string) (models.Thread, error)
	GetThreadStatus(slug string) (models.Thread, int)
	GetThreadById(id int) (models.Thread, int)
	GetPost(id int) (models.Post, int)
	GetUser(name string) (models.User, int)
	AddNewThread(newThread models.Thread) (uint64, error)
	GetForumUsers(slug string, limit int, since string, desc bool) (models.Users, error)
	AddForumUsers(slug string, author string) error
	GetForumThreads(slug string, limit int, since string, desc bool) (models.Threads, error)
	EditMessage(id int, message string) error
	Clear() error
	Status() (models.Status, error)
	CreateUser(nickname string, user models.User) error
	GetUserByNicknameOrEmail(nickname string, email string) (models.Users, error)
	UpdateUser(nickname string, user models.User) error
	IncrementThreads(forum string) error
	IncrementPosts(forum string) error
	AddPost(post *models.Post) error
	UpdateThread(slugOrId string, thread models.Thread) (models.Thread, error)
	UpdateThreadById(id int, thread models.Thread) (models.Thread, error)
	CheckVote(id int, nickname string) (bool, error)
	AddVote(id int, vote models.Vote) error
	UpdateVote(id int, vote models.Vote) error
	GetValueVote(id int, nickname string) (int, error)
	GetPostsFlat(id int ,limit int, since int) (models.Posts, error)
	GetPostsFlatDesc(id int ,limit int, since int) (models.Posts, error)
	GetPostsTree(id int ,limit int) (models.Posts, error)
	GetPostsTreeDesc(id int ,limit int) (models.Posts, error)
	GetPostsTreeSince(id int ,limit int, since int) (models.Posts, error)
	GetPostsTreeSinceDesc(id int ,limit int, since int) (models.Posts, error)
	GetPostsParentTree(id int ,limit int) (models.Posts, error)
	GetPostsParentTreeDesc(id int ,limit int) (models.Posts, error)
	GetPostsParentTreeSince(id int ,limit int, since int) (models.Posts, error)
	GetPostsParentTreeSinceDesc(id int ,limit int, since int) (models.Posts, error)
	GetPostNull(id int) (models.PostNullMessage, int)
}
