package event

import (
	"subd/models"
)

//go:generate mockgen -destination=./mock/usecase_mock.go -package=mock -source=./application/event/usecase.go

type UseCase interface {
	CreateNewForum(newForum *models.Forum) (models.Forum, int)
	CreateNewThread(newThread *models.Thread) (models.Thread, int)
	GetForum(slug string) (models.Forum, int)
	GetForumUsers(slug string, limit int, since string, desc bool) (models.Users, int)
	GetThreads(slug string, limit int, since string, desc bool) (models.Threads, int)
	GetPost(id int, related string) (models.FullPost, int)
	EditMessage(id int, message string) (models.Post, int)
	Clear() error
	Status() (models.Status, error)
	CreateUser(nickname string, user models.User) (models.Users, int)
	GetUser(nickname string) (models.User, int)
	UpdateUser(nickname string, user models.User) (models.User, int)
	CreateNewPosts(newPosts []*models.Post, slugOrId string)  int
	GetThread(slugOrId string) (models.Thread, int)
	UpdateThread(slugOrId string, newThread models.Thread) (models.Thread, int)
	Vote(slugOrId string, vote models.Vote) (models.Thread, int)
	GetThreadSortFlat(slugOrId string, limit int, since int, desc bool) (models.Posts, int)
	GetThreadSortTree(slugOrId string, limit int, since int, desc bool) (models.Posts, int)
	GetThreadSortParentTree(slugOrId string, limit int, since int, desc bool) (models.Posts, int)
	EditMessageNull(id int) (models.PostNullMessage, int)
}
