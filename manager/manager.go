package manager

import (
	model "example/api_go/model"
)

type DBManager interface {
	CloseDB()
	model.IUser
	model.IAccessTokenData
}

type DataManager struct {
	DBManager
}

// New creates DataManager object
func New(db DBManager) *DataManager {
	return &DataManager{
		db,
	}
}
