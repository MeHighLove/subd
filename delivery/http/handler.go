package http

import (
	"github.com/labstack/echo"
	"github.com/mailru/easyjson"
	"net/http"
	"strconv"
	smth "subd"
	"subd/constants"
	"subd/models"
)

type SmthHandler struct {
	UseCase   smth.UseCase
}

func CreateSmthHandler(e *echo.Echo, uc smth.UseCase) {
	handler := SmthHandler{UseCase: uc}


	e.POST("/api/forum/create", handler.CreateForum)
	e.GET("/api/:slug/details", handler.ForumDetails)
	e.POST("/api/:slug/create", handler.CreateThread)
	e.GET("api/:slug/users", handler.GetForumUsers)
	/*e.GET("/api/v1/", eventHandler.GetAllEvents, middleware.GetPage)
	e.GET("/api/v1/event/:id", eventHandler.GetOneEvent, middleware.GetId)
	e.GET("/link/event/:id", eventHandler.GetEventLink, middleware.GetId)
	e.GET("/api/v1/event/name/:id", eventHandler.GetOneEventName, middleware.GetId)
	e.POST("/api/v1/near", eventHandler.GetNear, middleware.GetPage)
	e.GET("/api/v1/event", eventHandler.GetEvents, middleware.GetPage)
	e.GET("/api/v1/search", eventHandler.FindEvents, middleware.GetPage)
	//create & delete & save вообще не должно быть, пользователь НИКАК не может создавать и удалять что-либо, только админ работает с БД
	e.POST("/api/v1/create", eventHandler.Create)
	e.DELETE("/api/v1/event/:id", eventHandler.Delete, middleware.GetId)
	e.POST("/api/v1/save/:id", eventHandler.Save, middleware.GetId)
	e.GET("api/v1/event/:id/image", eventHandler.GetImage, middleware.GetId)
	e.GET("/api/v1/recommend", eventHandler.Recommend, middleware.GetPage, auth.GetSession)*/
}

func (sd SmthHandler) GetForumUsers(c echo.Context) error {
	defer c.Request().Body.Close()

	slug := c.Param("slug")
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit == 0 {
		limit = 100
	}
	since := c.QueryParam("since")
	desc, err := strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		desc = false
	}

	users, status := sd.UseCase.GetForumUsers(slug, limit, since, desc)
	if status == constants.NotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find forum with slug " + slug)
	}

	if _, err := easyjson.MarshalToWriter(users, c.Response().Writer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Response().Status = status

	return nil
}

func (sd SmthHandler) ForumDetails(c echo.Context) error {
	defer c.Request().Body.Close()

	slug := c.Param("slug")

	forum, status := sd.UseCase.GetForum(slug)
	if status == constants.NotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find forum with slug " + slug)
	}

	if _, err := easyjson.MarshalToWriter(forum, c.Response().Writer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Response().Status = status

	return nil
}

