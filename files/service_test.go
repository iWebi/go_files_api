package files
import (
	"testing"
	"net/http/httptest"
	"io"
	"log"
)


var (
	server   *httptest.Server
	reader   io.Reader //Ignore this for now
	baseUrl string
)

func init() {
	session := Session("mongodb://127.0.0.1:27017")
	server = httptest.NewServer(AppHandlers(session))
	log.Println("server.URL = ", server.URL)
}

func TestUploadApi(t *testing.T) {

}
