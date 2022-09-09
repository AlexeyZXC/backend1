package main

import (
	"fmt"

	"github.com/AlexeyZXC/backend1/tree/CourseProject/api/handler"
	"github.com/AlexeyZXC/backend1/tree/CourseProject/api/handler/routerchi"
	"github.com/AlexeyZXC/backend1/tree/CourseProject/api/server"
	"github.com/AlexeyZXC/backend1/tree/CourseProject/app/repo/link"
	"github.com/AlexeyZXC/backend1/tree/CourseProject/db/postgres"
)

func main() {

	db, err := postgres.NewPgDB()
	if err != nil {
		fmt.Println("db connect error: ", err)
		return
	}

	ls := link.NewLinks(db)
	h := handler.NewHandlers(ls)

	rh := routerchi.NewRouterChi(h)

	srv := server.NewServer(":8000", rh)

	srv.Start(ls)
}
