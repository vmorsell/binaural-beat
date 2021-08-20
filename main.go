package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"syscall"

	"github.com/gordonklaus/portaudio"
)

const sampleRate = 44100

func main() {
	portaudio.Initialize()
	defer portaudio.Terminate()

	go play()

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan
	fmt.Println("Exiting...")
}

func play() {
	// Initialize an oscillator inducing a delta binaural beat of 4 Hz.
	o, err := NewOscillator(200, 204, sampleRate)
	if err != nil {
		log.Fatalf("new oscillator: %v", err)
	}
	o.Start()
	defer o.Close()

	// Block function from returning.
	select {}
}

type Oscillator struct {
	FreqL      float64
	FreqR      float64
	SampleRate float64

	*portaudio.Stream
	phaseL float64
	phaseR float64
}

func NewOscillator(freqL, freqR, sampleRate float64) (*Oscillator, error) {
	o := &Oscillator{
		FreqL: freqL,
		FreqR: freqR,
	}

	if err := o.openStream(); err != nil {
		return nil, fmt.Errorf("open stream: %w", err)
	}

	return o, nil
}

func (o *Oscillator) openStream() error {
	s, err := portaudio.OpenDefaultStream(0, 2, sampleRate, 256, o.processAudio)
	if err != nil {
		return fmt.Errorf("open default stream: %w", err)
	}
	o.Stream = s
	return nil
}

func (o *Oscillator) stepPhase() {
	_, o.phaseL = math.Modf(o.phaseL + o.FreqL/sampleRate)
	_, o.phaseR = math.Modf(o.phaseR + o.FreqR/sampleRate)
}

func (o *Oscillator) processAudio(out [][]float32) {
	for i := range out[0] {
		out[0][i] = float32(math.Sin(2 * math.Pi * o.phaseL))
		out[1][i] = float32(math.Sin(2 * math.Pi * o.phaseR))
		o.stepPhase()
	}
}
