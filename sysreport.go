package main

import (
  "fmt"
  "strconv"
  "net/http"
  "flag"
  "./plugins"
)

var facterMsg string
var ohaiMsg string
var debugPtr = flag.Bool("debug", false, "Enable debug output on console")
var portPtr  = flag.Int("port", 8080, "Port to listen on")
var sslPtr   = flag.Bool("ssl", false, "Enable HTTPS")
var cpemPtr  = flag.String("cpem", "cert.pem", "Path to SSL cert in PEM format")
var kpemPtr  = flag.String("ckey", "key.pem", "Path to SSL key in PEM format")

func main() {
  flag.Parse()

  if *debugPtr == true {fmt.Println("Console debug output enabled.")}

  urlport := ":" + strconv.Itoa(*portPtr)

  http.HandleFunc("/ohai", ohaiViewHandler)
  http.HandleFunc("/facter", facterViewHandler)
  if *sslPtr == true {
    http.ListenAndServeTLS(urlport, *cpemPtr, *kpemPtr, nil)
  } else {
    http.ListenAndServe(urlport, nil)
  }
}

func ohaiViewHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Content-type", "application/json")

  jsonMsg, err := getOhaiResponse()
  if err != nil {
    http.Error(w, "Oops", http.StatusInternalServerError)
  }
  fmt.Fprintf(w, jsonMsg)
}

func facterViewHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Content-type", "application/json")

  jsonMsg, err := getFacterResponse()
  if err != nil {
    http.Error(w, "Oops", http.StatusInternalServerError)
  }
  fmt.Fprintf(w, jsonMsg)
}

func getOhaiResponse() (string, error){

  // Note: We pull this once and only once.
  // Should have a timer of some sort.
  if ohaiMsg != "" {
    return ohaiMsg, nil
  }

  jsonMsg, err := plugins.Ohai(*debugPtr)

  if err != nil {
    return "", err
  } else {
    return jsonMsg, nil
  }

}

func getFacterResponse() (string, error){

  // Note: We pull this once and only once.
  // Should have a timer of some sort.
  if facterMsg != "" {
    return facterMsg, nil
  }

  jsonMsg, err := plugins.Facter(*debugPtr)

  if err != nil {
    return "", err
  } else {
    return jsonMsg, nil
  }

}
