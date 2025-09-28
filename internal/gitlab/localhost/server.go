// package localhost enable testing a fake local host server
package localhost

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
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

	//"/api/v4/projects:
	//"/api/v4/projects/{id} many endpoints
	mux.HandleFunc("/api/v4/projects/", LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetProjects(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Handle the GitLab API v4 events endpoint: /api/v4/events
	mux.HandleFunc("/api/v4/", LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetEvents(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Handle the GitLab API v4 groups endpoints: /api/v4/groups/{id}/merge_requests and /api/v4/groups/{id}/projects
	mux.HandleFunc("/api/v4/groups/", LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if strings.Contains(r.URL.Path, "/merge_requests") {
				handler.GetGroupsMergeRequests(w, r)
			} else if strings.Contains(r.URL.Path, "/projects") {
				handler.GetGroupsProjects(w, r)
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
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
