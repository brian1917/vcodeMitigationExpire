# Veracode Mitigation Expiring Utility

## Description
Designed to run a a scheduled frequency to automatically expire mitigations. The types of mitigations, expiration
references, and other settings are controlled in a json config file.

## Third-party Packages
1. github.com/brian1917/vcodeapi

## Parameters
1.  **-config**: path to json config file

## Configuration File
A sample config file (`sampleConfig.json`) is included in the repository. The configuration is explained with comments below.
```{
     "auth": {
       "user": "apiUserName",
       "password": "Pwd123"
     },
     "targetMitigations": {                     // Mitigation types set to true will be included for expiring (rejecting)
       "potentialFalsePositive": false,
       "mitigatedByDesign": false,
       "mitigationByOSEnv": false,
       "mitigatedByNetworkEnv": false,
       "reviewedNoActionTaken" : true,
       "remediatedByUser": false
     },
     "commentText": {
       "requireCommentText": true,              // If set to true, only mitigations with the text in the approval comments will be expired (rejected)
       "text": "INCLUDE IN EXPIRATION UTILITY"
     },
     "appScope": {                              // If set to false, all apps in account are used. If set to true, specify a text file with app IDs on each line
       "limitAppList": false,
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

## Executables
I've added the executables for Mac (vcodeMitigationExpire) and Windows (vcodeMitigationExpire.exe).
Building from source is preferred, but I'll try to keep these up-to-date for those that don't have Go installed.
* For Windows, users download the EXE and from the command line run *_vcodeMitigationExpire.exe -config config.json_*
* For Mac, download the executable, set it to be an executable: *_chmod +x vcodeMitigationExpire_* and then run *_./vcodeMitigationExpire -config config.json_*