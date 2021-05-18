package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"net/http"
	"runtime"
	"time"
	"x-ui/logger"
)

type ProcessState string

const (
	Running ProcessState = "running"
	Stop    ProcessState = "stop"
	Error   ProcessState = "error"
)

type Status struct {
	Cpu float64 `json:"cpu"`
	Mem struct {
		Current uint64 `json:"current"`
		Total   uint64 `json:"total"`
	} `json:"mem"`
	Swap struct {
		Current uint64 `json:"current"`
		Total   uint64 `json:"total"`
	} `json:"swap"`
	Disk struct {
		Current uint64 `json:"current"`
		Total   uint64 `json:"total"`
	} `json:"disk"`
	Xray struct {
		State    ProcessState `json:"state"`
		ErrorMsg string       `json:"errorMsg"`
		Version  string       `json:"version"`
	} `json:"xray"`
	Uptime   uint64    `json:"uptime"`
	Loads    []float64 `json:"loads"`
	TcpCount int       `json:"tcpCount"`
	UdpCount int       `json:"udpCount"`
	NetIO    struct {
		Up   uint64 `json:"up"`
		Down uint64 `json:"down"`
	} `json:"netIO"`
	NetTraffic struct {
		Sent uint64 `json:"sent"`
		Recv uint64 `json:"recv"`
	} `json:"netTraffic"`
}

type Release struct {
	TagName string `json:"tag_name"`
}

func stopServerController(a *ServerController) {
	a.stopTask()
}

type ServerController struct {
	BaseController

	ctx    context.Context
	cancel context.CancelFunc

	lastStatus        *Status
	lastRefreshTime   time.Time
	lastGetStatusTime time.Time
}

func NewServerController(g *gin.RouterGroup) *ServerController {
	ctx, cancel := context.WithCancel(context.Background())
	a := &ServerController{
		ctx:               ctx,
		cancel:            cancel,
		lastGetStatusTime: time.Now(),
	}
	a.initRouter(g)
	go a.runTask()
	runtime.SetFinalizer(a, stopServerController)
	return a
}

func (a *ServerController) initRouter(g *gin.RouterGroup) {
	g.POST("/server/status", a.status)
	g.POST("/server/getXrayVersion", a.getXrayVersion)
}

func (a *ServerController) refreshStatus() {
	status := &Status{}

	now := time.Now()

	percents, err := cpu.Percent(time.Second*2, false)
	if err != nil {
		logger.Warning("get cpu percent failed:", err)
	} else {
		status.Cpu = percents[0]
	}

	upTime, err := host.Uptime()
	if err != nil {
		logger.Warning("get uptime failed:", err)
	} else {
		status.Uptime = upTime
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		logger.Warning("get virtual memory failed:", err)
	} else {
		status.Mem.Current = memInfo.Used
		status.Mem.Total = memInfo.Total
	}

	swapInfo, err := mem.SwapMemory()
	if err != nil {
		logger.Warning("get swap memory failed:", err)
	} else {
		status.Swap.Current = swapInfo.Used
		status.Swap.Total = swapInfo.Total
	}

	distInfo, err := disk.Usage("/")
	if err != nil {
		logger.Warning("get dist usage failed:", err)
	} else {
		status.Disk.Current = distInfo.Used
		status.Disk.Total = distInfo.Total
	}

	avgState, err := load.Avg()
	if err != nil {
		logger.Warning("get load avg failed:", err)
	} else {
		status.Loads = []float64{avgState.Load1, avgState.Load5, avgState.Load15}
	}

	ioStats, err := net.IOCounters(false)
	if err != nil {
		logger.Warning("get io counters failed:", err)
	} else if len(ioStats) > 0 {
		ioStat := ioStats[0]
		status.NetTraffic.Sent = ioStat.BytesSent
		status.NetTraffic.Recv = ioStat.BytesRecv

		if a.lastStatus != nil {
			duration := now.Sub(a.lastRefreshTime)
			seconds := float64(duration) / float64(time.Second)
			up := uint64(float64(status.NetTraffic.Sent-a.lastStatus.NetTraffic.Sent) / seconds)
			down := uint64(float64(status.NetTraffic.Recv-a.lastStatus.NetTraffic.Recv) / seconds)
			status.NetIO.Up = up
			status.NetIO.Down = down
		}
	} else {
		logger.Warning("can not find io counters")
	}

	tcpConnStats, err := net.Connections("tcp")
	if err != nil {
		logger.Warning("get connections failed:", err)
	} else {
		status.TcpCount = len(tcpConnStats)
	}

	udpConnStats, err := net.Connections("udp")
	if err != nil {
		logger.Warning("get connections failed:", err)
	} else {
		status.UdpCount = len(udpConnStats)
	}

	// TODO 临时
	status.Xray.State = Running
	status.Xray.ErrorMsg = ""
	status.Xray.Version = "1.0.0"

	a.lastStatus = status
	a.lastRefreshTime = now
}

func (a *ServerController) runTask() {
	for {
		select {
		case <-a.ctx.Done():
			break
		default:
		}
		now := time.Now()
		if now.Sub(a.lastGetStatusTime) > time.Minute*3 {
			time.Sleep(time.Second * 2)
			continue
		}
		a.refreshStatus()
	}
}

func (a *ServerController) stopTask() {
	a.cancel()
}

func (a *ServerController) status(c *gin.Context) {
	a.lastGetStatusTime = time.Now()

	jsonMsgObj(c, "", a.lastStatus, nil)
}

var lastVersions []string
var lastGetReleaseTime time.Time

func (a *ServerController) getXrayVersion(c *gin.Context) {
	now := time.Now()
	if now.Sub(lastGetReleaseTime) <= time.Minute {
		jsonMsgObj(c, "", lastVersions, nil)
		return
	}
	url := "https://api.github.com/repos/XTLS/Xray-core/releases"
	resp, err := http.Get(url)
	if err != nil {
		jsonMsg(c, "获取版本失败，请稍后尝试", err)
		return
	}

	defer resp.Body.Close()
	buffer := bytes.NewBuffer(make([]byte, 8192))
	buffer.Reset()
	_, err = buffer.ReadFrom(resp.Body)
	if err != nil {
		jsonMsg(c, "获取版本失败，请稍后尝试", err)
		return
	}

	releases := make([]Release, 0)
	err = json.Unmarshal(buffer.Bytes(), &releases)
	if err != nil {
		jsonMsg(c, "获取版本失败，请向作者反馈此问题", err)
		return
	}
	versions := make([]string, 0, len(releases))
	for _, release := range releases {
		versions = append(versions, release.TagName)
	}
	lastVersions = versions
	lastGetReleaseTime = time.Now()

	jsonMsgObj(c, "", versions, nil)
}
