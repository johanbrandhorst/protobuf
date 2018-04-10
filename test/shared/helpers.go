package shared

import (
	"github.com/pkg/errors"
)

func reportError(property string, got, wanted interface{}) error {
	return errors.Errorf("Unexpected %s seen, got %v, wanted %v", property, got, wanted)
}

func unexpectedError(function string, err error) error {
	return errors.Errorf("Unexpected %s error seen: %v", function, err)
}
