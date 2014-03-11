package main

import (
  "fmt"
  "strconv"
  "net/http"
  "flag"
  "./plugins"
  auth "github.com/abbot/go-http-auth"
)

var debugPtr = flag.Bool("debug", false, "Enable debug output on console")
var portPtr  = flag.Int("port", 8080, "Port to listen on")
var sslPtr   = flag.Bool("ssl", false, "Enable HTTPS")
var authPtr  = flag.Bool("auth", false, "Enable Basic Authentication")
var cpemPtr  = flag.String("cpem", "cert.pem", "Path to SSL cert in PEM format")
var kpemPtr  = flag.String("ckey", "key.pem", "Path to SSL key in PEM format")

func main() {
  flag.Parse()

  if *debugPtr == true {fmt.Println("Console debug output enabled.")}

  urlport := ":" + strconv.Itoa(*portPtr)

  authenticator := auth.NewBasicAuthenticator("sysreport", Secret)

  if *authPtr == true {
    http.HandleFunc("/", authenticator.Wrap(authRootViewHandler))
    http.HandleFunc("/packages", authenticator.Wrap(authPackagesViewHandler))
    http.HandleFunc("/facter", authenticator.Wrap(authFacterViewHandler))
    http.HandleFunc("/ohai", authenticator.Wrap(authOhaiViewHandler))
  } else {
    http.HandleFunc("/", rootViewHandler)
    http.HandleFunc("/packages", packagesViewHandler)
    http.HandleFunc("/facter", facterViewHandler)
    http.HandleFunc("/ohai", ohaiViewHandler)
  }

  if *sslPtr == true {
    http.ListenAndServeTLS(urlport, *cpemPtr, *kpemPtr, nil)
  } else {
    http.ListenAndServe(urlport, nil)
  }
}

func Secret(user, realm string) string {
        if user == "test" {
                // password is "hello"
                return "$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1"
        }
        return ""
}

func Motd() string {

  msg := `<center><pre style="font-family:monospace;font-size:2em;">
  _________            __________                             __   
 /   _____/__.__. _____\______   \ ____ ______   ____________/  |_ 
 \_____  <   |  |/  ___/|       _// __ \\____ \ /  _ \_  __ \   __\
 /        \___  |\___ \ |    |   \  ___/|  |_> >  <_> )  | \/|  |  
/_______  / ____/____  >|____|_  /\___  >   __/ \____/|__|   |__|  
        \/\/         \/        \/     \/|__|                       
</pre></center>`

  return msg

}

func rootViewHandler(w http.ResponseWriter, r *http.Request) {

  headerSet(w, "text/html")

  msg := Motd()

  fmt.Fprintf(w, msg)
}

func authRootViewHandler(w http.ResponseWriter, r *auth.AuthenticatedRequest) {

  headerSet(w, "text/html")

  msg := Motd()

  fmt.Fprintf(w, msg)
}

func packagesViewHandler(w http.ResponseWriter, r *http.Request) {

  headerSet(w, "application/json")

  jsonMsg, err := getResponse("packages")
  if err != nil {
    http.Error(w, "Oops", http.StatusInternalServerError)
  }
  fmt.Fprintf(w, jsonMsg)
}

func authPackagesViewHandler(w http.ResponseWriter, r *auth.AuthenticatedRequest) {

  headerSet(w, "application/json")

  jsonMsg, err := getResponse("packages")
  if err != nil {
    http.Error(w, "Oops", http.StatusInternalServerError)
  }
  fmt.Fprintf(w, jsonMsg)
}

func facterViewHandler(w http.ResponseWriter, r *http.Request) {

  headerSet(w, "application/json")

  jsonMsg, err := getResponse("Facter")
  if err != nil {
    http.Error(w, "Oops", http.StatusInternalServerError)
  }
  fmt.Fprintf(w, jsonMsg)
}

func authFacterViewHandler(w http.ResponseWriter, r *auth.AuthenticatedRequest) {

  headerSet(w, "application/json")

  jsonMsg, err := getResponse("Facter")
  if err != nil {
    http.Error(w, "Oops", http.StatusInternalServerError)
  }
  fmt.Fprintf(w, jsonMsg)
}

func ohaiViewHandler(w http.ResponseWriter, r *http.Request) {

  headerSet(w, "application/json")

  jsonMsg, err := getResponse("ohai")
  if err != nil {
    http.Error(w, "Oops", http.StatusInternalServerError)
  }
  fmt.Fprintf(w, jsonMsg)
}

func authOhaiViewHandler(w http.ResponseWriter, r *auth.AuthenticatedRequest) {

  headerSet(w, "application/json")

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

func headerSet(w http.ResponseWriter, ctype string) http.ResponseWriter {
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Content-type", ctype)

  return w
}
