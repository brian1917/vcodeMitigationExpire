package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/brian1917/vcodeapi"
)

func main() {

	// SET UP LOGGING FILE
	f, err := os.OpenFile("vcodeMitigationExpire.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Printf("Started running")

	// SET SOME VARIABLES
	var appSkip bool
	var flaws []vcodeapi.Flaw
	var recentBuild string
	var errorCheck error
	var flawList []string
	var buildsBack int

	// PARSE CONFIG FILE AND LOG CONFIG SETTINGS
	config := parseConfig()
	log.Printf("Config Settings:\n"+
		"[*] Target Mitigations:%v\n"+
		"[*] Comment Text:%v\n"+
		"[*] App Scope:%v\n"+
		"[*] Expiration Details:%v\n",
		config.TargetMitigations, config.CommentText, config.AppScope, config.ExpirationDetails)

	// GET APP LIST
	appList := getApps(config.Auth.CredsFile, config.AppScope.LimitAppList, config.AppScope.AppListTextFile)
	appCounter := 0

	// CYCLE THROUGH EACH APP
	for _, appID := range appList {
		//ADJUST SOME VARIABLES
		flawList = []string{}
		appSkip = false
		appCounter++

		fmt.Printf("Processing App ID %v (%v of %v)\n", appID, appCounter, len(appList))

		//GET THE BUILD LIST
		buildList, err := vcodeapi.ParseBuildList(config.Auth.CredsFile, appID)
		if err != nil {
			log.Fatal(err)
		}

		// GET FOUR MOST RECENT BUILD IDS
		if len(buildList) == 0 {
			appSkip = true
			flaws = nil
			recentBuild = ""
		} else {
			//GET THE DETAILED RESULTS FOR MOST RECENT BUILD
			flaws, _, errorCheck = vcodeapi.ParseDetailedReport(config.Auth.CredsFile, buildList[len(buildList)-1].BuildID)
			recentBuild = buildList[len(buildList)-1].BuildID
			buildsBack = 1
			//IF THAT BUILD HAS AN ERROR, GET THE NEXT MOST RECENT (CONTINUE FOR 4 TOTAL BUILDS)
			for i := 1; i < 4; i++ {
				if len(buildList) > i && errorCheck != nil {
					flaws, _, errorCheck = vcodeapi.ParseDetailedReport(config.Auth.CredsFile, buildList[len(buildList)-(i+1)].BuildID)
					recentBuild = buildList[len(buildList)-(i+1)].BuildID
					buildsBack = i + 1
				}
			}
			// IF 4 MOST RECENT BUILDS HAVE ERRORS, THERE ARE NO RESULTS AVAILABLE
			if errorCheck != nil {
				appSkip = true
			}
		}

		//CHECK FLAWS AND
		if appSkip == false {
			for _, f := range flaws {
				// ONLY RUN ON MITIGATED REMEDIATION STATUS (TAKES INTO ACCOUNT ACCEPTED AND NOT FIXED)
				if f.RemediationStatus == "Mitigated" || f.RemediationStatus == "Reviewed - No Action Taken" || f.RemediationStatus == "Potential False Positive" || f.RemediationStatus == "Remediated by User" {
					//THE MOST RECENT MITIGATION ACTION IS THE ACCEPTANCE, PROPOSAL SHOULD BE SECOND LAST IN ARRAY
					recentProposal := f.Mitigations.Mitigation[len(f.Mitigations.Mitigation)-2]
					recentApproval := f.Mitigations.Mitigation[len(f.Mitigations.Mitigation)-1]
					// CHECK FOR MITIGATION TYPE
					if (recentProposal.Action == "Potential False Positive" && config.TargetMitigations.PotentialFalsePositive == true) ||
						(recentProposal.Action == "Mitigate by Design" && config.TargetMitigations.MitigatedByDesign == true) ||
						(recentProposal.Action == "Mitigate by Network Environment" && config.TargetMitigations.MitigatedByNetworkEnv == true) ||
						(recentProposal.Action == "Mitigate by OS Environment" && config.TargetMitigations.MitigationByOSEnv == true) ||
						(recentProposal.Action == "Reviewed - No Action Taken" && config.TargetMitigations.ReviewedNoActionTaken == true) ||
						(recentProposal.Action == "Remediated by User" && config.TargetMitigations.RemediatedByUser == true) {
						// CHECK FOR INCLUDING COMMENT TEXT
						if (config.CommentText.RequireCommentText == true && strings.Contains(recentApproval.Description, config.CommentText.Text)) ||
							(config.CommentText.RequireCommentText == false) {
							// CHECK IF EXPIRED AND BUILD ARRAY
							if expireCheck(f, config) == true {
								flawList = append(flawList, f.Issueid)
							}
						}
					}
				}
			}
			// IF WE HAVE FLAWS MEETING CRITERIA, RUN UPDATE MITIGATION API
			if len(flawList) > 0 {
				log.Printf("[*]Trying to mitigate IDs %v in Build ID %v in App ID %v", flawList, recentBuild, appID)

				// TRY TO EXPIRE FLAW
				expireError := vcodeapi.ParseUpdateMitigation(config.Auth.CredsFile, recentBuild,
					"rejected", config.ExpirationDetails.RejectionComment, strings.Join(flawList, ","))

				// IF WE HAVE AN ERROR, WE NEED TO TRY 2 BUILDS BACK FROM RESULTS BUILD
				// EXAMPLE = RESULTS IN BUILD 3 (MANUAL); DYNAMIC IS BUILD 2; STATIC IS BUILD 1 (BUILD WE NEED TO MITIGATE STATIC FLAW)
				for i := 0; i < 1; i++ {
					if expireError != nil {
						expireError = vcodeapi.ParseUpdateMitigation(config.Auth.CredsFile, buildList[len(buildList)-(buildsBack+i)].BuildID,
							"rejected", config.ExpirationDetails.RejectionComment, strings.Join(flawList, ","))
						if expireError != nil {
							log.Printf("[*] %v", expireError)
						}
					}
				}
				// IF EXPIRE ERROR IS STILL NOT NULL, NOW WE LOG THE ERROR AND EXIT
				if expireError != nil {
					log.Fatalf("[!] Could not reject Flaw IDs %v in App ID %v", flawList, appID)
				}
				log.Printf("App ID %v: Reject Flaw IDs %v\n", appID, strings.Join(flawList, ","))
			}
		}
	}
	log.Printf("Completed running")
}
