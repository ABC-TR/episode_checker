package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type series struct {
	name            string
	watched_episode string
	imdb_key        string
}

type series_again struct {
	series  series
	watched bool
	last_ep string
}

func main() {
	file := "series.txt"
	serieses := read_file(file)
	var serieslist []series_again
	var temp series_again

	for _, series := range serieses {
		lastep := get_episode(series.imdb_key)

		if lastep == series.watched_episode {

			temp = series_again{
				series:  series,
				watched: true,
				last_ep: lastep,
			}

		} else {

			temp = series_again{
				series:  series,
				watched: false,
				last_ep: lastep,
			}
		}

		serieslist = append(serieslist, temp)
	}

	output(serieslist)

}

func modify_file(file string, data []byte) {
	err := ioutil.WriteFile(file, data, 0777)

	if err != nil {
		log.Fatalln(err)
	}
}

func read_file(file string) []series {
	data, err := ioutil.ReadFile(file)
	var serieses []series
	if err != nil {
		log.Fatalln(err)
	}

	sdata := string(data)
	list := strings.Split(sdata, "\n")

	for _, line := range list {
		if line == "" {
			continue
		}

		if strings.Contains(line, "//") {
			continue
		}

		// 58 = :
		a_series := strings.Split(line, "::")

		for n, char := range []byte(a_series[2]) {
			if char == 10 {
				a_series[2] = string([]byte(a_series[2])[:n])
			}
		}

		series := series{
			name:            a_series[0],
			imdb_key:        a_series[1],
			watched_episode: a_series[2],
		}
		serieses = append(serieses, series)
	}
	return serieses
}

func get_episode(imdb_key string) string {
	url := "https://api.tvmaze.com/lookup/shows?imdb="
	//imdb_key := "tt1520211"
	var json_data map[string]interface{}
	resp := get_http(url + imdb_key)

	//get_latest_episode(resp, json_data)

	json.Unmarshal([]byte(resp), &json_data)

	previousepisode_url := json_data["_links"].(map[string]interface{})["previousepisode"].(map[string]interface{})["href"]
	previousepisode_url_string := fmt.Sprintf("%v", previousepisode_url)

	resp = get_http(previousepisode_url_string)
	json.Unmarshal([]byte(resp), &json_data)

	episode := json_data["number"]
	season := json_data["season"]
	result := fmt.Sprintf("%v-%v\n", season, episode)

	b_res := []byte(result)

	for n, char := range b_res {
		if char == 10 {
			result = string([]byte(result)[:n])
		}
	}

	return result
}

func get_http(url string) string {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	sb := string(body)
	return sb
}

func output(serieslist []series_again) {

	/*

	   |   NAME   |  ✓  |  X  | * |
	   +--------------------------+
	   | She-Hulk | 5-6 | 5-7 | X |
	   |  ANDOR   | 2-4 | 2-4 | ✓ |
	   +--------------------------+
	*/

	len1 := 4
	len2 := 1
	len3 := 1
	len4 := 1

	for _, series := range serieslist {
		if len([]rune(series.series.name)) > len1 {
			len1 = len(series.series.name)
		}
		if len([]rune(series.series.watched_episode)) > len2 {
			len2 = len(series.series.watched_episode)
		}
		if len([]rune(series.last_ep)) > len3 {
			len3 = len(series.last_ep)
		}
	}

	a := len1 + len2 + len3 + len4 + 2*4

	fmt.Print("\n+")
	for i := -2; i <= a; i++ {
		fmt.Print("-")
	}
	fmt.Print("+\n")

	fmt.Print("|")
	if len("NAME") <= len1 {
		fmt.Printf(" %v%v |", "NAME", strings.Repeat(" ", len1-len("NAME")))
	} else {
		fmt.Printf(" NAME |")
	}

	if len("✓") <= len2 {
		fmt.Printf(" %v%v |", " ✓", strings.Repeat(" ", len2-len("X")-1))
	} else {
		fmt.Printf(" ✓ |")
	}

	if len("X") <= len3 {
		fmt.Printf(" %v%v |", " X", strings.Repeat(" ", len3-len("X")-1))
	} else {
		fmt.Printf(" X |")
	}
	fmt.Printf(" * |")

	fmt.Print("\n+")
	for i := -2; i <= a; i++ {
		fmt.Print("-")
	}
	fmt.Print("+\n")

	for _, series := range serieslist {
		fmt.Print("|")
		if len([]rune(series.series.name)) <= len1 {
			fmt.Printf(" %v%v |", series.series.name, strings.Repeat(" ", len1-len([]rune(series.series.name))))
		} else {
			fmt.Printf(" %v |", series.series.name)
		}

		if len([]rune(series.series.watched_episode)) <= len2 {
			fmt.Printf(" %v%v |", series.series.watched_episode, strings.Repeat(" ", len2-len([]rune(series.series.watched_episode))))
		} else {
			fmt.Printf(" %v |", series.series.watched_episode)
		}

		if len([]rune(series.last_ep)) <= len3 {
			fmt.Printf(" %v%v |", series.last_ep, strings.Repeat(" ", len3-len([]rune(series.last_ep))))
		} else {
			fmt.Printf(" %v |", series.last_ep)
		}

		if series.watched == true {
			fmt.Printf(" ✓ |")
		} else {
			fmt.Printf(" X |")
		}
		fmt.Print("\n")
	}

	fmt.Print("+")
	for i := -2; i <= a; i++ {
		fmt.Print("-")
	}
	fmt.Print("+\n")
}
