package repository

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/go-openapi/strfmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
	event "subd"
	"subd/models"
	"time"
)

type SomeDatabase struct {
	pool   *pgxpool.Pool
}

func NewSomeDatabase(conn *pgxpool.Pool) event.Repository {
	return &SomeDatabase{pool: conn}
}

func (sd SomeDatabase) CheckThread(slug string) (bool, error) {
	var id []uint64
	err := pgxscan.Select(context.Background(), sd.pool, &id,
		`SELECT 1 FROM threads
	WHERE slug = $1 LIMIT 1`, slug)

	if errors.As(err, &pgx.ErrNoRows) || len(id) == 0 {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (sd SomeDatabase) CheckThreadById(id int) (bool, error) {
	var ids []uint64
	err := pgxscan.Select(context.Background(), sd.pool, &ids,
		`SELECT 1 FROM threads
	WHERE id = $1 LIMIT 1`, id)

	if errors.As(err, &pgx.ErrNoRows) || len(ids) == 0 {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (sd SomeDatabase) CheckForum(slug string) (bool, error) {
	var id []uint64
	err := pgxscan.Select(context.Background(), sd.pool, &id,
		`SELECT 1 FROM forums
	WHERE slug = $1 LIMIT 1`, slug)

	if errors.As(err, &pgx.ErrNoRows) || len(id) == 0 {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (sd SomeDatabase) CheckUser(user string) (bool, error) {
	var id []uint64
	err := pgxscan.Select(context.Background(), sd.pool, &id,
		`SELECT 1 FROM users
	WHERE nickname = $1 LIMIT 1`, user)

	if errors.As(err, &pgx.ErrNoRows) || len(id) == 0 {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (sd SomeDatabase) CheckUserByEmail(email string) (bool, error) {
	var id []uint64
	err := pgxscan.Select(context.Background(), sd.pool, &id,
		`SELECT 1 FROM users
	WHERE email = $1 LIMIT 1`, email)

	if errors.As(err, &pgx.ErrNoRows) || len(id) == 0 {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (sd SomeDatabase) CheckUserByNicknameOrEmail(nickname string, email string) (bool, error) {
	var id []uint64
	err := pgxscan.Select(context.Background(), sd.pool, &id,
		`SELECT 1 FROM users
	WHERE nickname = $1 OR email = $2 LIMIT 1`, nickname, email)

	if errors.As(err, &pgx.ErrNoRows) || len(id) == 0 {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (sd SomeDatabase) CheckPost(id int) (bool, error) {
	var ids []uint64
	err := pgxscan.Select(context.Background(), sd.pool, &ids,
		`SELECT 1 FROM posts
	WHERE id = $1 LIMIT 1`, id)

	if errors.As(err, &pgx.ErrNoRows) || len(ids) == 0 {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (sd SomeDatabase) CheckVote(id int, nickname string) (bool, error) {
	var ids []uint64
	err := pgxscan.Select(context.Background(), sd.pool, &ids,
		`SELECT 1 FROM votes
	WHERE thread = $1 AND nickname = $2 LIMIT 1`, id, nickname)

	if errors.As(err, &pgx.ErrNoRows) || len(ids) == 0 {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (sd SomeDatabase) AddNewForum(newForum *models.Forum) (error, bool) {
	resp, err := sd.pool.Exec(context.Background(),
		`INSERT INTO forums 
		VALUES (default, $1, $2, default, default, $3)`,
		newForum.Title, newForum.Owner, newForum.Slug)
	if err != nil {
		return err, false
	}
	if resp.RowsAffected() == 0 {
		return nil, true
	}

	return nil, false
}

func (sd SomeDatabase) AddVote(id int, vote models.Vote) error {
	_, err := sd.pool.Exec(context.Background(),
		`INSERT INTO votes 
		VALUES ($1, $2, $3)`,
		id, vote.Voice, vote.Nickname)
	if err != nil {
		return err
	}

	return nil
}

func (sd SomeDatabase) GetForumCounts(slug string) (uint64, uint64, error) {
	var posts, threads uint64
	err := sd.pool.QueryRow(context.Background(),
			`SELECT posts, threads FROM forums WHERE slug = $1`, slug).Scan(&posts, &threads)

	if err != nil {
		return posts, threads, err
	}

	return posts, threads, nil
}

func (sd SomeDatabase) GetForum(slug string) (models.Forum, int) {
	var forum []models.Forum
	err := pgxscan.Select(context.Background(), sd.pool, &forum,
		`SELECT title, owner, posts, threads, slug FROM forums WHERE slug = $1`, slug)

	if errors.As(err, &pgx.ErrNoRows) || len(forum) == 0 {
		return models.Forum{}, http.StatusNotFound
	}

	if err != nil {
		return models.Forum{}, http.StatusInternalServerError
	}

	return forum[0], http.StatusOK
}

func (sd SomeDatabase) GetPost(id int) (models.Post, int) {
	var post []models.Post
	err := pgxscan.Select(context.Background(), sd.pool, &post,
		`SELECT id, author, created, forum, is_edited, message, parent, thread FROM posts WHERE id = $1`, id)

	if errors.As(err, &pgx.ErrNoRows) || len(post) == 0 {
		return models.Post{}, http.StatusNotFound
	}

	if err != nil {
		return models.Post{}, http.StatusInternalServerError
	}

	return post[0], http.StatusOK
}

func (sd SomeDatabase) GetPostNull(id int) (models.PostNullMessage, int) {
	var post []models.PostNullMessage
	err := pgxscan.Select(context.Background(), sd.pool, &post,
		`SELECT id, author, created, forum, message, parent, thread FROM posts WHERE id = $1`, id)

	if errors.As(err, &pgx.ErrNoRows) || len(post) == 0 {
		return models.PostNullMessage{}, http.StatusNotFound
	}

	if err != nil {
		return models.PostNullMessage{}, http.StatusInternalServerError
	}

	return post[0], http.StatusOK
}

func (sd SomeDatabase) GetUser(name string) (models.User, int) {
	var user []models.User
	err := pgxscan.Select(context.Background(), sd.pool, &user,
		`SELECT nickname, fullname, about, email FROM users WHERE nickname = $1`, name)

	if errors.As(err, &pgx.ErrNoRows) || len(user) == 0 {
		return models.User{}, http.StatusNotFound
	}

	if err != nil {
		return models.User{}, http.StatusInternalServerError
	}

	return user[0], http.StatusOK
}

func (sd SomeDatabase) GetThreadById(id int) (models.Thread, int) {
	var thread []models.ThreadSQL
	err := pgxscan.Select(context.Background(), sd.pool, &thread,
		`SELECT * FROM threads WHERE id = $1`, id)

	if errors.As(err, &pgx.ErrNoRows) || len(thread) == 0 {
		return models.Thread{}, http.StatusNotFound
	}

	if err != nil {
		return models.Thread{}, http.StatusInternalServerError
	}

	return models.ConvertThread(thread[0]), http.StatusOK
}

func (sd SomeDatabase) GetThreadStatus(slug string) (models.Thread, int) {
	var thread []models.ThreadSQL
	err := pgxscan.Select(context.Background(), sd.pool, &thread,
		`SELECT * FROM threads WHERE slug = $1`, slug)

	if errors.As(err, &pgx.ErrNoRows) || len(thread) == 0 {
		return models.Thread{}, http.StatusNotFound
	}

	if err != nil {
		return models.Thread{}, http.StatusInternalServerError
	}

	return models.ConvertThread(thread[0]), http.StatusOK
}

func (sd SomeDatabase) GetThread(slug string) (models.Thread, error) {
	var ev []models.ThreadSQL
	err := pgxscan.Select(context.Background(), sd.pool, &ev,
		`SELECT * FROM threads WHERE slug = $1`, slug)

	if errors.As(err, &pgx.ErrNoRows) || len(ev) == 0 {
		return models.Thread{}, nil
	}

	if err != nil {
		return models.Thread{}, err
	}

	return models.ConvertThread(ev[0]), nil
}

func (sd SomeDatabase) GetValueVote(id int, nickname string) (int, error) {
	var num []int
	err := pgxscan.Select(context.Background(), sd.pool, &num,
		`SELECT voice FROM votes WHERE thread = $1 AND nickname = $2`, id, nickname)

	if err != nil {
		return 0, err
	}

	return num[0], nil
}

func (sd SomeDatabase) AddNewThread(newThread models.Thread) (uint64, error) {
	var id uint64
	var err error
	if newThread.Slug == ""{
		err = sd.pool.QueryRow(context.Background(),
			`INSERT INTO threads VALUES (default, $1, $2, $3, $4, null, $5, default) RETURNING id`,
			newThread.Author, newThread.Created, newThread.Forum, newThread.Message,
			newThread.Title).Scan(&id)
	} else {
		err = sd.pool.QueryRow(context.Background(),
			`INSERT INTO threads VALUES (default, $1, $2, $3, $4, $5, $6, default) RETURNING id`,
			newThread.Author, newThread.Created, newThread.Forum, newThread.Message,
			newThread.Slug, newThread.Title).Scan(&id)
	}
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (sd SomeDatabase) AddPost(newPosts []*models.Post, thread models.Thread, now time.Time) int {
	for i := range newPosts {
		newPosts[i].Thread = int(thread.Id)
		newPosts[i].Forum = thread.Forum
		newPosts[i].Created = strfmt.DateTime(now)
		err := sd.pool.QueryRow(context.Background(),
			`INSERT INTO posts VALUES (default, $1, $2, $3, default, $4, $5, $6) RETURNING id`,
			newPosts[i].Author, newPosts[i].Created, newPosts[i].Forum,
			newPosts[i].Message, newPosts[i].Parent, newPosts[i].Thread).Scan(&newPosts[i].Id)
		if err != nil {
			if err.Error() == "ERROR: 00404 (SQLSTATE 00404)" {
				return http.StatusConflict
			} else {
				return http.StatusNotFound
			}
		}

		_, err = sd.pool.Exec(context.Background(),
			`UPDATE forums SET posts = posts + 1 WHERE slug = $1`, newPosts[i].Forum)
		_, err = sd.pool.Exec(context.Background(),
			`INSERT INTO forum_users 
		VALUES ($1, $2)`,
			newPosts[i].Forum, newPosts[i].Author)
	}

	return http.StatusCreated
}

func (sd SomeDatabase) GetForumUsers(slug string, limit int, since string, desc bool) (models.Users, error) {
	var users models.Users
	var err error
	if since != "" {
		if desc == true {
			err = pgxscan.Select(context.Background(), sd.pool, &users,
				`SELECT users.nickname, users.fullname, users.email, users.about FROM forum_users JOIN users
			ON forum_users.nickname = users.nickname
			WHERE forum_users.forum = $1 AND users.nickname < $2 
			ORDER BY users.nickname DESC LIMIT $3`, slug, since, limit)
		} else {
			err = pgxscan.Select(context.Background(), sd.pool, &users,
				`SELECT users.nickname, users.fullname, users.email, users.about FROM forum_users JOIN users
			ON forum_users.nickname = users.nickname
			WHERE forum_users.forum = $1 AND users.nickname > $2 
			ORDER BY users.nickname LIMIT $3`, slug, since, limit)
		}
	} else {
		if desc == true {
			err = pgxscan.Select(context.Background(), sd.pool, &users,
				`SELECT users.nickname, users.fullname, users.email, users.about FROM forum_users JOIN users
			ON forum_users.nickname = users.nickname
			WHERE forum_users.forum = $1 
			ORDER BY users.nickname DESC LIMIT $2`, slug, limit)
		} else {
			err = pgxscan.Select(context.Background(), sd.pool, &users,
				`SELECT users.nickname, users.fullname, users.email, users.about FROM forum_users JOIN users
			ON forum_users.nickname = users.nickname
			WHERE forum_users.forum = $1 
			ORDER BY users.nickname LIMIT $2`, slug, limit)
		}
	}
	if errors.As(err, &pgx.ErrNoRows) || len(users) == 0 {
		return models.Users{}, nil
	}

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (sd SomeDatabase) GetPostsFlat(id int ,limit int, since int) (models.Posts, error) {
	var posts models.Posts
		err := pgxscan.Select(context.Background(), sd.pool, &posts,
			`SELECT id, author, created, forum, is_edited, message, parent, thread 
			FROM posts WHERE thread = $1 AND id > $2 
			ORDER BY created, id LIMIT $3`, id, since, limit)

	if errors.As(err, &pgx.ErrNoRows) || len(posts) == 0 {
		return models.Posts{}, nil
	}

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (sd SomeDatabase) GetPostsFlatDesc(id int ,limit int, since int) (models.Posts, error) {
	var posts models.Posts
	var err error
	if since != 0 {
		err = pgxscan.Select(context.Background(), sd.pool, &posts,
			`SELECT id, author, created, forum, is_edited, message, parent, thread
			FROM posts WHERE thread = $1 AND id < $2 
			ORDER BY created DESC, id DESC LIMIT $3`, id, since, limit)
	} else {
		err = pgxscan.Select(context.Background(), sd.pool, &posts,
			`SELECT id, author, created, forum, is_edited, message, parent, thread
			FROM posts WHERE thread = $1 
			ORDER BY created DESC, id DESC LIMIT $2`, id, limit)
	}

	if errors.As(err, &pgx.ErrNoRows) || len(posts) == 0 {
		return models.Posts{}, nil
	}

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (sd SomeDatabase) GetPostsParentTree(id int ,limit int) (models.Posts, error) {
	var posts models.Posts
	err := pgxscan.Select(context.Background(), sd.pool, &posts,
		`SELECT posts.id, posts.author, posts.created, posts.forum,
			posts.is_edited, posts.message, posts.parent, posts.thread 
			FROM (SELECT * FROM posts a WHERE a.parent = 0 AND a.thread = $1
			ORDER BY a.path LIMIT $2) AS b
			JOIN posts ON b.path[1] = posts.path[1]
			ORDER BY posts.path[1], posts.path`, id, limit)

	if errors.As(err, &pgx.ErrNoRows) || len(posts) == 0 {
		return models.Posts{}, nil
	}

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (sd SomeDatabase) GetPostsParentTreeDesc(id int ,limit int) (models.Posts, error) {
	var posts models.Posts
	err := pgxscan.Select(context.Background(), sd.pool, &posts,
		`SELECT posts.id, posts.author, posts.created, posts.forum,
			posts.is_edited, posts.message, posts.parent, posts.thread 
			FROM (SELECT * FROM posts a WHERE a.parent = 0 AND a.thread = $1
			ORDER BY a.path DESC LIMIT $2) AS b
			JOIN posts ON b.path[1] = posts.path[1]
			ORDER BY posts.path[1] DESC, posts.path`, id, limit)

	if errors.As(err, &pgx.ErrNoRows) || len(posts) == 0 {
		return models.Posts{}, nil
	}

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (sd SomeDatabase) GetPostsParentTreeSince(id int ,limit int, since int) (models.Posts, error) {
	var posts models.Posts
	err := pgxscan.Select(context.Background(), sd.pool, &posts,
		`SELECT posts.id, posts.author, posts.created, posts.forum,
			posts.is_edited, posts.message, posts.parent, posts.thread 
			FROM (SELECT * FROM posts a WHERE a.parent = 0 AND a.thread = $1
			AND a.path[1] > (SELECT path[1] FROM posts WHERE id = $2)
			ORDER BY a.path LIMIT $3) AS b
			JOIN posts ON b.path[1] = posts.path[1]
			ORDER BY posts.path[1], posts.path`, id, since, limit)

	if errors.As(err, &pgx.ErrNoRows) || len(posts) == 0 {
		return models.Posts{}, nil
	}

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (sd SomeDatabase) GetPostsParentTreeSinceDesc(id int ,limit int, since int) (models.Posts, error) {
	var posts models.Posts
	err := pgxscan.Select(context.Background(), sd.pool, &posts,
		`SELECT posts.id, posts.author, posts.created, posts.forum,
			posts.is_edited, posts.message, posts.parent, posts.thread 
			FROM (SELECT * FROM posts a WHERE a.parent = 0 AND a.thread = $1
			AND a.path[1] < (SELECT path[1] FROM posts WHERE id = $2)
			ORDER BY a.path DESC LIMIT $3) AS b
			JOIN posts ON b.path[1] = posts.path[1]
			ORDER BY posts.path[1] DESC, posts.path`, id, since, limit)

	if errors.As(err, &pgx.ErrNoRows) || len(posts) == 0 {
		return models.Posts{}, nil
	}

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (sd SomeDatabase) GetPostsTree(id int ,limit int) (models.Posts, error) {
	var posts models.Posts
	err := pgxscan.Select(context.Background(), sd.pool, &posts,
		`SELECT id, author, created, forum, is_edited, message, parent, thread
			FROM posts WHERE thread = $1
			ORDER BY path LIMIT $2`, id, limit)

	if errors.As(err, &pgx.ErrNoRows) || len(posts) == 0 {
		return models.Posts{}, nil
	}

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (sd SomeDatabase) GetPostsTreeDesc(id int ,limit int) (models.Posts, error) {
	var posts models.Posts
	err := pgxscan.Select(context.Background(), sd.pool, &posts,
		`SELECT id, author, created, forum, is_edited, message, parent, thread
			FROM posts WHERE thread = $1
			ORDER BY path DESC LIMIT $2`, id, limit)

	if errors.As(err, &pgx.ErrNoRows) || len(posts) == 0 {
		return models.Posts{}, nil
	}

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (sd SomeDatabase) GetPostsTreeSince(id int ,limit int, since int) (models.Posts, error) {
	var posts models.Posts
	err := pgxscan.Select(context.Background(), sd.pool, &posts,
		`SELECT id, author, created, forum, is_edited, message, parent, thread
			FROM posts WHERE thread = $1 AND path > (SELECT path FROM posts WHERE id = $2)
			ORDER BY path LIMIT $3`, id, since, limit)

	if errors.As(err, &pgx.ErrNoRows) || len(posts) == 0 {
		return models.Posts{}, nil
	}

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (sd SomeDatabase) GetPostsTreeSinceDesc(id int ,limit int, since int) (models.Posts, error) {
	var posts models.Posts
	err := pgxscan.Select(context.Background(), sd.pool, &posts,
		`SELECT id, author, created, forum, is_edited, message, parent, thread
			FROM posts WHERE thread = $1 AND path < (SELECT path FROM posts WHERE id = $2)
			ORDER BY path DESC LIMIT $3`, id, since, limit)

	if errors.As(err, &pgx.ErrNoRows) || len(posts) == 0 {
		return models.Posts{}, nil
	}

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (sd SomeDatabase) AddForumUsers(slug string, author string) error {
	_, err := sd.pool.Exec(context.Background(),
		`INSERT INTO forum_users 
		VALUES ($1, $2)`,
		slug, author)
	if err != nil {
		return err
	}

	return nil
}

func (sd SomeDatabase) GetForumThreads(slug string, limit int, since string, desc bool) (models.Threads, error) {
	var threads []models.ThreadSQL
	var err error
	if since == "" {
		if desc == true {
			err = pgxscan.Select(context.Background(), sd.pool, &threads,
				`SELECT * FROM threads WHERE forum = $1
				ORDER BY created DESC LIMIT $2`, slug, limit)
		} else {
			err = pgxscan.Select(context.Background(), sd.pool, &threads,
				`SELECT * FROM threads WHERE forum = $1
				ORDER BY created LIMIT $2`, slug, limit)
		}
	} else {
		if desc == true {
			err = pgxscan.Select(context.Background(), sd.pool, &threads,
				`SELECT * FROM threads WHERE forum = $1 AND created <= $2
				ORDER BY created DESC LIMIT $3`, slug, since, limit)
		} else {
			err = pgxscan.Select(context.Background(), sd.pool, &threads,
				`SELECT * FROM threads WHERE forum = $1 AND created >= $2
				ORDER BY created LIMIT $3`, slug, since, limit)
		}
	}
	if errors.As(err, &pgx.ErrNoRows) || len(threads) == 0 {
		return models.Threads{}, nil
	}

	if err != nil {
		return nil, err
	}

	var threads2 []models.Thread
	for i := range threads {
		threads2 = append(threads2, models.ConvertThread(threads[i]))
	}
	return threads2, nil
}

func (sd SomeDatabase) EditMessage(id int, message string) error {
	_, err := sd.pool.Exec(context.Background(),
		`UPDATE posts SET is_edited = true, message = $1 WHERE id = $2`, message, id)

	if err != nil {
		return err
	}

	return nil
}

func (sd SomeDatabase) UpdateThread(slugOrId string, thread models.Thread) (models.Thread, error) {
	err := sd.pool.QueryRow(context.Background(),
		`UPDATE threads SET message = $1, title = $2 WHERE slug = $3
			RETURNING threads.id, threads.author, threads.created, threads.forum,
			threads.message, threads.slug, threads.title, threads.votes`, thread.Message,
			thread.Title, slugOrId).Scan(&thread.Id, &thread.Author, &thread.Created,
				&thread.Forum, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)

	if err != nil {
		return models.Thread{}, err
	}

	return thread, nil
}

func (sd SomeDatabase) UpdateVote(id int, vote models.Vote) error {
	_, err := sd.pool.Exec(context.Background(),
		`UPDATE votes SET voice = $1 WHERE thread = $2 AND nickname = $3`, vote.Voice,
		id, vote.Nickname)

	if err != nil {
		return err
	}

	return nil
}

func (sd SomeDatabase) UpdateThreadById(id int, thread models.Thread) (models.Thread, error) {
	err := sd.pool.QueryRow(context.Background(),
		`UPDATE threads SET message = $1, title = $2 WHERE id = $3
			RETURNING threads.id, threads.author, threads.created, threads.forum,
			threads.message, threads.slug, threads.title, threads.votes`, thread.Message,
		thread.Title, id).Scan(&thread.Id, &thread.Author, &thread.Created,
		&thread.Forum, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)

	if err != nil {
		return models.Thread{}, err
	}

	return thread, nil
}

func (sd SomeDatabase) IncrementThreads(forum string) error {
	_, err := sd.pool.Exec(context.Background(),
		`UPDATE forums SET threads = threads + 1 WHERE slug = $1`, forum)

	if err != nil {
		return err
	}

	return nil
}

func (sd SomeDatabase) IncrementPosts(forum string) error {
	_, err := sd.pool.Exec(context.Background(),
		`UPDATE forums SET posts = posts + 1 WHERE slug = $1`, forum)

	if err != nil {
		return err
	}

	return nil
}

func (sd SomeDatabase) Clear() error {
	_, err := sd.pool.Exec(context.Background(),
		`TRUNCATE users, forums, threads, posts, votes, forum_users`)

	if err != nil {
		return err
	}

	return nil
}

func (sd SomeDatabase) Status() (models.Status, error) {
	status := models.Status{}
	err := sd.pool.QueryRow(context.Background(),
		`SELECT (SELECT count(id) FROM forums) as forums, 
			(SELECT count(id) FROM posts) as posts, 
			(SELECT count(id) FROM users) as users,
			(SELECT count(id) FROM threads) as threads`).Scan(&status.Forum, &status.Post, &status.User, &status.Thread)

	if err != nil {
		return models.Status{}, err
	}

	return status, nil
}

func (sd SomeDatabase) CreateUser(nickname string, user models.User) error {
	_, err := sd.pool.Exec(context.Background(),
		`INSERT INTO users VALUES(default, $1, $2, $3, $4)`, nickname, user.Fullname, user.About, user.Email)

	if err != nil {
		return err
	}

	return nil
}

func (sd SomeDatabase) UpdateUser(nickname string, user models.User) error {
	_, err := sd.pool.Exec(context.Background(),
		`UPDATE users SET fullname = $1, about = $2, email = $3 WHERE nickname = $4`,
		user.Fullname, user.About, user.Email, nickname)

	if err != nil {
		return err
	}

	return nil
}

func (sd SomeDatabase) GetUserByNicknameOrEmail(nickname string, email string) (models.Users, error) {
	var users models.Users
	err := pgxscan.Select(context.Background(), sd.pool, &users,
		`SELECT nickname, fullname, about, email FROM users WHERE nickname = $1 OR email = $2`, nickname, email)

	if errors.As(err, &pgx.ErrNoRows) || len(users) == 0 {
		return models.Users{}, nil
	}

	if err != nil {
		return models.Users{}, err
	}

	return users, nil
}