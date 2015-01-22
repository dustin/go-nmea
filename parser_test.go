package nmea

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestChecksum(t *testing.T) {
	tests := map[string]bool{
		"":     false,
		"*00":  false,
		"$*00": true,
		"$*01": false,
		"^0*0": false,
		"$0*0": false,
		"$*xx": false,
		"$GPRMC,162254.00,A,3723.02837,N,12159.39853,W,0.820,188.36,110706,,,A*74": true,
		"$GPRMC,162254.00,A,3723.02837,N,12159.39853,W,0.820,188.36,110706,,,A*72": false,
	}

	for in, exp := range tests {
		if checkChecksum(in) != exp {
			t.Errorf("Failed on %v/%v", in, exp)
		}
	}
}

func TestSampleParsing(t *testing.T) {
	for _, s := range strings.Split(ubloxSample, "\n") {
		parseMessage(s, nil)
	}
}

type rmcHandler struct {
	rmc RMC
}

func (r *rmcHandler) HandleRMC(rmc RMC) {
	r.rmc = rmc
}

func logJSON(t *testing.T, h interface{}) {
	j, err := json.Marshal(h)
	if err != nil {
		t.Errorf("Failed to marshal %v: %v", h, err)
	}
	t.Logf("%T: %s", h, j)
}

func TestRMCHandling(t *testing.T) {
	h := &rmcHandler{}
	for _, s := range strings.Split(ubloxSample, "\n") {
		parseMessage(s, h)
	}
	logJSON(t, h.rmc)
}
