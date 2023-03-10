package uptime

import (
	"errors"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
)

func TestAccumulateErrors(t *testing.T) {
	errs := []error{
		nil,
		nil,
		errors.New("foo"),
		nil,
		errors.New("bar"),
	}

	err := accumulateErrors(errs...)
	require.Error(t, err)

	merr := new(multierror.Error)
	if !errors.As(err, &merr) {
		t.Fatal("err is not a multierror")
	}
	require.Len(t, merr.Errors, 2)
}
