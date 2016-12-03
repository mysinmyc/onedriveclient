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

    //Perform authentication
    vAuthenticationError :=vClient.PerformNewAuthentication(auth.NewHttpAuthHelper("localhost:8080", "!!!clientid!!!", "!!!clientsecret!!!", []string{"offline_access", "onedrive.readonly"}),time.Second * 120	)

    //Check authentication
    if vAuthenticationError!=nil {
        log.Panic(vAuthenticationError)
    }
...

````

