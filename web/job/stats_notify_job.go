package job

import (
	"fmt"
	"net"
	"os"
	"x-ui/logger"
	"x-ui/web/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)
type StatsNotifyJob struct {
	enable 		   bool
	xrayService    service.XrayService
	inboundService service.InboundService
	settingService service.SettingService
}

func NewStatsNotifyJob() *StatsNotifyJob {
	return new(StatsNotifyJob)
}

func (j *StatsNotifyJob) Run() {
	if !j.xrayService.IsXrayRunning() {
		return
	}
	//Telegram bot basic info
	tgBottoken,err:=j.settingService.GetTgBotToken()
	if err != nil {
		logger.Warning("StatsNotifyJob run failed,GetTgBotToken fail:", err)
		return
	}
	tgBotid,err:=j.settingService.GetTgBotChatId()
	if err != nil {
		logger.Warning("StatsNotifyJob run failed,GetTgBotChatId fail:", err)
		return
	}
	//get traffic 
	inbouds,err := j.inboundService.GetAllInbounds()
	if err != nil {
		logger.Warning("StatsNotifyJob run failed:", err)
		return
	}
	var upTraffic int64
	var downTraffic int64
	var totalTraffic int64
	for _, inbound := range inbouds {
		upTraffic+=inbound.Up
		downTraffic+=inbound.Down
		totalTraffic+=inbound.Total
	}
	upTraffic=upTraffic/(1024*1024)
	downTraffic=downTraffic/(1024*1024)
	totalTraffic=totalTraffic/(1024*1024)
	//get hostname
	name, err := os.Hostname()
	if err != nil {
		fmt.Println("get hostname error:",err)
		return
	}
	//get ip address
	var ip string
    netInterfaces, err := net.Interfaces()
    if err != nil {
        fmt.Println("net.Interfaces failed, err:", err.Error())
        return 
	}
 
    for i := 0; i < len(netInterfaces); i++ {
        if (netInterfaces[i].Flags & net.FlagUp) != 0 {
            addrs, _ := netInterfaces[i].Addrs()
 
            for _, address := range addrs {
                if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
                    if ipnet.IP.To4() != nil {
                        ip=ipnet.IP.String()
						break
                    }else{
						ip=ipnet.IP.String()
						break
					}
                }
            }
        }
    }

	bot, err := tgbotapi.NewBotAPI(tgBottoken)
	if err != nil {
		fmt.Println("get tgbot error:",err)
	}
	bot.Debug = true
	fmt.Printf("Authorized on account %s", bot.Self.UserName)
	var info string
	info=fmt.Sprintf("主机名称:%s\r\n",name)
	info+=fmt.Sprintf("IP地址:%s\r\n",ip)
	info+=fmt.Sprintf("上行流量↑:%dM\r\n下行流量↓:%dM\r\n总流量:%dM\r\n",upTraffic,downTraffic,totalTraffic)
	msg := tgbotapi.NewMessage(int64(tgBotid),info)
	//msg.ReplyToMessageID = int(tgBotid)
	bot.Send(msg)
}
