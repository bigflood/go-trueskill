package gaussian

import (
	"math"
	"fmt"
)

type Gaussian struct {
	Pi  float64
	Tau float64
}

func MuSigma(mu, sigma float64) Gaussian {
	pi := 1 / (sigma * sigma)

	return Gaussian{
		Pi:  pi,
		Tau: pi * mu,
	}
}

func MuVariance(mu, variance float64) Gaussian {
	return MuSigma(mu, math.Sqrt(variance))
}

func (g Gaussian) String() string {
	return fmt.Sprintf("(mu=%f sigma=%f pi=%f tau=%f)",
		g.Mu(), g.Sigma(), g.Pi, g.Tau)
}

func (g Gaussian) Sigma() float64 {
	return math.Sqrt(1 / g.Pi)
}

func (g Gaussian) Var() float64 {
	sigma := g.Sigma()
	return sigma * sigma
}

func (g Gaussian) Mu() float64 {
	return g.Tau / g.Pi
}

func (g Gaussian) Mul(h Gaussian) Gaussian {
	pi := g.Pi + h.Pi
	tau := g.Tau + h.Tau
	return Gaussian{
		Pi:  pi,
		Tau: tau,
	}
}

func (g Gaussian) Div(h Gaussian) Gaussian {
	pi := g.Pi - h.Pi
	tau := g.Tau - h.Tau
	return Gaussian{
		Pi:  pi,
		Tau: tau,
	}
}

func (g Gaussian) Add(h Gaussian) Gaussian {
	mu := g.Mu() + h.Mu()
	pi := g.Pi * h.Pi / (g.Pi + h.Pi)

	return Gaussian{
		Pi:  pi,
		Tau: pi * mu,
	}
}

func (g Gaussian) Sub(h Gaussian) Gaussian {
	mu := g.Mu() - h.Mu()
	pi := g.Pi * h.Pi / (g.Pi + h.Pi)

	return Gaussian{
		Pi:  pi,
		Tau: pi * mu,
	}
}

func Cdf(mu, sigma, x float64) float64 {
	return 0.5 * math.Erfc(-(x - mu)/(sigma*math.Sqrt2))
}

func Pdf(mu, sigma, x float64) float64 {
	a := (x - mu) / math.Abs(sigma)
	return 1 / math.Sqrt(2*math.Pi) * math.Abs(sigma) * math.Exp(-a*a/2)
}

func Ppf(mu, sigma, x float64) float64 {
	erfcinv := math.Erfcinv(2 * x)
	return mu - sigma*math.Sqrt2*erfcinv
}
