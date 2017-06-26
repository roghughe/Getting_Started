package language

import (
	"encoding/json"
	"fmt"
	"math"
)

type Abser interface {
	Abs() float64
}

type MyFloat float64

func (f MyFloat) Abs() float64 {
	fmt.Println("Calling Abs() - 1")
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

type Vertex struct {
	X, Y float64
}

func (v *Vertex) Abs() float64 {
	fmt.Println("Calling Abs() - 2")
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func DemoInterfaces() {
	fmt.Println("Demoing Interfaces")
	var a Abser
	f := MyFloat(-math.Sqrt2)
	v := Vertex{3, 4}

	a = f  // a MyFloat implements Abser
	a = &v // a *Vertex implements Abser

	// In the following line, v is a Vertex (not *Vertex)
	// and does NOT implement Abser. If uncommented, this willnot compile
	// a = v

	fmt.Println(a.Abs())
}

//===================================================================================
//
//===================================================================================

/** This is an interface defintion for a Dummy struct */
type HttpResponseFetcher interface {
	Fetch(url string) ([]byte, error)
}

// Dummy
type Fetcher struct{}

func (Fetcher) Fetch(url string) ([]byte, error) {
	return nil, nil
}

func populateInfo(fetcher HttpResponseFetcher, parsedInfo *Info) error {
	response, err := fetcher.Fetch("http://example.com/info")

	if err == nil {
		err = json.Unmarshal(response, parsedInfo)

		if err == nil {
			return nil
		}
	}

	return err
}

type Info struct {
	str string
}
