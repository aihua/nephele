package web

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	cat "github.com/ctripcorp/cat.go"
	"github.com/ctripcorp/nephele/imgsvr/handler/proc"
	"github.com/ctripcorp/nephele/imgsvr/handler/proc/command"
	"github.com/ctripcorp/nephele/imgsvr/web/route"
	"github.com/ctripcorp/nephele/util"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

type ImageProcessHandler struct {
	listenPort  int
	listenAddr  string
	router      *route.Router
	listenErrCh chan error
	reloadCh    chan struct{}
}

func NewImageProcessWebHandler(listenPort int) *ImageProcessHandler {
	h := &ImageProcessHandler{
		listenPort:  listenPort,
		listenAddr:  fmt.Sprintf("127.0.0.1:%d", listenPort),
		router:      route.New(),
		listenErrCh: make(chan error),
		reloadCh:    make(chan struct{}),
	}

	h.router.Get("/*imagepath", h.handle)
	h.router.Get("/reload/", h.reload)

	return h
}

func (h *ImageProcessHandler) Run() {
	log.Infof("Listening on %s", h.listenAddr)
	proc.Run()
	h.listenErrCh <- http.ListenAndServe(h.listenAddr, h.router)
}

func (h *ImageProcessHandler) ListenError() <-chan error {
	return h.listenErrCh
}

func (h *ImageProcessHandler) Reload() <-chan struct{} {
	return h.reloadCh
}

func (h *ImageProcessHandler) handle(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		commandArg *command.CommandArgument
		handleOk   = true
		rootCtx    = route.Context(r)
		imagePath  = route.Param(rootCtx, "imagepath")
		catVar     = cat.Instance()
		catTran    = catVar.NewTransaction("URL", getShortImagePath(imagePath))
		ctxWithCat = context.WithValue(rootCtx, "cat", catVar)
	)
	defer postHandle(catVar, catTran, imagePath, w, &err, &handleOk)
	h.recUrlEvt(catVar, imagePath, r)
	commandArg, err = command.Parse(ctxWithCat)
	if err != nil {
		return
	}

	//arg,  := command.Parse(catCtx)
}

func (h *ImageProcessHandler) reload(w http.ResponseWriter, r *http.Request) {
	//do reload
	h.reloadCh <- struct{}{}
}

func (h *ImageProcessHandler) recUrlEvt(c cat.Cat, imagePath string, r *http.Request) {
	util.LogEvent(c, "URL", "URL.Method", map[string]string{
		"Http": fmt.Sprintf("%s %s", r.Method, imagePath),
	})

	util.LogEvent(c, "image", fmt.Sprintf("%s:%d", util.LocalIP(), h.listenPort), nil)

	util.LogEvent(c, "URL", "URL.Client", map[string]string{
		"clientip": util.HttpClietIP(r),
		"serverip": util.LocalIP(),
		"proto":    r.Proto,
		"referer":  r.Referer(),
	})
}

func postHandle(c cat.Cat, tran cat.Transaction, imagePath string, w http.ResponseWriter, err *error, handleOk *bool) {
	p := recover()
	if p != nil {
		c.LogPanic(p)
		tran.SetStatus(p)
		log.WithFields(log.Fields{"uri": imagePath}).Error(*err)
	}

	if *handleOk {
		tran.SetStatus("0")
		tran.Complete()
	} else {
		tran.SetStatus(*err)
		tran.Complete()
	}
	if p != nil || *err != nil {
		http.Error(w, http.StatusText(404), 404)
	}
}

func getShortImagePath(imagePath string) string {
	segments := strings.Split(imagePath, "/")
	if len(segments) < 4 {
		return imagePath
	}
	if segments[2] == "fd" || segments[2] == "t1" {
		return util.JoinString("/", segments[1], "/", segments[2], "/", segments[3])
	} else {
		return util.JoinString("/", segments[1], "/", segments[2])
	}
}
