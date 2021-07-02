package apperr

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Exit codes are int values for the exit code that shell interpreter can interpret
const (
	ExitCodeOK    int = 0
	ExitCodeError int = iota
)

// ErrorFormatter is the interface for format
type ErrorFormatter interface {
	Format(s fmt.State, verb rune)
}

// ExitCoder is the wrapper interface for urfave/cli
type ExitCoder interface {
	error
	ExitCode() int
}

// ExitError is the wrapper struct for urfave/cli
type ExitError struct {
	exitCode int
	err      error
}

// NewExitError makes a new ExitError
func NewExitError(exitCode int, err error) *ExitError {
	return &ExitError{
		exitCode: exitCode,
		err:      err,
	}
}

// Error returns the string message, fulfilling the interface required by `error`
func (ee *ExitError) Error() string {
	if ee.err == nil {
		return ""
	}
	return fmt.Sprintf("%v", ee.err)
}

// ExitCode returns the exit code, fulfilling the interface required by `ExitCoder`
func (ee *ExitError) ExitCode() int {
	return ee.exitCode
}

// HandleExit returns int value that shell interpreter can interpret as the exit code
// If err has error message, it will be displayed to stderr
// This function is heavily inspired by urfave/cli.HandleExitCoder
func HandleExit(err error) int {
	if err == nil {
		return ExitCodeOK
	}

	if exitErr, ok := err.(ExitCoder); ok { //nolint:errorlint
		if err.Error() != "" {
			if _, ok := exitErr.(ErrorFormatter); ok {
				logrus.Errorf("%+v", err)
			} else {
				logrus.Error(err)
			}
		}
		return exitErr.ExitCode()
	}

	logrus.Error(err)
	return ExitCodeError
}
