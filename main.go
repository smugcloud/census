package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const api = "https://www.broadbandmap.gov/broadbandmap/"

const data = "https://www.broadbandmap.gov/broadbandmap/demographic/jun2014/state/ids/"

var cs, averages bool

var split []string

//Geography holds the State detail so we can get the fip
type Geography struct {
	Results struct {
		State []State `json:"state"`
	} `json:"Results"`
}

//State holds the fip from the json response
type State struct {
	Fips string `json:"fips"`
}

//StateDetails is the full list so we can create the CSV, or average
type StateDetails struct {
	Results []struct {
		GeographyID                 string  `json:"geographyId"`
		GeographyName               string  `json:"geographyName"`
		LandArea                    float64 `json:"landArea"`
		Population                  uint64  `json:"population"`
		Households                  uint64  `json:"households"`
		RaceWhite                   float64 `json:"raceWhite"`
		RaceBlack                   float64 `json:"raceBlack"`
		RaceHispanic                float64 `json:"raceHispanic"`
		RaceAsian                   float64 `json:"raceAsian"`
		RaceNativeAmerican          float64 `json:"raceNativeAmerican"`
		IncomeBelowPoverty          float64 `json:"incomeBelowPoverty"`
		MedianIncome                float64 `json:"medianIncome"`
		IncomeLessThan25            float64 `json:"incomeLessThan25"`
		IncomeBetween25To50         float64 `json:"incomeBetween25to50"`
		IncomeBetween50To100        float64 `json:"incomeBetween50to100"`
		IncomeBetween100To200       float64 `json:"incomeBetween100to200"`
		IncomeGreater200            float64 `json:"incomeGreater200"`
		EducationHighSchoolGraduate float64 `json:"educationHighSchoolGraduate"`
		EducationBachelorOrGreater  float64 `json:"educationBachelorOrGreater"`
		AgeUnder5                   float64 `json:"ageUnder5"`
		AgeBetween5To19             float64 `json:"ageBetween5to19"`
		AgeBetween20To34            float64 `json:"ageBetween20to34"`
		AgeBetween35To59            float64 `json:"ageBetween35to59"`
		AgeGreaterThan60            float64 `json:"ageGreaterThan60"`
		MyAreaIndicator             bool    `json:"myAreaIndicator"`
	} `json:"Results"`
}

func main() {
	flag.Parse()
	//Create a slightly better usage description
	flag.Usage = func() {
		fmt.Printf("Usage: census [params] [comma separated list of states]\n\ne.g. census --averages oregon,washington,california\n\n")
		flag.PrintDefaults()
	}
	// Verify that a flag has been provided
	if len(os.Args) < 2 {
		fmt.Printf("One command line flag is required.\n\n")
		flag.Usage()
		os.Exit(-1)
	}
	//store all positional arguments
	p := os.Args[2:]

	var fipIds []string
	//Happy path to keep things simple.  Ideally, spaces between the comma
	//and the next state should be allowed (TODO)
	if len(p) == 1 {
		//split on the comma to build our slice
		split := strings.Split(p[0], ",")
		//sort it so the output stays alphabetically sorted
		sort.Strings(split)
		//Get the fip ID's to do the remaining actions
		for _, v := range split {
			fipIds = append(fipIds, getFIPS(v))
		}

	} else {
		fmt.Printf("State values should be comma separated (e.g. oregon,washington,california).\n")
		os.Exit(-1)
	}

	// Switch on the flag
	switch {
	case averages == true:
		fmt.Println(getAverageIncomeBelowPoverty(fipIds))
	case cs == true:
		var finalFile [][]string
		finalFile = printCSVStates(fipIds)
		w := csv.NewWriter(os.Stdout)
		w.WriteAll(finalFile)

	default:
		flag.PrintDefaults()
		os.Exit(-1)
	}

}

//printCSVStates gets us a two dimensional slice so we can use the CSV WriteAll function
func printCSVStates(fips []string) [][]string {
	var allStates []StateDetails
	for _, v := range fips {
		resp, err := http.Get(data + v + "?format=json")
		if err != nil {
			log.Fatalf("Error getting state details: %v\n", err)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		var tmpState StateDetails
		err = json.Unmarshal(body, &tmpState)
		if err != nil {
			log.Fatalf("Something went wrong getting state details: %v", err)
		}
		allStates = append(allStates, tmpState)

	}
	var csvfile []string
	finalFile := make([][]string, len(allStates))

	//More concise and scalable way to convert the values to strings for writing
	//to the CSV functions
	// for i, conv := range

	//Using reflect and strconv to dynamically capture the values observed
	//in the json response, and convert it to a string for the two-dimensional slice
	//used by the CSV writer.
	for i, conv := range allStates {
		t := reflect.ValueOf(&conv.Results[0]).Elem()
		for p := 0; p < t.NumField(); p++ {
			switch t.Field(p).Interface().(type) {
			case bool:
				csvfile = append(csvfile, strconv.FormatBool(t.Field(p).Bool()))
			case float64:
				csvfile = append(csvfile, strconv.FormatFloat(t.Field(p).Float(), 'g', -1, 64))
			case uint64:
				csvfile = append(csvfile, strconv.FormatUint(t.Field(p).Uint(), 10))
			case string:
				csvfile = append(csvfile, t.Field(p).Interface().(string))
			}

		}

		finalFile[i] = append(finalFile[i], csvfile...)
		csvfile = csvfile[:0]
	}

	return finalFile
}

//getAverageIncomeBelowPoverty gives us the integer showing the average income
func getAverageIncomeBelowPoverty(fips []string) int {
	var sum float64
	var income StateDetails
	for _, v := range fips {
		resp, err := http.Get(data + v + "?format=json")
		if err != nil {
			log.Fatalf("Error getting Income: %v\n", err)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &income)
		if err != nil {
			log.Fatalf("Something went wrong getting income: %v", err)
		}

		sum += income.Results[0].IncomeBelowPoverty

	}
	return int((sum / float64(len(fips))) * 100)
}

//Get the fip ID for a given state
func getFIPS(state string) string {
	var fips Geography
	resp, err := http.Get(api + "census/state/" + state + "?format=json")
	if err != nil {
		log.Fatalf("Error getting FIPS: %v\n", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &fips)

	return fips.Results.State[0].Fips
}

//Set our CLI flags
func init() {
	flag.BoolVar(&cs, "csv", false, "Print CSV output of all state information.")
	flag.BoolVar(&averages, "averages", false, "Return average income below poverty across\n\tthe states specified.")

}
