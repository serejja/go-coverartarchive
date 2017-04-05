package coverartarchive

import "errors"

var ErrMaxRedirectsReached = errors.New("Maximum amount of consecutive redirects reached")

var ErrNotFound = errors.New("Not found")
