package repository

import (
	"context"
	"errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
	event "subd"
	"subd/models"
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

func (sd SomeDatabase) AddPost(post models.Post) (models.Post, error) {
	err := sd.pool.QueryRow(context.Background(),
		`INSERT INTO posts VALUES (default, $1, $2, $3, default, $4, $5, $6) RETURNING id`,
		post.Author, post.Created, post.Forum, post.Message, post.Parent, post.Thread).Scan(&post.Id)
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
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
			ORDER BY a.path LIMIT $3) AS b
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
				`select * from threads where forum = $1
				order by created DESC LIMIT $2`, slug, limit)
		} else {
			err = pgxscan.Select(context.Background(), sd.pool, &threads,
				`select * from threads where forum = $1
				order by created LIMIT $2`, slug, limit)
		}
	} else {
		if desc == true {
			err = pgxscan.Select(context.Background(), sd.pool, &threads,
				`select * from threads where forum = $1 AND created <= $2
				order by created DESC LIMIT $3`, slug, since, limit)
		} else {
			err = pgxscan.Select(context.Background(), sd.pool, &threads,
				`select * from threads where forum = $1 AND created >= $2
				order by created LIMIT $3`, slug, since, limit)
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


/*func (ed EventDatabase) GetNearEvents(now time.Time, coord models.Coordinates, page int) ([]models.EventCardWithCoordsSQL, error) {
	var events []models.EventCardWithCoordsSQL
	err := pgxscan.Select(context.Background(), ed.pool, &events,
		`SELECT id, title, place, description, start_date, end_date, (
		2*6371*asin(sqrt(abs(sin(radians((latitude- $2)/2))*sin(radians((latitude- $2)/2))+
                     sin(radians((longitude- $3)/2))*sin(radians((longitude- $3)/2))
                         *cos(latitude)*cos($2))))) AS distance FROM events

		WHERE end_date > $1
		ORDER BY distance, id
		LIMIT $4 OFFSET $5`, now, coord.Latitude, coord.Longitude, constants.EventsPerPage, (page-1)*constants.EventsPerPage)

	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		ed.logger.Debug("no rows in method GetNearEvents")
		return []models.EventCardWithCoordsSQL{}, nil
	}

	if err != nil {
		ed.logger.Warn(err)
		return nil, err
	}

	return events, nil
}

func (ed EventDatabase) GetAllEvents(now time.Time, page int) ([]models.EventCardWithDateSQL, error) {
	var events []models.EventCardWithDateSQL
	err := pgxscan.Select(context.Background(), ed.pool, &events,
		`SELECT id, title, place, description, start_date, end_date FROM events
		WHERE end_date > $1
		ORDER BY id DESC
		LIMIT $2 OFFSET $3`, now, constants.EventsPerPage, (page-1)*constants.EventsPerPage)

	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		ed.logger.Debug("no rows in method GetAllEvents")
		return []models.EventCardWithDateSQL{}, nil
	}

	if err != nil {
		ed.logger.Warn(err)
		return nil, err
	}

	return events, nil
}

func (ed EventDatabase) GetEventsByCategory(typeEvent string, now time.Time, page int) ([]models.EventCardWithDateSQL, error) {
	var events []models.EventCardWithDateSQL
	err := pgxscan.Select(context.Background(), ed.pool, &events,
		`SELECT id, title, place, description, start_date, end_date FROM events
		WHERE category = $1 AND end_date > $2
		ORDER BY id DESC
		LIMIT $3 OFFSET $4`, typeEvent, now, constants.EventsPerPage, (page-1)*constants.EventsPerPage)

	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		ed.logger.Debug("no rows in method GetEventsByCategory")
		return []models.EventCardWithDateSQL{}, nil
	}

	if err != nil {
		ed.logger.Warn(err)
		return []models.EventCardWithDateSQL{}, err
	}

	return events, nil
}

func (ed EventDatabase) GetOneEventByID(eventId uint64) (models.EventSQL, error) {
	var ev []models.EventSQL
	err := pgxscan.Select(context.Background(), ed.pool, &ev,
		`SELECT * FROM events WHERE id = $1`, eventId)

	if errors.As(err, &pgx.ErrNoRows) || len(ev) == 0 {
		ed.logger.Debug("no event with id " + fmt.Sprint(eventId))
		return models.EventSQL{}, echo.NewHTTPError(http.StatusNotFound, errors.New("Event with id "+fmt.Sprint(eventId)+" not found"))
	}

	if err != nil {
		ed.logger.Warn(err)
		return models.EventSQL{}, err
	}

	return ev[0], nil
}

func (ed EventDatabase) GetOneEventNameByID(eventId uint64) (string, error) {
	var ev []string
	err := pgxscan.Select(context.Background(), ed.pool, &ev,
		`SELECT title FROM events WHERE id = $1`, eventId)

	if errors.As(err, &pgx.ErrNoRows) || len(ev) == 0 {
		ed.logger.Debug("no event with id " + fmt.Sprint(eventId))
		return "", echo.NewHTTPError(http.StatusNotFound, errors.New("Event with id "+fmt.Sprint(eventId)+" not found"))
	}

	if err != nil {
		ed.logger.Warn(err)
		return "", err
	}

	return ev[0], nil
}

func (ed EventDatabase) GetTags(eventId uint64) (models.Tags, error) {
	var parameters models.Tags
	err := pgxscan.Select(context.Background(), ed.pool, &parameters,
		`SELECT t.id AS id, t.name AS name
		FROM tags t
		JOIN event_tag e on e.tag_id = t.id
        WHERE e.event_id = $1`, eventId)

	if errors.As(err, &pgx.ErrNoRows) || len(parameters) == 0 {
		ed.logger.Debug("no rows in method GetTags")
		return models.Tags{}, nil
	}

	if err != nil {
		ed.logger.Warn(err)
		return models.Tags{}, err
	}

	return parameters, err
}

func (ed EventDatabase) AddEvent(newEvent *models.Event) error {
	_, err := ed.pool.Exec(context.Background(),
		`INSERT INTO events (title, place, subway, street, description, category, start_date, end_date, image) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		newEvent.Title, newEvent.Place, newEvent.Subway, newEvent.Street, newEvent.Description,
		newEvent.Category, newEvent.StartDate, newEvent.EndDate, newEvent.Image)
	if err != nil {
		ed.logger.Warn(err)
		return err
	}

	return nil
}

func (ed EventDatabase) DeleteById(eventId uint64) error {
	resp, err := ed.pool.Exec(context.Background(),
		`DELETE FROM events WHERE id = $1`, eventId)

	if err != nil {
		ed.logger.Warn(err)
		return err
	}

	if resp.RowsAffected() == 0 {
		return echo.NewHTTPError(http.StatusNotFound, errors.New("Event with id "+fmt.Sprint(eventId)+" not found"))
	}

	return nil
}

func (ed EventDatabase) UpdateEventAvatar(eventId uint64, path string) error {
	_, err := ed.pool.Exec(context.Background(),
		`UPDATE events SET image = $1 WHERE id = $2`, path, eventId)

	if err != nil {
		ed.logger.Warn(err)
		return err
	}

	return nil
}

func (ed EventDatabase) FindEvents(str string, now time.Time, page int) ([]models.EventCardWithDateSQL, error) {
	var events []models.EventCardWithDateSQL
	err := pgxscan.Select(context.Background(), ed.pool, &events,
		`SELECT DISTINCT ON(e.id) e.id, e.title, e.place,
		e.description, e.start_date, e.end_date FROM
        events e JOIN event_tag et on e.id = et.event_id
        JOIN tags t on et.tag_id = t.id
		WHERE (LOWER(title) LIKE '%' || $1 || '%' OR LOWER(description) LIKE '%' || $1 || '%'
		OR LOWER(category) LIKE '%' || $1 || '%' OR LOWER(t.name) LIKE '%' || $1 || '%')
		AND end_date > $2
		ORDER BY e.id DESC
		LIMIT $3 OFFSET $4`, str, now, constants.EventsPerPage, (page-1)*constants.EventsPerPage)

	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		ed.logger.Debug("no rows in method FindEvents with string " + str)
		return []models.EventCardWithDateSQL{}, nil
	}

	if err != nil {
		ed.logger.Warn(err)
		return []models.EventCardWithDateSQL{}, err
	}

	return events, nil
}

func (ed EventDatabase) RecomendSystem(uid uint64, category string) error {

	_, err := ed.pool.Exec(context.Background(),
		`UPDATE user_preference SET `+constants.Category[category]+`=`+constants.Category[category]+`+1 `+`WHERE user_id = $1`, uid)

	if errors.As(err, &pgx.ErrNoRows) {
		ed.logger.Debug("no rows in method RecomendSystem with id " + fmt.Sprint(uid))
		return err
	}

	if err != nil {
		ed.logger.Warn(err)
		return err
	}

	return nil
}

func (ed EventDatabase) GetPreference(uid uint64) (models.Recomend, error) {
	var recomend []models.Recomend
	err := pgxscan.Select(context.Background(), ed.pool, &recomend,
		`SELECT entertainment, education, cinema, exhibition, festival, tour, concert
		FROM user_preference
		WHERE user_id = $1`, uid)

	if errors.As(err, &pgx.ErrNoRows) {
		ed.logger.Debug("no rows in method GetPreference with id " + fmt.Sprint(uid))
		return models.Recomend{}, err
	}

	if err != nil {
		ed.logger.Warn(err)
		return models.Recomend{}, err
	}

	return recomend[0], nil
}

func (ed EventDatabase) CategorySearch(str string, category string, now time.Time, page int) ([]models.EventCardWithDateSQL, error) {
	var events []models.EventCardWithDateSQL
	err := pgxscan.Select(context.Background(), ed.pool, &events,
		`SELECT DISTINCT ON(e.id) e.id, e.title, e.place,
		e.description, e.start_date, e.end_date FROM
        events e JOIN event_tag et on e.id = et.event_id
        JOIN tags t on et.tag_id = t.id
		WHERE (LOWER(title) LIKE '%' || $1 || '%' OR LOWER(description) LIKE '%' || $1 || '%'
		OR LOWER(t.name) LIKE '%' || $1 || '%') AND e.category = $2
		AND end_date > $3
		ORDER BY e.id DESC
		LIMIT $4 OFFSET $5`, str, category, now, constants.EventsPerPage, (page-1)*constants.EventsPerPage)

	if errors.As(err, &pgx.ErrNoRows) || len(events) == 0 {
		ed.logger.Debug("no rows in method CategorySearch with searchstring " + str)
		return []models.EventCardWithDateSQL{}, nil
	}

	if err != nil {
		ed.logger.Warn(err)
		return []models.EventCardWithDateSQL{}, err
	}

	return events, nil
}

func (ed EventDatabase) GetRecommended(uid uint64, now time.Time, page int) ([]models.EventCardWithDateSQL, error) {
	recommend, err := ed.GetPreference(uid)
	if err != nil {
		ed.logger.Debug(string(err.Error()))
		return ed.GetAllEvents(now, page)
	}
	var eventsFirst, eventsSecond, eventsThird []models.EventCardWithDateSQL
	var s []float64
	s = append(s, float64(recommend.Concert), float64(recommend.Tour), float64(recommend.Entertainment),
		float64(recommend.Festival), float64(recommend.Exhibition),
		float64(recommend.Cinema), float64(recommend.Education))

	var first string
	if float64(recommend.Concert) == floats.Max(s) {
		first = "Концерт"
	}
	if float64(recommend.Cinema) == floats.Max(s) {
		first = "Кино"
	}
	if float64(recommend.Exhibition) == floats.Max(s) {
		first = "Выставка"
	}
	if float64(recommend.Education) == floats.Max(s) {
		first = "Образование"
	}
	if float64(recommend.Tour) == floats.Max(s) {
		first = "Экскурсия"
	}
	if float64(recommend.Festival) == floats.Max(s) {
		first = "Фестиваль"
	}
	if float64(recommend.Entertainment) == floats.Max(s) {
		first = "Развлечения"
	}

	err = pgxscan.Select(context.Background(), ed.pool, &eventsFirst,
		`SELECT id, title, place, description, start_date, end_date FROM events
			WHERE category = $1 AND end_date > $2
			ORDER BY id DESC
			LIMIT 3 OFFSET $3`, first, now, (page-1)*3)
	if errors.As(err, &pgx.ErrNoRows) || len(eventsFirst) == 0 {
		ed.logger.Debug("debug")
	} else {
		if err != nil {
			ed.logger.Warn(err)
			return nil, err
		}
	}

	if first != "Концерт" {
		err = pgxscan.Select(context.Background(), ed.pool, &eventsSecond,
			`SELECT id, title, place, description, start_date, end_date FROM events
			WHERE category = 'Концерт' AND end_date > $1
			ORDER BY id DESC
			LIMIT 2 OFFSET $2`, now, (page-1)*2)

		if errors.As(err, &pgx.ErrNoRows) || len(eventsSecond) == 0 {
			ed.logger.Debug("debug")
		} else {
			if err != nil {
				ed.logger.Warn(err)
				return nil, err
			}
		}
		if first != "Выставка" {
			err = pgxscan.Select(context.Background(), ed.pool, &eventsThird,
				`SELECT id, title, place, description, start_date, end_date FROM events
			WHERE category = 'Выставка' AND end_date > $1
			ORDER BY id DESC
			LIMIT 1 OFFSET $2`, now, (page-1)*1)
			if errors.As(err, &pgx.ErrNoRows) || len(eventsThird) == 0 {
				ed.logger.Debug("debug")
			} else {
				if err != nil {
					ed.logger.Warn(err)
					return nil, err
				}
			}
		} else {
			err = pgxscan.Select(context.Background(), ed.pool, &eventsThird,
				`SELECT id, title, place, description, start_date, end_date FROM events
			WHERE category = 'Фестиваль' AND end_date > $1
			ORDER BY id DESC
			LIMIT 1 OFFSET $2`, now, (page-1)*1)
			if errors.As(err, &pgx.ErrNoRows) || len(eventsThird) == 0 {
				ed.logger.Debug("debug")
			} else {
				if err != nil {
					ed.logger.Warn(err)
					return nil, err
				}
			}
		}
	} else {
		err = pgxscan.Select(context.Background(), ed.pool, &eventsSecond,
			`SELECT id, title, place, description, start_date, end_date FROM events
			WHERE category = 'Выставка' AND end_date > $1
			ORDER BY id DESC
			LIMIT 2 OFFSET $2`, now, (page-1)*2)
		if errors.As(err, &pgx.ErrNoRows) || len(eventsSecond) == 0 {
			ed.logger.Debug("debug")
		} else {
			if err != nil {
				ed.logger.Warn(err)
				return nil, err
			}
		}
		err = pgxscan.Select(context.Background(), ed.pool, &eventsThird,
			`SELECT id, title, place, description, start_date, end_date FROM events
			WHERE category = 'Фестиваль' AND end_date > $1
			ORDER BY id DESC
			LIMIT 1 OFFSET $2`, now, (page-1)*1)
		if errors.As(err, &pgx.ErrNoRows) || len(eventsThird) == 0 {
			ed.logger.Debug("debug")
		} else {
			if err != nil {
				ed.logger.Warn(err)
				return nil, err
			}
		}
	}
	eventsFirst = append(eventsFirst, eventsSecond...)
	eventsFirst = append(eventsFirst, eventsThird...)

	return eventsFirst, nil
}

/*func (ed EventDatabase) GetRecommended(uid uint64, now time.Time, page int) ([]models.EventCardWithDateSQL, error) {
	recomend, err := ed.GetPreference(uid)
	if err != nil || (recomend.Concert == recomend.Movie && recomend.Movie == recomend.Show && recomend.Show == 0) {
		ed.logger.Debug(string(err.Error()))
		return ed.GetAllEvents(now, 1)
	}
	var eventsPrefer []models.EventCardWithDateSQL
	var param string
	if recomend.Concert >= recomend.Movie && recomend.Concert >= recomend.Show {
		param = "Музей"
	}
	if recomend.Movie >= recomend.Concert && recomend.Movie >= recomend.Show {
		param = "Кино"
	}
	if recomend.Show >= recomend.Concert && recomend.Show >= recomend.Movie {
		param = "Выставка"
	}
	err = pgxscan.Select(context.Background(), ed.pool, &eventsPrefer,
		`SELECT id, title, description, image, start_date, end_date FROM events
			WHERE category = $1 AND end_date > $2
			UNION
			SELECT id, title, description, image, start_date, end_date FROM events
			WHERE category != $1 AND end_date > $2
			ORDER BY id DESC
			LIMIT 6 OFFSET $3`, param, now, (page-1)*6)
	if err != nil {
		ed.logger.Warn(err)
		return nil, err
	}

	if len(eventsPrefer) == 0 {
		ed.logger.Debug("no rows in method GetRecomended")
		return []models.EventCardWithDateSQL{}, nil
	}

	return eventsPrefer, nil
}*/
