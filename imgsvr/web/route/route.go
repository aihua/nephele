package route

import (
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
	"sync"
)

var (
	mtx   = sync.RWMutex{}
	ctxts = map[*http.Request]context.Context{}
)

func handle(h http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for _, p := range params {
			ctx = context.WithValue(ctx, param(p.Key), p.Value)
		}

		mtx.Lock()
		ctxts[r] = ctx
		mtx.Unlock()

		h(w, r)

		mtx.Lock()
		delete(ctxts, r)
		mtx.Unlock()

	}
}

func Context(r *http.Request) context.Context {
	mtx.RLock()
	defer mtx.RUnlock()
	return ctxts[r]
}

type param string

func Param(ctx context.Context, p string) string {
	return ctx.Value(param(p)).(string)
}

type Router struct {
	rtr    *httprouter.Router
	prefix string
}

func New() *Router {
	return &Router{rtr: httprouter.New()}
}

func (r *Router) WithPrefix(prefix string) *Router {
	return &Router{rtr: r.rtr, prefix: r.prefix + prefix}
}

func (r *Router) Get(path string, h http.HandlerFunc) {
	r.rtr.GET(r.prefix+path, handle(h))
}

func (r *Router) Del(path string, h http.HandlerFunc) {
	r.rtr.DELETE(r.prefix+path, handle(h))
}

func (r *Router) Put(path string, h http.HandlerFunc) {
	r.rtr.PUT(r.prefix+path, handle(h))
}

func (r *Router) Post(path string, h http.HandlerFunc) {
	r.rtr.POST(r.prefix+path, handle(h))
}

func (r *Router) Redirect(w http.ResponseWriter, req *http.Request, path string, code int) {
	http.Redirect(w, req, r.prefix+path, code)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.rtr.ServeHTTP(w, req)
}
