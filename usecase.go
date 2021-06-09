package event

import (
	"subd/models"
)

//go:generate mockgen -destination=./mock/usecase_mock.go -package=mock -source=./application/event/usecase.go

type UseCase interface {
	/*GetAllEvents(page int) (models.EventCards, error)
	GetOneEvent(eventId uint64) (models.Event, error)
	GetOneEventName(eventId uint64) (string, error)
	Delete(eventId uint64) error
	CreateNewEvent(newEvent *models.Event) error
	SaveImage(eventId uint64, img *multipart.FileHeader) error
	GetEventsByCategory(typeEvent string, page int) (models.EventCards, error)
	GetImage(eventId uint64) ([]byte, error)
	FindEvents(str string, category string, page int) (models.EventCards, error)
	RecomendSystem(uid uint64, category string) error
	GetRecommended(uid uint64, page int) (models.EventCards, error)
	GetNear(coord models.Coordinates, page int) (models.EventCardsWithCoords, error)*/
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
	CreateNewPosts(newPosts models.Posts, slugOrId string) (models.Posts, int)
	GetThread(slugOrId string) (models.Thread, int)
	UpdateThread(slugOrId string, newThread models.Thread) (models.Thread, int)
}
