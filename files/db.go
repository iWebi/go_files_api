package files

import (
	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

type DataStore struct {
	session *mgo.Session
}

func Session(ip string) *mgo.Session {
	s, err := mgo.Dial(ip)
	if err != nil {
		panic(err)
	}
	return s
}

func getStore(r *http.Request) *DataStore {
	return context.Get(r, "store").(*DataStore)
}

func (ds *DataStore) db() *mgo.Database {
	return ds.session.DB("GoFiles")
}

func (ds *DataStore) filesCol() *mgo.Collection {
	return ds.db().C("files")
}

func (s *DataStore) getById(id bson.ObjectId) (*FileMetadata, error) {
	file := &FileMetadata{}
	e := s.filesCol().FindId(id).One(file)
	return file, e
}

func (ds *DataStore) create(file *FileMetadata) (bson.ObjectId, error) {
	file.Id = bson.NewObjectId()
	file.Uploaded = time.Now()
	return file.Id, ds.filesCol().Insert(file)
}
