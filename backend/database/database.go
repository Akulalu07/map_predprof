package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Mars struct {
	ID     int `gorm:"primary_key"`
	API    string
	Points []byte //`gorm:""`
}

var (
	Db *gorm.DB
)

func SaveDb(conn string, points [][]byte) error {
	newdots := make([]byte, 0)
	for _, item := range points {
		newdots = append(newdots, item...)
	}
	mars := Mars{API: conn, Points: newdots}
	err := Db.Create(&mars).Error
	return err
}

func Init(conn string) (*gorm.DB, error) {
	var err error
	Db, err = gorm.Open(sqlite.Open(conn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = Db.AutoMigrate(&Mars{})
	if err != nil {
		return nil, err
	}
	return Db, nil
}
