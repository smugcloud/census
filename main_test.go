package main

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"testing"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	defer log.SetOutput(os.Stderr)
	return buf.String()
}
func TestAverage(t *testing.T) {
	split := []string{"oregon", "washington", "california"}
	var fipIds []string

	for _, v := range split {
		fipIds = append(fipIds, getFIPS(v))
	}
	output := getAverageIncomeBelowPoverty(fipIds)
	if output != 0.15046666666666667 {
		t.Fatalf("Average output does not match: %v", output)
	}
}

func TestCSV(t *testing.T) {
	fipIds := []string{"06", "41", "53"}
	var csvfile [][]string
	csvfile = printCSVStates(fipIds)
	// fmt.Println(csvfile)
	compare := [][]string{{"06", "California", "158180.02042103", "38660952", "14450824", "0.5088", "0.044", "0.329", "0.1151", "0.0031", "0.1572", "69823.7016", "0.1994", "0.2171", "0.3046", "0.2137", "0.0652", "0.7512", "0.2587", "0.0615", "0.2091", "0.202", "0.3332", "0.1942", "false"},
		{"41", "Oregon", "96098.56583654", "3996309", "1779290", "0.8444", "0.0115", "0.1105", "0.0283", "0.0053", "0.1594", "53775.8649", "0.2422", "0.2633", "0.3191", "0.1457", "0.0298", "0.8444", "0.2524", "0.0536", "0.1945", "0.1919", "0.3265", "0.2334", "false"},
		{"53", "Washington", "70555.17981912", "7077005", "3091503", "0.7885", "0.028", "0.1076", "0.0686", "0.0073", "0.1348", "63192.7444", "0.1986", "0.2379", "0.3333", "0.1894", "0.0408", "0.8624", "0.2743", "0.0556", "0.2023", "0.1933", "0.3315", "0.2172", "false"}}

	if reflect.DeepEqual(csvfile, compare) == false {
		t.Fatalf("2D slices don't match: %v", csvfile)
	}
}
