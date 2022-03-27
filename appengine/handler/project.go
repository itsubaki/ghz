package handler

import (
	"context"
	"fmt"

	"golang.org/x/oauth2/google"
)

var ProjectID = func() string {
	creds, err := google.FindDefaultCredentials(context.Background())
	if err != nil {
		panic(fmt.Sprintf("find default credentials: %v", err))
	}

	return creds.ProjectID
}()
