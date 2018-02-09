package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const api = "https://www.broadbandmap.gov/broadbandmap/"

const data = "https://www.broadbandmap.gov/broadbandmap/demographic/jun2014/state/ids/"

var split []string

type Geography struct {
	Results struct {
		State []State `json:"state"`
	} `json:"Results"`
}

type State struct {
	Fips string `json:"fips"`
}

type AverageResults struct {
	Results []struct {
		IncomeBelowPoverty float32 `json:"incomeBelowPoverty"`
	} `json:"Results"`
}

func main() {

	//store all positional arguments
	p := os.Args[1:]
	var fipIds []string
	// fmt.Printf("Length of p: %v\n", len(p))
	// fmt.Printf("args = %v\n", p)
	if len(p) == 1 {
		split := strings.Split(p[0], ",")
		for _, v := range split {
			fipIds = append(fipIds, getFIPS(v))
		}

	} else {
		fmt.Printf("State values should be comma separated (e.g. oregon,washington,california)")
		os.Exit(-1)
	}
	fmt.Println(fipIds)
	fmt.Println(getAverageIncomeBelowPoverty(fipIds))

}

func getAverageIncomeBelowPoverty(fips []string) float32 {
	var sum float32
	var income AverageResults
	for _, v := range fips {
		fmt.Printf("printing v: %v\n", v)
		resp, err := http.Get(data + v + "?format=json")
		if err != nil {
			log.Fatalf("Error getting Income: %v\n", err)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("%v\n", string(body))
		err = json.NewDecoder(resp.Body).Decode(&income.Results[0].IncomeBelowPoverty)
		if err != nil {
			log.Fatalf("Something went wrong getting income: %v", err)
		}

		sum += income.Results[0].IncomeBelowPoverty
		fmt.Printf("Sum is %v\n", sum)

	}
	fmt.Printf("Sum is %v", sum)
	return sum / float32(len(fips))
}

func getFIPS(state string) string {
	var fips Geography
	resp, err := http.Get(api + "census/state/" + state + "?format=json")
	if err != nil {
		log.Fatalf("Error getting FIPS: %v\n", err)
	}
	defer resp.Body.Close()
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Printf("%v\n", string(body))
	json.NewDecoder(resp.Body).Decode(&fips)
	return fips.Results.State[0].Fips
	// fmt.Printf("Fips value: %s\n", fips.Results.State[0].Fips)
}
