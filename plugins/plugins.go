package plugins

import (
  "encoding/json"
  "os/exec"
  "os"
  "errors"
  "strings"
  "bufio"
  "bytes"
  "github.com/spkane/go-utils/strutils"
  "github.com/spkane/go-utils/jsonutils"
  "github.com/spkane/go-utils/debugtools"
)

func Call(plugin string, debug bool) (string, error) {

  ltype := strings.ToLower(plugin)

  out, err := "", error(nil)

  switch ltype {
    default: out, err = Facter(debug)
    case "packages": out, err = Packages(debug)
    case "facter": out, err = Facter(debug)
    case "ohai": out, err = Ohai(debug)
  }

  return out, err

}

func DetermineDistro() (string) {

  file, err := os.Open("/etc/debian_version") // For read access.
  if err == nil {
    return "debian"
  } else {
    file.Close()
  }

  file2, err2 := os.Open("/etc/redhat-release") // For read access.
  if err2 == nil {
    return "redhat"
  } else {
    file2.Close()
  }

  return "unknown"

}

func Packages(debug bool) (string, error) {

  distro := DetermineDistro()

  out, err := []byte{}, error(nil)

  switch distro {
    default: out, err = []byte{0}, errors.New("unknown")
    case "debian":
      out, err = exec.Command("dpkg-query", "-W").Output()
    case "redhat":
      out, err = exec.Command("rpm", "-qa", "--queryformat", "\"%{NAME}\t%{VERSION}\n\"").Output()
  }

  if err != nil {
    return "", err
  }

  reader  := bytes.NewBuffer(out)
  scanner := bufio.NewScanner(reader)
  jsonMsg := "{\n"

  for scanner.Scan() {
    s := strings.Split(scanner.Text(), "\t")
    strutils.QuoteString(s[0])
    strutils.QuoteString(s[1])
    jsonMsg += "  " + strutils.QuoteString(s[0]) + " : " + strutils.QuoteString(s[1]) + ",\n"

	}

  jsonMsg = strutils.TrimSuffix(jsonMsg, ",\n")
  jsonMsg += "\n}"
  return jsonMsg, nil

}

func Facter(debug bool) (string, error) {

  // We could get this as yaml output (--yaml) instead...
  out, err := exec.Command("facter","-j","-p").Output()
  if err != nil {
    out2, err2 := exec.Command("facter","-j").Output()
    debugtools.CheckError(err2)
    out = out2
  }

  var input interface{}
  err3 := json.Unmarshal(out, &input)
  debugtools.CheckError(err3)

  jsonMsg := jsonutils.JsonBuild(input, debug)

  if err3 != nil {
    return "", err
  } else {
    return jsonMsg, nil
  }

}

func Ohai(debug bool) (string, error) {

  out, err := exec.Command("ohai", "-l", "error").Output()

  if err != nil {
    return "", err
  } else {
    return string(out), nil
  }

}
