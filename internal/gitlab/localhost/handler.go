package localhost

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleEvents(w http.ResponseWriter, r *http.Request) {
	// CLADE implement the handler specified in openapi.yaml
	events := 123 // CLAUDE implement this
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
