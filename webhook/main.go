package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	admission "github.com/scanyourkube/webhook/admitter"

	log "github.com/sirupsen/logrus"
	admissionv1 "k8s.io/api/admission/v1"
)

func init() {
	log.SetLevel(log.DebugLevel)

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {

	cert, err := tls.LoadX509KeyPair("/etc/opt/tls.crt", "/etc/opt/tls.key")
	if err != nil {
		panic(err)
	}
	// handle our core application
	log.Info("starting server")

	setupHealthCheck()
	server := http.Server{
		Addr: fmt.Sprintf(":%d", 443),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}
	http.HandleFunc("/mutate", ServeMutate)
	log.Info("get server server")
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func setupHealthCheck() {
	go func() {
		log.Info("starting health check")
		healthMux := http.NewServeMux()
		healthMux.HandleFunc("/health", ServeHealth)
		http.ListenAndServe(":8080", healthMux)
	}()
}

// ServeHealth returns 200 when things are good
func ServeHealth(w http.ResponseWriter, r *http.Request) {
	log.Info("Health check")
	fmt.Fprint(w, "OK")
}

// ServeMutate returns an admission review with deployment or statefulset mutations as a json patch
// in the review response
func ServeMutate(w http.ResponseWriter, r *http.Request) {

	in, err := parseRequest(*r)
	if err != nil {
		log.WithFields(log.Fields{
			"in":    in,
			"error": err,
		}).Fatal("Failed to parse request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Info("parsed request")

	var adm admission.IAdmitter

	switch in.Request.Kind.Kind {
	case "Deployment":
		adm = admission.DeploymentAdmitter{}
	case "StatefulSet":
		adm = admission.StatefulSetAdmitter{}
	}

	out, err := adm.MutateAdmissionReview(in.Request)
	if err != nil {
		log.WithFields(log.Fields{
			"out":   out,
			"error": err,
		}).Fatalf("Failed to mutate %s", in.Kind)
		e := fmt.Sprintf("could not generate admission response: %v", err)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	log.Infof("Mutated %s", in.Kind)

	w.Header().Set("Content-Type", "application/json")
	jout, err := json.Marshal(out)
	if err != nil {
		log.WithFields(log.Fields{
			"jout":  jout,
			"error": err,
		}).Fatal("Failed to parse response")
		e := fmt.Sprintf("could not parse admission response: %v", err)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", jout)
}

// parseRequest extracts an AdmissionReview from an http.Request if possible
func parseRequest(r http.Request) (*admissionv1.AdmissionReview, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("Content-Type: %q should be %q",
			r.Header.Get("Content-Type"), "application/json")
	}

	bodybuf := new(bytes.Buffer)
	bodybuf.ReadFrom(r.Body)
	body := bodybuf.Bytes()

	if len(body) == 0 {
		return nil, fmt.Errorf("admission request body is empty")
	}

	var a admissionv1.AdmissionReview

	if err := json.Unmarshal(body, &a); err != nil {
		return nil, fmt.Errorf("could not parse admission review request: %v", err)
	}

	if a.Request == nil {
		return nil, fmt.Errorf("admission review can't be used: Request field is nil")
	}

	return &a, nil
}
