package nmea

import (
	"fmt"
	"time"
)

// FixQuality represents the quality of a position fix in a GGA packet.
type FixQuality int

const (
	InvalidFix = FixQuality(iota)
	GPSFix
	DGPSFix
	PPSFix
	RealTimeKinematicFix
	FloatRealTimeKinematicFix
	EstimatedFix
	ManualInputModeFix
	SimulationModeFix
)

var fixNames = []string{
	InvalidFix:                "invalid fix",
	GPSFix:                    "gps",
	DGPSFix:                   "dgps",
	PPSFix:                    "pps",
	RealTimeKinematicFix:      "rt kinematic",
	FloatRealTimeKinematicFix: "float rt kinematic",
	EstimatedFix:              "estimated",
	ManualInputModeFix:        "manual mode",
	SimulationModeFix:         "sim mode",
}

func (q FixQuality) String() string {
	if q < 0 || int(q) >= len(fixNames) {
		return fmt.Sprintf("[Invalid Fix Value: %d]", q)
	}
	return fixNames[q]
}

// GGA represents a Fix information message.
type GGA struct {
	Taken               time.Time
	Latitude, Longitude float64
	Quality             FixQuality
	NumSats             int
	HorizontalDilution  float64
	Altitude            float64
	GeoidHeight         float64
}

// A GGAHandler handles GGA messages from a stream.
type GGAHandler interface {
	HandleGGA(GGA)
}

// GLL represents a Lat/Lon data message.
type GLL struct {
	Latitude, Longitude float64
	Taken               time.Time
	Active              bool
}

// A GLLHandler handles GLL messages from a stream.
type GLLHandler interface {
	HandleGLL(GLL)
}

type GSAFix int

const (
	_ = GSAFix(iota)
	NoFix
	Fix2D
	Fix3D
)

func (g GSAFix) String() string {
	if g < NoFix || g > Fix3D {
		return fmt.Sprintf("[Invalid GSA Fix: %d]", g)
	}
	return []string{"", "no fix", "2D fix", "3D fix"}[g]
}

// GSA represents a Overall Satellite data message.
type GSA struct {
	Auto             bool
	Fix              GSAFix
	SatsUsed         []int
	PDOP, HDOP, VDOP float64
}

// A GSAHandler handles GSA messages from a stream.
type GSAHandler interface {
	HandleGSA(GSA)
}

type GSVSatInfo struct {
	PRN       int
	Elevation int
	Azimuth   int
	SNR       int
}

// GSV represents a Detailed Satellite data message.
type GSV struct {
	InView         int
	SentenceNum    int
	TotalSentences int
	SatInfo        []GSVSatInfo
}

// A GSVHandler handles GSV messages from a stream.
type GSVHandler interface {
	HandleGSV(GSV)
}

// RMC represents a recommended minimum data for gps message.
type RMC struct {
	Timestamp           time.Time
	Status              rune
	Latitude, Longitude float64
	Speed               float64
	Angle               float64
	Magvar              float64
}

// A RMCHandler handles RMC messages from a stream.
type RMCHandler interface {
	HandleRMC(RMC)
}

// VTG represents a Vector track an Speed over the Ground message.
type VTG struct {
	True, Magnetic float64
	Knots, KMH     float64
}

// A VTGHandler handles VTG messages from a stream.
type VTGHandler interface {
	HandleVTG(VTG)
}

// ZDA represents a Date and Time message.
type ZDA struct {
	Timestamp time.Time
}

// A ZDAHandler handles ZDA messages from a stream.
type ZDAHandler interface {
	HandleZDA(ZDA)
}
