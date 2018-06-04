package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"net/http"
	"time"
)

type quakeFeatures struct {
	Features []quakeFeature
}

type quakeFeature struct {
	Properties quakeProperties
}

type quakeProperties struct {
	Time      string
	Depth     float32
	Magnitude float32
	Locality  string
	Mmi       int
}

// Response contains the message for the world
type Response struct {
	Version string  `json:"version"`
	Body    ResBody `json:"response"`
}

// ResBody is the actual body of the response
type ResBody struct {
	OutputSpeech     Payload ` json:"outputSpeech,omitempty"`
	ShouldEndSession bool    `json:"shouldEndSession"`
}

// Payload ...
type Payload struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
}

func NewResponse(speech string) Response {
	return Response{
		Version: "1.0",
		Body: ResBody{
			OutputSpeech: Payload{
				Type: "PlainText",
				Text: speech,
			},
			ShouldEndSession: true,
		},
	}
}

func Handler(request http.Request) (Response, error) {
	response, err := http.Get("https://api.geonet.org.nz/quake?MMI=3")
	if err != nil {
		return Response{}, err
	}

	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Response{}, err
	}

	var q quakeFeatures

	err = json.Unmarshal(b, &q)
	if err != nil {
		return Response{}, err
	}

	if len(q.Features) > 0 {
		p := q.Features[0].Properties
		t, err := time.Parse(time.RFC3339, p.Time)
		if err != nil {
			return Response{}, err
		}

		ts := t.Format("Monday 2 January 2006, 3 4 PM")

		return NewResponse(fmt.Sprintf("The latest earthquake was a magnitude %0.1f earthquake near %s at %s", p.Magnitude, p.Locality, ts)), nil
	}

	return Response{}, nil
}

func main() {
	lambda.Start(Handler)
}
