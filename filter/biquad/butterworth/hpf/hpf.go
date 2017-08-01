package hpf

import (
	"audio-effects/filter/biquad"
	"math"
)

// HPF --
type HPF struct {
	L    biquad.BiQuad
	R    biquad.BiQuad
	Freq float64
}

// Process processes a single sample
func (l *HPF) Process(in float64, channel int) float64 {
	var output float64
	if channel == 0 {
		output = l.L.DoBiQuad(in)
		output = output*l.L.C0 + in*l.L.D0
		return output
	}
	output = l.R.DoBiQuad(in)
	output = output*l.R.C0 + in*l.R.D0
	return output
}

// UpdateCoefficients --
func (l *HPF) UpdateCoefficients(samplerate float64) {
	C := math.Tan(l.Freq / samplerate)
	l.L.A0 = 1 / (1 + math.Sqrt(2)*C + math.Pow(C, 2))
	l.L.A1 = -2 * l.L.A0
	l.L.A2 = l.L.A0
	l.L.B1 = 2 * l.L.A0 * (math.Pow(C, 2) - 1)
	l.L.B2 = l.L.A0 * (1 - math.Sqrt(2)*C + math.Pow(C, 2))

	l.L.C0 = 1.0
	l.L.D0 = 0.0

	l.R.A0 = 1 / (1 + math.Sqrt(2)*C + math.Pow(C, 2))
	l.R.A1 = -2 * l.L.A0
	l.R.A2 = l.L.A0
	l.R.B1 = 2 * l.L.A0 * (math.Pow(C, 2) - 1)
	l.R.B2 = l.L.A0 * (1 - math.Sqrt(2)*C + math.Pow(C, 2))

	l.R.C0 = 1.0
	l.R.D0 = 0.0
}
