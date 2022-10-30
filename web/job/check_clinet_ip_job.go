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

)

type CheckClientIpJob struct {
	xrayService    service.XrayService
	inboundService service.InboundService
}
var job *CheckClientIpJob
  
func NewCheckClientIpJob() *CheckClientIpJob {
	job = new(CheckClientIpJob)
	return job
}

func (j *CheckClientIpJob) Run() {
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
	for clientEmail, ips := range InboundClientIps {
		inboundClientIps,err := GetInboundClientIps(clientEmail)
		if(err != nil){
			addInboundClientIps(clientEmail,ips)

		}else{
			updateInboundClientIps(inboundClientIps,clientEmail,ips)
		}
			
	}


}
func GetAccessLogPath() string {
	
    config, err := os.ReadFile("bin/config.json")
    checkError(err)

	jsonConfig := map[string]interface{}{}
    err = json.Unmarshal([]byte(config), &jsonConfig)
    if err != nil {
        logger.Warning(err)
    }
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
// https://codereview.stackexchange.com/a/192954
func Unique(slice []string) []string {
    // create a map with all the values as key
    uniqMap := make(map[string]struct{})
    for _, v := range slice {
        uniqMap[v] = struct{}{}
    }

    // turn the map keys into a slice
    uniqSlice := make([]string, 0, len(uniqMap))
    for v := range uniqMap {
        uniqSlice = append(uniqSlice, v)
    }
    return uniqSlice
}

func GetInboundClientIps(clientEmail string) (*model.InboundClientIps, error) {
	db := database.GetDB()
	InboundClientIps := &model.InboundClientIps{}
	err := db.Model(model.InboundClientIps{}).Where("client_email = ?", clientEmail).First(InboundClientIps).Error
	if err != nil {
		return nil, err
	}
	return InboundClientIps, nil
}
func addInboundClientIps(clientEmail string,ips []string) error {
	inboundClientIps := &model.InboundClientIps{}
    jsonIps, err := json.Marshal(ips)
	checkError(err)

	inboundClientIps.ClientEmail = clientEmail
	inboundClientIps.Ips = string(jsonIps)
	

	db := database.GetDB()
	tx := db.Begin()

	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	err = tx.Save(inboundClientIps).Error
	if err != nil {
		return err
	}
	return nil
}
func updateInboundClientIps(inboundClientIps *model.InboundClientIps,clientEmail string,ips []string) error {

    jsonIps, err := json.Marshal(ips)
	checkError(err)

	inboundClientIps.ClientEmail = clientEmail
	inboundClientIps.Ips = string(jsonIps)
	
	// check inbound limitation
	inbound, _ := GetInboundByEmail(clientEmail)

	limitIpRegx, _ := regexp.Compile(`"limitIp": .+`)

	limitIpMactch := limitIpRegx.FindString(inbound.Settings)
	limitIpMactch =  ss.Split(limitIpMactch, `"limitIp": `)[1]
    limitIp, err := strconv.Atoi(limitIpMactch)


	if(limitIp < len(ips) && limitIp != 0 && inbound.Enable) {

		DisableInbound(inbound.Id)
	}

	db := database.GetDB()
	err = db.Save(inboundClientIps).Error
	if err != nil {
		return err
	}
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
func DisableInbound(id int) error{
	db := database.GetDB()
	result := db.Model(model.Inbound{}).
		Where("id = ? and enable = ?", id, true).
		Update("enable", false)
	err := result.Error
	logger.Warning("disable inbound with id:",id)

	if err == nil {
		job.xrayService.SetToNeedRestart()
	}

	return err
}
