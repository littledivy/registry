package handler

import (
	"fmt"
	"net/http"
	"io/ioutil"
	 "strings"
)

const registryURL = "https://x.nest.land"

var modules = []string{"autopilot", "sass", "css", "audio"}

func containsModule(slice []string, item string) bool {
    set := make(map[string]struct{}, len(slice))
    for _, s := range slice {
        set[s] = struct{}{}
    }

    _, ok := set[item] 
    return ok
}

func validateModule(path string) string {
  parts := strings.Split(path, "/")
  if len(parts) > 1 {
	name := strings.Split(parts[1], "@")[0]
	return name
  }
  return ""
}

func Handler(w http.ResponseWriter, r *http.Request) {
	 if r.Method != http.MethodGet {
        http.Error(w, "only GET allowed", http.StatusMethodNotAllowed)
        return
    }
    moduleName := validateModule(r.URL.Path)
    if moduleName == "" || !containsModule(modules, moduleName) {
     	w.Write([]byte("Not found"))
     	return
    } 

    urlstr := fmt.Sprintf("%s%s", registryURL, r.URL.Path)
    fmt.Println(urlstr)
    // request the proxy url
    resp, err := http.Get(urlstr)
    if err != nil {
        http.Error(w, fmt.Sprintf("error creating request to %s", urlstr), http.StatusInternalServerError)
        return
    }
    // make sure body gets closed when this function exits
    defer resp.Body.Close()

    // read entire response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        http.Error(w, "error reading response body", http.StatusInternalServerError)
        return
    }

    // write status code and body from proxy request into the answer
    w.WriteHeader(resp.StatusCode)
    w.Write(body)
}