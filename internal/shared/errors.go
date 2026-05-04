package shared

import "errors"

var ErrUserNotFound = errors.New("ErrUserNotFound")
var ErrAppNotFound = errors.New("ErrAppNotFound")
var ErrValueNotFound = errors.New("ErrValueNotFound")
var ErrUserExistsInApp = errors.New("ErrUserExistsInApp")
var ErrSessionNotFound = errors.New("ErrSessionNotFound")
var ErrInvalidPassword = errors.New("ErrInvalidPassword")
var ErrRequiredFieldMissing = errors.New("ErrRequiredFieldMissing")
