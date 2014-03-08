package plugins

import (
  "encoding/json"
  "os/exec"
  "strings"
  "bufio"
  "bytes"
  "github.com/spkane/go-utils/strutils"
  "github.com/spkane/go-utils/jsonutils"
  "github.com/spkane/go-utils/debugtools"
)

func Call(plugin string, debug bool) (string, error) {

  ltype := strings.ToLower(plugin)

  switch ltype {
    default: return Facter(debug)
    case "dpkg": return Dpkg(debug)
    case "facter": return Facter(debug)
    case "ohai": return Ohai(debug)
  }

}

func Dpkg(debug bool) (string, error) {

  out, err := exec.Command("dpkg-query", "-W").Output()
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
    jsonMsg += "  " + strutils.QuoteString(s[0]) + " : " + strutils.QuoteString(s[1]) + "\n"

	}

  jsonMsg += "}"
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
