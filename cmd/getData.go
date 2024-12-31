/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/spf13/cobra"
)

// getDataCmd represents the getData command
var getDataCmd = &cobra.Command{
	Use:   "getData",
	Short: "get data from url which you enter and you can also set id to get specific data",
	Long:  `get data from url which you enter and you can also set id to get specific data`,
	Run: func(cmd *cobra.Command, args []string) {
		urlString, _ := cmd.Flags().GetString("url")
		key, _ := cmd.Flags().GetString("key")
		if urlString == "" {
			fmt.Println("Please provide a URL")
			return
		}

		parsedURL, err := url.Parse(urlString)

		if err != nil {
			fmt.Println("Invalid URL")
			return
		}

		response, err := http.Get(parsedURL.String())
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

		if key == "" {
			fmt.Println("data: ", string(body))
			return
		}

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
	},
}

func init() {
	rootCmd.AddCommand(getDataCmd)

	getDataCmd.Flags().StringP("url", "U", "", "URL to get data from")
	getDataCmd.Flags().StringP("key", "K", "", "key you want to get data for")
}
