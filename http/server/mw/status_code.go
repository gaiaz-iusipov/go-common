package httpservermw

import "net/http"

var _ http.ResponseWriter = (*responseWriterWrapper)(nil)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) { rw.statusCode = code }

func StatusCodeFn(handlerFunc func(rw http.ResponseWriter, code int)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rww := &responseWriterWrapper{ResponseWriter: rw, statusCode: http.StatusOK}
			next.ServeHTTP(rww, req)
			handlerFunc(rw, rww.statusCode)
			rw.WriteHeader(rww.statusCode)
		})
	}
}
