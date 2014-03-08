package main

import (
  "fmt"
  "strconv"
  "net/http"
  "flag"
  "./plugins"
)

var debugPtr = flag.Bool("debug", false, "Enable debug output on console")
var portPtr  = flag.Int("port", 8080, "Port to listen on")
var sslPtr   = flag.Bool("ssl", false, "Enable HTTPS")
var cpemPtr  = flag.String("cpem", "cert.pem", "Path to SSL cert in PEM format")
var kpemPtr  = flag.String("ckey", "key.pem", "Path to SSL key in PEM format")

func main() {
  flag.Parse()

  if *debugPtr == true {fmt.Println("Console debug output enabled.")}

  urlport := ":" + strconv.Itoa(*portPtr)

  http.HandleFunc("/dpkg", dpkgViewHandler)
  http.HandleFunc("/facter", facterViewHandler)
  http.HandleFunc("/ohai", ohaiViewHandler)
  if *sslPtr == true {
    http.ListenAndServeTLS(urlport, *cpemPtr, *kpemPtr, nil)
  } else {
    http.ListenAndServe(urlport, nil)
  }
}

func dpkgViewHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Content-type", "application/json")

  jsonMsg, err := getResponse("dpkg")
  if err != nil {
    http.Error(w, "Oops", http.StatusInternalServerError)
  }
  fmt.Fprintf(w, jsonMsg)
}

func facterViewHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Content-type", "application/json")

  jsonMsg, err := getResponse("Facter")
  if err != nil {
    http.Error(w, "Oops", http.StatusInternalServerError)
  }
  fmt.Fprintf(w, jsonMsg)
}

func ohaiViewHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Content-type", "application/json")

  jsonMsg, err := getResponse("ohai")
  if err != nil {
    http.Error(w, "Oops", http.StatusInternalServerError)
  }
  fmt.Fprintf(w, jsonMsg)
}

func getResponse(plugin string) (string, error){

  jsonMsg, err := plugins.Call(plugin, *debugPtr)

  if err != nil {
    return "", err
  } else {
    return jsonMsg, nil
  }

}
