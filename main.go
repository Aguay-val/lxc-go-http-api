package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/lxc/go-lxc.v2"
)

const lxcpath string = "/var/lib/lxc"

type version struct {
	Version string `json:"version"`
}

type containers struct {
	Containers []string `json:"containers"`
}

// HTTPClientResp format API client response to JSON
type HTTPClientResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type apiError struct {
	Error   error
	Message string
	Code    int
}

// APIOptions define all client API options
type APIOptions struct {
	Force        string `json:"force"`
	Start        string
	Name         string
	TemplateOpts lxc.TemplateOptions
}

type apiHandler func(http.ResponseWriter, *http.Request) *apiError

func (fn apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		log.Printf("ERROR: %s\n", e.Error)
		jsonError := &HTTPClientResp{
			Status:  "error",
			Message: e.Message}
		js, _ := json.Marshal(jsonError)
		fmt.Println(string(js))
		http.Error(w, string(js), e.Code)
	}
}

// GetVersion return LXC version
func GetVersion(w http.ResponseWriter, r *http.Request) *apiError {
	lxcVersion := &version{
		Version: lxc.Version()}

	js, err := json.Marshal(lxcVersion)

	// err = errors.New("oops")

	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return &apiError{err, "error message", 500}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(js)))
	return nil
}

// GetContainers return names list of active containers
func GetContainers(w http.ResponseWriter, r *http.Request) *apiError {
	lxcContainers := &containers{
		Containers: lxc.ContainerNames(lxcpath)}

	js, err := json.Marshal(lxcContainers)

	if err != nil {
		return &apiError{err, err.Error(), 500}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(js)))

	return nil
}

// CreateContainer create a new container
func CreateContainer(w http.ResponseWriter, r *http.Request) *apiError {
	var opts APIOptions
	err := json.NewDecoder(r.Body).Decode(&opts)

	if err != nil {
		return &apiError{err, err.Error(), 400}
	}

	c, err := lxc.NewContainer(opts.Name, lxcpath)

	if err != nil {
		return &apiError{err, err.Error(), 500}
	}
	defer c.Release()

	if err := c.Create(opts.TemplateOpts); err != nil {
		return &apiError{err, err.Error(), 500}
	}

	if opts.Start == "yes" {
		if err := c.Start(); err != nil {
			return &apiError{err, err.Error(), 500}
		}
	}

	jsonResp := &HTTPClientResp{
		Status:  "success",
		Message: "container created"}
	js, _ := json.Marshal(jsonResp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(string(js)))

	return nil
}

// DestroyContainer destroy a container
func DestroyContainer(w http.ResponseWriter, r *http.Request) *apiError {
	vars := mux.Vars(r)

	var opts APIOptions
	err := json.NewDecoder(r.Body).Decode(&opts)

	if err != nil {
		return &apiError{err, err.Error(), 400}
	}

	if vars["container"] == "" {
		var err error
		return &apiError{err, "no container name passed", 500}
	}

	c, err := lxc.NewContainer(vars["container"], lxcpath)

	if err != nil {
		return &apiError{err, err.Error(), 400}
	}

	if opts.Force == "yes" {
		err := c.Stop()
		if err != nil {
			return &apiError{err, err.Error(), 500}
		}
	}

	err = c.Destroy()

	if err != nil {
		return &apiError{err, err.Error(), 400}
	} else {
		jsonResp := &HTTPClientResp{
			Status:  "success",
			Message: "container destroyed"}
		js, _ := json.Marshal(jsonResp)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(string(js)))

		return nil
	}

	return nil
}

func main() {
	r := mux.NewRouter()
	r.Handle("/version", apiHandler(GetVersion)).Methods("GET")
	r.Handle("/containers", apiHandler(GetContainers)).Methods("GET")
	r.Handle("/create", apiHandler(CreateContainer)).Methods("POST")
	r.Handle("/destroy/", apiHandler(DestroyContainer)).Methods("DELETE")
	r.Handle("/destroy/{container}", apiHandler(DestroyContainer)).Methods("DELETE")
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
