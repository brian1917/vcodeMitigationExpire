package main

import (
	"github.com/brian1917/vcodeapi"
	"log"
	"time"
)

func expireCheck(flaw vcodeapi.Flaw, config config) bool {

	var expCheck bool
	var expDate time.Time
	var refDate time.Time
	var err error

	if config.ExpirationDetails.SpecificDate == true {
		expDate, err = time.Parse("2006-01-02", config.ExpirationDetails.Date)
	} else if config.ExpirationDetails.DateFlawFound == true {
		refDate, err = time.Parse("2006-01-02 15:04:05 MST", flaw.DateFirstOccurrence)
		if err != nil {
			log.Fatal(err)
		}
		expDate = refDate.AddDate(0, 0, config.ExpirationDetails.DaysToExpire)
	} else {
		refDate, err = time.Parse("2006-01-02 15:04:05 MST", flaw.Mitigations.Mitigation[len(flaw.Mitigations.Mitigation)-2].Date)
		if err != nil {
			log.Fatal(err)
		}
		expDate = refDate.AddDate(0, 0, config.ExpirationDetails.DaysToExpire)
	}

	diff := int64(expDate.Sub(time.Now()))
	if diff <= 0 {
		expCheck = true
	} else {
		expCheck = false
	}
	return expCheck
}
