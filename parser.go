package nmea

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var notHandled = errors.New("not handled")

var parsers = map[string]func([]string, interface{}) error{
	"$GPRMC": rmcParser,
	"$GPVTG": vtgParser,
}

type cumulativeFloatParser struct {
	err error
}

func (c *cumulativeFloatParser) parse(s string) float64 {
	if s == "" {
		return 0
	}
	rv, err := strconv.ParseFloat(s, 64)
	if err != nil {
		c.err = err
	}
	return rv
}

func parseDMS(s, ref string) (float64, error) {
	n := 2
	m := 1.0
	if ref == "E" || ref == "W" {
		n = 3
	}
	if ref == "S" || ref == "W" {
		m = -1
	}

	cp := &cumulativeFloatParser{}
	deg := cp.parse(s[:n])
	min := cp.parse(s[n:])
	deg += (min / 60.0)
	deg *= m

	return deg, cp.err
}

/*
   0:   RMC          Recommended Minimum sentence C
   1:   123519       Fix taken at 12:35:19 UTC
   2:   A            Status A=active or V=Void.
   3,4: 4807.038,N   Latitude 48 deg 07.038' N
   5,6: 01131.000,E  Longitude 11 deg 31.000' E
   7:   022.4        Speed over the ground in knots
   8:   084.4        Track angle in degrees True
   9:   230394       Date - 23rd of March 1994
   10,11:  003.1,W      Magnetic Variation
*/
func rmcParser(parts []string, handler interface{}) error {
	h, ok := handler.(RMCHandler)
	if !ok {
		return notHandled
	}

	t, err := time.Parse("150405.99 020106 UTC", parts[1]+" "+parts[9]+" UTC")
	if err != nil {
		return err
	}

	lat, err := parseDMS(parts[3], parts[4])
	if err != nil {
		return err
	}
	lon, err := parseDMS(parts[5], parts[6])
	if err != nil {
		return err
	}

	cp := &cumulativeFloatParser{}

	speed := cp.parse(parts[7])
	angle := cp.parse(parts[8])
	magvar := 0.0
	if parts[10] != "" {
		magvar = cp.parse(parts[10])
		if parts[11] == "W" {
			magvar *= -1
		}
	}

	if cp.err != nil {
		return cp.err
	}

	h.HandleRMC(RMC{
		Timestamp: t,
		Status:    rune(parts[2][0]),
		Latitude:  lat,
		Longitude: lon,
		Speed:     speed,
		Angle:     angle,
		Magvar:    magvar,
	})

	return nil
}

/*
VTG - Velocity made good. The gps receiver may use the LC prefix
instead of GP if it is emulating Loran output.

  $GPVTG,054.7,T,034.4,M,005.5,N,010.2,K*48

where:
        // 0:    VTG          Track made good and ground speed
        // 1,2:  054.7,T      True track made good (degrees)
        // 3,4:  034.4,M      Magnetic track made good
        // 5,6:  005.5,N      Ground speed, knots
        // 7,8:  010.2,K      Ground speed, Kilometers per hour
*/
func vtgParser(parts []string, handler interface{}) error {
	h, ok := handler.(VTGHandler)
	if !ok {
		return notHandled
	}

	if parts[2] != "T" || parts[4] != "M" || parts[6] != "N" || parts[8] != "K" {
		return fmt.Errorf("Unexpected VTG packet: %#v", parts)
	}

	cp := &cumulativeFloatParser{}
	vtg := VTG{
		True:     cp.parse(parts[1]),
		Magnetic: cp.parse(parts[3]),
		Knots:    cp.parse(parts[5]),
		KMH:      cp.parse(parts[7]),
	}

	if cp.err != nil {
		return cp.err
	}

	h.HandleVTG(vtg)

	return nil
}

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

func parseMessage(line string, handler interface{}) {
	if !checkChecksum(line) {
		// skip bad checksums
		return
	}

	parts := strings.Split(line[:len(line)-3], ",")

	if p, ok := parsers[parts[0]]; ok {
		if err := p(parts, handler); err != nil {
			log.Printf("Error parsing %#v: %v", parts, err)
		}
	}
}
