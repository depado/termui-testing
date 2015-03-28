package utils

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func GetRAMUsage() (total, free, available uint64) {
	contents, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		return
	}
	lines := strings.Split(string(contents), "\n")
	total, err = strconv.ParseUint(strings.Fields(lines[0])[1], 10, 64)
	if err != nil {
		fmt.Println("Error parsing Uint", err)
	}
	free, err = strconv.ParseUint(strings.Fields(lines[1])[1], 10, 64)
	if err != nil {
		fmt.Println("Error parsing Uint", err)
	}
	available, err = strconv.ParseUint(strings.Fields(lines[2])[1], 10, 64)
	if err != nil {
		fmt.Println("Error parsing Uint", err)
	}
	return
}
