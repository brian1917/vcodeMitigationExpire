package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
)

type config struct {
	Auth struct {
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"auth"`

	TargetMitigations struct {
		PotentialFalsePositive bool `json:"potentialFalsePositive"`
		MitigatedByDesign      bool `json:"mitigatedByDesign"`
		MitigationByOSEnv      bool `json:"mitigationByOSEnv"`
		MitigatedByNetworkEnv  bool `json:"mitigatedByNetworkEnv"`
		ReviewedNoActionTaken  bool `json:"reviewedNoActionTaken"`
		RemediatedByUser       bool `json:"remediatedByUser"`
	} `json:"targetMitigations"`

	CommentText struct {
		RequireCommentText bool   `json:"requireCommentText"`
		Text               string `json:"text"`
	} `json:"commentText"`

	AppScope struct {
		LimitAppList    bool   `json:"limitAppList"`
		AppListTextFile string `json:"appListTextFile"`
	} `json:"appScope"`

	ExpirationDetails struct {
		DateFlawFound            bool   `json:"DateFlawFound"`
		DateOfMitigationApproval bool   `json:"dateOfMitigationApproval"`
		SpecificDate             bool   `json:"specificDate"`
		Date                     string `json:"date"`
		DaysToExpire             int    `json:"daysToExpire"`
	} `json:"expirationDetails"`
}

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "Veracode username")
}

func parseConfig() config {

	flag.Parse()

	//READ CONFIG FILE
	var config config

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
	}

	// VALIDATE AT LEAST ONE TARGET MITIGATION IS SET
	if config.TargetMitigations.PotentialFalsePositive == false &&
		config.TargetMitigations.MitigatedByDesign == false &&
		config.TargetMitigations.MitigationByOSEnv == false &&
		config.TargetMitigations.MitigatedByNetworkEnv == false &&
		config.TargetMitigations.ReviewedNoActionTaken == false &&
		config.TargetMitigations.RemediatedByUser == false {
		log.Fatal("at least one target mitigation must be set to true")
	}

	// VALIDATE EXPIRATION CONFIG
	counter := 0
	if config.ExpirationDetails.DateFlawFound == true {
		counter += 1
	}
	if config.ExpirationDetails.DateOfMitigationApproval == true {
		counter += 1
	}
	if config.ExpirationDetails.SpecificDate == true {
		counter += 1
	}
	if counter == 0 {
		log.Fatal("One expiration trigger needs to be set to true")
	}
	if counter > 1 {
		log.Fatal("Only one expiration trigger is allowed")
	}

	return config
}
