package hosts

import (
	"github.com/ctripcorp/nephele/imgsvr/config"
)

type StartOpt struct {
	NewHosting bool
}

type Host interface {
	GetPort() int
	Start(opt StartOpt) <-chan error
	Stop()
	ReloadConf()
}

func New(cfg config.Cfg) Host {
	var host Host
	if cfg.MasterMode {
		host = NewMaster(cfg)
	} else if cfg.ImgProcWorkerMode {
		host = NewImageProcessWorker(cfg)
	}
	return host
}
