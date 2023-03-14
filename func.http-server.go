package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func convertYAMLtoJSON(yamlFileName string) ([]byte, error) {
	// Read the contents of the YAML file
	yamlData, err := ioutil.ReadFile(yamlFileName)
	if err != nil {
		return nil, err
	}

	// Convert the YAML data to JSON
	jsonData, err := yaml.YAMLToJSON(yamlData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func convertJSONtoYAML(jsonPath string) error {
	// Read JSON data from file
	jsonData, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	// Convert JSON data to YAML
	yamlData, err := yaml.JSONToYAML(jsonData)
	if err != nil {
		return err
	}

	// Write YAML data to file
	err = ioutil.WriteFile("/etc/go-zones/server.yml", yamlData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func updateJSONHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the incoming JSON data from the request body.
	var newData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&newData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Read the YAML file into a byte slice
	jsonData, err := convertYAMLtoJSON("/etc/go-zones/server.yml")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonData))

	// Merge the parsed JSON data with the existing JSON data.
	var existingData map[string]interface{}
	err = json.Unmarshal(jsonData, &existingData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for k, v := range newData {
		existingData[k] = v
	}

	// Write the updated JSON data back to the file.
	updatedData, err := json.Marshal(existingData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ioutil.WriteFile("/etc/go-zones/dns.temp.json", updatedData, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	convertJSONtoYAML("/etc/go-zones/dns.temp.json")

	// Send a response indicating that the JSON was successfully updated.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("JSON updated successfully"))
}

func dnsConfigsHandler(w http.ResponseWriter, r *http.Request) {
	// Read the YAML data from a file
	yamlData, err := ioutil.ReadFile("/etc/go-zones/server.yml")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Convert YAML to JSON
	jsonData, err := yaml.YAMLToJSON(yamlData)
	if err != nil {
		http.Error(w, "Failed to convert YAML to JSON", http.StatusInternalServerError)
		return
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON data to the response writer
	w.Write(jsonData)
}

// NewRouter generates the router used in the HTTP Server
func NewRouter(basePath string) *http.ServeMux {
	if basePath == "" {
		basePath = "/" + appName
	}
	// Create router and define routes and return that router
	router := http.NewServeMux()

	// Version Output - reads from variables.go
	router.HandleFunc(basePath+"/version", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+basePath+"/version", r)
		fmt.Fprintf(w, appName+" version: %s\n", appVersion)
	})

	// Healthz endpoint for kubernetes platforms
	router.HandleFunc(basePath+"/healthz", func(w http.ResponseWriter, r *http.Request) {
		logNeworkRequestStdOut(r.Method+" "+basePath+"/healthz", r)
		fmt.Fprintf(w, "OK")
	})

	router.HandleFunc("/dns-config", dnsConfigsHandler)
	router.HandleFunc("/update-dns-config", updateJSONHandler)

	return router
}

// RunHTTPServer will run the HTTP Server
func (config Config) RunHTTPServer() {
	// Set up a channel to listen to for interrupt signals
	var runChan = make(chan os.Signal, 1)

	// Set up a context to allow for graceful server shutdowns in the event
	// of an OS interrupt (defers the cancel just in case)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		config.Application.Server.Timeout.Server,
	)
	defer cancel()

	// Define server options
	server := &http.Server{
		Addr:         config.Application.Server.Host + ":" + config.Application.Server.Port,
		Handler:      NewRouter(config.Application.Server.BasePath),
		ReadTimeout:  config.Application.Server.Timeout.Read * time.Second,
		WriteTimeout: config.Application.Server.Timeout.Write * time.Second,
		IdleTimeout:  config.Application.Server.Timeout.Idle * time.Second,
	}

	// Only listen on IPV4
	l, err := net.Listen("tcp4", config.Application.Server.Host+":"+config.Application.Server.Port)
	check(err)

	// Handle ctrl+c/ctrl+x interrupt
	signal.Notify(runChan, os.Interrupt, syscall.SIGTSTP)

	// Alert the user that the server is starting
	log.Printf("Server is starting on %s\n", server.Addr)

	// Run the server on a new goroutine
	go func() {
		//if err := server.ListenAndServe(); err != nil {
		if err := server.Serve(l); err != nil {
			if err == http.ErrServerClosed {
				// Normal interrupt operation, ignore
			} else {
				log.Fatalf("Server failed to start due to err: %v", err)
			}
		}
	}()

	// Block on this channel listeninf for those previously defined syscalls assign
	// to variable so we can let the user know why the server is shutting down
	interrupt := <-runChan

	// If we get one of the pre-prescribed syscalls, gracefully terminate the server
	// while alerting the user
	log.Printf("Server is shutting down due to %+v\n", interrupt)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server was unable to gracefully shutdown due to err: %+v", err)
	}
}
