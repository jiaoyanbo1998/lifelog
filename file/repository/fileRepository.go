package repository

import "lifelog-grpc/file/repository/dao"

type FileRepository interface {
}

type FileRepositoryV1 struct {
	fileDao dao.FileDao
}

func NewFileService(fileDao dao.FileDao) FileRepository {
	return &FileRepositoryV1{
		fileDao: fileDao,
	}
}
