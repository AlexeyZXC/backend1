package main

import (
	"github.com/AlexeyZXC/backend1/tree/CourseProject/api/handler"
	"github.com/AlexeyZXC/backend1/tree/CourseProject/api/handler/routerchi"
	"github.com/AlexeyZXC/backend1/tree/CourseProject/api/server"
	"github.com/AlexeyZXC/backend1/tree/CourseProject/app/repo/link"
)

func main() {

	ls := link.NewLinks()
	h := handler.NewHandlers(ls)

	rh := routerchi.NewRouterChi(h)

	srv := server.NewServer(":8000", rh)

	srv.Start(ls)
}
