// Copyright © 2017 Chris Holmes chris@holmser.net
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/cobra"
)

var (
	v           bool
	radius      float64
	resultCount int
)

// City only takes the fields we need to make our query
type City struct {
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	State     string  `json:"state"`
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

var iterations int

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Run queries against the Redis data",
	Long:  `Query ingests a list of cities and makes 1000 queries as fast as it can against redis`,
	Run: func(cmd *cobra.Command, args []string) {
		client, c := initRedis()
		cities := loadCities()
		iterations = iterations / 100
		for i := 0; i < iterations; i++ {
			run100(cities, client, c)
		}
	},
}

func run100(cities []City, client *redis.Client, c chan []redis.GeoLocation) {
	var wg sync.WaitGroup

	rand.Seed(time.Now().UTC().UnixNano())
	if v {
		defer timeTrack(time.Now(), "100 queries")
	}
	for i := 0; i < 100; i++ {
		// Waitgroup ensures that all routines have completed be
		wg.Add(1)
		go doGeoSearch(cities[rand.Intn(len(cities))], client, c)
		go doPipeHM(client, c, &wg)
	}
	wg.Wait()

}

func loadCities() []City {

	jsonFile, _ := os.Open("data/cities_short.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var city []City
	json.Unmarshal(byteValue, &city)
	fmt.Printf("%v cities loaded successfully\n", len(city))
	return city
}

// initRedis() initializes the redis client and returns a pointer to the Client
// as well as a channel for the processes to communicate on.
func initRedis() (*redis.Client, chan []redis.GeoLocation) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	c := make(chan []redis.GeoLocation)

	return client, c

}

func doGeoSearch(loc City, client *redis.Client, c chan<- []redis.GeoLocation) {
	// TODO:  fix hardcoded query
	// fmt.Println("doGeoSearch(loc, client, c)")

	q := redis.GeoRadiusQuery{Unit: "mi", Radius: radius, Count: resultCount}
	res, err := client.GeoRadiusRO("places", loc.Longitude, loc.Latitude, &q).Result()
	if err != nil {
		fmt.Println("Geoquery Err: ", err)
	} else {
		c <- res
	}
}

func doPipeHM(client *redis.Client, c <-chan []redis.GeoLocation, wg *sync.WaitGroup) {
	defer wg.Done()
	list := <-c
	var res []*redis.StringStringMapCmd
	client.Pipelined(func(pipe redis.Pipeliner) error {
		for i := range list {
			res = append(res, pipe.HGetAll(list[i].Name))
		}
		return nil
	})

}

func init() {
	// Cobra Flags
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().IntVarP(&iterations, "num", "n", 100, "Number of random queries to execute")
	queryCmd.Flags().Float64VarP(&radius, "radius", "r", 25.0, "Search radius in miles")
	queryCmd.Flags().IntVarP(&resultCount, "limit", "l", 50, "Limit results sent back by Redis")
	queryCmd.Flags().BoolVarP(&v, "verbose", "v", true, "verbose")

}
