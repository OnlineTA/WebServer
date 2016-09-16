package main

import (
  "encoding/json"
  "fmt"
  "io"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "os/exec"
  "path"
  "strings"
  "strconv"
)

type File struct {
  Name string
  Filename string
}

type Conf struct {
  Require []File
  Try []File
}

var conf Conf

type Uploader struct {
  w http.ResponseWriter
  r *http.Request
  id string
  outdirpath string
}

// http://stackoverflow.com/a/10510783
func exists(path string) (bool, error) {
  _, err := os.Stat(path)

  if err == nil {
    return true, nil
  }

  if os.IsNotExist(err) {
    return false, nil
  }

  return true, err
}

func (u *Uploader) upload(argument, filename string) int {

  file, header, err := u.r.FormFile(argument)
  if err != nil {
    if err.Error() == "http: no such file" {
      fmt.Fprintf(u.w,
        "You forgot to append `%s`.\n", filename)
      return http.StatusBadRequest
    } else {
      log.Println(err)
      return http.StatusInternalServerError
    }
  }
  defer file.Close()

  if header.Filename != filename {
    fmt.Fprintf(u.w,
      "The argument `%s` should have filename `%s`.\n",
      argument, filename)
    return http.StatusBadRequest
  }

  out_file, err := os.Create(path.Join(u.outdirpath, filename))
  if err != nil {
    log.Println(err)
    return http.StatusInternalServerError
  }
  defer out_file.Close()

  _, err = io.Copy(out_file, file)
  if err != nil {
    log.Println(err)
    return http.StatusInternalServerError
  }

  fmt.Fprintf(u.w, "* `%s` uploaded successfully.\n", filename)
  return http.StatusOK
}

func (u *Uploader) grade(id string) {
  outdirpath, err := ioutil.TempDir("uploads", "submission")
  if err != nil {
    log.Println(err)
    u.w.WriteHeader(http.StatusForbidden)
    return
  }
  defer os.RemoveAll(outdirpath)

  u.outdirpath = outdirpath

  var status int

  for _, file := range conf.Require {
    status = u.upload(file.Name, file.Filename)
    if status != http.StatusOK {
      u.w.WriteHeader(status)
      return
    }
  }

  for _, file := range conf.Try {
    _ = u.upload(file.Name, file.Filename)
  }

  cmd := exec.Command(path.Join("tests", id, "run.sh"), outdirpath)

  stdoutPipe, err := cmd.StdoutPipe()
  if err != nil {
    u.w.WriteHeader(http.StatusInternalServerError)
    return
  }

  stderrPipe, err := cmd.StderrPipe()
  if err != nil {
    u.w.WriteHeader(http.StatusInternalServerError)
    return
  }

  go io.Copy(u.w, stdoutPipe)
  go io.Copy(u.w, stderrPipe)

  err = cmd.Start()
  if err != nil {
    log.Println(err)
    return
  }

  err = cmd.Wait()
  if err != nil {
    log.Println(err)
    return
  }
}

func getId(path string) (string, int) {
  path = strings.Trim(path, "/")
  prefix := "grade/"

  if !strings.HasPrefix(path, prefix) {
    return "", http.StatusForbidden
  }

  id := path[len(prefix):]
  _, err := strconv.Atoi(id)
  if err != nil {
    return "", http.StatusForbidden
  }

  return id, http.StatusOK
}

func handler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Server", "OnlineTA")

  u := Uploader{ w, r, "", "" }

  id, code := getId(r.URL.Path)
  if code != http.StatusOK {
    w.WriteHeader(code)
    return
  }

  ok, err := exists(path.Join("tests", id))

  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  if ok {
    u.grade(id)
  } else if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
  } else {
    w.WriteHeader(http.StatusBadRequest)
  }
}

func main() {

  file, err := os.Open("files.json")
  if err != nil {
    panic(err)
  }
  decoder := json.NewDecoder(file)
  err = decoder.Decode(&conf)
  if err != nil {
    panic(err)
  }

  errorLog, err := os.OpenFile("error.log",
    os.O_CREATE | os.O_APPEND | os.O_WRONLY, 0600)
  if err != nil {
    panic(err)
  }
  defer errorLog.Close()

  log.SetOutput(errorLog)

  http.HandleFunc("/", handler)
  http.ListenAndServe(":8080", nil)
}
