CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS forum_users CASCADE;
DROP TABLE IF EXISTS votes CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS forums CASCADE;
DROP TABLE IF EXISTS users CASCADE;

DROP FUNCTION IF EXISTS insert_votes();
DROP FUNCTION IF EXISTS update_votes();
DROP FUNCTION IF EXISTS post_path();

DROP TRIGGER IF EXISTS insert_votes ON votes;
DROP TRIGGER IF EXISTS update_votes ON votes;
DROP TRIGGER IF EXISTS post_path ON posts;

CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    nickname CITEXT COLLATE "C" UNIQUE NOT NULL,
    fullname CITEXT        NOT NULL,
    about    TEXT                      NOT NULL,
    email    CITEXT UNIQUE             NOT NULL
);

CREATE INDEX user_nickname ON users using hash (nickname);
CREATE INDEX user_email ON users using hash (email);

CREATE TABLE forums
(
    id      SERIAL PRIMARY KEY,
    title   TEXT                      NOT NULL,
    owner   CITEXT REFERENCES users(nickname) ON DELETE CASCADE NOT NULL,
    posts   INT DEFAULT 0,
    threads INT DEFAULT 0,
    slug    CITEXT UNIQUE NOT NULL
);

CREATE INDEX forums_slug ON forums USING hash (slug);
CREATE INDEX forums_owners on forums (owner);

CREATE TABLE threads
(
    id      SERIAL PRIMARY KEY,
    author  CITEXT REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    created TIMESTAMP WITH TIME ZONE DEFAULT now(),
    forum   CITEXT REFERENCES forums (slug) ON DELETE CASCADE NOT NULL,
    message TEXT               NOT NULL,
    slug    CITEXT UNIQUE,
    title   CITEXT NOT NULL,
    votes   INT                      DEFAULT 0
);

create index threads_forum_created on threads (forum, created);
create index threads_created on threads (created);
create index threads_slug on threads using hash (slug);

CREATE TABLE posts
(
    id        BIGSERIAL PRIMARY KEY,
    author    CITEXT REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    created   TIMESTAMP WITH TIME ZONE DEFAULT now(),
    forum     CITEXT REFERENCES forums (slug) ON DELETE CASCADE NOT NULL,
    is_edited BOOLEAN                  DEFAULT FALSE,
    message   TEXT               NOT NULL,
    parent    INT                NOT NULL,
    thread    INT REFERENCES threads (id) ON DELETE CASCADE NOT NULL,
    path      BIGINT[]
);

create index posts_id on posts (id);
create index posts_thread_created_id on posts (thread, created, id);
create index posts_thread_id on posts (thread, id);
create index posts_thread_path on posts (thread, path);
create index posts_path_1_path on posts ((path[1]));

CREATE TABLE votes
(
    thread   INT REFERENCES threads (id) NOT NULL,
    voice    INT                NOT NULL,
    nickname CITEXT REFERENCES users (nickname) NOT NULL,
    UNIQUE (thread, nickname)
);

create unique index votes_user_thread on votes (thread, nickname);

CREATE TABLE forum_users
(
    forum    CITEXT REFERENCES forums (slug) ON DELETE CASCADE NOT NULL,
    nickname CITEXT COLLATE "C" REFERENCES users (nickname) ON DELETE CASCADE NOT NULL,
    UNIQUE (forum, nickname)
);

create index forum_users_nickname on forum_users (nickname);

CREATE OR REPLACE FUNCTION insert_votes()
    RETURNS TRIGGER AS
$insert_votes$
BEGIN
    IF new.voice > 0 THEN
        UPDATE threads SET votes = (votes + 1)
        WHERE id = new.thread;
    ELSE
        UPDATE threads SET votes = (votes - 1)
        WHERE id = new.thread;
    END IF;
    RETURN new;
END;
$insert_votes$ language plpgsql;

CREATE TRIGGER insert_votes
    BEFORE INSERT ON votes FOR EACH ROW
EXECUTE PROCEDURE insert_votes();



CREATE OR REPLACE FUNCTION update_votes()
    RETURNS TRIGGER AS
$update_votes$
BEGIN
    IF new.voice > 0 THEN
        UPDATE threads
        SET votes = (votes + 2)
        WHERE threads.id = new.thread;
    else
        UPDATE threads
        SET votes = (votes - 2)
        WHERE threads.id = new.thread;
    END IF;
    RETURN new;
END;
$update_votes$ LANGUAGE plpgsql;

CREATE TRIGGER update_votes
    BEFORE UPDATE ON votes FOR EACH ROW
EXECUTE PROCEDURE update_votes();



CREATE OR REPLACE FUNCTION post_path()
    RETURNS TRIGGER AS
$post_path$
DECLARE
    parent_thread BIGINT;
    parent_path   BIGINT[];
BEGIN
    IF (new.parent = 0) THEN
        new.path := new.path || new.id;
    ELSE
        SELECT thread, path
        FROM posts p
        WHERE p.thread = new.thread
        AND p.id = new.parent
        INTO parent_thread, parent_path;
        IF parent_thread != new.thread OR NOT FOUND THEN
            RAISE EXCEPTION USING ERRCODE = '00404';
        END IF;
        new.path := parent_path || new.id;
    END IF;
    RETURN new;
END;
$post_path$ LANGUAGE plpgsql;

CREATE TRIGGER post_path
    BEFORE INSERT ON posts FOR EACH ROW
EXECUTE PROCEDURE post_path();