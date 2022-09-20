// The package servers as a router.
package routerchi

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/AlexeyZXC/backend1/tree/CourseProject/api/handler"
	"github.com/AlexeyZXC/backend1/tree/CourseProject/app/repo/pages"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type RouterChi struct {
	*chi.Mux
	ls  *handler.Handlers
	log *zap.SugaredLogger
}

// NewRouterChi returns a new object of Router with the embedded a Chi mux.
func NewRouterChi(ls *handler.Handlers, log *zap.SugaredLogger) *RouterChi {
	r := chi.NewRouter()
	ret := &RouterChi{
		ls:  ls,
		log: log,
	}

	mwl := &mwLog{log}
	r.Use(mwl.Logger)

	r.Get("/stat/{id}", ret.OpenStatLink)
	r.Get("/sl/{id}", ret.OpenShortLink)
	r.Post("/", ret.CreateShortLink)
	r.Get("/", ret.defaultPage)

	ret.Mux = r
	return ret
}

type Link handler.Link

func (l *Link) Bind(r *http.Request) error {
	return nil
}

func (Link) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type mwLog struct {
	log *zap.SugaredLogger
}

// Logger serves as a middleware for logging purpose.
func (mwl *mwLog) Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		body, _ := io.ReadAll(r.Body)
		defer r.Body.Close()

		r.Body = io.NopCloser(bytes.NewBuffer(body))
		mwl.log.Infoln(r.Method + ": " + string(body))

		handler.ServeHTTP(w, r)
	})
}

// CreateShortLink processes Post requests on "/" path.
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

	_, err = fmt.Fprintf(w, pages.DefaultPageContent, l.LongLink, pages.ShortLinkUrl+l.ShortLink, pages.StatUrl+l.ShortLink)
	if err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

// OpenShortLink processes Get requests on "/sl/{id}" path.
func (rt *RouterChi) OpenShortLink(w http.ResponseWriter, r *http.Request) {
	sl := chi.URLParam(r, "id")
	if sl == "" {
		render.Render(w, r, ErrRender(errors.New("empty short link")))
	}
	sli, err := strconv.Atoi(sl)
	if err != nil {
		render.Render(w, r, ErrRender(errors.New("wrong short link")))
		return
	}

	ll, err := rt.ls.GetLongLink(r.Context(), sli)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	if ll.LongLink != "" {
		http.Redirect(w, r, ll.LongLink, http.StatusMovedPermanently)

		if err := rt.ls.UpdateStat(r.Context(), sli, r.RemoteAddr); err != nil {
			fmt.Printf("failed to update stat for short Link(%v), err: %v \n", sli, err)
		}

	} else {
		http.Redirect(w, r, ll.LongLink, http.StatusNotFound)
	}
}

// defaultPage processes Get requests on "/" path.
func (rt *RouterChi) defaultPage(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, pages.DefaultPageContent, "", "", "")
	if err != nil {
		render.Render(w, r, ErrRender(err))
	}
}

// OpenStatLink processes Get requests on "/stat/{id}" path.
func (rt *RouterChi) OpenStatLink(w http.ResponseWriter, r *http.Request) {
	sl := chi.URLParam(r, "id")
	if sl == "" {
		render.Render(w, r, ErrRender(errors.New("empty short link")))
	}
	sli, err := strconv.Atoi(sl)
	if err != nil {
		render.Render(w, r, ErrRender(errors.New("wrong short link")))
		return
	}

	stat, err := rt.ls.GetStat(r.Context(), sli)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	ll, err := rt.ls.GetLongLink(r.Context(), sli)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	var rl Link
	rl.ShortLink = pages.ShortLinkUrl + sl
	rl.LongLink = ll.LongLink

	for _, s := range stat {
		var st handler.Stat
		st.UserIP = s.UserIP
		st.PassTime = s.PassTime

		rl.StatData = append(rl.StatData, st)
	}

	render.Render(w, r, rl)
}
