package files

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"bufio"
	"mime/multipart"
	"io"
	"gopkg.in/mgo.v2"
	"errors"
	"strconv"
)

func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

func bytesTostring(b []byte) string {
	return string(b[:])
}

func getHostAndPort(r *http.Request) string {
	proto := "http"
	if strings.Contains(r.Proto, "HTTPS") {
		proto = "https"
	}
	return proto + "://" + r.Host
}

func getRequestUrl(r *http.Request) string {
	return getHostAndPort(r) + r.URL.Path
}

func decode(r *http.Request, v interface{}) (err error) {
	err = json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		log.Println("Failed to decode request body")
		err = &InvalidInputError{}
	} else {
		log.Println("Decoded object = ", v)
	}
	return
}

func errorResponse(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	b, _ := json.Marshal(map[string]interface{}{
		"message": err.Error(),
	})
	log.Println("Sending error response for \"" + err.Error() + "\" error")
	fmt.Fprintf(w, "%s", bytesTostring(b))
}

func jsonResponse(w http.ResponseWriter, v interface{}, headers map[string]string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	//Any custom headers passed in
	for k, v := range headers {
		w.Header().Set(k, v)
	}

	w.WriteHeader(statusCode)

	if v != nil {
		b, _ := json.Marshal(v)
		fmt.Fprintf(w, "%s", bytesTostring(b))
	}
}

func requestId() string {
	// Go does not have the concept of thread id.
	// Below serves as a Unique request Id to track all the log statements pertaining to current http request processing.
	b10 := []byte("")
	b10 = strconv.AppendInt(b10, time.Now().Unix(), 10)
	return bytesTostring(b10)
}

func writeToGridFile(file multipart.File, gridFile *mgo.GridFile) error {
	reader := bufio.NewReader(file)
	defer func() { file.Close() }()
	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return errors.New("Could not read the input file")
		}
		if n == 0 {
			break
		}
		// write a chunk
		if _, err := gridFile.Write(buf[:n]); err != nil {
			return errors.New("Could not write to GridFs for "+ gridFile.Name())
		}
	}
	gridFile.Close()
	return nil
}