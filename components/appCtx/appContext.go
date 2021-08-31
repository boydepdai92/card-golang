package appCtx

import "gorm.io/gorm"

type AppContext interface {
	GetMainDatabaseConnection() *gorm.DB
	GetPepper() string
	GetSecret() string
	GetExpiresIn() int
}

type appContext struct {
	db        *gorm.DB
	pepper    string
	secret    string
	expiresIn int
}

func NewAppContext(db *gorm.DB, pepper string, secret string, expiresIn int) *appContext {
	return &appContext{db: db, pepper: pepper, secret: secret, expiresIn: expiresIn}
}

func (appCtx *appContext) GetMainDatabaseConnection() *gorm.DB {
	return appCtx.db
}

func (appCtx *appContext) GetPepper() string {
	return appCtx.pepper
}

func (appCtx *appContext) GetSecret() string {
	return appCtx.secret
}

func (appCtx *appContext) GetExpiresIn() int {
	return appCtx.expiresIn
}
