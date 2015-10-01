package vagrant

import (
	"bufio"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// MachineReadableItem - you can find more information
//  about the format at: https://docs.vagrantup.com/v2/cli/machine-readable.html
type MachineReadableItem struct {
	Timestamp int64
	Target    string
	Type      string
	Data      string
}

// ParseMachineReadableItemsFromString ...
func ParseMachineReadableItemsFromString(str, targetFilter, typeFilter string) []MachineReadableItem {
	res := []MachineReadableItem{}

	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		line := scanner.Text()
		splits := strings.Split(line, ",")
		if len(splits) < 4 {
			log.Debugf("Skipping line: Not enough components: %s", line)
			continue
		}
		time, err := strconv.ParseInt(splits[0], 10, 64)
		if err != nil {
			log.Debugf("Skipping line: Failed to parse timestamp from line: %s", line)
			continue
		}
		targetStr := splits[1]
		if targetFilter != "" && targetStr != targetFilter {
			log.Debugf("Skipping line: target (%s) doesn't match the provided filter (%s)", targetStr, targetFilter)
			continue
		}
		typeStr := splits[2]
		if typeFilter != "" && typeStr != typeFilter {
			log.Debugf("Skipping line: type (%s) doesn't match the provided filter (%s)", typeStr, typeFilter)
			continue
		}
		dataStr := strings.Join(splits[3:], ",")
		//
		res = append(res, MachineReadableItem{time, targetStr, typeStr, dataStr})
	}
	return res
}
