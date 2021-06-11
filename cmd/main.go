package main

import (
	"subd/server"
)

func main() {
	s := server.NewServer()
	s.ListenAndServe()
}
