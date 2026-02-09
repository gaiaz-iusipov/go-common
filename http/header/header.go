package httpheader

const (
	// WWWAuthenticate response header advertises the HTTP authentication methods (or challenges)
	// that might be used to gain access to a specific resource.
	WWWAuthenticate = "Www-Authenticate"
)

const (
	CacheControl          = "Cache-Control"
	CacheControlImmutable = "public, max-age=31536000, immutable"
)

const (
	ContentType     = "Content-Type"
	ContentTypeJSON = "application/json"
	ContentTypeHTML = "text/html; charset=utf-8"
	ContentTypeJS   = "text/javascript; charset=utf-8"
	ContentTypeIcon = "image/x-icon"
)

const (
	TraceID = "X-Trace-Id"
)
