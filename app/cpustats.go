package main

import (
	"fmt"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

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


func (s *CpuacctGroup) GetStats(path string, stats *cgroups.Stats) error {
	userModeUsage, kernelModeUsage, err := getCpuUsageBreakdown(path)
	if err != nil {
		return err
	}

	totalUsage, err := getCgroupParamUint(path, "cpuacct.usage")
  fmt.Println("totalUsage", totalUsage)
	if err != nil {
		return err
	}

	percpuUsage, err := getPercpuUsage(path)
  fmt.Println("percpuUsage", percpuUsage)
	if err != nil {
		return err
	}

	stats.CpuStats.CpuUsage.TotalUsage = totalUsage
	stats.CpuStats.CpuUsage.PercpuUsage = percpuUsage
	stats.CpuStats.CpuUsage.UsageInUsermode = userModeUsage
	stats.CpuStats.CpuUsage.UsageInKernelmode = kernelModeUsage
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

// ===========================

var (
	ErrNotSupportStat = errors.New("stats are not supported for subsystem")
	ErrNotValidFormat = errors.New("line is not a valid key value format")
)
// Saturates negative values at zero and returns a uint64.
// Due to kernel bugs, some of the memory cgroup stats can be negative.
func parseUint(s string, base, bitSize int) (uint64, error) {
	value, err := strconv.ParseUint(s, base, bitSize)
	if err != nil {
		intValue, intErr := strconv.ParseInt(s, base, bitSize)
		// 1. Handle negative values greater than MinInt64 (and)
		// 2. Handle negative values lesser than MinInt64
		if intErr == nil && intValue < 0 {
			return 0, nil
		} else if intErr != nil && intErr.(*strconv.NumError).Err == strconv.ErrRange && intValue < 0 {
			return 0, nil
		}

		return value, err
	}

	return value, nil
}

// Parses a cgroup param and returns as name, value
//  i.e. "io_service_bytes 1234" will return as io_service_bytes, 1234
func getCgroupParamKeyValue(t string) (string, uint64, error) {
	parts := strings.Fields(t)
	switch len(parts) {
	case 2:
		value, err := parseUint(parts[1], 10, 64)
		if err != nil {
			return "", 0, fmt.Errorf("Unable to convert param value (%q) to uint64: %v", parts[1], err)
		}

		return parts[0], value, nil
	default:
		return "", 0, ErrNotValidFormat
	}
}

// Gets a single uint64 value from the specified cgroup file.
func getCgroupParamUint(cgroupPath, cgroupFile string) (uint64, error) {
	contents, err := ioutil.ReadFile(filepath.Join(cgroupPath, cgroupFile))
	if err != nil {
		return 0, err
	}

	return parseUint(strings.TrimSpace(string(contents)), 10, 64)
}


func main()  {

  sysfs_path := "/sys/fs/cgroup/cpuacct/system.slice/docker-8f32af9ace20e3001195dc6388f0f7d00fd6ada516bbcc645dd49f7997b23b61.scope"
  var stats cgroups.Stats
  var x CpuacctGroup

	err := x.GetStats(sysfs_path, &stats)

  fmt.Println("stats", stats)
  fmt.Println("err", err)

}
