package coverartarchive

import "errors"

var ErrNotFound = errors.New("Not found")

var ErrInvalidMBID = errors.New("MBID cannot be parsed as a valid UUID")

var ErrRateLimitReached = errors.New("You have exceeded your rate limit")
