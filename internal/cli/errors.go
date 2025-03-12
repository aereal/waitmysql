package cli

import (
	"errors"
	"fmt"
)

var (
	ErrMissingDSN = &MissingRequiredFlagError{FlagName: "-dsn"}
)

type MissingRequiredFlagError struct {
	FlagName string
}

func (e *MissingRequiredFlagError) Error() string {
	return fmt.Sprintf("%s must be given", e.FlagName)
}

func (e *MissingRequiredFlagError) Is(err error) bool {
	mrfErr := new(MissingRequiredFlagError)
	if !errors.As(err, &mrfErr) {
		return false
	}
	return e.FlagName == mrfErr.FlagName
}
