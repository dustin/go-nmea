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

type rmcHandler struct{ t *testing.T }

func (r rmcHandler) HandleRMC(rmc RMC) {
	j, err := json.Marshal(rmc)
	if err != nil {
		panic(err)
	}
	r.t.Logf("%s", j)
}

func TestRMCHandling(t *testing.T) {
	for _, s := range strings.Split(ubloxSample, "\n") {
		parseMessage(s, rmcHandler{t})
	}
}
