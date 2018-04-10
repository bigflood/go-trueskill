package trueskill

import (
	"math"
	"github.com/bigflood/go-trueskill/gaussian"
)

type TrueSkill struct {
	Mu              float64 // mean of ratings
	Sigma           float64 // standard deviation of ratings
	Beta            float64
	Tau             float64 // dynamic factor
	DrawProbability float64
}

func (ts *TrueSkill) Init() *TrueSkill {
	ts.InitWithMu(25)
	return ts
}

func (ts *TrueSkill) InitWithMu(mu float64) *TrueSkill {
	ts.Mu = mu
	ts.Sigma = ts.Mu / 3
	ts.Beta = ts.Sigma / 2
	ts.Tau = ts.Sigma / 100
	ts.DrawProbability = 0.1
	return ts
}

func (ts *TrueSkill) CalcDrawProbability(drawMargin float64, size int) float64 {
	return ts.Cdf(drawMargin/(math.Sqrt(float64(size))*ts.Beta))*2 - 1
}

func (ts *TrueSkill) CalcDrawMargin(drawProbability float64, size int) float64 {
	return ts.Ppf((drawProbability+1)/2.) * math.Sqrt(float64(size)) * ts.Beta
}

func (ts *TrueSkill) Cdf(x float64) float64 {
	return gaussian.Cdf(0, 1, x)
}

func (ts *TrueSkill) Pdf(x float64) float64 {
	return gaussian.Pdf(0, 1, x)
}

func (ts *TrueSkill) Ppf(x float64) float64 {
	return gaussian.Ppf(0, 1, x)
}

func (ts *TrueSkill) Quality1vs1(a, b gaussian.Gaussian) float64 {
	d := a.Mu() - b.Mu()
	beta22 := ts.Beta * ts.Beta * 2
	v := beta22 + a.Var() + b.Var()
	e := - 0.5 * d * d / v
	s := beta22 / v

	return math.Exp(e) * math.Sqrt(s)
}

func (ts *TrueSkill) Rate1vs1(a, b gaussian.Gaussian, draw bool) (gaussian.Gaussian, gaussian.Gaussian) {
	tau2 := ts.Tau * ts.Tau
	beta2 := ts.Beta * ts.Beta

	// prior factors
	a2 := gaussian.MuVariance(a.Mu(), a.Var()+tau2)
	b2 := gaussian.MuVariance(b.Mu(), b.Var()+tau2)

	// likelyhood factors
	a3 := gaussian.MuVariance(a2.Mu(), a2.Var()+beta2)
	b3 := gaussian.MuVariance(b2.Mu(), b2.Var()+beta2)

	// sum: diff.
	d := a3.Sub(b3)

	t := ts.truncateFactor(d, draw)

	// sum factors (up)
	a4 := b3.Add(t)
	b4 := a3.Sub(t)

	a5 := a4.Mul(a3)
	b5 := b4.Mul(b3)

	a6 := a5.Div(a3)
	b6 := b5.Div(b3)

	// likelyhood factors (up)
	a7 := gaussian.MuVariance(a6.Mu(), a6.Var()+beta2)
	b7 := gaussian.MuVariance(b6.Mu(), b6.Var()+beta2)

	// player skills
	a8 := a7.Mul(a2)
	b8 := b7.Mul(b2)

	return a8, b8
}

func (ts *TrueSkill) truncateFactor(d gaussian.Gaussian, draw bool) gaussian.Gaussian {
	sqrtPi := math.Sqrt(d.Pi)

	var v, w float64
	if draw {
		v, w = ts.vwDraw(d)
	} else {
		v, w = ts.vwWin(d)
	}

	denom := 1 - w

	pi, tau := d.Pi/denom, (d.Tau+sqrtPi*v)/denom

	t := gaussian.Gaussian{Pi: pi, Tau: tau}

	return t.Div(d)
}

func (ts *TrueSkill) vwWin(d gaussian.Gaussian) (float64, float64) {
	sqrtPi := math.Sqrt(d.Pi)

	dm := ts.CalcDrawMargin(ts.DrawProbability, 2)

	x := (d.Mu() - dm) * sqrtPi
	cdf := ts.Cdf(x)
	pdf := ts.Pdf(x)
	v := pdf / cdf
	w := v * (v + x)

	return v, w
}

func (ts *TrueSkill) vwDraw(d gaussian.Gaussian) (float64, float64) {
	dm := ts.CalcDrawMargin(ts.DrawProbability, 2)

	mu := d.Mu()
	sqrtPi := math.Sqrt(d.Pi)

	absMuSqrtPi := math.Abs(mu) * sqrtPi
	dmSqrtPi := dm * sqrtPi

	a := dmSqrtPi - absMuSqrtPi
	b := -dmSqrtPi - absMuSqrtPi

	pdfa := ts.Pdf(a)
	pdfb := ts.Pdf(b)

	denom := ts.Cdf(a) - ts.Cdf(b)
	numer := pdfb - pdfa

	v := numer / denom
	if mu < 0 {
		v = -v
	}

	w := (v * v) + (a*pdfa-b*pdfb)/denom

	return v, w
}
