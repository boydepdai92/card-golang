package main

import (
	"card-warhouse/components/appCtx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
)

func main() {
	dsn := os.Getenv("DB_CONNECTION_STR")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if nil != err {
		log.Fatalln("Can not connect to database")
		return
	}

	pepper := os.Getenv("PEPPER")
	jwtSecret := os.Getenv("JWT_SECRET")
	expiresIn, _ := strconv.Atoi(os.Getenv("JWT_EXPIRES_IN"))

	appContext := appCtx.NewAppContext(db, pepper, jwtSecret, expiresIn)

	Register(appContext)
}
