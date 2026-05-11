package main

import (
	"net/http"
)

type readinessResponse struct {
	Status string `json:"status"`
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
	res := readinessResponse{
		Status: "ok",
	}
	w.Write([]byte(res.Status))
}
