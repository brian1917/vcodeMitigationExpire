# Veracode Mitigation Expiring Utility
[![Go Report Card](https://goreportcard.com/badge/github.com/brian1917/vcodeMitigationExpire)](https://goreportcard.com/report/github.com/brian1917/vcodeMitigationExpire)

## Description
Utility designed to be run on a regular cadence (e.g., weekly cron job) to expire mitigations. 
The types of mitigations, expiration references, and other settings are controlled in a JSON config file. Executables
for Windows and Mac are available for users that don't want to build from source.

## Third-party Packages
github.com/brian1917/vcodeapi (https://godoc.org/github.com/brian1917/vcodeapi)

## Parameters
`-config`: path to JSON config file

## Configuration File
A sample config file (`sampleConfig.json`) is included in the repository. An annotated version is below:
```
{
     "auth": {
       "credsFile": "/Users/name/.veracode/credentials"
     },
    "mode":{
      "logOnly": true,                          // True will just log which mitigations would be rejected with current config
      "rejectMitigations": false                // True will actually reject the mitigations
  },
     "targetMitigations": {                     // Mitigation types set to true will
       "potentialFalsePositive": false,         // be included for expiring (rejecting)
       "mitigatedByDesign": false,
       "mitigationByOSEnv": false,
       "mitigatedByNetworkEnv": false,
       "reviewedNoActionTaken" : true,
       "remediatedByUser": false
     },
     "commentText": {                           // If set to true, only mitigations with the text
       "requireCommentText": true,              // in the approval comments will be expired (rejected)
       "text": "INCLUDE IN EXPIRATION UTILITY"
     },
     "appScope": {                              // If set to false, all apps in account are used.
       "limitAppList": false,                   // If set to true, specify a text file with app IDs on each line
       "appListTextFile": ""
     },
     "expirationDetails": {
       "DateFlawFound": false,                  // Expiration Date = Date of first occurrence + Days to expire
       "dateOfMitigationApproval": true,        // Expiration Date = Date of Mitigation Approval + Days to expire
       "specificDate": false,                   // Expiration Date = Date provided below in yyyy-mm-dd format
       "date":"",
       "DaysToExpire": 90
       "rejectionComment": "Expired automatically by utility."      // Comment left by utility when rejecting expired mitigation
     }

   }
 ```
* At least one and only one `mode` must be set to `true`.
* At least one `targetMitigation` must be set to `true`
* One (and only one) of `dateFlawFound`, `dateOfMitigationApproval`, `specificDate` must be set to true.

## Logging
Each run will create a log file. The file naming convention is `vcodeMitigationExpire_YYYYMMDD_HHMMSS.log` with the time stamp based on the start of execution. Logging captures when the utility starts running, config settings, mitigations that were expired or would be expired depending on mode, and when the utility stops running.

## Advanced Usage
There might be cases where a single use-case isn't applicable for an organization (e.g.,some apps or
flaws require different expiration times. In those cases, run multiple instances of the utility with different config files.
For example, run the utility weekly (or daily) for each config file that is needed for each scenario.

## Credentials File
Must be structured like the following:
```
[DEFAULT]
veracode_api_key_id = ID HERE
veracode_api_key_secret = SECRET HERE
```

## Executables
I've added the executables for Mac (vcodeMitigationExpire) and Windows (vcodeMitigationExpire.exe).
Building from source is preferred, but I'll try to keep these up-to-date for those that don't have Go installed.
* For Windows, download the executable and from the command line run `vcodeMitigationExpire.exe -config config.json`
* For Mac, download the executable, set it to be an executable: `chmod +x vcodeMitigationExpire` and then run `./vcodeMitigationExpire -config config.json`
