package filter

// Filter implements a common eq interface
type Filter interface {
	// UpdateCoefficients updates the struct
	UpdateCoefficients(samplerate float64)
	Process(in float64, channel int) float64
}
