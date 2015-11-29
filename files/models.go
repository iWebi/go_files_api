package files

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type FileMetadata struct {
	Id       bson.ObjectId `json:"id" bson:"_id"`
	Name     string        `json:"name"`
	Md5      string        `json:"md5"`
	Uploaded time.Time     `json:"uploaded"`
}

func (f *FileMetadata) String() string {
	b, err := json.Marshal(f)
	if err != nil {
		return ""
	}
	return string(b[:])
}

func (f *FileMetadata) OK() (err error) {
	switch {
	case len(f.Name) == 0:
		err = &MandatoryError{field: "name"}
	case len(f.Md5) == 0:
		err = &MandatoryError{field: "md5"}
	}
	return
}
