package config

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	cat "github.com/ctripcorp/cat.go"
	"os"
	"runtime"
)

type Env string

type Cfg struct {
	fs                 *flag.FlagSet
	verbose            bool
	logPath            string
	MasterMode         bool
	ImgProcWorkerMode  bool
	ImgProcWorkerCount int
	ImgUpWorkerMode    bool
	Port               int
	Env                Env
	Appid              string
	GoMaxProc          int
}

var (
	prod Env = "Prod"
	uat  Env = "Uat"
	cfg  Cfg = Cfg{}
)

func init() {
	flag.CommandLine.Init(os.Args[0], flag.ContinueOnError)
	cfg.fs = flag.CommandLine

	cfg.fs.BoolVar(
		&cfg.verbose, "v", false,
		"Enable verbose mode",
	)

	cfg.fs.StringVar(
		&cfg.logPath, "log.path", "nephele.log",
		"Set log path")

	cfg.fs.StringVar(
		&cfg.Appid, "id", "nephele", "set appid")

	cfg.fs.IntVar(
		&cfg.GoMaxProc, "gmp", 1, "set GOMAXPROCS")

	cfg.fs.BoolVar(
		&cfg.MasterMode, "master", false, "run as master mode")

	cfg.fs.BoolVar(
		&cfg.ImgProcWorkerMode, "process.worker", false, "run as image process worker mode")

	cfg.fs.BoolVar(
		&cfg.ImgUpWorkerMode, "upload.worker", false, "run as image upload worker mode")

	cfg.fs.IntVar(
		&cfg.Port, "port", 9001, "set server port")

	cfg.fs.IntVar(
		&cfg.ImgProcWorkerCount, "worker.count", runtime.NumCPU(), "set image process worker count")
}

func Parse(args []string) (Cfg, error) {
	err := cfg.fs.Parse(args)
	if err != nil {
		if err != flag.ErrHelp {
			log.Errorf("Invalid command line arguments. Help: %s -h", os.Args[0])
		}
		return Cfg{}, err
	}

	if Env(os.Getenv("NEPHELE_ENV")) == prod {
		cfg.Env = prod
	} else {
		cfg.Env = uat
	}

	return cfg, nil
}

func ApplyCatConfig() {
	if cfg.Env == prod {
		cat.CAT_HOST = cat.PROD
	} else {
		cat.CAT_HOST = cat.UAT
	}
	cat.DOMAIN = cfg.Appid
}

func ApplyLogConfig() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	//output, _ := os.OpenFile(cfg.logPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	output := os.Stdout
	log.SetOutput(output)
	if cfg.verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func GetImageProcConfig() ProcConfig {
	return NewLocalConfig(cfg.Env)
}
