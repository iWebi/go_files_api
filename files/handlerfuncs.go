package files

import (
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"log"
	"fmt"
	"errors"
	"io"
)

const _24K = (1 << 20) * 24 //memory settings

func GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Println("Received request for downloading file with id ", id)

	gridFs := getStore(r).db().GridFS("fs")
	if gridFsFile, err := gridFs.OpenId(bson.ObjectIdHex(id)); err != nil {
		errorResponse(w, err, http.StatusNotFound)
		return
	} else {
		if written, err := io.Copy(w, gridFsFile); err != nil {
			errorResponse(w, err, http.StatusInternalServerError)
			return
		} else {
			log.Println("Sending ", written, " bytes for ", gridFsFile.Name())
		}
		gridFsFile.Close()
		w.Header().Set("Content-Disposition", "attachment; filename="+gridFsFile.Name())
		w.Header().Set("Content-Type", "application/x-download")
	}
}

func GetByPath(w http.ResponseWriter, r *http.Request) {

}

func GetAll(w http.ResponseWriter, r *http.Request) {

}

func GetContent(w http.ResponseWriter, r *http.Request) {

}

func GetByName(w http.ResponseWriter, r *http.Request) {

}


func Upload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["path"]

	if err := r.ParseMultipartForm(_24K); err != nil {
		log.Println("Error in ParseMultipartForm. ", err.Error())
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	formData := r.MultipartForm.Value
	log.Println(formData)
	inputMd5FormData := formData["md5"]
	if inputMd5FormData == nil || len(inputMd5FormData[0]) == 0 {
		errorResponse(w, &MandatoryError{field:"md5"}, http.StatusBadRequest)
		return
	}
	inputMd5 := inputMd5FormData[0] //MultipartForm.Value stores slices in maps.

	fileMetadata := make(map[string]string)
	store := getStore(r)
	for _, fileHeaders := range r.MultipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()
			fileMetadata["path"] = path
			if gridFile, err := store.db().GridFS("fs").Create(fileHeader.Filename); err != nil {
				errorResponse(w, err, http.StatusInternalServerError)
				return
			} else {
				gridFile.SetMeta(fileMetadata)
				gridFile.SetName(fileHeader.Filename)
				if err := writeToGridFile(file, gridFile); err != nil {
					errorResponse(w, err, http.StatusInternalServerError)
					return
				}
				//if input md5 did not match with the Mongo GridFs md5, return error
				if inputMd5 != gridFile.MD5() {
					log.Println("Input MD5 " + inputMd5 + " did not match with uploaded MD5 " + gridFile.MD5())
					errorResponse(w, errors.New("Upload failed"), http.StatusInternalServerError)
					return
				}
				headers := make(map[string]string)
				if objId, ok := gridFile.Id().(bson.ObjectId); ok {
					headers["location"] = getHostAndPort(r) + fmt.Sprintf(`/files/%x`, string(objId))
				}
				log.Println("File successfully uploaded ", gridFile.Id(), gridFile.Name(), gridFile.MD5())
				jsonResponse(w, nil, headers, http.StatusCreated)
			}
		}
	}
}

func UpdateById(w http.ResponseWriter, r *http.Request) {

}

func UpdateContent(w http.ResponseWriter, r *http.Request) {

}

func UpdateByPath(w http.ResponseWriter, r *http.Request) {

}

func DeleteById(w http.ResponseWriter, r *http.Request) {

}

func DeleteByPath(w http.ResponseWriter, r *http.Request) {

}