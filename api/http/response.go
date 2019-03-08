package http

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) responseJSON(w http.ResponseWriter, data interface{}) {
	res, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (h *Handler) responseUnauthorized(w http.ResponseWriter) {
	res, err := json.Marshal(map[string]string{
		"error": "Unauthorized",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (h *Handler) responseErr(w http.ResponseWriter, err error) {
	res, err := json.Marshal(map[string]interface{}{
		"error": err.Error(),
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
