package main    

import (
    "fmt"
    "path/filepath"
    "regexp"
)

func main() {
  sysfs_path := "/sys/fs/cgroup/cpuacct/system.slice/docker-*.scope"
  files, _ := filepath.Glob(sysfs_path)
	for _,element := range files {
    fmt.Println("container=", element) // contains a list of all files in the current directory
    r, _ := regexp.Compile("/sys/fs/cgroup/cpuacct/system.slice/docker-([0-9a-fA-F]+).scope")
    s := r.FindStringSubmatch(element)
    fmt.Println("match=", s[1])
	}
}
