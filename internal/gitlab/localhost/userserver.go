// package localhost enable testing a fake local host server
package localhost

import (
	"net/http"
	"net/http/httptest"
)

// Integration test helper
type UserServer struct {
	server  *httptest.Server
	handler *UserHandler
}

func NewUserServer() *UserServer {
	repo := NewInMemoryUserRepository()
	service := NewUserService(repo)
	handler := NewUserHandler(service)
	router := SetupUserRouter(handler)
	server := httptest.NewServer(router)

	return &UserServer{
		server:  server,
		handler: handler,
	}
}

func (ts *UserServer) Close() {
	ts.server.Close()
}

func (ts *UserServer) URL() string {
	return ts.server.URL
}

// Router setup with middleware
func SetupUserRouter(handler *UserHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// Apply logging middleware
	mux.HandleFunc("/users", LoggingMiddleware(handler.CreateUser))
	mux.HandleFunc("/users/", LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if r.URL.Path == "/users/" {
				handler.ListUsers(w, r)
			} else {
				handler.GetUser(w, r)
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	return mux
}
