package forms

import "regexp"

// UsernameRegexp regexp.
var UsernameRegexp = regexp.MustCompile(`^[a-zA-Z0-9\\.\\-_]+$`)

// EmailRegexp regexp.
var EmailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
