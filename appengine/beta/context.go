package beta

import (
	"context"

	"google.golang.org/appengine/v2"
)

func Context() context.Context {
	return appengine.BackgroundContext()
}
