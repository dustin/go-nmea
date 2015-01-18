package nmea

import "testing"

func TestChecksum(t *testing.T) {
	tests := map[string]bool{
		"$GPRMC,162254.00,A,3723.02837,N,12159.39853,W,0.820,188.36,110706,,,A*74": true,
		"$GPRMC,162254.00,A,3723.02837,N,12159.39853,W,0.820,188.36,110706,,,A*72": false,
	}

	for in, exp := range tests {
		if checkChecksum(in) != exp {
			t.Errorf("Failed on %v/%v", in, exp)
		}
	}
}
