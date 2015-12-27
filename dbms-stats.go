package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type QueryResponse struct {
	Batchcomplete string
	Query         struct {
		Normalized []map[string]string
		Pages      map[string]struct {
			Pageid    int
			Ns        int
			Title     string
			Revisions []map[string]string
		}
	}
}

func query_wikipedia_api(plcontinue string, url string) (string, []byte) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var body_json interface{}
	err = json.Unmarshal(body, &body_json)
	if err != nil {
		fmt.Println("error:", err)
	}

	body_map := body_json.(map[string]interface{})

	if continue_map, ok := body_map["continue"].(map[string]interface{}); ok {
		plcontinue = continue_map["plcontinue"].(string)
	} else {
		plcontinue = ""
	}

	return plcontinue, body
}

func get_page_titles() []string {
	const root_page_title = "List_of_relational_database_management_systems"
	plcontinue := "init"
	var url string
	var query_response QueryResponse
	var body []byte
	var pages []string
	var titles []string

	for plcontinue != "" {
		if plcontinue != "init" {
			url = fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&format=json&titles=%s&prop=revisions&rvprop=content&rvsection=1&plcontinue=%s", root_page_title, plcontinue)
		} else {
			url = fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=query&format=json&titles=%s&prop=revisions&rvprop=content&rvsection=1", root_page_title)
		}
		plcontinue, body = query_wikipedia_api(plcontinue, url)
		err := json.Unmarshal(body, &query_response)
		if err != nil {
			fmt.Println("error:", err)
		}
		text := query_response.Query.Pages["1568820"].Revisions[0]["*"]
		re := regexp.MustCompile(`\[\[.*\]\]`)
		pages = re.FindAllString(text, -1)
		titles = make([]string, len(pages))
		for i, page := range pages {
			page = strings.Replace(page, "[", "", -1)
			page = strings.Replace(page, "]", "", -1)
			page = strings.Split(page, "|")[0]
			titles[i] = page
		}
	}
	return titles
}

func main() {
	for _, page := range get_page_titles() {
		fmt.Println(page)
	}
}
