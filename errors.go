package tessera

import "errors"

var (
	ErrBadRequest           = errors.New("Tessera: Bad request")
	ErrAppNotFound          = errors.New("Tessera: App not found")
	ErrConflict             = errors.New("Tessera: User already exists in app")
	ErrServer               = errors.New("Tessera: Internal server error")
	ErrUserNotFound         = errors.New("Tessera: User not found")
	ErrInvalidPassword      = errors.New("Tessera: Invalid password")
	ErrUserExistsInApp      = errors.New("Tessera: User already exists in app")
	ErrRequiredFieldMissing = errors.New("Tessera: Some required fields are missing")
)

var ErrMap = map[string]error{
	"ErrUserNotFound":         ErrUserNotFound,
	"ErrAppNotFound":          ErrAppNotFound,
	"ErrInvalidPassword":      ErrInvalidPassword,
	"ErrUserExistsInApp":      ErrUserExistsInApp,
	"ErrRequiredFieldMissing": ErrRequiredFieldMissing,
}
