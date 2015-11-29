package files

import "gopkg.in/mgo.v2/bson"

//interface to be implemented by DAO
type Service interface {
	allFiles() ([]FileMetadata, error)
	getById(id bson.ObjectId) (*FileMetadata, error)
	create(FileMetadata) (bson.ObjectId, error)
	delete(id bson.ObjectId) (bool, error)
}
