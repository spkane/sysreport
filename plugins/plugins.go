package plugins

import (
  "encoding/json"
  "os/exec"
  "github.com/spkane/go-utils/jsonutils"
  "github.com/spkane/go-utils/debugtools"
)

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
