package files

import (
	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
	"io"
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

//MDbHandler wraps database session to create a copy of the session for every http request
func MDbHandler(out io.Writer, h http.Handler, s *mgo.Session) http.Handler {
	return mdbHandler{out, h, s}
}

type mdbHandler struct {
	writer  io.Writer
	handler http.Handler
	session *mgo.Session
}

func (h mdbHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s := h.session.Copy()
	defer s.Close()
	store := &DataStore{s.Copy()}
	context.Set(req, "store", store)

	h.handler.ServeHTTP(w, req)
}

// Log req Handler attaches a unique request Id to the log statements.
// Go does not have the concept of thread id.
// This handler attaches Unique request Id to track all the log statements pertaining to current http request processing.
func LogReqHandler(h http.Handler) http.Handler {
	return logReqIdHandler{h}
}

type logReqIdHandler struct {
	handler http.Handler
}

func (h logReqIdHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	reqIdPrefix := "[" + requestId() + "]" //attach the request id to all log statements
	log.SetPrefix(reqIdPrefix)
	h.handler.ServeHTTP(w, req)
}

// Final app handlers to serve on http server
// Handlers are wrapped into more handlers for adding respective functionality
func AppHandlers(session *mgo.Session) http.Handler {
	log.SetFlags(log.LstdFlags | log.Lshortfile)


	logger := &lumberjack.Logger{
		Filename:   "/tmp/go_files_api.log",
		MaxSize:    500, // megabytes
		MaxBackups: 20,
		MaxAge:     30, //days
	}
	log.SetOutput(logger)

	r := mux.NewRouter()
	r.HandleFunc("/files/{id}", GetById).Methods("GET")
	r.HandleFunc("/files/path/{path}", Upload).Methods("POST")

	//Adapters
	// Add a RecoveryOption to send 500 error
	recoveryHandler := handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)) //for suppressing panics
	logHandler := handlers.LoggingHandler(logger, r) //fix this. Should log to a file
	dbHandler := MDbHandler(os.Stdout, logHandler, session) //wrap MongoDB session for each handler request
	return LogReqHandler(recoveryHandler(dbHandler))
}