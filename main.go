// Package classification LXC HTTP API
//
// That is to provide a detailed overview of the lxc http api specs
//
// Terms Of Service:
//
//     Schemes: http, https
//     Host: localhost:8080
//     Base path: /
//     Version: 0.1
//     License: MIT http://opensource.org/licenses/MIT
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - basic:
//
//     SecurityDefinitions:
//     basic:
//          type: basic
//
// swagger:meta
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux" // http-swagger middleware
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

// Version model
// swagger:model Version
type Version struct {
	// Currently LXC version
	// example: 4.0
	Version string `json:"version"`
}

// Containers model
// swagger:model Containers
type Containers struct {
	// List of active containers
	Containers []string `json:"containers"`
}

// HTTPClientResp format API client response to JSON
// swagger:model HTTPClientResp
type HTTPClientResp struct {
	// Request status
	// example: success
	Status string `json:"status"`

	// Request message
	// example: container created successfully
	Message string `json:"message"`
}

type apiError struct {
	Error   error
	Message string
	Code    int
}

// ContainerTemplate model
// swagger:model ContainerTemplate
type ContainerTemplate struct {
	// Container name
	// required: true
	// example: dummy
	Name string `json:"name"`

	// Defined if container need to be started after created
	// example: true
	Started bool `json:"started"`

	// Container template
	// required: true
	TemplateOpts lxc.TemplateOptions `json:"template"`
}

// DestroyOptions model
// swagger:model DestroyOptions
type DestroyOptions struct {
	// Defined if container need to be stopped before destroy
	// example: true
	Force bool `json:"force"`
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
	var opts ContainerTemplate
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

	if opts.Started {
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

	var opts DestroyOptions
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

	a := middleware.RedocOpts{
		SpecURL: "/swagger/swagger.json",
	}
	configuredRouter := middleware.Redoc(a, r)

	// swagger:operation GET /version general version
	//
	// Return current LXC version
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: Version response
	//     schema:
	//       "$ref": "#/definitions/Version"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/HTTPClientResp"
	r.Handle("/version", apiHandler(GetVersion)).Methods("GET")

	// swagger:operation GET /containers containers containers
	//
	// Return containers list
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: Containers response
	//     schema:
	//       "$ref": "#/definitions/Containers"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/HTTPClientResp"
	r.Handle("/containers", apiHandler(GetContainers)).Methods("GET")

	// Serve swagger json file
	// r.Path("/swagger.json").Handler(http.FileServer(http.Dir("./swagger")))
	r.PathPrefix("/swagger/").Handler(
		http.StripPrefix("/swagger/", http.FileServer(http.Dir("./docs"))))

	// swagger:operation POST /create container create
	//
	// Create a new container
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: template
	//   in: body
	//   description: container template
	//   required: true
	//   schema:
	//     "$ref": "#/definitions/ContainerTemplate"
	// responses:
	//   '200':
	//     description: API response
	//     schema:
	//       "$ref": "#/definitions/HTTPClientResp"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/HTTPClientResp"
	r.Handle("/create", apiHandler(CreateContainer)).Methods("POST")

	// Handle request to destroy endpoint when container name is missing
	r.Handle("/destroy/", apiHandler(DestroyContainer)).Methods("DELETE")

	// swagger:operation DELETE /delete/{container} container delete
	//
	// Delete a container
	// ---
	// produces:
	// - application/json
	// parameters:
	// - name: container
	//   in: path
	//   type: string
	//   required: true
	//   description: Container name
	// - name: force
	//   in: body
	//   schema:
	//     "$ref": "#/definitions/DestroyOptions"
	// responses:
	//   '200':
	//     description: API response
	//     schema:
	//       "$ref": "#/definitions/HTTPClientResp"
	//   default:
	//     description: unexpected error
	//     schema:
	//       "$ref": "#/definitions/HTTPClientResp"
	r.Handle("/destroy/{container}", apiHandler(DestroyContainer)).Methods("DELETE")
	http.Handle("/", r)

	srv := &http.Server{
		Handler: configuredRouter,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
