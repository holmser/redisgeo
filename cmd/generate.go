// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/cobra"
)

var (
	cityCount int
	left      = [...]string{
		"Prancing",
		"Dancing",
		"Filthy",
		"Hot",
		"Sparkly",
		"Fierce",
		"Jostling",
		"Bruising",
		"Blazing",
		"Speedy",
		"Flatulant",
		"admiring",
		"adoring",
		"affectionate",
		"agitated",
		"amazing",
		"angry",
		"awesome",
		"blissful",
		"boring",
		"brave",
		"clever",
		"cocky",
		"compassionate",
		"competent",
		"condescending",
		"confident",
		"cranky",
		"dazzling",
		"determined",
		"distracted",
		"dreamy",
		"eager",
		"ecstatic",
		"elastic",
		"elated",
		"elegant",
		"eloquent",
		"epic",
		"fervent",
		"festive",
		"flamboyant",
		"focused",
		"friendly",
		"frosty",
		"gallant",
		"gifted",
		"goofy",
		"gracious",
		"happy",
		"hardcore",
		"heuristic",
		"hopeful",
		"hungry",
		"infallible",
		"inspiring",
		"jolly",
		"jovial",
		"keen",
		"kind",
		"laughing",
		"loving",
		"lucid",
		"mystifying",
		"modest",
		"musing",
		"naughty",
		"nervous",
		"nifty",
		"nostalgic",
		"objective",
		"optimistic",
		"peaceful",
		"pedantic",
		"pensive",
		"practical",
		"priceless",
		"quirky",
		"quizzical",
		"relaxed",
		"reverent",
		"romantic",
		"sad",
		"serene",
		"sharp",
		"silly",
		"sleepy",
		"stoic",
		"stupefied",
		"suspicious",
		"tender",
		"thirsty",
		"trusting",
		"unruffled",
		"upbeat",
		"vibrant",
		"vigilant",
		"vigorous",
		"wizardly",
		"wonderful",
		"xenodochial",
		"youthful",
		"zealous",
		"zen",
		"Cory's",
		"Jenn's",
		"Matthew's",
		"Craig's",
		"Cathy's",
		"Grandmas",
	}
	right = [...]string{
		"unicorn",
		"chocolate",
		"christmas",
		"new years",
		"mardi gras",
		"spring",
		"summer",
		"fall",
		"winter",
		"St. Paddy's Day",
	}
)

type RaceCategory struct {
	Category string  `json:"category"`
	Distance float32 `json:"distance"`
	Count    int     `json:"count"`
}
type Location struct {
	Country   string
	State     string
	City      string
	Latitude  float64
	Longitude float64
}
type Race struct {
	Uuid     int
	Name     string
	Distance float32
	City     string
	State    string
	Country  string
}

// generateCmd represents the generate command
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func raceNameGenerator(first, second []string, third string) string {
	rstring := fmt.Sprintf("%s %s %s", first[rand.Intn(len(left))], second[rand.Intn(len(right))], third)
	return rstring
}

func getRandName(names []string) string {
	name := names[rand.Intn(len(names))]
	return name
}

func getNames() []string {
	file, err := os.Open("data/names.txt")
	check(err)
	defer file.Close()
	var names []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}
	return names
}

func getCities() []City {
	jsonFile, err := os.Open("data/cities.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Successfully Opened cities.json")
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}
	var cities []City
	// fmt.Println(byteValue)
	json.Unmarshal(byteValue, &cities)
	return cities

}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate test data",
	Long:  `This command generates statistically relevant test data based on a few inputs`,
	Run: func(cmd *cobra.Command, args []string) {
		// names := getNames()
		connString := host + ":" + port
		rand.Seed(time.Now().UTC().UnixNano())
		cities := getCities()

		jsonFile, err := os.Open("races.json")
		// if we os.Open returns an error then handle it
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Successfully opened races.json")
		}
		defer jsonFile.Close()

		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			fmt.Println(err)
		}
		client := redis.NewClient(&redis.Options{
			Addr:     connString,
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		pipe := client.Pipeline()

		var races []RaceCategory
		// fmt.Println(byteValue)
		json.Unmarshal(byteValue, &races)
		raceId := 0
		for i := 0; i < len(races); i++ {
			for j := 0; j < races[i].Count*5; j++ {
				raceId++
				// name := getRandName(names)
				city := cities[rand.Intn(cityCount)]
				// location := Location{
				// 	Country: "US",
				// 	City:    city,
				// }
				raceMeta := map[string]interface{}{
					"name":      strings.Title(raceNameGenerator(left[:], right[:], races[i].Category)),
					"distance":  races[i].Distance,
					"country":   "US",
					"city":      city.City,
					"state":     city.State,
					"latitude":  city.Latitude,
					"longitude": city.Longitude}

				uuid := strconv.Itoa(raceId)
				// fmt.Println(city)
				loc := &redis.GeoLocation{
					Name:      uuid,
					Latitude:  city.Latitude,
					Longitude: city.Longitude}

				pipe.HMSet(uuid, raceMeta).Err()
				if err != nil {
					fmt.Println(err)
				}

				pipe.GeoAdd("places", loc).Err()
				if err != nil {
					fmt.Println(err)
				}

				if raceId%1000 == 0 {
					_, err := pipe.Exec()
					if err != nil {
						fmt.Println(err)
					}
				}
				// pipe.HM
				// fmt.Println(name, city, races[i].Distance, races[i].Category)
				// fmt.Println(strings.Title(raceNameGenerator(left[:], right[:], races[i].Category)))

			}
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&host, "host", "r", "localhost", "Redis hostname")
	generateCmd.Flags().StringVarP(&port, "port", "p", "6379", "Redis port")
	generateCmd.Flags().IntVarP(&cityCount, "city-count", "c", 1000, "Number of cities to use")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
