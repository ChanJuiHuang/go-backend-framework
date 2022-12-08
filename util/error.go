package util

import (
	"fmt"
)

const ErrorDelimiter string = "---"

func WrapError(err error) error {
	return fmt.Errorf("%s%s%w", Stack(2), ErrorDelimiter, err)
}
