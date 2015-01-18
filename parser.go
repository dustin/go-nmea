package nmea

import (
	"log"
	"strconv"
)

func checkChecksum(line string) bool {
	cs := 0
	if len(line) < 4 {
		return false
	}

	if line[0] != '$' {
		return false
	}
	if line[len(line)-3] != '*' {
		return false
	}
	exp, err := strconv.ParseInt(line[len(line)-2:], 16, 64)
	if err != nil {
		log.Printf("Failed to parse checksum: %v", err)
		return false
	}

	for _, c := range line[1:] {
		if c == '*' {
			break
		}
		cs = cs ^ int(c)
	}

	return cs == int(exp)
}
