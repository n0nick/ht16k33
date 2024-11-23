package main

import (
	"flag"
	"log"
	"time"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/devices/v3/ht16k33"
	"periph.io/x/host/v3"
)

func main() {
	title := flag.String("title", "", "Optional title to flash before value")
	value := flag.String("value", "null", "Value to present")
	flag.Parse()

	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatal(err)
	}
	defer bus.Close()

	display, err := ht16k33.NewAlphaNumericDisplay(bus, ht16k33.I2CAddr)
	if err != nil {
		log.Fatal(err)
	}
	// defer display.Halt()

	callOnSubstrings(*title, 4, func(s string) {
		if _, err := display.WriteString(s); err != nil {
			log.Fatal(err)
		}
		time.Sleep(150 * time.Millisecond)
	})
	time.Sleep(200 * time.Millisecond)

	if _, err := display.WriteString(*value); err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
}

// callOnSubstrings calls the function `f` on each n-length substring of `st`, rotating one character at a time.
func callOnSubstrings(st string, n int, f func(string)) {
	// Iterate through the string one character at a time
	for i := 0; i < len(st); i++ {
		// Get the substring of length `n` starting from index `i`
		end := i + n
		if end > len(st) {
			break // Stop if the substring exceeds the string length
		}
		f(st[i:end])
	}
}
