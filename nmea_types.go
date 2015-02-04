package nmea

import "time"

// AAM represents a Waypoint Arrival Alarm message.
type AAM struct {
	// TODO
}

// A AAMHandler handles AAM messages from a stream.
type AAMHandler interface {
	HandleAAM(AAM)
}

// ALM represents a Almanac data message.
type ALM struct {
	// TODO
}

// A ALMHandler handles ALM messages from a stream.
type ALMHandler interface {
	HandleALM(ALM)
}

// APA represents a Auto Pilot A sentence message.
type APA struct {
	// TODO
}

// A APAHandler handles APA messages from a stream.
type APAHandler interface {
	HandleAPA(APA)
}

// APB represents a Auto Pilot B sentence message.
type APB struct {
	// TODO
}

// A APBHandler handles APB messages from a stream.
type APBHandler interface {
	HandleAPB(APB)
}

// BOD represents a Bearing Origin to Destination message.
type BOD struct {
	// TODO
}

// A BODHandler handles BOD messages from a stream.
type BODHandler interface {
	HandleBOD(BOD)
}

// BWC represents a Bearing using Great Circle route message.
type BWC struct {
	// TODO
}

// A BWCHandler handles BWC messages from a stream.
type BWCHandler interface {
	HandleBWC(BWC)
}

// DTM represents a Datum being used. message.
type DTM struct {
	// TODO
}

// A DTMHandler handles DTM messages from a stream.
type DTMHandler interface {
	HandleDTM(DTM)
}

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

// GRS represents a GPS Range Residuals message.
type GRS struct {
	// TODO
}

// A GRSHandler handles GRS messages from a stream.
type GRSHandler interface {
	HandleGRS(GRS)
}

type GSAFix int

const (
	_ = GSAFix(iota)
	NoFix
	Fix2D
	Fix3D
)

func (g GSAFix) String() string {
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

// GST represents a GPS Pseudorange Noise Statistics message.
type GST struct {
	// TODO
}

// A GSTHandler handles GST messages from a stream.
type GSTHandler interface {
	HandleGST(GST)
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

// MSK represents a send control for a beacon receiver message.
type MSK struct {
	// TODO
}

// A MSKHandler handles MSK messages from a stream.
type MSKHandler interface {
	HandleMSK(MSK)
}

// MSS represents a Beacon receiver status information message.
type MSS struct {
}

// A MSSHandler handles MSS messages from a stream.
type MSSHandler interface {
	HandleMSS(MSS)
}

// RMA represents a recommended Loran data message.
type RMA struct {
}

// A RMAHandler handles RMA messages from a stream.
type RMAHandler interface {
	HandleRMA(RMA)
}

// RMB represents a recommended navigation data for gps message.
type RMB struct {
	// TODO
}

// A RMBHandler handles RMB messages from a stream.
type RMBHandler interface {
	HandleRMB(RMB)
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

// RTE represents a route message message.
type RTE struct {
	// TODO
}

// A RTEHandler handles RTE messages from a stream.
type RTEHandler interface {
	HandleRTE(RTE)
}

// TRF represents a Transit Fix Data message.
type TRF struct {
	// TODO
}

// A TRFHandler handles TRF messages from a stream.
type TRFHandler interface {
	HandleTRF(TRF)
}

// STN represents a Multiple Data ID message.
type STN struct {
	// TODO
}

// A STNHandler handles STN messages from a stream.
type STNHandler interface {
	HandleSTN(STN)
}

// VBW represents a dual Ground / Water Spped message.
type VBW struct {
	// TODO
}

// A VBWHandler handles VBW messages from a stream.
type VBWHandler interface {
	HandleVBW(VBW)
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

// WCV represents a Waypoint closure velocity (Velocity Made Good) message.
type WCV struct {
	// TODO
}

// A WCVHandler handles WCV messages from a stream.
type WCVHandler interface {
	HandleWCV(WCV)
}

// WPL represents a Waypoint Location information message.
type WPL struct {
	// TODO
}

// A WPLHandler handles WPL messages from a stream.
type WPLHandler interface {
	HandleWPL(WPL)
}

// XTC represents a cross track error message.
type XTC struct {
	// TODO
}

// A XTCHandler handles XTC messages from a stream.
type XTCHandler interface {
	HandleXTC(XTC)
}

// XTE represents a measured cross track error message.
type XTE struct {
	// TODO
}

// A XTEHandler handles XTE messages from a stream.
type XTEHandler interface {
	HandleXTE(XTE)
}

// ZTG represents a Zulu (UTC) time and time to go (to destination) message.
type ZTG struct {
	// TODO
}

// A ZTGHandler handles ZTG messages from a stream.
type ZTGHandler interface {
	HandleZTG(ZTG)
}

// ZDA represents a Date and Time message.
type ZDA struct {
	Timestamp time.Time
}

// A ZDAHandler handles ZDA messages from a stream.
type ZDAHandler interface {
	HandleZDA(ZDA)
}
