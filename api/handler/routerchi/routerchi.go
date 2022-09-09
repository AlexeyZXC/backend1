package routerchi

import (
	"fmt"
	"io"
	"net/http"

	"github.com/AlexeyZXC/backend1/tree/CourseProject/api/handler"
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

	//r.Post("/genShortLink/{longLink}", ret.CreateShortLink)
	r.Post("/", ret.CreateShortLink)
	r.Get("/", ret.defaultPage)

	ret.Mux = r
	return ret
}

type Link handler.Link

func (l *Link) Bind(r *http.Request) error { //todo get longlink from request
	if err := r.ParseForm(); err != nil {
		fmt.Println("fail parse form: ", err)
		return err
	}
	l.LongLink = r.FormValue("lurl")

	fmt.Println("bind: r.Form:", r.Form)
	// data := link.Stat{
	// 	UserIP:   r.RemoteAddr,
	// 	PassTime: time.Now(),
	// }
	//l.StatData = append(l.StatData, data)
	return nil
}

func (Link) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (rt *RouterChi) CreateShortLink(w http.ResponseWriter, r *http.Request) {

	method := r.Method
	fmt.Printf("CreateShortLink: method: %v\n", method)

	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	fmt.Println("body: ", string(b))

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

	render.Render(w, r, Link(l))
}

func (rt *RouterChi) defaultPage(w http.ResponseWriter, r *http.Request) {

	q := chi.URLParam(r, "fname")
	method := r.Method
	fmt.Printf("CreateShortLink: method: %v, fname: %v\n", method, q)

	bodyContent := `
	<h1>URL shortener</h1>

	<form action="/" method="post" enctype="text/json">
	  <label for="lurl">Long URL:</label>
	  <input type="text" id="lurl" name="lurl"><br><br>
	  <label for="surl">Short URL:</label>
	  <input type="text" id="surl" name="surl"><br><br>
	  <input type="submit" value="Submit">
	</form>
	`
	n, err := w.Write([]byte(bodyContent))
	if err != nil {
		fmt.Println("defaultPage: err:", err)
	}
	if n != len(bodyContent) {
		fmt.Printf("defaultPage: wrote: %v instead of %v\n", n, len(bodyContent))
	}

	//render.Render(w, r, nil)
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

	for _, s := range stat {
		var st handler.Stat
		st.UserIP = s.UserIP
		st.PassTime = s.PassTime
		rl.StatData = append(rl.StatData, st)
	}

	render.Render(w, r, rl)
}
