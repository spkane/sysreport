package main

import (
  "fmt"
  "strconv"
  "net/http"
  "encoding/json"
  "os/exec"
  "flag"
  "github.com/spkane/go-utils/jsonutils"
  "github.com/spkane/go-utils/debugtools"
)

var msg string
var debugPtr = flag.Bool("debug", false, "Enable debug output on console")
var portPtr  = flag.Int("port", 8080, "Port to listen on")
var sslPtr   = flag.Bool("ssl", false, "Enable HTTPS")
var cpemPtr  = flag.String("cpem", "cert.pem", "Path to SSL cert in PEM format")
var kpemPtr  = flag.String("ckey", "key.pem", "Path to SSL key in PEM format")

func main() {
  flag.Parse()

  if *debugPtr == true {fmt.Println("Console debug output enabled.")}

  urlport := ":" + strconv.Itoa(*portPtr)

  http.HandleFunc("/facter", viewHandler)
  if *sslPtr == true {
    http.ListenAndServeTLS(urlport, *cpemPtr, *kpemPtr, nil)
  } else {
    http.ListenAndServe(urlport, nil)
  }
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Content-type", "application/json")

  jsonMsg, err := getResponse()
  if err != nil {
    http.Error(w, "Oops", http.StatusInternalServerError)
  }
  fmt.Fprintf(w, jsonMsg)
}

func getResponse() (string, error){

  // Note: We pull this once and only once.
  // Should have a timer of some sort.
  if msg != "" {
    return msg, nil
  }

  out, err := exec.Command("facter","-j","-p").Output()
  if err != nil {
    out2, err2 := exec.Command("facter","-j").Output()
    debugtools.CheckError(err2)
    out = out2
  }

  var input interface{}
  err3 := json.Unmarshal(out, &input)
  debugtools.CheckError(err3)

  msg = jsonutils.JsonBuild(input, *debugPtr)

  if err3 != nil {
    return "", err
  }

  return msg, nil
}
