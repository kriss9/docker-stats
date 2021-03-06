package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"
  "stats"
	"github.com/docker/libcontainer/cgroups"
	"github.com/docker/libcontainer/system"
)

const (
	cgroupCpuacctStat   = "cpuacct.stat"
	nanosecondsInSecond = 1000000000
)

var clockTicks = uint64(system.GetClockTicks())

type CpuacctGroup struct {
}


func (s *CpuacctGroup) GetStats(path string, metric *cgroups.Stats) error {
	userModeUsage, kernelModeUsage, err := getCpuUsageBreakdown(path)
	if err != nil {
		return err
	}

	totalUsage, err := stats.GetCgroupParamUint(path, "cpuacct.usage")
	if err != nil {
		return err
	}

	percpuUsage, err := getPercpuUsage(path)
	if err != nil {
		return err
	}

	metric.CpuStats.CpuUsage.TotalUsage = totalUsage
	metric.CpuStats.CpuUsage.PercpuUsage = percpuUsage
	metric.CpuStats.CpuUsage.UsageInUsermode = userModeUsage
	metric.CpuStats.CpuUsage.UsageInKernelmode = kernelModeUsage
	return nil
}

// Returns user and kernel usage breakdown in nanoseconds.
func getCpuUsageBreakdown(path string) (uint64, uint64, error) {
	userModeUsage := uint64(0)
	kernelModeUsage := uint64(0)
	const (
		userField   = "user"
		systemField = "system"
	)

	// Expected format:
	// user <usage in ticks>
	// system <usage in ticks>



	data, err := ioutil.ReadFile(filepath.Join(path, cgroupCpuacctStat))
  fmt.Println("path=", filepath.Join(path, cgroupCpuacctStat))
	if err != nil {
		return 0, 0, err
	}
	fields := strings.Fields(string(data))
	if len(fields) != 4 {
		return 0, 0, fmt.Errorf("failure - %s is expected to have 4 fields", filepath.Join(path, cgroupCpuacctStat))
	}
	if fields[0] != userField {
		return 0, 0, fmt.Errorf("unexpected field %q in %q, expected %q", fields[0], cgroupCpuacctStat, userField)
	}
	if fields[2] != systemField {
		return 0, 0, fmt.Errorf("unexpected field %q in %q, expected %q", fields[2], cgroupCpuacctStat, systemField)
	}
	if userModeUsage, err = strconv.ParseUint(fields[1], 10, 64); err != nil {
		return 0, 0, err
	}
	if kernelModeUsage, err = strconv.ParseUint(fields[3], 10, 64); err != nil {
		return 0, 0, err
	}

  fmt.Println("userModeUsage", userModeUsage)
  fmt.Println("kernelModeUsage", kernelModeUsage)

	return (userModeUsage * nanosecondsInSecond) / clockTicks, (kernelModeUsage * nanosecondsInSecond) / clockTicks, nil
}

func getPercpuUsage(path string) ([]uint64, error) {
	percpuUsage := []uint64{}
	data, err := ioutil.ReadFile(filepath.Join(path, "cpuacct.usage_percpu"))
	if err != nil {
		return percpuUsage, err
	}
	for _, value := range strings.Fields(string(data)) {
		value, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return percpuUsage, fmt.Errorf("Unable to convert param value to uint64: %s", err)
		}
		percpuUsage = append(percpuUsage, value)
	}
	return percpuUsage, nil
}

func main()  {

  sysfs_path := "/sys/fs/cgroup/cpuacct/system.slice/docker-b9d1fa3b40ac74244cd7f3cc3893d70bf2c62fa65ebfdb7d6a62142fddb8c04b.scope"
  var metrics cgroups.Stats
  var x CpuacctGroup
  var curCpu, prevCpu uint64 = 0,0
  var curTime, prevTime time.Time
 
	prevCpu = metrics.CpuStats.CpuUsage.TotalUsage
	prevTime = time.Now()
  time.Sleep(time.Second)

  for {
		curCpu = metrics.CpuStats.CpuUsage.TotalUsage
		curTime = time.Now()

		x.GetStats(sysfs_path, &metrics)
    fmt.Println("prevCpu=", prevCpu)
    fmt.Println("curCpu=", curCpu)
    fmt.Println("diff Time=", curTime.Sub(prevTime))
    fmt.Println("diff Cpu=", curCpu - prevCpu)
    fmt.Println(" Cpu(%)=", (float64(curCpu - prevCpu) / float64(curTime.Sub(prevTime))))

    prevCpu = curCpu
    prevTime = curTime
    time.Sleep(time.Second)
  }
}
