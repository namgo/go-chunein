package main

import (
	"os"
	"io/ioutil"
	"strings"
	"log"
	"errors"
	"fmt"
	"net/url"
	"net/http"
	"github.com/manifoldco/promptui"
)

type Station struct {
	title string
	bitrate int
	url string
	asString string
	
}

var (
	statusError = errors.New("Bad Status")
)

func _formatQuery(searchFmt string, searchString string) string {
	searchEscaped := url.QueryEscape(searchString)
	return fmt.Sprintf(searchFmt, searchEscaped)
	
}

func search(searchFmt string, searchString string) ([]Station, error) {
	stations := make([]Station, 0)
	opml, err := newOPMLFromURL(_formatQuery(searchFmt, searchString))
	if err != nil {
		return nil, err
	}
	if opml.Head.Status != 200 {
		return nil, statusError
	}
	for _, a := range opml.Body.Outlines {
		if a.Item == "station" {
			station := Station{
				bitrate: a.Bitrate,
				title: a.Text,
				url: a.URL,
				asString: fmt.Sprintf("%s[%d]", a.Text, a.Bitrate),
			}
			stations = append(stations, station)
		}
	}
	return stations, nil
}

func prompt(stations []Station) Station {
	m := make(map[string]Station)
	keys := make([]string, 0)
	for _, s := range(stations) {
		m[s.asString] = s
		keys = append(keys, s.asString)
	}
	prompt := promptui.Select{
		Label: "Select Station",
		Items: keys,
		Size: 20,
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Panic(err)
	}
	return m[result]
}

func play(serverUrl string, station Station) error {
	res, err := http.Get(station.url)
	if err != nil {
		log.Panic(err)
	}
	data, err := ioutil.ReadAll(res.Body)
	for _, u  := range strings.Split(string(data), "\n") {
		if len(u) > 1 {
			_, err = http.PostForm(serverUrl,
				url.Values{"url": {u}})
		}
	}
	return err
}

func main() {
	stations, err := search("https://opml.tunein.com/Search.ashx?query=%s", os.Args[1])
	if err != nil {
		log.Panic(err)
	}
	station := prompt(stations)
	if err = play(os.Args[2], station); err != nil {
		log.Panic(err)
	}
	
}
