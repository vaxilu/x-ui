package service

import (
	_ "embed"
	"strconv"
	"strings"
	"x-ui/database"
	"x-ui/database/model"
	"x-ui/logger"
	"x-ui/util/random"
)

//go:embed config.json
var xrayTemplateConfig string

type SettingService struct {
}

func (s *SettingService) ClearSetting() error {
	db := database.GetDB()
	return db.Delete(model.Setting{}).Error
}

func (s *SettingService) getSetting(key string) (*model.Setting, error) {
	db := database.GetDB()
	setting := &model.Setting{}
	err := db.Model(model.Setting{}).Where("key = ?", key).First(setting).Error
	if err != nil {
		return nil, err
	}
	return setting, nil
}

func (s *SettingService) saveSetting(key string, value string) error {
	setting, err := s.getSetting(key)
	db := database.GetDB()
	if database.IsNotFound(err) {
		return db.Create(&model.Setting{
			Key:   key,
			Value: value,
		}).Error
	} else if err != nil {
		return err
	}
	setting.Key = key
	setting.Value = value
	return db.Save(setting).Error
}

func (s *SettingService) getString(key string, defaultValue string) (string, error) {
	setting, err := s.getSetting(key)
	if database.IsNotFound(err) {
		return defaultValue, nil
	} else if err != nil {
		return "", err
	}
	return setting.Value, nil
}

func (s *SettingService) getInt(key string, defaultValue int) (int, error) {
	str, err := s.getString(key, strconv.Itoa(defaultValue))
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(str)
}

func (s *SettingService) GetXrayConfigTemplate() (string, error) {
	return s.getString("xray_template_config", xrayTemplateConfig)
}

func (s *SettingService) GetListen() (string, error) {
	return s.getString("web_listen", "")
}

func (s *SettingService) GetPort() (int, error) {
	return s.getInt("web_port", 65432)
}

func (s *SettingService) GetCertFile() (string, error) {
	return s.getString("web_cert_file", "")
}

func (s *SettingService) GetKeyFile() (string, error) {
	return s.getString("web_key_file", "")
}

func (s *SettingService) GetSecret() ([]byte, error) {
	seq := random.Seq(32)
	secret, err := s.getString("secret", seq)
	if secret == seq {
		err := s.saveSetting("secret", secret)
		if err != nil {
			logger.Warning("save secret failed:", err)
		}
	}
	return []byte(secret), err
}

func (s *SettingService) GetBasePath() (string, error) {
	basePath, err := s.getString("web_base_path", "/")
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}
	if !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}
	return basePath, nil
}
