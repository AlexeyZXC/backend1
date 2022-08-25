package main

// Выбираю chi роутер т.к. он предоставляет дополнительную функциональность, что позволяет быстрее разрабатывать и писать более читабельный код.
// С другой стороны он использует сигнатуры методов сходные со сигнатцурами из стандартных пакетов.
import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("chi router is working."))
	})
	http.ListenAndServe(":8000", r)
}
