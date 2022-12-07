package capstone_test

import (
	"github.com/gilmored/pivottech/capstone"
	"testing"
)

func TestCapstone(t *testing.T) {
	tests := map[string]struct {
		vin    string
		expect bool
		f      func(string) bool
	}{
		"Pass": {vin: "1GNALDEK9FZ108495", expect: true, f: capstone.GetVehicleInfo},
		"Fail": {vin: "kdia9ej483092d", expect: false, f: capstone.GetVehicleInfo},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.f != nil {
				result := test.f(test.vin)
				if result != test.expect {
					t.Errorf("got %v, expected %v", result, test.expect)
				}
				return
			}
		})
	}
}
