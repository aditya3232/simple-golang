package util

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func GetCPUSample() (idle, total uint64) {
	contents, err := ioutil.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseUint(fields[i], 10, 64)
				if err != nil {
					fmt.Println("Error: ", i, fields[i], err)
				}
				total += val // tally up all the numbers to get total ticks
				if i == 4 {  // idle is the 5th field in the cpu line
					idle = val
				}
			}
			return
		}
	}
	return
}

func GetMemorySample() (total, free, buffers, cached uint64) {
	contents, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "MemTotal:":
			total, _ = strconv.ParseUint(fields[1], 10, 64)
		case "MemFree:":
			free, _ = strconv.ParseUint(fields[1], 10, 64)
		case "Buffers:":
			buffers, _ = strconv.ParseUint(fields[1], 10, 64)
		case "Cached:":
			cached, _ = strconv.ParseUint(fields[1], 10, 64)
		}
	}
	return
}

func GetCoreSample() (coreCount int) {
	contents, err := ioutil.ReadFile("/proc/cpuinfo")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if fields[0] == "processor" {
			coreCount++
		}
	}
	return
}
