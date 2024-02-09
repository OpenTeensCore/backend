package main

import (
	"OpenTeens/dao"
	"OpenTeens/model"
	"OpenTeens/router"
)

func main() {
	err := dao.InitMySql()
	if err != nil {
		panic(err)
	}

	err = dao.DB.AutoMigrate(
		&model.UserAccount{},
		&model.UserInfo{},
		&model.UserToken{},
		&model.EmailVerification{},
	)
	if err != nil {
		return
	}

	r := router.SetUpRoute()

	r.Run("0.0.0.0:8088")
}
