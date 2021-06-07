package usecase

import (
	"net/http"
	"strings"
	smth "subd"
	"subd/constants"
	"subd/models"
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

	s.repo.AddForumUsers(newThread.Forum, newThread.Author)

	return *newThread, http.StatusCreated
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
