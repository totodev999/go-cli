package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func main() {
	// arg[1] url
	// arg[2] key optional
	args := os.Args[1:]

	if args[0] == "" {
		fmt.Println("Please provide a URL")
		return
	}

	urlString, err := url.Parse(args[0])

	if err != nil {
		fmt.Println("Invalid URL")
		return
	}

	response, err := http.Get(urlString.String())
	if err != nil {
		fmt.Println("Error fetching URL")
		return
	}

	fmt.Println("Status Code: ", response.Status)
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading body")
		return
	}

	if len(args) <= 1 || args[1] == "" {
		fmt.Println("data: ", string(body))
		return
	}

	key := args[1]
	var data interface{}

	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error parsing JSON")
		return
	}

	switch v := data.(type) {
	case map[string]interface{}:
		fmt.Println("Value: ", v[key])

	case []interface{}:
		if len(v) == 0 {
			fmt.Println("Empty array")
			return
		}
		for _, item := range v {
			if itemMap, ok := item.(map[string]interface{}); ok {
				fmt.Println("Value: ", itemMap[key])
			} else {
				fmt.Println("data is array but not json in it", item)
			}
		}
	}
}
