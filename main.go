package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
)

const (
	lang             = "en-US"
	maxRequestLength = 1500
	outputFormat     = "mp3"
)

func main() {
	var (
		voice = flag.String("voice", "Amy", "the voice to use")
		text  = flag.String("text", "", "the text to convert")
		file  = flag.String("file", "", "the text file to convert")
		out   = flag.String("out", "output.mp3", "the name of the output file")
	)
	flag.Parse()

	p := polly.New(session.Must(session.NewSession()))

	textToSynthesize := getTextToSynthesize(text, file)
	batches := splitTextIntoBatches(textToSynthesize)

	f, err := os.Create(*out)
	must(err)
	defer f.Close()

	now := time.Now()
	err = synthesizeText(p, voice, batches, f)
	must(err)
	seconds := time.Now().Sub(now) / time.Second

	fmt.Printf("audio file written to %s in %d seconds\n", *out, seconds)
}

func synthesizeText(p *polly.Polly, voice *string, batches []string, f *os.File) error {
	numBatches := len(batches)
	fmt.Printf("converting in %d batch(es)\n", numBatches)

	for i, s := range batches {

		speech, err := p.SynthesizeSpeech(&polly.SynthesizeSpeechInput{
			Text:         aws.String(s),
			VoiceId:      voice,
			OutputFormat: aws.String(outputFormat),
		})
		if err != nil {
			return err
		}

		_, err = io.Copy(f, speech.AudioStream)
		if err != nil {
			return err
		}
		speech.AudioStream.Close()

		fmt.Printf("batch %d/%d complete\n", i+1, numBatches)

	}

	return nil
}

func splitTextIntoBatches(text string) []string {
	batches := []string{}
	if len(text) <= maxRequestLength {
		batches = append(batches, text)
	} else {
		for len(text) > maxRequestLength {
			batch := text[:maxRequestLength]
			batches = append(batches, batch)
			text = text[maxRequestLength:]
		}
		batches = append(batches, text)
	}
	return batches
}

func getTextToSynthesize(text, file *string) string {
	if *text != "" && *file != "" {
		exit("both file and text can't be used")
	}
	if *text == "" && *file == "" {
		exit("file or text are required")
	}

	if *file == "" {
		return *text
	}

	b, err := ioutil.ReadFile(*file)
	must(err)
	contents := string(b)
	return contents
}

func must(err error) {
	if err != nil {
		exit(err.Error())
	}
}

func exit(err string) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
