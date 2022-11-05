package job

import (
	"x-ui/logger"
	"x-ui/web/service"
	"x-ui/database"
	"x-ui/database/model"
    "os"
 	ss "strings"
	"regexp"
    "encoding/json"
	"gorm.io/gorm"
    "strconv"
	"strings"
	"time"
	"net"
 	"github.com/go-cmd/cmd"
	"sort"
)

type CheckClientIpJob struct {
	xrayService    service.XrayService
	inboundService service.InboundService
}
var job *CheckClientIpJob
var disAllowedIps []string 

func NewCheckClientIpJob() *CheckClientIpJob {
	job = new(CheckClientIpJob)
	return job
}

func (j *CheckClientIpJob) Run() {
	logger.Debug("Check Client IP Job...")
	processLogFile()
}

func processLogFile() {
	accessLogPath := GetAccessLogPath()
	if(accessLogPath == "") {
		logger.Warning("xray log not init in config.json")
		return
	}

    data, err := os.ReadFile(accessLogPath)
	InboundClientIps := make(map[string][]string)
    checkError(err)

	// clean log
	if err := os.Truncate(GetAccessLogPath(), 0); err != nil {
		checkError(err)
	}
	
	lines := ss.Split(string(data), "\n")
	for _, line := range lines {
		ipRegx, _ := regexp.Compile(`[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+`)
		emailRegx, _ := regexp.Compile(`email:.+`)

		matchesIp := ipRegx.FindString(line)
		if(len(matchesIp) > 0) {
			ip := string(matchesIp)
			if( ip == "127.0.0.1" || ip == "1.1.1.1") {
				continue
			}

			matchesEmail := emailRegx.FindString(line)
			if(matchesEmail == "") {
				continue
			}
			matchesEmail = ss.Split(matchesEmail, "email: ")[1]
	
			if(InboundClientIps[matchesEmail] != nil) {
				if(contains(InboundClientIps[matchesEmail],ip)){
					continue
				}
				InboundClientIps[matchesEmail] = append(InboundClientIps[matchesEmail],ip)

				

			}else{
			InboundClientIps[matchesEmail] = append(InboundClientIps[matchesEmail],ip)
		}
		}

	}
	err = ClearInboudClientIps()
	if err != nil {
		return
	}

	var inboundsClientIps []*model.InboundClientIps
	disAllowedIps = []string{}

	for clientEmail, ips := range InboundClientIps {
		inboundClientIps := GetInboundClientIps(clientEmail, ips)
		if inboundClientIps != nil {
			inboundsClientIps = append(inboundsClientIps, inboundClientIps)
		}
	}

	err = AddInboundsClientIps(inboundsClientIps)
	checkError(err)

	// check if inbound connection is more than limited ip and drop connection
	LimitDevice := func() { LimitDevice() }

	stop := schedule(LimitDevice, 700 *time.Millisecond)
	time.Sleep(60 * time.Second)
	stop <- true
 
}
func GetAccessLogPath() string {
	
    config, err := os.ReadFile("bin/config.json")
    checkError(err)

	jsonConfig := map[string]interface{}{}
    err = json.Unmarshal([]byte(config), &jsonConfig)
	checkError(err)
	if(jsonConfig["log"] != nil) {
		jsonLog := jsonConfig["log"].(map[string]interface{})
		if(jsonLog["access"] != nil) {

			accessLogPath := jsonLog["access"].(string)

			return accessLogPath
		}
	}
	return ""

}
func checkError(e error) {
    if e != nil {
		logger.Warning("client ip job err:", e)
	}
}
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func ClearInboudClientIps() error {
	db := database.GetDB()
	err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.InboundClientIps{}).Error
	checkError(err)
	return err
}

func GetInboundClientIps(clientEmail string, ips []string) *model.InboundClientIps {
	jsonIps, err := json.Marshal(ips)
	if err != nil {
		return nil
	}

	inboundClientIps := &model.InboundClientIps{}
	inboundClientIps.ClientEmail = clientEmail
	inboundClientIps.Ips = string(jsonIps)

	inbound, err := GetInboundByEmail(clientEmail)
	if err != nil {
		return nil
	}
	limitIpRegx, _ := regexp.Compile(`"limitIp": .+`)
	limitIpMactch := limitIpRegx.FindString(inbound.Settings)
	limitIpMactch =  ss.Split(limitIpMactch, `"limitIp": `)[1]
    limitIp, err := strconv.Atoi(limitIpMactch)
	if err != nil {
		return nil
	}

	if(limitIp < len(ips) && limitIp != 0 && inbound.Enable) {

		if(limitIp == 1){
			limitIp = 2
		}
		disAllowedIps = append(disAllowedIps,ips[limitIp - 1:]...)

	}
	logger.Debug("disAllowedIps ",disAllowedIps)
    sort.Sort(sort.StringSlice(disAllowedIps))

	return inboundClientIps
}

func AddInboundsClientIps(inboundsClientIps []*model.InboundClientIps) error {
	if inboundsClientIps == nil || len(inboundsClientIps) == 0 {
		return nil
	}
	db := database.GetDB()
	tx := db.Begin()

	err := tx.Save(inboundsClientIps).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func GetInboundByEmail(clientEmail string) (*model.Inbound, error) {
	db := database.GetDB()
	var inbounds *model.Inbound
	err := db.Model(model.Inbound{}).Where("settings LIKE ?", "%" + clientEmail + "%").Find(&inbounds).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return inbounds, nil
}

func LimitDevice(){
	
	localIp,err := LocalIP()
	checkError(err)

	c := cmd.NewCmd("bash","-c","ss --tcp | grep -E '" + IPsToRegex(localIp) + "'| awk '{if($1==\"ESTAB\") print $4,$5;}'","| sort | uniq -c | sort -nr | head")

	<-c.Start()
	if len(c.Status().Stdout) > 0 {
		ipRegx, _ := regexp.Compile(`[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+`)
		portRegx, _ := regexp.Compile(`(?:(:))([0-9]..[^.][0-9]+)`)

		for _, row := range c.Status().Stdout {
			
			data := strings.Split(row," ")
			
			destIp,destPort,srcIp,srcPort := "","","",""
 

			destIp = string(ipRegx.FindString(data[0]))

			destPort = portRegx.FindString(data[0])
			destPort = strings.Replace(destPort,":","",-1)
			
			
			srcIp = string(ipRegx.FindString(data[1]))

			srcPort = portRegx.FindString(data[1])
			srcPort = strings.Replace(srcPort,":","",-1)

			if(contains(disAllowedIps,srcIp)){
				dropCmd := cmd.NewCmd("bash","-c","ss -K dport = " + srcPort)
				dropCmd.Start()

				logger.Debug("request droped : ",srcIp,srcPort,"to",destIp,destPort)
			} 
		}
	}

}

func LocalIP() ([]string, error) {
	// get machine ips

	ifaces, err := net.Interfaces()
	ips := []string{}
	if err != nil {
		return ips, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return ips, err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			ips = append(ips,ip.String())
			
		}
	}
	logger.Debug("System IPs : ",ips)

	return ips, nil
}


func IPsToRegex(ips []string) (string){

	regx := ""
	for _, ip := range ips {
		regx += "(" + strings.Replace(ip, ".", "\\.", -1) + ")"

	}
	regx = "(" + strings.Replace(regx, ")(", ")|(.", -1)  + ")"

	return regx
}

func schedule(LimitDevice func(), delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			LimitDevice()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}
