package main

import (
	"bufio"
	"github.com/brian1917/vcodeapi"
	"log"
	"os"
)

func getApps(user, password string, limit bool, txtfile string) []string {
	var apps []string

	if limit == false {
		appList, err := vcodeapi.ParseAppList(user, password)
		if err != nil {
			log.Fatal(err)
		}
		for _, app := range appList {
			apps = append(apps, app.AppID)
		}
	} else {
		file, err := os.Open(txtfile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			apps = append(apps, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	return apps
}
