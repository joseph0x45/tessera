package shared

import "errors"

var ErrUserNotFound = errors.New("UserNotFound")
var ErrAppNotFound = errors.New("ErrAppNotFound")
var ErrValueNotFound = errors.New("ErrValueNotFound")
var ErrUserExistsInApp = errors.New("ErrUserExistsInApp")
