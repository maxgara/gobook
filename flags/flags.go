package main

import (
	"flag"
	"fmt"

	"github.com/syscll/tempconv"
)

func main() {
	// var period = flag.Duration("period", 1*time.Second, "sleep period")
	// flag.Parse()
	// fmt.Printf("sleeping for %v\n", *period)
	// time.Sleep(*period)
	// fmt.Println("done\n")
	var temp = CelsiusFlag("temp", 32, "the temperature")
	flag.Parse()
	fmt.Println(*temp)
}

// Flag for temperatures eg. 25C. Allows C or F. satisfies flag.Value.
type celsiusFlag struct {
	tempconv.Celsius
}

var x flag.Value

// set temp for f
func (f *celsiusFlag) Set(s string) error {
	var temp float64
	var unit string
	_, err := fmt.Sscanf(s, "%f%s", &temp, &unit)
	if err != nil {
		return err
	}
	switch unit {
	case "C":
		f.Celsius = tempconv.Celsius(temp)
	case "F":
		f.Celsius = tempconv.FahrenheitToCelsius(tempconv.Fahrenheit(temp))
	default:
		return fmt.Errorf("invalid temp unit %s (must be C or F)", unit)
	}
	return nil
}

// Create new celsiusFlag, return ptr to Celsius value
func CelsiusFlag(name string, dvalue float64, usage string) *tempconv.Celsius {
	f := celsiusFlag{tempconv.Celsius(dvalue)}
	flag.CommandLine.Var(&f, name, usage)
	return &f.Celsius
}
