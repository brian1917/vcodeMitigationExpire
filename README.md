# Veracode Mitigation Expiring Utility

## Description
Designed to run on a scheduled frequency to expire mitigations. The types of mitigations, expiration
references, and other settings are controlled in a JSON config file.

## Third-party Packages
github.com/brian1917/vcodeapi

## Parameters
`-config`: path to json config file

## Configuration File
A sample config file (`sampleConfig.json`) is included in the repository. An annotated version is below:
```{
     "auth": {
       "user": "apiUserName",
       "password": "Pwd123"
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
     }

   }
 ```
* At least one `targetMitigation` must be set to `true`
* One (and only one) of `dateFlawFound`, `dateOfMitigationApproval`, `specificDate` must be set to true.

## Executables
I've added the executables for Mac (vcodeMitigationExpire) and Windows (vcodeMitigationExpire.exe).
Building from source is preferred, but I'll try to keep these up-to-date for those that don't have Go installed.
* For Windows, download the EXE and from the command line run `vcodeMitigationExpire.exe -config config.json`
* For Mac, download the executable, set it to be an executable: `chmod +x vcodeMitigationExpire` and then run `./vcodeMitigationExpire -config config.json`