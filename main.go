package main

import (
	"fmt"
	"github.com/brian1917/vcodeapi"
	"log"
	"strings"
	"os"
	"reflect"
)

func main() {

	f, err := os.OpenFile("vcodeMitigationExpire.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	fmt.Println(reflect.TypeOf(f))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Printf("Started running")

	var appSkip bool
	var flaws []vcodeapi.Flaw
	var recentBuild string
	var errorCheck error
	var flawList []string

	config := parseConfig()

	appList := getApps(config.Auth.User, config.Auth.Password, config.AppScope.LimitAppList, config.AppScope.AppListTextFile)

	appCounter := 0

	// CYCLE THROUGH EACH APP
	for _, appID := range appList {
		//ADJUST SOME VARIABLES
		flawList = []string{}
		appSkip = false
		appCounter += 1

		fmt.Printf("Processing App ID %v (%v of %v)\n", appID, appCounter, len(appList))

		//GET THE BUILD LIST
		buildList, err := vcodeapi.ParseBuildList(config.Auth.User, config.Auth.Password, appID)
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
			flaws, _, errorCheck = vcodeapi.ParseDetailedReport(config.Auth.User, config.Auth.Password, buildList[len(buildList)-1].BuildID)
			recentBuild = buildList[len(buildList)-1].BuildID
			//IF THAT BUILD HAS AN ERROR, GET THE NEXT MOST RECENT (CONTINUE FOR 4 TOTAL BUILDS)
			for i := 1; i < 4; i++ {
				if len(buildList) > i && errorCheck != nil {
					flaws, _, errorCheck = vcodeapi.ParseDetailedReport(config.Auth.User, config.Auth.Password, buildList[len(buildList)-(i+1)].BuildID)
					recentBuild = buildList[len(buildList)-(i+1)].BuildID
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
				if f.Remediation_status == "Mitigated" {
					//THE MOST RECENT MITIGATION ACTION IS THE ACCEPTANCE, PROPOSAL SHOULD BE SECOND LAST IN ARRAY
					recentProposal := f.Mitigations.Mitigation[len(f.Mitigations.Mitigation)-2]
					recentApproval := f.Mitigations.Mitigation[len(f.Mitigations.Mitigation)-1]
					// RUN CHECKS FOR MITIGATION TYPE
					if (recentProposal.Action == "Potential False Positive" && config.TargetMitigations.PotentialFalsePositive == true) ||
						(recentProposal.Action == "Mitigate by Design" && config.TargetMitigations.MitigatedByDesign == true) ||
						(recentProposal.Action == "Mitigate by Network Environment" && config.TargetMitigations.MitigatedByNetworkEnv == true) ||
						(recentProposal.Action == "Mitigate by OS Environment" && config.TargetMitigations.MitigationByOSEnv == true) ||
						(recentProposal.Action == "Reviewed - No Action Taken" && config.TargetMitigations.ReviewedNoActionTaken == true) ||
						(recentProposal.Action == "Remediated by User" && config.TargetMitigations.RemediatedByUser == true) {
						// RUN CHECK FOR INCLUDING COMMENT TEXT
						if (config.CommentText.RequireCommentText == true && strings.Contains(recentApproval.Description, config.CommentText.Text)) ||
							(config.CommentText.RequireCommentText == false){
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
				/*vcodeapi.ParseUpdateMitigation(config.Auth.User, config.Auth.Password, recentBuild,
					"reject", "Expired", strings.Join(flawList, ","))*/
				log.Printf("App ID %v: Reject Flaw IDs %v\n", appID, strings.Join(flawList, ","))
			}
		}
	}
	log.Printf("Completed running")
	if 2==3{
		fmt.Println(recentBuild)
	}
}
