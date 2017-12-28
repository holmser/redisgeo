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
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"
)

func lineCounter() int {
	file, _ := os.Open("US.txt")
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
		// result := strings.Split(fileScanner.Text(), "\t")
		// fmt.Println(result)
	}
	return lineCount
}

// parseRecord takes a line of text and returns a pointer to a redis GeoLocation
// for insertion into the database
func parseRecord(record string) *redis.GeoLocation {
	res := strings.Split(record, "\t")
	lat, err := strconv.ParseFloat(res[4], 64)
	lon, err := strconv.ParseFloat(res[5], 64)
	if err != nil {
		fmt.Println(err)
	}
	location := redis.GeoLocation{
		Name:      res[0] + " " + res[1],
		Latitude:  lat,
		Longitude: lon}
	// name := strings.Replace(res[1], "\"", "\\\"", -1)
	// fmt.Printf("GEOADD \"places\" %v %v \"%v\"\r\n", res[5], res[4], res[0]+" "+name)
	return &location
}

// loadCmd represents the load command
var host string
var port string

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load geo data into redis",
	Long:  `Loads geo data into redois`,
	Run: func(cmd *cobra.Command, args []string) {
		connString := host + ":" + port
		totalItems := lineCounter()
		fmt.Printf("Loading %v locations to %v\n", totalItems, connString)
		uiprogress.Start()                          // start rendering
		bar := uiprogress.AddBar(totalItems / 1000) // Add a new bar

		// optionally, append and prepend completion and elapsed time
		bar.AppendCompleted()
		bar.PrependElapsed()

		client := redis.NewClient(&redis.Options{
			Addr:     connString,
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		file, _ := os.Open("US.txt")
		fileScanner := bufio.NewScanner(file)
		lineCount := 0
		pipe := client.Pipeline()
		for fileScanner.Scan() {
			lineCount++
			place := parseRecord(fileScanner.Text())
			pipe.GeoAdd("places", place)

			if lineCount%1000 == 0 {
				_, err := pipe.Exec()
				if err != nil {
					fmt.Println(err)
				} else {
					bar.Incr()
				}
			}
		}
		_, err := pipe.Exec()

		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loadCmd.PersistentFlags().String(&host, "host", "localhost", "Redis hostname")
	// loadCmd.Flags().String(&port, "port", "6379", "Redis port")
	// loadCmd.Flag().st
	loadCmd.Flags().StringVarP(&host, "host", "r", "localhost", "Redis hostname")
	loadCmd.Flags().StringVarP(&port, "port", "p", "6379", "Redis port")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
