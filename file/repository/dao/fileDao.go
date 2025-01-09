package dao

import "gorm.io/gorm"

type FileDao interface {
}

type FileDaoV1 struct {
	db *gorm.DB
}

func NewFileService(db *gorm.DB) FileDao {
	return &FileDaoV1{
		db: db,
	}
}
