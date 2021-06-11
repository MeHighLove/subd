package models

import (
	"database/sql"
	"github.com/go-openapi/strfmt"
)

type Vote struct {
	Nickname string `json:"nickname"`
	Voice int `json:"voice"`
}

type PostNullMessage struct {
	Author   string          `json:"author"`
	Created  strfmt.DateTime `json:"created"`
	Forum    string          `json:"forum"`
	Id       int           `json:"id"`
	Message  string          `json:"message"`
	Parent   int           `json:"parent"`
	Thread   int           `json:"thread"`
}

type Forum struct {
	Title string `json:"title"`
	Owner string `json:"user"`
	Posts uint64 `json:"posts"`
	Threads uint64 `json:"threads"`
	Slug string `json:"slug"`
}

type FullPost struct {
	Author *User   `json:"author"`
	Forum  *Forum  `json:"forum"`
	Post   *Post   `json:"post"`
	Thread *Thread `json:"thread"`
}

type Status struct {
	Forum uint64 `json:"forum"`
	Post uint64 `json:"post"`
	Thread uint64 `json:"thread"`
	User uint64 `json:"user"`
}

type NewMessage struct {
	Message string `json:"message"`
}

type Post struct {
	Author   string          `json:"author"`
	Created  strfmt.DateTime `json:"created"`
	Forum    string          `json:"forum"`
	Id       int           `json:"id"`
	IsEdited bool            `json:"isEdited"`
	Message  string          `json:"message"`
	Parent   int           `json:"parent"`
	Thread   int           `json:"thread"`
}

type ThreadSQL struct {
	Id uint64 `json:"id"`
	Author string `json:"author"`
	Created strfmt.DateTime `json:"created"`
	Forum string `json:"forum"`
	Message string `json:"message"`
	Slug sql.NullString `json:"slug"`
	Title string `json:"title"`
	Votes int `json:"votes"`
}

type Thread struct {
	Id uint64 `json:"id"`
	Author string `json:"author"`
	Created strfmt.DateTime `json:"created"`
	Forum string `json:"forum"`
	Message string `json:"message"`
	Slug string `json:"slug"`
	Title string `json:"title"`
	Votes int `json:"votes"`
}

type User struct {
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About string `json:"about"`
	Email strfmt.Email `json:"email"`
}

//easyjson:json
type Users []User

//easyjson:json
type Threads []Thread

//easyjson:json
type Posts []Post

func ConvertPostToNullMessage(post Post) (PostNullMessage) {
	var newPost PostNullMessage
	newPost.Message = post.Message
	newPost.Id = post.Id
	newPost.Forum = post.Forum
	newPost.Thread = post.Thread
	newPost.Author = post.Author
	newPost.Created = post.Created
	newPost.Parent = post.Parent
	return newPost
}

func ConvertThread(old ThreadSQL) (Thread) {
	var newThread Thread
	newThread.Author = old.Author
	newThread.Id = old.Id
	newThread.Forum = old.Forum
	newThread.Slug = old.Slug.String
	newThread.Title = old.Title
	newThread.Created = old.Created
	newThread.Message = old.Message
	newThread.Votes = old.Votes
	return newThread
}