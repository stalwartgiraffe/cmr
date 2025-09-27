// package localhost enable testing a fake local host server
package localhost

import (
	"log"
	"net/http"
	"net/http/httptest"
)

// Server local host server for integration testing
type Server struct {
	server  *httptest.Server
	handler *Handler
}

func NewServer() *Server {
	events := NewEventsRepoMem()
	service := NewService(events)
	handler := NewHandler(service)

	s := &Server{
		handler: handler,
	}
	s.Start()
	return s
}

func (ts *Server) Start() {
	router := SetupRouter(ts.handler)
	ts.server = httptest.NewServer(router)
}

func (ts *Server) Close() {
	ts.server.Close()
}

func (ts *Server) URL() string {
	return ts.server.URL
}

// SetupRouter creates the route handlers.
// see the swagger doc
// https://gitlab.com/gitlab-org/gitlab/-/tree/master
// https://gitlab.com/gitlab-org/gitlab/-/blob/master/doc/api/openapi/openapi_v2.yaml
// Router setup with middleware
func SetupRouter(handler *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// Handle the GitLab API v4 events endpoint: /api/v4/users/{id}/events
	mux.HandleFunc("/api/v4/users/", LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.HandleEvents(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Handle the GitLab API v4 groups merge requests endpoint: /api/v4/groups/{id}/merge_requests
	mux.HandleFunc("/api/v4/groups/", LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.HandleMergeRequests(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	return mux
}

// LoggingMiddleware logs events simply
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}
