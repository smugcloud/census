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
	"strconv"
	"strings"
)

const api = "https://www.broadbandmap.gov/broadbandmap/"

const data = "https://www.broadbandmap.gov/broadbandmap/demographic/jun2014/state/ids/"

var cs, averages bool

var split []string

type Geography struct {
	Results struct {
		State []State `json:"state"`
	} `json:"Results"`
}

type State struct {
	Fips string `json:"fips"`
}

type StateDetails struct {
	Results []struct {
		GeographyID                 string  `json:"geographyId"`
		GeographyName               string  `json:"geographyName"`
		LandArea                    float64 `json:"landArea"`
		Population                  int     `json:"population"`
		Households                  int     `json:"households"`
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

	// Verify that a flag has been provided
	if len(os.Args) < 2 {
		fmt.Printf("One command line flag is required.\n")
		flag.Usage()
		os.Exit(-1)
	}
	//store all positional arguments
	p := os.Args[2:]
	var fipIds []string
	if len(p) == 1 {
		split := strings.Split(p[0], ",")
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
		printCSVStates(fipIds)
	default:
		flag.PrintDefaults()
		os.Exit(-1)
	}

}

func printCSVStates(fips []string) {
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

	w := csv.NewWriter(os.Stdout)
	for _, conv := range allStates {
		var record []string
		record = append(record, conv.Results[0].GeographyID)
		record = append(record, conv.Results[0].GeographyName)
		record = append(record, strconv.FormatFloat(conv.Results[0].LandArea, 'g', -1, 64))
		record = append(record, strconv.Itoa(conv.Results[0].Population))
		record = append(record, strconv.Itoa(conv.Results[0].Households))
		record = append(record, strconv.FormatFloat(conv.Results[0].RaceWhite, 'g', -1, 64))
		record = append(record, strconv.FormatFloat(conv.Results[0].RaceBlack, 'g', -1, 64))
		record = append(record, strconv.FormatFloat(conv.Results[0].RaceHispanic, 'g', -1, 64))
		record = append(record, strconv.FormatFloat(conv.Results[0].RaceAsian, 'g', -1, 64))
		record = append(record, strconv.FormatFloat(conv.Results[0].RaceNativeAmerican, 'g', -1, 64))
		record = append(record, strconv.FormatFloat(conv.Results[0].IncomeBelowPoverty, 'g', -1, 64))
		record = append(record, strconv.FormatFloat(conv.Results[0].MedianIncome, 'g', -1, 64))
		record = append(record, strconv.FormatFloat(conv.Results[0].IncomeLessThan25, 'g', -1, 64))
		record = append(record, strconv.FormatFloat(conv.Results[0].IncomeBetween25To50, 'g', -1, 64))
		record = append(record, strconv.FormatFloat(conv.Results[0].IncomeBetween50To100, 'g', -1, 64))
		record = append(record, strconv.FormatFloat(conv.Results[0].IncomeBetween100To200, 'g', -1, 64))
		record = append(record, strconv.FormatFloat(conv.Results[0].IncomeGreater200, 'g', -1, 64))
		record = append(record, strconv.FormatFloat(conv.Results[0].EducationHighSchoolGraduate, 'g', -1, 64))
		record = append(record, strconv.FormatFloat(conv.Results[0].EducationBachelorOrGreater, 'g', -1, 64))
		record = append(record, strconv.FormatFloat(conv.Results[0].AgeUnder5, 'g', -1, 64))
		record = append(record, strconv.FormatFloat(conv.Results[0].AgeBetween5To19, 'g', -1, 64))
		record = append(record, strconv.FormatFloat(conv.Results[0].AgeBetween20To34, 'g', -1, 64))
		record = append(record, strconv.FormatFloat(conv.Results[0].AgeBetween35To59, 'g', -1, 64))
		record = append(record, strconv.FormatFloat(conv.Results[0].AgeGreaterThan60, 'g', -1, 64))
		record = append(record, strconv.FormatBool(conv.Results[0].MyAreaIndicator))
		w.Write(record)
	}
	w.Flush()
	//}

}

func getAverageIncomeBelowPoverty(fips []string) float64 {
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
		// fmt.Printf("Sum is %v\n", sum)

	}
	// fmt.Printf("Sum is %v", sum)
	return sum / float64(len(fips))
}

func getFIPS(state string) string {
	var fips Geography
	resp, err := http.Get(api + "census/state/" + state + "?format=json")
	if err != nil {
		log.Fatalf("Error getting FIPS: %v\n", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &fips)
	// fmt.Printf("%v\n", string(body))
	// json.NewDecoder(resp.Body).Decode(&fips)
	return fips.Results.State[0].Fips
	// fmt.Printf("Fips value: %s\n", fips.Results.State[0].Fips)
}

func init() {
	flag.BoolVar(&cs, "csv", false, "Print CSV output.")
	flag.BoolVar(&averages, "averages", false, "Return average ")

}
