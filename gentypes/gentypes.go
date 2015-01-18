package main

import "fmt"

var types = [][]string{
	{"AAM", "Waypoint Arrival Alarm"},
	{"ALM", "Almanac data"},
	{"APA", "Auto Pilot A sentence"},
	{"APB", "Auto Pilot B sentence"},
	{"BOD", "Bearing Origin to Destination"},
	{"BWC", "Bearing using Great Circle route"},
	{"DTM", "Datum being used."},
	{"GGA", "Fix information"},
	{"GLL", "Lat/Lon data"},
	{"GRS", "GPS Range Residuals"},
	{"GSA", "Overall Satellite data"},
	{"GST", "GPS Pseudorange Noise Statistics"},
	{"GSV", "Detailed Satellite data"},
	{"MSK", "send control for a beacon receiver"},
	{"MSS", "Beacon receiver status information"},
	{"RMA", "recommended Loran data"},
	{"RMB", "recommended navigation data for gps"},
	{"RMC", "recommended minimum data for gps"},
	{"RTE", "route message"},
	{"TRF", "Transit Fix Data"},
	{"STN", "Multiple Data ID"},
	{"VBW", "dual Ground / Water Spped"},
	{"VTG", "Vector track an Speed over the Ground"},
	{"WCV", "Waypoint closure velocity (Velocity Made Good)"},
	{"WPL", "Waypoint Location information"},
	{"XTC", "cross track error"},
	{"XTE", "measured cross track error"},
	{"ZTG", "Zulu (UTC) time and time to go (to destination)"},
	{"ZDA", "Date and Time"},
}

func main() {
	for _, t := range types {
		fmt.Printf("// %v represents a %v message.\n", t[0], t[1])
		fmt.Printf("type %v struct {\n", t[0])
		fmt.Printf("}\n\n")
		fmt.Printf("// A %vHandler handles %v messages from a stream.\n", t[0], t[0])
		fmt.Printf("type %vHandler interface {\n", t[0])
		fmt.Printf("\tHandle%v(%v)\n", t[0], t[0])
		fmt.Printf("}\n\n")
	}
}
