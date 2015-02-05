package nmea

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

var (
	errBadChecksum = errors.New("bad checksum")

	parsers = map[string]func([]string, interface{}) error{
		"$GPRMC": rmcParser,
		"$GPVTG": vtgParser,
		"$GPGGA": ggaParser,
		"$GPGSA": gsaParser,
		"$GPGLL": gllParser,
		"$GPZDA": zdaParser,
		"$GPGSV": gsvParser,
		"$GPAAM": aamParser,
		"$GPGST": gstParser,
	}
)

type cumulativeErrorParser struct {
	err error
}

func (c *cumulativeErrorParser) parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	rv, err := strconv.ParseFloat(s, 64)
	if err != nil {
		c.err = err
	}
	return rv
}

func (c *cumulativeErrorParser) parseInt(s string) int {
	if s == "" {
		return 0
	}
	rv, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		c.err = err
	}
	return int(rv)
}

func (c *cumulativeErrorParser) parseDMS(s, ref string) float64 {
	n := 2
	m := 1.0
	if ref == "E" || ref == "W" {
		n = 3
	}
	if ref == "S" || ref == "W" {
		m = -1
	}

	deg := c.parseFloat(s[:n])
	min := c.parseFloat(s[n:])
	deg += (min / 60.0)
	deg *= m

	return deg
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
		return nil
	}

	t, err := time.Parse("150405.99 020106 UTC", parts[1]+" "+parts[9]+" UTC")
	if err != nil {
		return err
	}

	cp := &cumulativeErrorParser{}

	lat := cp.parseDMS(parts[3], parts[4])
	lon := cp.parseDMS(parts[5], parts[6])
	speed := cp.parseFloat(parts[7])
	angle := cp.parseFloat(parts[8])
	magvar := 0.0
	if parts[10] != "" {
		magvar = cp.parseFloat(parts[10])
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
		return nil
	}

	if parts[2] != "T" || parts[4] != "M" || parts[6] != "N" || parts[8] != "K" {
		return fmt.Errorf("Unexpected VTG packet: %#v", parts)
	}

	cp := &cumulativeErrorParser{}
	vtg := VTG{
		True:     cp.parseFloat(parts[1]),
		Magnetic: cp.parseFloat(parts[3]),
		Knots:    cp.parseFloat(parts[5]),
		KMH:      cp.parseFloat(parts[7]),
	}

	if cp.err != nil {
		return cp.err
	}

	h.HandleVTG(vtg)

	return nil
}

/*
 $GPGGA,123519,4807.038,N,01131.000,E,1,08,0.9,545.4,M,46.9,M,,*47

Where:
     1:    123519       Fix taken at 12:35:19 UTC
     2,3:  4807.038,N   Latitude 48 deg 07.038' N
     4,5:  01131.000,E  Longitude 11 deg 31.000' E
     6:    1            Fix quality: 0 = invalid
                               1 = GPS fix (SPS)
                               2 = DGPS fix
                               3 = PPS fix
			       4 = Real Time Kinematic
			       5 = Float RTK
                               6 = estimated (dead reckoning) (2.3 feature)
			       7 = Manual input mode
			       8 = Simulation mode
     7:     08           Number of satellites being tracked
     8:     0.9          Horizontal dilution of position
     9,10:  545.4,M      Altitude, Meters, above mean sea level
     11,12: 46.9,M       Height of geoid (mean sea level) above WGS84
                      ellipsoid
     (empty field) time in seconds since last DGPS update
     (empty field) DGPS station ID number
     *47          the checksum data, always begins with *

*/
func ggaParser(parts []string, handler interface{}) error {
	h, ok := handler.(GGAHandler)
	if !ok {
		return nil
	}

	if len(parts) < 13 || parts[10] != "M" || parts[12] != "M" {
		return fmt.Errorf("Unexpected GGA packet: %#v", parts)
	}

	t, err := time.Parse("150405 UTC", parts[1]+" UTC")
	if err != nil {
		return err
	}

	cp := &cumulativeErrorParser{}
	h.HandleGGA(GGA{
		Taken:              t,
		Latitude:           cp.parseDMS(parts[2], parts[3]),
		Longitude:          cp.parseDMS(parts[4], parts[5]),
		Quality:            FixQuality(cp.parseInt(parts[6])),
		HorizontalDilution: cp.parseFloat(parts[8]),
		NumSats:            cp.parseInt(parts[7]),
		Altitude:           cp.parseFloat(parts[9]),
		GeoidHeight:        cp.parseFloat(parts[11]),
	})

	return cp.err
}

