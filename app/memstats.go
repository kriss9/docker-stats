package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
  "stats"
	"github.com/docker/libcontainer/cgroups"
)

type MemoryGroup struct {
}


func (s *MemoryGroup) GetStats(path string, metric *cgroups.Stats) error {
	// Set stats from memory.stat.
	statsFile, err := os.Open(filepath.Join(path, "memory.stat"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer statsFile.Close()

	sc := bufio.NewScanner(statsFile)
	for sc.Scan() {
		t, v, err := stats.GetCgroupParamKeyValue(sc.Text())
		if err != nil {
			return fmt.Errorf("failed to parse memory.stat (%q) - %v", sc.Text(), err)
		}
		metric.MemoryStats.Stats[t] = v
	}

	// Set memory usage and max historical usage.
	value, err := stats.GetCgroupParamUint(path, "memory.usage_in_bytes")
	if err != nil {
		return fmt.Errorf("failed to parse memory.usage_in_bytes - %v", err)
	}
	metric.MemoryStats.Usage = value
	value, err = stats.GetCgroupParamUint(path, "memory.max_usage_in_bytes")
	if err != nil {
		return fmt.Errorf("failed to parse memory.max_usage_in_bytes - %v", err)
	}
	metric.MemoryStats.MaxUsage = value
	value, err = stats.GetCgroupParamUint(path, "memory.failcnt")
	if err != nil {
		return fmt.Errorf("failed to parse memory.failcnt - %v", err)
	}
	metric.MemoryStats.Failcnt = value

	return nil
}


func main()  {
  sysfs_path := "/sys/fs/cgroup/memory/system.slice/docker-8f32af9ace20e3001195dc6388f0f7d00fd6ada516bbcc645dd49f7997b23b61.scope"
	actualStats := *cgroups.NewStats()
  var x MemoryGroup

	err := x.GetStats(sysfs_path, &actualStats)

  fmt.Println("stats", actualStats)
  fmt.Println("err", err)

}
