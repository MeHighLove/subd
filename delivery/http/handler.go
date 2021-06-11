package http

import (
	"fmt"
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

//TODO Разобраться с чеками(они все не нужны). Передавать можно все по ссылке.
func CreateSmthHandler(e *echo.Echo, uc smth.UseCase) {
	handler := SmthHandler{UseCase: uc}


	e.POST("/api/forum/create", handler.CreateForum)
	e.GET("/api/forum/:slug/details", handler.ForumDetails)
	e.POST("/api/forum/:slug/create", handler.CreateThread)
	e.GET("api/forum/:slug/users", handler.GetForumUsers)
	e.GET("/api/forum/:slug/threads", handler.GetThreads)
	e.GET("/api/post/:id/details", handler.GetPostDetails)
	e.POST("/api/post/:id/details", handler.EditMessage)
	e.POST("/api/service/clear", handler.Clear)
	e.GET("/api/service/status", handler.Status)
	e.POST("/api/thread/:slug_or_id/create", handler.CreatePosts)
	e.GET("/api/thread/:slug_or_id/details", handler.GetThreadDetails)
	e.POST("/api/thread/:slug_or_id/details", handler.UpdateThread)
	e.GET("/api/thread/:slug_or_id/posts", handler.GetThreadSort)
	e.POST("/api/thread/:slug_or_id/vote", handler.Vote)
	e.POST("/api/user/:nickname/create", handler.CreateUser)
	e.GET("/api/user/:nickname/profile", handler.GetUser)
	e.POST("/api/user/:nickname/profile", handler.UpdateUser)
}

func (sd SmthHandler) GetThreadSort(c echo.Context) error {
	defer c.Request().Body.Close()

	slugOrId := c.Param("slug_or_id")
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit == 0 {
		limit = 100
	}
	since, _ := strconv.Atoi(c.QueryParam("since"))
	desc, err := strconv.ParseBool(c.QueryParam("desc"))
	if err != nil {
		desc = false
	}
	sort := c.QueryParam("sort")

	var posts models.Posts
	var status int
	if sort == "tree" {
		posts, status = sd.UseCase.GetThreadSortTree(slugOrId, limit, since, desc)
		if status == constants.NotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Can't find post with id ")
		}

		return c.JSON(status, posts)
	}
	if sort == "parent_tree" {
		posts, status = sd.UseCase.GetThreadSortParentTree(slugOrId, limit, since, desc)
		if status == constants.NotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Can't find post with id ")
		}

		return c.JSON(status, posts)
	}

	posts, status = sd.UseCase.GetThreadSortFlat(slugOrId, limit, since, desc)
	if status == constants.NotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find post with id ")
	}

	return c.JSON(status, posts)
}

func (sd SmthHandler) Vote(c echo.Context) error {
	defer c.Request().Body.Close()

	slugOrId := c.Param("slug_or_id")

	vote := &models.Vote{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, vote); err != nil {
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	thread, status := sd.UseCase.Vote(slugOrId, *vote)

	if status == constants.NotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find thread with slug " + slugOrId)
	}

	return c.JSON(status, thread)
}

func (sd SmthHandler) UpdateThread(c echo.Context) error {
	defer c.Request().Body.Close()

	slugOrId := c.Param("slug_or_id")

	newThread := &models.Thread{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, newThread); err != nil {
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	thread, status := sd.UseCase.UpdateThread(slugOrId, *newThread)

	if status == http.StatusConflict {
		return echo.NewHTTPError(http.StatusConflict, "Can't find thread with slug " + slugOrId)
	}
	if status == constants.NotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find thread with slug " + slugOrId)
	}

	return c.JSON(status, thread)
}

func (sd SmthHandler) UpdateUser(c echo.Context) error {
	defer c.Request().Body.Close()

	nickname := c.Param("nickname")

	newUser := &models.User{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, newUser); err != nil {
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	user, status := sd.UseCase.UpdateUser(nickname, *newUser)

	if status == http.StatusConflict {
		return echo.NewHTTPError(http.StatusConflict, "Can't find user with nickname " + nickname)
	}
	if status == constants.NotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find user with nickname " + nickname)
	}

	return c.JSON(status, user)
}

