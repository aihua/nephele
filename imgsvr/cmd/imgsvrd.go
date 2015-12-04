package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ctripcorp/nephele/imgsvr/config"
	"github.com/ctripcorp/nephele/imgsvr/hosts"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	os.Exit(Main())
}

func Main() int {
	var (
		err error
		cfg config.Cfg
	)
	if cfg, err = config.Parse(os.Args[1:]); err != nil {
		return 2
	}

	runtime.GOMAXPROCS(cfg.GoMaxProc)

	config.ApplyLogConfig()
	config.ApplyCatConfig()

	log.Debugf("start svr in %s", cfg.Env)

	host := hosts.New(cfg)
	startErrCh := host.Start(hosts.StartOpt{})
	defer host.Stop()

	hup := make(chan os.Signal)
	signal.Notify(hup, syscall.SIGHUP)
	waitForReloading(host, hup)

	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case <-term:
		log.Warn("Received SIGTERM, exiting gracefully...")
	case err := <-startErrCh:
		log.Errorln("Error starting web server, exiting gracefully:", err)
	}

	log.Info("See you next time!")
	return 0
}

func waitForReloading(h hosts.Host, hup chan os.Signal) {
	go func() {
		for {
			<-hup
			h.ReloadConf()
		}
	}()
}
