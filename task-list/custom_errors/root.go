package custom_errors

import (
	"errors"
	"fmt"
)

var InvalidFlag = errors.New("Invalid Flag:")

var InvalidArgument = errors.New("Invalid Argument:")

var CreateInvalidFlagErrorWithMessage = func(message string) error {

	return fmt.Errorf("%w %s", InvalidFlag, message)

}

var CreateInvalidArgumentErrorWithMessage = func(message string) error {

	return fmt.Errorf("%w %s", InvalidFlag, message)

}
