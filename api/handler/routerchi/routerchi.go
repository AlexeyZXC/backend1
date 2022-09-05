package routerchi

import (
	"net/http"
	"time"

	"github.com/AlexeyZXC/backend1/tree/CourseProject/api/handler"
	"github.com/AlexeyZXC/backend1/tree/CourseProject/app/repo/link"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type RouterChi struct {
	*chi.Mux
	ls *handler.Handlers
}

func NewRouterChi(ls *handler.Handlers) *RouterChi {
	r := chi.NewRouter()
	ret := &RouterChi{
		ls: ls,
	}

	r.Post("/genShortLink/{longLink}", ret.CreateShortLink)

	ret.Mux = r
	return ret
}

type Link link.Link

func (l *Link) Bind(r *http.Request) error { //todo get longlink from request
	l.LongLink = r.FormValue("longlink")
	data := link.Stat{
		UserIP:   r.RemoteAddr,
		PassTime: time.Now(),
	}
	l.StatData = append(l.StatData, data)
	return nil
}

func (Link) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (rt *RouterChi) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	rl := Link{}
	if err := render.Bind(r, &rl); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	l, err := rt.ls.CreateShortLink(r.Context(), rl.LongLink)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Render(w, r, Link(*l))
}

func (rt *RouterChi) UpdateStat(w http.ResponseWriter, r *http.Request) {
	rl := Link{}
	if err := render.Bind(r, &rl); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	if err := rt.ls.UpdateStat(r.Context(), rl.ShortLink, rl.StatData[0].UserIP); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Render(w, r, Link{})
}

func (rt *RouterChi) GetStat(w http.ResponseWriter, r *http.Request) {
	rl := Link{}
	if err := render.Bind(r, &rl); err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	stat, err := rt.ls.GetStat(r.Context(), rl.ShortLink)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	rl.StatData = stat

	render.Render(w, r, rl)
}
