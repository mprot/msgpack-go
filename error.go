package msgpack

import "fmt"

const errLengthLimitExceeded = errorString("length limit exceeded")

type errorString string

func errorf(format string, args ...interface{}) error {
	return errorString(fmt.Sprintf(format, args...))
}

func (e errorString) Error() string {
	return string(e)
}

// TypeError specifies an error type for unexpected wire types.
type TypeError struct {
	Actual   Type
	Expected Type
}

// Error returns the error message of the error.
func (e TypeError) Error() string {
	return "unexpected type: " + string(e.Actual) + " (expected " + string(e.Expected) + ")"
}

type intOverflowError struct{}

func (e intOverflowError) Error() string {
	return "integer overflow"
}

type floatOverflowError struct{}

func (e floatOverflowError) Error() string {
	return "floating-point overflow"
}

type invalidExtensionError struct {
	typ int8
}

func newInvalidExtensionError(typ int8) error {
	return invalidExtensionError{typ}
}

func (e invalidExtensionError) Error() string {
	return fmt.Sprintf("invalid extension type %d", e.typ)
}
