package event

import "subd/models"

type Repository interface {
	/*AddEvent(newEvent *models.Event) error
	GetAllEvents(now time.Time, page int) ([]models.EventCardWithDateSQL, error)
	GetNearEvents(now time.Time, coord models.Coordinates, page int) ([]models.EventCardWithCoordsSQL, error)
	GetOneEventByID(eventId uint64) (models.EventSQL, error)
	GetOneEventNameByID(eventId uint64) (string, error)
	DeleteById(eventId uint64) error
	GetTags(eventId uint64) (models.Tags, error)
	UpdateEventAvatar(eventId uint64, path string) error
	GetEventsByCategory(typeEvent string, now time.Time, page int) ([]models.EventCardWithDateSQL, error)
	FindEvents(str string, now time.Time, page int) ([]models.EventCardWithDateSQL, error)
	RecomendSystem(uid uint64, category string) error
	GetPreference(uid uint64) (models.Recomend, error)
	GetRecommended(uid uint64, now time.Time, page int) ([]models.EventCardWithDateSQL, error)
	CategorySearch(str string, category string, now time.Time, page int) ([]models.EventCardWithDateSQL, error)*/
	CheckUser(user string) (bool, error)
	CheckUserByEmail(email string) (bool, error)
	CheckUserByNicknameOrEmail(nickname string, email string) (bool, error)
	AddNewForum(newForum *models.Forum) (error, bool)
	GetForumCounts(slug string) (uint64, uint64, error)
	GetForum(slug string) (models.Forum, int)
	CheckForum(slug string) (bool, error)
	CheckThread(slug string) (bool, error)
	CheckPost(id int) (bool, error)
	GetThread(slug string) (models.Thread, error)
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
}
