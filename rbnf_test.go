package rbnf

import (
	"testing"
)

var fs = "Expected '%v', got '%v'"

func Test_exp(t *testing.T) {
	if w, g := 10, exp(10, 1); w != g {
		t.Errorf(fs, w, g)
	}
	if w, g := 100, exp(10, 2); w != g {
		t.Errorf(fs, w, g)
	}
	if w, g := 1, exp(10, 0); w != g {
		t.Errorf(fs, w, g)
	}
}

func Test_Divisor(t *testing.T) {
	var r BaseRule

	r = BaseRule{10, "", "", "tio", "", ""}
	if w, g := 10, r.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = BaseRule{100, "", "", "hundra", "", ""}
	if w, g := 100, r.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = BaseRule{200, "", "", "hundra", "", ""}
	if w, g := 100, r.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = BaseRule{2000, "", "", "tusen", "", ""}
	if w, g := 1000, r.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = BaseRule{2000, "", "", "tusen", "", ">"}
	if w, g := 1000, r.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

	r = BaseRule{2000, "", "", "tusen", "", ">>"}
	if w, g := 1000, r.Divisor(); w != g {
		t.Errorf(fs, w, g)
	}

}
