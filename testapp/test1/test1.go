package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"quebic-faas/common"
	"quebic-faas/types"
	"sort"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	serverStop := make(chan bool)
	var routes = []types.Resource{
		types.Resource{URL: "/u1"},
		types.Resource{URL: "/u2"},
	}
	server(routes, serverStop)

	time.Sleep(30 * time.Second)
	serverStop <- true
	fmt.Printf("New routes got.....\n")

	routes = []types.Resource{
		types.Resource{URL: "/u1"},
		types.Resource{URL: "/u2"},
		types.Resource{URL: "/u3"},
	}
	serverStop <- false
	server(routes, serverStop)

	time.Sleep(30 * time.Second)
	serverStop <- true
	fmt.Printf("New routes got.....\n")

	routes = []types.Resource{
		types.Resource{URL: "/b1"},
		types.Resource{URL: "/b2"},
		types.Resource{URL: "/b3"},
	}
	serverStop <- false
	server(routes, serverStop)

}

func server(routes []types.Resource, serverStop chan bool) {
	router := mux.NewRouter()

	//isWorking := make(chan bool)
	var isWorking bool

	for _, route := range routes {
		log.Printf("setup %v", route.URL)
		router.HandleFunc(route.URL, func(w http.ResponseWriter, r *http.Request) {

			isWorking = true

			log.Printf("start %v", route.URL)
			data := []types.Function{}
			writeResponse(w, &data, 200)
			log.Printf("finish %v", route.URL)

			isWorking = false

		}).Methods("GET")
	}

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {

		log.Printf("start %v", "/health")

		isWorking = true
		data := []types.Function{}
		writeResponse(w, &data, 200)
		isWorking = false

		log.Printf("finish %v", "/health")

	}).Methods("GET")

	go func() {
		log.Printf("server running...")
		err := http.ListenAndServe("0.0.0.0:3000", router)
		if err != nil {
			log.Printf("quebic-faas-manager failed. error : %v", err)
		}
	}()

	select {
	case msg := <-serverStop:
		if msg {
			for isWorking {
				time.Sleep(2 * time.Second)
			}
			log.Printf("server stopped")
		}
	}
}

func _main() {

	router := mux.NewRouter()

	router.HandleFunc("/sort", func(w http.ResponseWriter, r *http.Request) {

		var routes = []types.Resource{
			//users
			{
				ID:            "/users:POST",
				URL:           "/users",
				RequestMethod: "POST",
			},
			{
				ID:            "/users:GET",
				URL:           "/users",
				RequestMethod: "GET",
			},
			{
				ID:            "/users:DELETE",
				URL:           "/users",
				RequestMethod: "DELETE",
			},
			{
				ID:            "/users:PUT",
				URL:           "/users",
				RequestMethod: "PUT",
			},
			{
				ID:            "/users:PATCH",
				URL:           "/users",
				RequestMethod: "PATCH",
			},
			// orders
			{
				ID:            "/orders:POST",
				URL:           "/orders",
				RequestMethod: "POST",
			},
			{
				ID:            "/orders:GET",
				URL:           "/orders",
				RequestMethod: "GET",
			},
			{
				ID:            "/orders:DELETE",
				URL:           "/orders",
				RequestMethod: "DELETE",
			},
			{
				ID:            "/orders:PUT",
				URL:           "/orders",
				RequestMethod: "PUT",
			},
			{
				ID:            "/orders:PATCH",
				URL:           "/orders",
				RequestMethod: "PATCH",
			},
		}

		sortByID := func(r1, r2 *types.Resource) bool {
			return r1.ID < r2.ID
		}

		By(sortByID).Sort(routes)

		writeResponse(w, routes, 200)

	}).Methods("GET")

	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s", r.FormValue("id"))
		writeResponse(w, "H1", 201)
	}).Methods("GET")

	router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")

		const MEMORY = 5 * 1024 * 1024 //5mb
		r.ParseMultipartForm(MEMORY)

		//Artifact
		artifactFile, handler, err := r.FormFile(common.FunctionSaveField_SOURCE)
		if err != nil {
			makeErrorResponse(w, 500, makeError("FormFile %v", err))
			return
		}
		defer artifactFile.Close()

		fmt.Printf("\ngot file %s\n", handler.Filename)

		f, err := os.OpenFile("/home/tharanga/file-upload/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			makeErrorResponse(w, 500, makeError("OpenFile %v", err))
			return
		}
		defer f.Close()

		io.Copy(f, artifactFile)

		//Spec
		specJSON := r.Form.Get(common.FunctionSaveField_SPEC)
		if specJSON == "" {
			makeErrorResponse(w, 500, makeError("specJSON empty", nil))
			return
		}

		functionDTO := &types.FunctionDTO{}
		err = json.Unmarshal([]byte(specJSON), functionDTO)
		if err != nil {
			makeErrorResponse(w, 500, makeError("specJSON Unmarshal failed %v", err))
			return
		}

		writeResponse(w, functionDTO, 201)

	}).Methods("POST")

	err := http.ListenAndServe("127.0.0.1:3000", router)
	if err != nil {
		log.Fatalf("api failed. error : %v", err)
	}

}

func makeError(format string, err error) error {

	if err != nil {
		return fmt.Errorf(format, err)
	}

	return fmt.Errorf(format)

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

// By is the type of a "less" function that defines the ordering of its Planet arguments.
type By func(r1, r2 *types.Resource) bool

type routesSorter struct {
	routes []types.Resource
	by     func(r1, r2 *types.Resource) bool // Closure used in the Less method.
}

// Len is part of sort.Interface.
func (s *routesSorter) Len() int {
	return len(s.routes)
}

// Swap is part of sort.Interface.
func (s *routesSorter) Swap(i, j int) {
	s.routes[i], s.routes[j] = s.routes[j], s.routes[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s *routesSorter) Less(i, j int) bool {
	return s.by(&s.routes[i], &s.routes[j])
}

// Sort is a method on the function type, By, that sorts the argument slice according to the function.
func (by By) Sort(routes []types.Resource) {
	rs := &routesSorter{
		routes: routes,
		by:     by, // The Sort method's receiver is the function (closure) that defines the sort order.
	}
	sort.Sort(rs)
}
