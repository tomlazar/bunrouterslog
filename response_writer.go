package bunrouterslog

import (
	"net/http"

	"github.com/felixge/httpsnoop"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func wrap(w http.ResponseWriter) responseWriter {
	rw := responseWriter{}
	rw.ResponseWriter = httpsnoop.Wrap(w, httpsnoop.Hooks{
		WriteHeader: func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return func(statusCode int) {
				if rw.statusCode == 0 {
					rw.statusCode = statusCode
				}
				next(statusCode)
			}
		},
		Write: func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return func(b []byte) (int, error) {
				n, err := next(b)
				rw.bytesWritten += n
				return n, err
			}
		},
	})

	return rw
}
