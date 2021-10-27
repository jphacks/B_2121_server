package session

import "net/http"

// ResponseWriter for extracting http header.
type fakeResponseWriter struct {
	headers http.Header
}

func (f *fakeResponseWriter) Header() http.Header {
	return f.headers
}

func (f *fakeResponseWriter) Write(_ []byte) (int, error) {
	panic("implement me")
}

func (f *fakeResponseWriter) WriteHeader(_ int) {
	panic("implement me")
}
