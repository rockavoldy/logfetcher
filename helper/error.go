package helper

import "errors"

var (
	ErrNotFound = errors.New("not found")
)

func HandleErr(err error, desc string) error {
	return errors.New(err.Error() + ": " + desc)
}
