package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"quebic-faas/common"
	"quebic-faas/types"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	release := make(chan bool)

	apiRouter := router.PathPrefix("/api/").Subrouter()
	apiRouter.HandleFunc("/functions", func(w http.ResponseWriter, r *http.Request) {
		data := []types.Function{}
		writeResponse(w, &data, 200)
	}).Methods("GET")

	apiRouter.HandleFunc("/routes", func(w http.ResponseWriter, r *http.Request) {
		data := []types.Resource{}
		release <- true
		writeResponse(w, &data, 200)
	}).Methods("GET")

	router.PathPrefix(common.HTTPStaticPathPrefix).Handler(http.FileServer(http.Dir(common.HTTPStaticDir)))

	fmt.Println("server running...")
	err := http.ListenAndServe("0.0.0.0:3000", router)
	if err != nil {
		release <- false
		fmt.Printf("quebic-faas-manager failed. error : %v\n", err)
	}
}

func makeErrorResponse(w http.ResponseWriter, status int, cause error) {
	errorResponse := types.ErrorResponse{Status: status, Cause: cause.Error()}
	writeResponse(w, &errorResponse, status)
}

func writeResponse(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeResponse failed %v", err)
	}
}