/*
  $GPGSA,A,3,04,05,,09,12,,,24,,,,,2.5,1.3,2.1*39

Where:
     1. A        Auto selection of 2D or 3D fix (M = manual)
     2. 3        3D fix - values include: 1 = no fix
                                       2 = 2D fix
                                       3 = 3D fix
     3-15.  04,05... PRNs of satellites used for fix (space for 12)
     15.    2.5      PDOP (dilution of precision)
     16. 1.3      Horizontal dilution of precision (HDOP)
     17. 2.1      Vertical dilution of precision (VDOP)
*/
func gsaParser(parts []string, handler interface{}) error {
	h, ok := handler.(GSAHandler)
	if !ok {
		return nil
	}

	if len(parts) != 18 {
		return fmt.Errorf("Unexpected GSA packet: %#v (len=%v)", parts, len(parts))
	}

	cp := &cumulativeErrorParser{}
	sats := []int{}
	for _, s := range parts[3:15] {
		if s != "" {
			sats = append(sats, cp.parseInt(s))
		}
	}

	h.HandleGSA(GSA{
		Auto:     parts[1] == "A",
		Fix:      GSAFix(cp.parseInt(parts[2])),
		SatsUsed: sats,
		PDOP:     cp.parseFloat(parts[15]),
		HDOP:     cp.parseFloat(parts[16]),
		VDOP:     cp.parseFloat(parts[17]),
	})

	return cp.err
}

/*
  $GPGLL,4916.45,N,12311.12,W,225444,A,*1D

Where:
     0,   GLL          Geographic position, Latitude and Longitude
     1,2: 4916.46,N    Latitude 49 deg. 16.45 min. North
     3,4: 12311.12,W   Longitude 123 deg. 11.12 min. West
     5:   225444       Fix taken at 22:54:44 UTC
     6:   A            Data Active or V (void)
*/
func gllParser(parts []string, handler interface{}) error {
	h, ok := handler.(GLLHandler)
	if !ok {
		return nil
	}

	t, err := time.Parse("150405 UTC", parts[5]+" UTC")
	if err != nil {
		return err
	}

	cp := &cumulativeErrorParser{}
	h.HandleGLL(GLL{
		Taken:     t,
		Latitude:  cp.parseDMS(parts[1], parts[2]),
		Longitude: cp.parseDMS(parts[3], parts[4]),
		Active:    parts[6] == "A",
	})
	return nil
}

/*
ZDA - Data and Time

  $GPZDA,hhmmss.ss,dd,mm,yyyy,xx,yy*CC
  $GPZDA,201530.00,04,07,2002,00,00*60

where:
	1.     hhmmss    HrMinSec(UTC)
        2,3,4. dd,mm,yyy Day,Month,Year
        5.     xx        local zone hours -13..13
        6.     yy        local zone minutes 0..59
*/
func zdaParser(parts []string, handler interface{}) error {
	h, ok := handler.(ZDAHandler)
	if !ok {
		return nil
	}

	if len(parts) != 7 || len(parts[1]) < 6 {
		return fmt.Errorf("Unexpected ZDA packet: %#v (len=%v)", parts, len(parts))
	}

	cp := &cumulativeErrorParser{}
	tz := time.UTC
	tzh := cp.parseInt(parts[5])
	tzm := cp.parseInt(parts[6])
	if tzh != 0 || tzm != 0 {
		// Has timezone.  Do timezone stuff
		tz = time.FixedZone("GPS", (tzh*3600)+(tzm*60))
	}

	ts := time.Date(
		cp.parseInt(parts[4]),
		time.Month(cp.parseInt(parts[3])),
		cp.parseInt(parts[2]),
		cp.parseInt(parts[1][:2]),
		cp.parseInt(parts[1][2:4]),
		cp.parseInt(parts[1][4:6]),
		int(float64(time.Second)*cp.parseFloat(parts[1][6:])),
		tz)

	h.HandleZDA(ZDA{ts})

	return cp.err
}

/*
  $GPGSV,2,1,08,01,40,083,46,02,17,308,41,12,07,344,39,14,22,228,45*75

Where:
      GSV          Satellites in view
      1: 2            Number of sentences for full data
      2: 1            sentence 1 of 2
      3: 08           Number of satellites in view

      01           Satellite PRN number
      40           Elevation, degrees
      083          Azimuth, degrees
      46           SNR - higher is better
           for up to 4 satellites per sentence
      *75          the checksum data, always begins with *

*/
func gsvParser(parts []string, handler interface{}) error {
	h, ok := handler.(GSVHandler)
	if !ok {
		return nil
	}

	cp := &cumulativeErrorParser{}
	gsv := GSV{
		InView:         cp.parseInt(parts[3]),
		SentenceNum:    cp.parseInt(parts[2]),
		TotalSentences: cp.parseInt(parts[1]),
	}

	for i := 4; i+4 <= len(parts); i += 4 {
		gsv.SatInfo = append(gsv.SatInfo, GSVSatInfo{
			cp.parseInt(parts[i]),
			cp.parseInt(parts[i+1]),
			cp.parseInt(parts[i+2]),
			cp.parseInt(parts[i+3]),
		})
	}

	h.HandleGSV(gsv)

	return cp.err
}

