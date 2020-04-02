package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	supnews "github.com/aiomonitors/supnewsmonitor/supnews"

	"github.com/fatih/color"
)

func initializeData() (*supnews.Config, error) {
	file, fileErr := ioutil.ReadFile("config.json")
	if fileErr != nil {
		return &supnews.Config{}, fileErr
	}
	var config supnews.Config
	unmarshalError := json.Unmarshal(file, &config)
	if unmarshalError != nil {
		return &supnews.Config{}, unmarshalError
	}
	if !strings.HasPrefix(config.Webhook, "https") {
		return &supnews.Config{}, errors.New("Invalid webhook passed, must be a URL")
	}
	if !strings.HasPrefix(config.Icon, "https") {
		if len(config.Icon) == 0 {
			config.Icon = ""
		} else {
			return &supnews.Config{}, errors.New("Invalid icon passed, must be a URL")
		}
	}
	if !strings.HasPrefix(config.Hexcode, "#") {
		config.Hexcode = "#000000"
	}
	if config.Groupname == "name of your group here" {
		config.Groupname = ""
	}
	return &config, nil
}

func sendLatest(c supnews.Config) error {
	res, resErr := supnews.GetPage("")
	if resErr != nil {
		return resErr
	}
	supnews.SendToWebhook(res, c)
	return nil
}

func main() {
	config, initErr := initializeData()
	if initErr != nil {
		color.Red("Error initializing data, please check your config.json file is in the correct format")
		os.Exit(1)
	}
	color.Blue("Options:\n[ 1 ] Send Current News to Webhook\n[ 2 ] Start Monitor")
	color.Magenta("Please enter 1 or 2 to continue")
	var input int
	_, err := fmt.Scanf("%d", &input)
	if err != nil {
		color.Red("Error in program %v", err)
		os.Exit(1)
	}
	switch input {
	case 1:
		sendErr := sendLatest(*config)
		if sendErr != nil {
			color.Red("Error in program %v", sendErr)
			os.Exit(1)
		}
	case 2:
		m, monitorErr := supnews.NewMonitor("", *config)
		if monitorErr != nil {
			color.Red("Error in program %v", monitorErr)
			os.Exit(1)
		}
		color.Yellow("Press ⌃C at any point to exit")
		m.StartMonitor()
	default:
		color.Red("Unknown input")
		os.Exit(1)
	}
}