func (sd SmthHandler) CreateForum(c echo.Context) error {
	defer c.Request().Body.Close()

	newForum := &models.Forum{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, newForum); err != nil {
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	forum, status := sd.UseCase.CreateNewForum(newForum)
	if status == constants.NotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find user with name " + newForum.Owner)
	}

	if _, err := easyjson.MarshalToWriter(forum, c.Response().Writer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Response().Status = status

	return nil
}

func (sd SmthHandler) CreateThread(c echo.Context) error {
	defer c.Request().Body.Close()

	newThread := &models.Thread{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, newThread); err != nil {
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	newThread.Forum = c.Param("slug")

	thread, status := sd.UseCase.CreateNewThread(newThread)
	if status == constants.NotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find user with name " + newThread.Author)
	}

	if _, err := easyjson.MarshalToWriter(thread, c.Response().Writer); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.Response().Status = status

	return nil
}

/*func (eh EventHandler) GetNear(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	page := c.Get(constants.PageKey).(int)

	coord := &models.Coordinates{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, coord); err != nil {
		middleware.ErrResponse(c, http.StatusTeapot)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	events, err := eh.UseCase.GetNear(*coord, page)
	events = eh.sanitizer.SanitizeEventCardsWithCoords(events)
	if err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		return err
	}

	if _, err = easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	eh.Logger.LogInfo(c, start, requestId)
	middleware.OkResponse(c)
	return nil
}

func (eh EventHandler) Recommend(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	page := c.Get(constants.PageKey).(int)
	uid := c.Get(constants.UserIdKey).(uint64)

	events, err := eh.UseCase.GetRecommended(uid, page)
	events = eh.sanitizer.SanitizeEventCards(events)
	if err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		return err
	}

	if _, err = easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	eh.Logger.LogInfo(c, start, requestId)
	middleware.OkResponse(c)
	return nil
}

func (eh EventHandler) GetAllEvents(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	page := c.Get(constants.PageKey).(int)

	events, err := eh.UseCase.GetAllEvents(page)
	events = eh.sanitizer.SanitizeEventCards(events)
	if err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		return err
	}

	if _, err = easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	eh.Logger.LogInfo(c, start, requestId)
	middleware.OkResponse(c)
	return nil
}

func (eh EventHandler) GetUserID(c echo.Context) (uint64, error) {
	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	cookie, err := c.Cookie(constants.SessionCookieName)
	if err != nil && cookie != nil {
		eh.Logger.LogWarn(c, start, requestId, err)
		middleware.ErrResponse(c, http.StatusForbidden)
		return 0, errors.New("user is not authorized")
	}

	var uid uint64
	var exists bool
	var code int

	if cookie != nil {
		exists, uid, err, code = eh.rpcAuth.Check(cookie.Value)
		if err != nil {
			eh.Logger.LogWarn(c, start, requestId, err)
			middleware.ErrResponse(c, code)
			return 0, err
		}

		if !exists {
			eh.Logger.LogWarn(c, start, requestId, err)
			middleware.ErrResponse(c, http.StatusForbidden)
			return 0, errors.New("user is not authorized")
		}
		middleware.OkResponse(c)
		return uid, nil
	}
	eh.Logger.LogWarn(c, start, requestId, err)
	middleware.ErrResponse(c, http.StatusForbidden)
	return 0, errors.New("user is not authorized")
}

func (eh EventHandler) GetOneEvent(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	id := c.Get(constants.IdKey).(int)

	ev, err := eh.UseCase.GetOneEvent(uint64(id))
	if err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return err
	}
	eh.sanitizer.SanitizeEvent(&ev)
	if uid, err := eh.GetUserID(c); err == nil {
		if err := eh.UseCase.RecomendSystem(uid, ev.Category); err != nil {
			eh.Logger.LogWarn(c, start, requestId, err)
		}
	} else {
		eh.Logger.LogWarn(c, start, requestId, err)
	}

	if _, err = easyjson.MarshalToWriter(ev, c.Response().Writer); err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	eh.Logger.LogInfo(c, start, requestId)
	middleware.OkResponse(c)
	return nil
}

func (eh EventHandler) GetOneEventName(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	id := c.Get(constants.IdKey).(int)

	name, err := eh.UseCase.GetOneEventName(uint64(id))
	if err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return err
	}
	sanName := eh.sanitizer.SanitizeEventName(name)

	eh.Logger.LogInfo(c, start, requestId)
	middleware.OkResponse(c)
	return c.String(http.StatusOK, sanName)
}

func (eh EventHandler) GetEvents(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	category := c.QueryParam("category")
	page := c.Get(constants.PageKey).(int)

	events, err := eh.UseCase.GetEventsByCategory(category, page)
	events = eh.sanitizer.SanitizeEventCards(events)

	if err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if _, err := easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	eh.Logger.LogInfo(c, start, requestId)
	middleware.OkResponse(c)
	return nil
}

//Эти функции будут удалены, поэтому почти не изменялись с переноса архитектуры
func (eh EventHandler) Create(c echo.Context) error {
	defer c.Request().Body.Close()

	newEvent := &models.Event{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, newEvent); err != nil {
		middleware.ErrResponse(c, http.StatusTeapot)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	if err := eh.UseCase.CreateNewEvent(newEvent); err != nil {
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return err
	}

	middleware.OkResponse(c)
	return c.JSON(http.StatusOK, *newEvent)
}

func (eh EventHandler) Delete(c echo.Context) error {
	defer c.Request().Body.Close()

	id := c.Get(constants.IdKey).(int)

	err := eh.UseCase.Delete(uint64(id))
	if err != nil {
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return err
	}

	middleware.OkResponse(c)
	return c.String(http.StatusOK, "Event with id "+fmt.Sprint(id)+" successfully deleted \n")
}

func (eh EventHandler) Save(c echo.Context) error {
	defer c.Request().Body.Close()

	id := c.Get(constants.IdKey).(int)

	img, err := c.FormFile("image")
	if err != nil {
		middleware.ErrResponse(c, http.StatusTeapot)
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	err = eh.UseCase.SaveImage(uint64(id), img)

	if err != nil {
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return err
	}
	middleware.OkResponse(c)
	return c.JSON(http.StatusOK, "Picture changed successfully")
}

func (eh EventHandler) GetImage(c echo.Context) error {
	defer c.Request().Body.Close()

	id := c.Get(constants.IdKey).(int)

	file, err := eh.UseCase.GetImage(uint64(id))
	if err != nil {
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return err
	}

	_, err = c.Response().Write(file)
	if err != nil {
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return err
	}
	middleware.OkResponse(c)
	return nil
}

func (eh EventHandler) FindEvents(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	str := c.QueryParam("find")
	category := c.QueryParam("category")
	page := c.Get(constants.PageKey).(int)

	events, err := eh.UseCase.FindEvents(str, category, page)
	events = eh.sanitizer.SanitizeEventCards(events)
	if err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		return err
	}

	if _, err := easyjson.MarshalToWriter(events, c.Response().Writer); err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	eh.Logger.LogInfo(c, start, requestId)
	middleware.OkResponse(c)
	return nil
}

func (eh EventHandler) GetEventLink(c echo.Context) error {
	defer c.Request().Body.Close()

	start := time.Now()
	requestId := fmt.Sprintf("%016x", rand.Int())
	id := c.Get(constants.IdKey).(int)

	ev, err := eh.UseCase.GetOneEvent(uint64(id))
	if err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		middleware.ErrResponse(c, http.StatusInternalServerError)
		return err
	}

	data, err := ioutil.ReadFile("2021_1_Fyvaoldzh/dist/index.html")
	if err != nil {
		eh.Logger.LogError(c, start, requestId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	input := models.ViewData{
		Id:    ev.ID,
		Title: ev.Title,
	}

	tmpl, _ := template.New("input").Parse(string(data))
	return tmpl.Execute(c.Response(), input)
}
*/