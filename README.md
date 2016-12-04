# onedriveclient 

A golang implementation of onedrive client




##Usage


````go


...

import (
    "time"
    "log"
    ...
	"github.com/mysinmyc/onedriveclient"
	"github.com/mysinmyc/onedriveclient/auth"
)

...

    //Create a new client instance
	vClient := onedriveclient.NewOneDriveClient()

    //Build an online authentication helper
    vAuthenticationHelper := auth.NewHttpAuthHelper("localhost:8080", "!!!clientid!!!", "!!!clientsecret!!!", []string{"offline_access", "onedrive.readonly"})

    //Perform authentication
    vAuthenticationError :=vClient.PerformNewAuthentication(vAuthenticationHelper,time.Second * 120	)

    //Check authentication
    if vAuthenticationError!=nil {
        log.Panic(vAuthenticationError)
    }
...

````


## Offline authentication

Sometimes there is no direct connectivity between client and server or httpd cannot be used. In these cases it can be used and offline authentication helper

In the following example program show an authentication url to copy and paste to/from and external browser.

````go

    //Build an offline authentication helper
    vAuthenticationHelper := auth.NewOfflineAuthHelper("!!!clientid!!!", "!!!clientsecret!!!", []string{"offline_access", "onedrive.readonly"})


````

The offline authentication supports also token reedim by using OfflineAuthHelper.ReedimTokenFromRedirectURI(string) instance method
 