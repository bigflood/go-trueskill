package gaussian

import (
	"testing"
	"math"
)

func TestGaussian_Mul(t *testing.T) {
	actual := MuSigma(25, 25/3.).Mul(MuSigma(30, 25/3.))
	expected := Gaussian{Pi: 0.02879999999999999574, Tau: 0.79199999999999981526}

	if d := expected.Mu() - actual.Mu(); math.Abs(d) > 0.00000000001 {
		t.Errorf("mu: expected %s, actual %s (delta=%e)", expected, actual, d)
	}

	if d := expected.Pi - actual.Pi; math.Abs(d) > 0.00000000001 {
		t.Errorf("pi: expected %s, actual %s (delta=%e)", expected, actual, d)
	}
}

func TestGaussian_Div(t *testing.T) {
	actual := MuSigma(25, 25/3.).Div(MuSigma(30, 25/3.))
	expected := Gaussian{Pi: 0, Tau: -0.07200000000000000844}

	if d := expected.Mu() - actual.Mu(); math.Abs(d) > 0.00000000001 {
		t.Errorf("mu: expected %s, actual %s (delta=%e)", expected, actual, d)
	}

	if d := expected.Pi - actual.Pi; math.Abs(d) > 0.00000000001 {
		t.Errorf("pi: expected %s, actual %s (delta=%e)", expected, actual, d)
	}
}
