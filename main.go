package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

//go:embed "ring.wav"
var b []byte

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Missing first argument: Duration to wait (example: 1s 2m 1h)")
		return
	}
	duration, err := time.ParseDuration(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	start := time.Now()

	streamer, format, err := wav.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	for time.Since(start) < duration {
		timeLeft := duration - time.Since(start)
		fmt.Printf("%c[2K\r%v", 27, timeLeft.Round(time.Duration(1*time.Second)))
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("%c[2K\rTIME IS UP!\n", 27)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
}
