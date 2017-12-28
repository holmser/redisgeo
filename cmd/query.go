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

// City only takes the fields we need to make our query
type City struct {
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	State     string  `json:"state"`
}

var v bool

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
		var wg sync.WaitGroup

		rand.Seed(time.Now().UTC().UnixNano())
		for i := 0; i < iterations; i++ {
			wg.Add(1)

			go doGeoSearch(cities[rand.Intn(len(cities))], client, c)
			doPipeHM(client, c, &wg)
		}
		wg.Wait()
	},
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
	q := redis.GeoRadiusQuery{Unit: "mi", Radius: 50, Count: 20}
	res, err := client.GeoRadiusRO("places", loc.Longitude, loc.Latitude, &q).Result()
	if err != nil {
		fmt.Println(err)
	} else {
		c <- res
	}
}

func doPipeHM(client *redis.Client, c <-chan []redis.GeoLocation, wg *sync.WaitGroup) {
	defer wg.Done()

	// fmt.Println("DoPipe Started")
	list := <-c
	var res []*redis.StringStringMapCmd
	client.Pipelined(func(pipe redis.Pipeliner) error {
		for i := range list {
			res = append(res, pipe.HGetAll(list[i].Name))
		}
		return nil
	})

}

// if v {
// 	for i := range res {
// 		fmt.Println(res[i].Val())
// 	}
// }
// fmt.Println("external", err)

// pipe.Exec()
// fmt.Println(pipe.Exec())
// }

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.Flags().IntVarP(&iterations, "num", "n", 100, "number of random queries to execute")
	queryCmd.Flags().BoolVarP(&v, "verbose", "v", false, "verbose")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// queryCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// queryCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