// GSVAccumulator combines several GSV structures into a single value.
type GSVAccumulator struct {
	InView  int
	Parts   int
	prev    int
	SatInfo []GSVSatInfo
}

// Add a GSV to the accumulating GSV state.  Returns true if
// this is the final state.
func (g *GSVAccumulator) Add(a GSV) bool {
	if a.TotalSentences != g.Parts || a.SentenceNum != g.prev+1 {
		g.InView = a.InView
		g.Parts = a.TotalSentences
		g.prev = a.SentenceNum
		g.SatInfo = a.SatInfo

		if a.SentenceNum != 1 {
			g.prev = 0
			g.SatInfo = nil
		}
		return a.TotalSentences == 1
	}

	g.prev = a.SentenceNum
	g.SatInfo = append(g.SatInfo, a.SatInfo...)

	return g.prev == g.Parts
}

/*
  	$GPAAM,A,A,0.10,N,WPTNME*32

Where:
    AAM    Arrival Alarm
    1:A          Arrival circle entered
    2:A          Perpendicular passed
    3:0.10       Circle radius
    4:N          Nautical miles
    5:WPTNME     Waypoint name
    *32          Checksum data

*/
func aamParser(parts []string, handler interface{}) error {
	h, ok := handler.(AAMHandler)
	if !ok {
		return nil
	}

	cp := &cumulativeErrorParser{}
	aam := AAM{
		Arrival:       parts[1] == "A",
		Perpendicular: parts[2] == "A",
		Radius:        cp.parseFloat(parts[3]),
	}

	h.HandleAAM(aam)

	return cp.err
}

/*
  	$GPGST,024603.00,3.2,6.6,4.7,47.3,5.8,5.6,22.0*58

Where:
    GST    pseudorange noise statistics
    1:024603.00  UTC time of associated GGA fix
    2:3.2        Total RMS standard deviation of ranges inputs to the navigation solution
    3:6.6        Standard deviation (meters) of semi-major axis of error ellipse
    4:4.7        Standard deviation (meters) of semi-minor axis of error ellipse
    5:47.3       Orientation of semi-major axis of error ellipse (true north degrees)
    6:5.8        Standard deviation (meters) of latitude error
    7:5.6        Standard deviation (meters) of longitude error
    8:22.0       Standard deviation (meters) of latitude error
    *32          Checksum data

*/
func gstParser(parts []string, handler interface{}) error {
	h, ok := handler.(GSTHandler)
	if !ok {
		return nil
	}

	t, err := time.Parse("150405 UTC", parts[1][:6]+" UTC")
	if err != nil {
		return err
	}

	cp := &cumulativeErrorParser{}
	gst := GST{
		Timestamp:             t,
		Deviation:             cp.parseFloat(parts[2]),
		MajorDeviceation:      cp.parseFloat(parts[3]),
		MinorDeviation:        cp.parseFloat(parts[4]),
		MajorOrientation:      cp.parseFloat(parts[5]),
		MinorOrientation:      cp.parseFloat(parts[6]),
		LatitudeErrDeviation:  cp.parseFloat(parts[7]),
		LongitudeErrDeviation: cp.parseFloat(parts[8]),
	}

	h.HandleGST(gst)

	return cp.err
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

func parseMessage(line string, handler interface{}) error {
	if !checkChecksum(line) {
		// skip bad checksums
		return errBadChecksum
	}

	parts := strings.Split(line[:len(line)-3], ",")

	if p, ok := parsers[parts[0]]; ok {
		return p(parts, handler)
	}
	return nil
}

// ErrorHandler handles error in processing individual messages.  If
// the error handler returns nil, the processor will keep executing,
// else Process will return the error the ErrorHandler returned.
type ErrorHandler func(err error) error

func defaultErrorHandler(err error) error {
	return nil
}

// Process all of the NMEA messages from the given reader.
//
// The handler satisfies any *Handler interfaces the application
// requires.  Any unhandled message will be ignored.
//
// An optional error handler can decide how to handle any errors that
// arise in parsing.  The default will ignore parser errors.
//
// Process returns nil on EOF.
func Process(r io.Reader, handler interface{}, errh ErrorHandler) error {
	if errh == nil {
		errh = defaultErrorHandler
	}
	s := bufio.NewScanner(r)
	for s.Scan() {
		err := parseMessage(s.Text(), handler)
		if err != nil {
			if e := errh(err); e != nil {
				return e
			}
		}
	}
	return s.Err()
}
