// package localhost enable testing a fake local host server
package localhost

import (
	"log"
	"net/http"
	"net/http/httptest"
)

// Integration test helper
type Server struct {
	server  *httptest.Server
	handler *Handler
	//events  *EventsRepoMem
}

func NewServer() *Server {
	events := NewEventsRepoMem()
	service := NewService(events)
	handler := NewHandler(service)

	return &Server{
		handler: handler,
	}
}

func (ts *Server) Close() {
	ts.server.Close()
}

func (ts *Server) URL() string {
	return ts.server.URL
}

// see the swagger doc
// https://gitlab.com/gitlab-org/gitlab/-/tree/master
// https://gitlab.com/gitlab-org/gitlab/-/blob/master/doc/api/openapi/openapi_v2.yaml
// Router setup with middleware
func SetupRouter(handler *Handler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/events/", LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.HandleEvents(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	return mux
}

// Simple logging middleware
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}
