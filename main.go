package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "./docs"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger" // http-swagger middleware
	lxc "gopkg.in/lxc/go-lxc.v2"
)

// @title LXC HTTP API
// @version 0.1
// @description An api to create and manage lxc thourgh HTTP
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host server.clerc.im:8000
// @BasePath /

const lxcpath string = "/var/lib/lxc"

type Version struct {
	Version string `json:"version" example:"4.0"`
}

type Containers struct {
	Containers []string `json:"containers"`
}

// HTTPClientResp format API client response to JSON
type HTTPClientResp struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"a successfully message"`
}

type apiError struct {
	Error   error
	Message string
	Code    int
}

// APIOptions define all client API options
type APIOptions struct {
	Start        string
	Name         string              `json:"name" example:"dummy"`
	TemplateOpts lxc.TemplateOptions `json:"opts"`
}

// APIDestroyOptions define destroy container options
type APIDestroyOptions struct {
	Force bool `json:"force" example:"true"`
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

// GetVersion godoc
// @Summary Get LXC version
// @Description Return LXC version in used
// @Tags general
// @Produce json
// @Success 200 {object} Version
// @Failure 500 {object} HTTPClientResp
// @Router /version [get]
func GetVersion(w http.ResponseWriter, r *http.Request) *apiError {
	lxcVersion := &Version{
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

// GetContainers godoc
// @Summary Get containers list
// @Description Return list of containers
// @Tags containers
// @Produce json
// @Success 200 {object} Containers
// @Failure 500 {object} HTTPClientResp
// @Router /containers [get]
func GetContainers(w http.ResponseWriter, r *http.Request) *apiError {
	lxcContainers := &Containers{
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

// CreateContainer godoc
// @Summary Create a new container
// @Description Create a new container
// @Accept json
// @Tags container
// @Produce json
// @Param options body lxc.TemplateOptions false "Creation parameters"
// @Success 200 {object} HTTPClientResp
// @Failure 400 {object} HTTPClientResp
// @Failure 500 {object} HTTPClientResp
// @Router /create [post]
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

// DestroyContainer godoc
// @Summary Destroy a container
// @Description Destroy a container
// @Tags container
// @Produce json
// @Param container path string true "Container name"
// @Param force body APIDestroyOptions false "Destroy container even if it running"
// @Success 200 {object} HTTPClientResp
// @Failure 400 {object} HTTPClientResp
// @Failure 500 {object} HTTPClientResp
// @Router /destroy/{container} [delete]
func DestroyContainer(w http.ResponseWriter, r *http.Request) *apiError {
	vars := mux.Vars(r)

	var opts APIDestroyOptions
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

	if opts.Force {
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
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://server.clerc.im:8000/swagger/doc.json"),
	))
	r.Handle("/version", apiHandler(GetVersion)).Methods("GET")
	r.Handle("/containers", apiHandler(GetContainers)).Methods("GET")
	r.Handle("/create", apiHandler(CreateContainer)).Methods("POST")
	r.Handle("/destroy/", apiHandler(DestroyContainer)).Methods("DELETE")
	r.Handle("/destroy/{container}", apiHandler(DestroyContainer)).Methods("DELETE")
	http.Handle("/", r)

	srv := &http.Server{
		Handler: r,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
