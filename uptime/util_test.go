package uptime

import (
	"errors"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
)

func TestAccumulateErrors(t *testing.T) {
	t.Run("all errors", func(t *testing.T) {
		errs := []error{
			errors.New("foo"),
			errors.New("bar"),
		}

		err := accumulateErrors(errs...)
		require.Error(t, err)

		merr := new(multierror.Error)
		if !errors.As(err, &merr) {
			t.Fatal("err is not a multierror")
		}
		require.Len(t, merr.Errors, 2)
	})
	t.Run("some errors", func(t *testing.T) {
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
	})
	t.Run("no errors", func(t *testing.T) {
		errs := []error{
			nil,
			nil,
			nil,
		}
		err := accumulateErrors(errs...)
		require.NoError(t, err)
	})
}
