package hosts

import (
	log "github.com/Sirupsen/logrus"
	cat "github.com/ctripcorp/cat.go"
	"github.com/ctripcorp/nephele/imgsvr/config"
	"github.com/ctripcorp/nephele/util"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Worker struct {
	host                Host
	tryReconnectCounter int
	heartbeatInfo       heartbeatValues
}

func (w *Worker) isDead() bool {
	return w.tryReconnectCounter >= 3
}

func (w *Worker) sick() {
	w.tryReconnectCounter = w.tryReconnectCounter + 1
}

func (w *Worker) healthy() {
	w.tryReconnectCounter = 0
}

type heartbeatValues struct {
	mapValues map[string]string
	urlValues url.Values
}

func (h *heartbeatValues) Get(key string) string {
	var val string
	if h.mapValues != nil {
		val = h.mapValues[key]
	} else {
		val = h.urlValues.Get(key)
	}
	return val
}

type Master struct {
	port        int
	workerCount int
	workers     []Worker
	startErrCh  chan error
}

func NewMaster(cfg config.Cfg) *Master {
	return &Master{port: cfg.Port, workerCount: cfg.ImgProcWorkerCount, startErrCh: make(chan error)}
}

func (m *Master) GetPort() int {
	return m.port
}

func (m *Master) Start(opt StartOpt) <-chan error {
	m.recRebootEvt()
	log.Debugf("start master process. port=%d and workers=%d", m.port, m.workerCount)
	m.workers = m.getWorkers()
	for _, worker := range m.workers {
		worker.host.Start(StartOpt{NewHosting: true})
	}
	go m.doHouseKeeping()
	go m.recHeartbeatEvt()
	return m.startErrCh
}

func (m *Master) Stop() {
	cmd := exec.Command("pkill", "imgsvrd")
	cmd.Output()
}

func (m *Master) ReloadConf() {

}

func (m *Master) getWorkers() []Worker {
	workers := make([]Worker, m.workerCount)
	for i := 1; i <= m.workerCount; i++ {
		workers[i] = Worker{host: &ImageProcessWorker{Port: m.port + i}}
	}
	return workers
}

func (m *Master) recRebootEvt() {
	cat := cat.Instance()
	tran := cat.NewTransaction("System", "Reboot")
	defer func() {
		tran.SetStatus("0")
		tran.Complete()
	}()
	util.LogEvent(cat, "Reboot", util.JoinString(util.GetIP(), ":", strconv.Itoa(m.port)), nil)
}

func (m *Master) doHouseKeeping() {
	t1 := time.NewTimer(time.Duration(time.Second * 8))
	for {
		select {
		case <-t1.C:
			log.Debug("start housekeeping")
			for _, worker := range m.workers {
				if !m.killDeadWorkerAndStartNew(&worker) {
					m.checkWorkerHealth(&worker)
				}
			}
			t1.Reset(time.Duration(time.Second * 2))
		}
	}
}

func (m *Master) killDeadWorkerAndStartNew(w *Worker) bool {
	isDead := w.isDead()
	if isDead {
		log.Debugf("find worker=%d is dead, prepared to restart", w.host.GetPort())
		err := killPort(w.host.GetPort())
		if err != nil {
			log.WithFields(log.Fields{
				"port": w.host.GetPort(),
				"type": "DaemonProcess.KillProcessError",
			}).Error(err.Error())
			util.LogErrorEvent(cat.Instance(), "DaemonProcess.KillProcessError", err.Error())
		} else {
			log.Infof("kill worker=%d ok.", w.host.GetPort())
		}
		w.host.Start(StartOpt{NewHosting: true})
		log.Debugf("restart worker=%d ok.", w.host.GetPort())
		w.healthy()
	}
	return isDead
}

func killPort(port int) error {
	cmd := exec.Command("sh", "-c", util.JoinString("lsof -i:", strconv.Itoa(port), "|grep LISTEN|awk '{print $2}'"))
	bts, err := cmd.Output()
	if err != nil {
		return err
	}
	pid := strings.TrimSpace(string(bts))
	if pid == "" {
		return nil
	}
	id, err := strconv.Atoi(pid)
	if err != nil {
		return err
	}
	p, err := os.FindProcess(id)
	if err != nil {
		return err
	}
	return p.Kill()
}

func (m *Master) checkWorkerHealth(w *Worker) {
	val, err := doHttpCall(util.JoinString(
		"http://127.0.0.1:", strconv.Itoa(w.host.GetPort()), "/heartbeat/"))
	if err != nil {
		w.sick()
		log.WithFields(log.Fields{
			"port": w.host.GetPort(),
			"type": "WorkerProcess.HeartbeatError",
		}).Error(err.Error())
		util.LogErrorEvent(cat.Instance(), "WorkerProcess.HeartbeatError", err.Error())
	} else {
		w.healthy()
		url, err := url.Parse(string(val))
		if err != nil {
			log.Errorf("parse heartbeat from worker=%d failed.", w.host.GetPort())
		} else {
			w.heartbeatInfo = heartbeatValues{mapValues: nil, urlValues: url.Query()}
		}
	}
}

func (m *Master) recHeartbeatEvt() {
	var nextDur time.Duration
	nowSec := time.Now().Second()
	if nowSec > 30 {
		nextDur = time.Duration(90-nowSec) * time.Second
	} else {
		nextDur = time.Duration(30-nowSec) * time.Second
	}

	t1 := time.NewTimer(nextDur)
	for {
		select {
		case <-t1.C:
			catInst := cat.Instance()
			tran := catInst.NewTransaction("System", "Status")
			h := catInst.NewHeartbeat("HeartBeat", util.GetIP())
			combineHeartbeat(h, strconv.Itoa(m.port), heartbeatValues{mapValues: util.GetStatus(), urlValues: nil})
			for _, w := range m.workers {
				combineHeartbeat(h, strconv.Itoa(w.host.GetPort()), w.heartbeatInfo)
			}
			h.SetStatus("0")
			h.Complete()
			tran.SetStatus("0")
			tran.Complete()
			t1.Reset(time.Duration(time.Minute))
		}
	}
}

func combineHeartbeat(h cat.Heartbeat, port string, val heartbeatValues) {
	h.Set("System", util.JoinString("Alloc_", port), val.Get("Alloc"))
	h.Set("System", util.JoinString("TotalAlloc_", port), val.Get("TotalAlloc"))
	h.Set("System", util.JoinString("Sys_", port), val.Get("Sys"))
	h.Set("System", util.JoinString("Mallocs_", port), val.Get("Mallocs"))
	h.Set("System", util.JoinString("Frees_", port), val.Get("Frees"))
	h.Set("System", util.JoinString("OtherSys_", port), val.Get("OtherSys"))
	h.Set("System", util.JoinString("PauseNs_", port), val.Get("PauseNs"))
	h.Set("HeapUsage", util.JoinString("HeapAlloc_", port), val.Get("HeapAlloc"))
	h.Set("HeapUsage", util.JoinString("HeapSys_", port), val.Get("HeapSys"))
	h.Set("HeapUsage", util.JoinString("HeapIdle_", port), val.Get("HeapIdle"))
	h.Set("HeapUsage", util.JoinString("HeapInuse_", port), val.Get("HeapInuse"))
	h.Set("HeapUsage", util.JoinString("HeapReleased_", port), val.Get("HeapReleased"))
	h.Set("HeapUsage", util.JoinString("HeapObjects_", port), val.Get("HeapObjects"))
	h.Set("GC", util.JoinString("NextGC_", port), val.Get("NextGC"))
	h.Set("GC", util.JoinString("LastGC_", port), val.Get("LastGC"))
	h.Set("GC", util.JoinString("NumGC_", port), val.Get("NumGC"))
}

func doHttpCall(url string) ([]byte, error) {
	client := http.Client{Timeout: time.Duration(time.Second)}
	resp, err := client.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bts, nil
}
