package hosts

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	cat "github.com/ctripcorp/cat.go"
	"github.com/ctripcorp/nephele/imgsvr/config"
	"github.com/ctripcorp/nephele/imgsvr/web"
	"github.com/ctripcorp/nephele/util"
	"os/exec"
	"strconv"
)

type ImageProcessWorker struct {
	Port       int
	startErrCh chan error
}

func NewImageProcessWorker(cfg config.Cfg) *ImageProcessWorker {
	return &ImageProcessWorker{Port: cfg.Port, startErrCh: make(chan error)}
}

func (w *ImageProcessWorker) GetPort() int {
	return w.Port
}

func (w *ImageProcessWorker) Start(opt StartOpt) <-chan error {
	if !opt.NewHosting {
		opts := &web.ImageProcessOptions{ListenAddress: fmt.Sprintf("127.0.0.1:%d", w.Port)}
		webHandler := web.NewImageProcessWebHandler(opts)
		util.LogEvent(cat.Instance(), "Reboot", fmt.Sprintf("%s:%d", util.LocalIP(), w.Port), nil)
		go webHandler.Run()
		return webHandler.ListenError()
	} else {
		cmd := exec.Command("go", "run", "imgsvrd.go", "-process.worker", "-port", strconv.Itoa(w.Port))
		err := cmd.Start()
		if err != nil {
			log.Errorf("start worker=%d failed. error=%s", w.Port, err.Error())
			util.LogErrorEvent(cat.Instance(), "DaemonProcess.StartWorkerError", err.Error())
			w.startErrCh <- err
		}
	}
	return w.startErrCh
}

func (w *ImageProcessWorker) Stop() {

}

func (w *ImageProcessWorker) ReloadConf() {

}
