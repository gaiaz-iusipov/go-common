package httpservermw

import (
	"net/http"
	"path/filepath"

	httpheader "github.com/gaiaz-iusipov/go-common/http/header"
)

type Header struct{}

func (Header) Add(key, val string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Header().Add(key, val)
			next.ServeHTTP(rw, req)
		})
	}
}

func (mw Header) CacheImmutable(next http.Handler) http.Handler {
	return mw.Add(httpheader.CacheControl, httpheader.CacheControlImmutable)(next)
}

func (Header) ContentTypeByExt(contentTypeByExt map[string]string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			ext := filepath.Ext(req.URL.Path)
			if v, ok := contentTypeByExt[ext]; ok {
				rw.Header().Set(httpheader.ContentType, v)
			}

			next.ServeHTTP(rw, req)
		})
	}
}