func (sd SmthHandler) GetUser(c echo.Context) error {
	defer c.Request().Body.Close()

	nickname := c.Param("nickname")

	user, status := sd.UseCase.GetUser(nickname)

	if status == constants.NotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find user with nickname " + nickname)
	}

	return c.JSON(status, user)
}

func (sd SmthHandler) CreateUser(c echo.Context) error {
	defer c.Request().Body.Close()

	nickname := c.Param("nickname")

	newUser := &models.User{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, newUser); err != nil {
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	users, status := sd.UseCase.CreateUser(nickname, *newUser)

	if status == http.StatusConflict {
		return c.JSON(status, users)
	}

	return c.JSON(status, users[0])
}

func (sd SmthHandler) Status(c echo.Context) error {
	defer c.Request().Body.Close()

	status, err := sd.UseCase.Status()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, status)
}

func (sd SmthHandler) Clear(c echo.Context) error {
	defer c.Request().Body.Close()

	err := sd.UseCase.Clear()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}

func (sd SmthHandler) EditMessage(c echo.Context) error {
	defer c.Request().Body.Close()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	newMessage := &models.NewMessage{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, newMessage); err != nil {
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	if newMessage.Message == "" {
		post, status := sd.UseCase.EditMessageNull(id)
		if status == constants.NotFound {
			return echo.NewHTTPError(http.StatusNotFound, "Can't find post with id "+fmt.Sprint(id))
		}
		return c.JSON(status, post)
	}

	post, status := sd.UseCase.EditMessage(id, newMessage.Message)
	if status == constants.NotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find post with id "+fmt.Sprint(id))
	}
	if status == http.StatusConflict {
		postNull := models.ConvertPostToNullMessage(post)
		return c.JSON(http.StatusOK, postNull)
	}

	return c.JSON(status, post)
}

func (sd SmthHandler) GetThreadDetails(c echo.Context) error {
	defer c.Request().Body.Close()

	slugOrId := c.Param("slug_or_id")

	thread, status := sd.UseCase.GetThread(slugOrId)
	if status == constants.NotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find thread with id " + slugOrId)
	}

	return c.JSON(status, thread)
}

func (sd SmthHandler) GetPostDetails(c echo.Context) error {
	defer c.Request().Body.Close()

	related := c.QueryParam("related")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil{
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	post, status := sd.UseCase.GetPost(id, related)
	if status == constants.NotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find post with id " + fmt.Sprint(id))
	}

	return c.JSON(status, post)
}

func (sd SmthHandler) GetThreads(c echo.Context) error {
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

	threads, status := sd.UseCase.GetThreads(slug, limit, since, desc)
	if status == constants.NotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find forum with slug " + slug)
	}

	return c.JSON(status, threads)
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

	return c.JSON(status, users)
}

func (sd SmthHandler) ForumDetails(c echo.Context) error {
	defer c.Request().Body.Close()

	slug := c.Param("slug")

	forum, status := sd.UseCase.GetForum(slug)
	if status == constants.NotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find forum with slug " + slug)
	}

	return c.JSON(status, forum)
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

	return c.JSON(status, forum)
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

	return c.JSON(status, thread)
}

func (sd SmthHandler) CreatePosts(c echo.Context) error {
	defer c.Request().Body.Close()

	newPosts := &models.Posts{}

	if err := easyjson.UnmarshalFromReader(c.Request().Body, newPosts); err != nil {
		return echo.NewHTTPError(http.StatusTeapot, err.Error())
	}

	slugOrId := c.Param("slug_or_id")

	posts, status := sd.UseCase.CreateNewPosts(*newPosts, slugOrId)
	if status == constants.NotFound {
		return echo.NewHTTPError(http.StatusNotFound, "Can't find user with name ")
	}
	if status == http.StatusConflict {
		return echo.NewHTTPError(http.StatusConflict, "Can't find user with name ")
	}

	return c.JSON(status, posts)
}