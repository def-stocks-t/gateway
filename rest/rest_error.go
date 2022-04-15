package rest

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
)

// JSON is a map alias, just for convenience
type JSON map[string]interface{}

// All error codes for UI mapping and translation
const (
	ErrInternal       = 0 // any internal error
	ErrObjectNotFound = 1 // can't find object
	ErrDecode         = 2 // failed to unmarshal incoming request
	ErrNoAccess       = 3 // rejected by auth
	ErrBadRequest     = 4 // bad request
	ErrValidation     = 5 // validation error
)

// SendErrorJSON makes {error: '', details: ''} json body and responds with error code
func SendErrorJSON(w http.ResponseWriter, r *http.Request, httpStatusCode int, err error, details string, errCode int) {
	log.WithError(err).Errorf("%s", errDetailsMsg(r, httpStatusCode, err, details, errCode))
	render.Status(r, httpStatusCode)
	render.JSON(w, r, JSON{"error": err.Error(), "details": details, "code": errCode})
}

func errDetailsMsg(r *http.Request, httpStatusCode int, err error, details string, errCode int) string {
	q := r.URL.String()
	if qun, e := url.QueryUnescape(q); e == nil {
		q = qun
	}

	srcFileInfo := ""
	if pc, file, line, ok := runtime.Caller(2); ok {
		filenameElems := strings.Split(file, "/")
		funcNameElems := strings.Split(runtime.FuncForPC(pc).Name(), "/")
		srcFileInfo = fmt.Sprintf(" [caused by %s:%d %s]", strings.Join(filenameElems[len(filenameElems)-3:], "/"),
			line, funcNameElems[len(funcNameElems)-1])
	}

	remoteIP := r.RemoteAddr
	if pos := strings.Index(remoteIP, ":"); pos >= 0 {
		remoteIP = remoteIP[:pos]
	}
	return fmt.Sprintf("%s - %v - %d (%d) - %s - %s%s", details, err, httpStatusCode, errCode, remoteIP, q, srcFileInfo)
}
