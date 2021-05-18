package service

import (
	"gorm.io/gorm"
	"x-ui/database"
	"x-ui/database/model"
	"x-ui/logger"
)

type UserService struct {
}

func (s *UserService) CheckUser(username string, password string) *model.User {
	db := database.GetDB()

	user := &model.User{}
	err := db.Model(model.User{}).
		Where("username = ? and password = ?", username, password).
		First(user).
		Error
	if err == gorm.ErrRecordNotFound {
		return nil
	} else if err != nil {
		logger.Warning("check user err:", err)
		return nil
	}
	return user
}
