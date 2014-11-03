package main

import (
	"encoding/json"
	"fmt"
	"log"
  "path/filepath"
  "regexp"
	"net/http"
	"os"
	"stats"

	"github.com/docker/libcontainer/cgroups"
	"github.com/gorilla/mux"
)

type Metric struct {
	Name string `json:"name"`
	Body string `json:"body"`
	Time int64  `json:"time"`
}

type ContainerList struct {
	Node       string  `json:"node"`
	Containers []string `json:"containers"`
}

func GetContainerList(resp http.ResponseWriter, req *http.Request) {
	body := ContainerList{}

  body.Node = "1.1.1.1"

  sysfs_path := "/sys/fs/cgroup/cpuacct/system.slice/docker-*.scope"
  files, _ := filepath.Glob(sysfs_path)
  for _,element := range files {
    r, _ := regexp.Compile("/sys/fs/cgroup/cpuacct/system.slice/docker-([0-9a-fA-F]+).scope")
    s := r.FindStringSubmatch(element)
    body.Containers = append(body.Containers, s[1])
  }

  fmt.Println(body)

	bodyJson, _ := json.Marshal(body)
	fmt.Fprintf(resp, "%s", string(bodyJson))
}

func GetAllContainerStats(resp http.ResponseWriter, req *http.Request) {
}

func GetContainerStats(resp http.ResponseWriter, req *http.Request) {
}

func pushStats(w http.ResponseWriter, r *http.Request) {
	sysfs_path := "/sys/fs/cgroup/cpuacct/system.slice/docker-1a60793ae082a3672f542c2e79c404a1905a03dafb4e72e631f09acb7c604a1c.scope"
	var metrics cgroups.Stats
	var x stats.CpuacctGroup

	err := x.GetStats(sysfs_path, &metrics)
	if err != nil {
		log.Fatal("GetStats: ", err)
	}
	b, _ := json.Marshal(metrics)
	fmt.Fprintf(w, "%s", string(b))
}


func main() {
	// configure logger
	f, err := os.OpenFile("urls.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	router := mux.NewRouter()

	router.HandleFunc("/containers", GetContainerList)
	router.HandleFunc("/containers/stats", GetAllContainerStats)
	router.HandleFunc("/containers/stats/{id:[0-9a-fA-F]+}", GetContainerStats)

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}
