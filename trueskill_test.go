package trueskill

import (
	"testing"
	"math"

	"github.com/bigflood/go-trueskill/gaussian"
)

func TestDrawProbability(t *testing.T) {
	ts := new(TrueSkill).Init()
	actual := ts.CalcDrawProbability(0.1, 2)
	expected := 0.013539900086822598

	AssertEqualFloat64(t, expected, actual)
}

func TestDrawMargin(t *testing.T) {
	ts := new(TrueSkill).Init()
	actual := ts.CalcDrawMargin(0.1, 2)
	expected := 0.74046658745214821717

	AssertEqualFloat64(t, expected, actual)
}

func TestQuality1vs1(t *testing.T) {
	actual := new(TrueSkill).Init().Quality1vs1(gaussian.MuSigma(2, 1), gaussian.MuSigma(4, 1))
	expected := 0.92084456003560233306

	AssertEqualFloat64(t, expected, actual)
}

func TestRate1vs1(t *testing.T) {
	for _, d := range []struct{ mu1, sigma1, mu2, sigma2, q, postMu1, postSigma1, postMu2, postSigma2 float64 }{
		{2.000, 1.000,
			4.000, 1.000,
			9.20844560035602333059e-01,
			2.18378742176196238844, 0.99346471679507120101,
			3.81621257823803761156, 0.99346471679507120101},
		{80, 1, 2, 1,
			1.02730431161284280829e-36,
			80, 1.00346621489935783345,
			2, 1.00346621489935783345},
		{1000, 1, 2, 1,
			0,
			80, 1.00346621489935783345,
			2, 1.00346621489935783345},
	} {
		q := new(TrueSkill).Init().Quality1vs1(gaussian.MuSigma(d.mu1, d.sigma1), gaussian.MuSigma(d.mu2, d.sigma2))

		AssertEqualFloat64(t, d.q, q)

		a, b := new(TrueSkill).Init().Rate1vs1(gaussian.MuSigma(d.mu1, d.sigma1), gaussian.MuSigma(d.mu2, d.sigma2), false)

		AssertEqualGaussian(t, gaussian.MuSigma(d.postMu1, d.postSigma1), a)
		AssertEqualGaussian(t, gaussian.MuSigma(d.postMu2, d.postSigma2), b)
	}
}

func TestRate1vs1Draw(t *testing.T) {
	for _, d := range []struct{ mu1, sigma1, mu2, sigma2, q, postMu1, postSigma1, postMu2, postSigma2 float64 }{
		{2.000, 1.000,
			4.000, 1.000,
			9.20844560035602333059e-01,
			2.05454825390105799698, 0.98968726588704969416,
			3.94545174609894200302, 0.98968726588704969416},
		{80, 1, 2, 1,
			1.02730431161284280829e-36,
			77.87120465878915354097, 0.98966304877719968314,
			4.12879534121084912357, 0.98966304877719968314},
		{1000, 1, 2, 1,
			0,
			972.66389918255993052298, 0.98961856648318713425,
			29.33610081744010500415, 0.98961856648318713425},
	} {
		q := new(TrueSkill).Init().Quality1vs1(gaussian.MuSigma(d.mu1, d.sigma1), gaussian.MuSigma(d.mu2, d.sigma2))

		AssertEqualFloat64(t, d.q, q)

		a, b := new(TrueSkill).Init().Rate1vs1(gaussian.MuSigma(d.mu1, d.sigma1), gaussian.MuSigma(d.mu2, d.sigma2), true)

		AssertEqualGaussian(t, gaussian.MuSigma(d.postMu1, d.postSigma1), a)
		AssertEqualGaussian(t, gaussian.MuSigma(d.postMu2, d.postSigma2), b)
	}
}

func AssertEqualGaussian(t *testing.T, expected, actual gaussian.Gaussian) {
	AssertEqualFloat64(t, expected.Mu(), actual.Mu())
	AssertEqualFloat64(t, expected.Sigma(), actual.Sigma())
}

func AssertEqualFloat64(t *testing.T, expected, actual float64) {
	t.Helper()
	d := expected - actual
	if math.Abs(d) > 0.00000000001 {
		t.Errorf("expected %.20f, actual %.20f (delta=%e)", expected, actual, d)
	}
}

func BenchmarkQuality1vs1(b *testing.B) {
	ts := new(TrueSkill).Init()
	for i := 0; i < b.N; i++ {
		ts.Quality1vs1(gaussian.MuSigma(25, 1), gaussian.MuSigma(30, 1))
	}
}

func BenchmarkRate1vs1(b *testing.B) {
	ts := new(TrueSkill).Init()
	for i := 0; i < b.N; i++ {
		ts.Rate1vs1(gaussian.MuSigma(30, 1), gaussian.MuSigma(25, 1), i%2 == 0)
	}
}
