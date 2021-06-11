package server

import (
	"context"
	"log"

	"subd/constants"
	"subd/delivery/http"
	"subd/repository"
	"subd/usecase"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

type Server struct {
	e       *echo.Echo
}

func NewServer() *Server {
	var server Server

	e := echo.New()

	pool, err := pgxpool.Connect(context.Background(), constants.DBConnect)
	if err != nil {
		log.Fatal(err)
	}
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	newRepository := repository.NewSomeDatabase(pool)

	newUC := usecase.NewSmth(newRepository)

	http.CreateSmthHandler(e, newUC)

	server.e = e
	return &server
}

func (s Server) ListenAndServe() {
	s.e.Logger.Fatal(s.e.Start(":5000"))
}
