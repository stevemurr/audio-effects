package main

import (
	"audio-effects/filter"
	"audio-effects/filter/biquad/butterworth/hpf"
	"audio-effects/filter/biquad/butterworth/lpf"
	"audio-effects/filter/biquad/second-order/parametric"
	"audio-effects/filter/biquad/shelving/highshelf"
	"audio-effects/filter/biquad/shelving/lowshelf"
	"flag"
	"io"
	"log"
	"math"
	"os"

	"github.com/naoina/toml"
	wav "github.com/youpy/go-wav"
)

var (
	inFile    = flag.String("in", "tmp.wav", "the wave file you want to process")
	inRate    = flag.Float64("r", 48000.0, "sample rate of input file")
	inDepth   = flag.Int("b", 24, "the bit depth of the input file")
	inConfig  = flag.String("c", "configs/example.toml", "toml file with eq settings")
	outFile   = flag.String("out", "out.wav", "the named output file")
	normalize = flag.Bool("norm", true, "normalize will scale the output to avoid clipping")
)

func processSample(in float64, fx []filter.Filter, bitDepth float64, channel int) float64 {
	for _, filt := range fx {
		in = filt.Process(in, channel)
	}
	return in
}

func process(w *wav.Reader, fx []filter.Filter, bitDepth float64, channels uint16, gain float64, results *[]wav.Sample) {
	for {
		samples, err := w.ReadSamples()
		if err == io.EOF {
			break
		}
		for _, sample := range samples {
			y := wav.Sample{}
			filtered := processSample(w.FloatValue(sample, 0), fx, bitDepth, 0)
			filtered *= gain
			y.Values[0] = int(filtered * math.Pow(2, bitDepth))
			if channels == 2 {
				filtered = processSample(w.FloatValue(sample, 1), fx, bitDepth, 1)
				filtered *= gain
				y.Values[1] = int(filtered * math.Pow(2, bitDepth))
			} else {
				y.Values[1] = y.Values[0]
			}
			*results = append(*results, y)
		}
	}
}

type config struct {
	Master struct {
		Gain float64
	}
	Parametric []*parametric.Parametric
	LowShelf   []*lowshelf.LowShelf
	HighShelf  []*highshelf.HighShelf
	HPF        []*hpf.HPF
	LPF        []*lpf.LPF
}

func parseConfig(c config, rate float64) []filter.Filter {
	fx := []filter.Filter{}
	for _, filt := range c.Parametric {
		filt.UpdateCoefficients(rate)
		fx = append(fx, filt)
	}
	for _, filt := range c.LowShelf {
		filt.UpdateCoefficients(rate)
		fx = append(fx, filt)
	}
	for _, filt := range c.HighShelf {
		filt.UpdateCoefficients(rate)
		fx = append(fx, filt)
	}
	for _, filt := range c.HPF {
		filt.UpdateCoefficients(rate)
		fx = append(fx, filt)
	}
	for _, filt := range c.LPF {
		filt.UpdateCoefficients(rate)
		fx = append(fx, filt)
	}
	return fx
}

func readConfig(val interface{}) error {
	fc, err := os.Open(*inConfig)
	if err != nil {
		return err
	}
	defer fc.Close()
	if err := toml.NewDecoder(fc).Decode(val); err != nil {
		return err
	}
	return nil
}

func writeWav(outFile string, results []wav.Sample, inRate float64, inDepth int, channels uint16) {
	out, err := os.Create(outFile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	wr := wav.NewWriter(out, uint32(len(results)), channels, uint32(inRate), uint16(inDepth))
	wr.WriteSamples(results)
}

func main() {
	flag.Parse()
	f, err := os.Open(*inFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var c config
	if err := readConfig(&c); err != nil {
		panic(err)
	}
	fx := parseConfig(c, *inRate)
	w := wav.NewReader(f)
	format, err := w.Format()
	if err != nil {
		log.Fatal(err)
	}
	results := &[]wav.Sample{}
	process(w, fx, float64(*inDepth), format.NumChannels, c.Master.Gain, results)
	writeWav(*outFile, *results, *inRate, *inDepth, format.NumChannels)
}
