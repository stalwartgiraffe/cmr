package app_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stalwartgiraffe/cmr/internal/app"
	//appfixtures "github.com/stalwartgiraffe/cmr/internal/app/fixtures"
)

func TestOnStartErr(t *testing.T) {
	fatalErr := fmt.Errorf("fatal error")
	type TC struct {
		ShutdownErrorCount int
		Shutdowns          *app.Shutdowns
		Req                *require.Assertions
	}
	newCase := func(t *testing.T) *TC {
		t.Parallel()
		tc := &TC{
			Req: require.New(t),
		}
		tc.Shutdowns = &app.Shutdowns{
			OnShutdownError: func(string, error) {
				tc.ShutdownErrorCount++
			},
		}
		return tc
	}
	const testPhase = "nope"
	t.Run("shutdown with nil start err", func(t *testing.T) {
		tc := newCase(t)
		tc.Shutdowns.ShutdownAll(testPhase, nil)
		tc.Req.Equal(0, tc.ShutdownErrorCount)
	})

	t.Run("Startup with one start err", func(t *testing.T) {
		tc := newCase(t)
		tc.Shutdowns.ShutdownAll(testPhase, fmt.Errorf("test error"))
		tc.Req.Equal(0, tc.ShutdownErrorCount)
	})
	t.Run("shutdown with one shutdown func and fatal err", func(t *testing.T) {
		tc := newCase(t)
		firstShutdown := 0
		tc.Shutdowns.AddShutdown(func(ctx context.Context) error {
			firstShutdown++
			return nil
		})
		tc.Shutdowns.ShutdownAll(testPhase, fatalErr)
		tc.Req.Equal(0, tc.ShutdownErrorCount)
		tc.Req.Equal(1, firstShutdown)
	})
	t.Run("shutdown with two shutdown func and fatal err", func(t *testing.T) {
		tc := newCase(t)
		firstShutdown := 0
		tc.Shutdowns.AddShutdown(func(ctx context.Context) error {
			firstShutdown++
			return nil
		})
		secondShutdown := 0
		tc.Shutdowns.AddShutdown(func(ctx context.Context) error {
			secondShutdown++
			return nil
		})
		tc.Shutdowns.ShutdownAll(testPhase, fatalErr)
		tc.Req.Equal(0, tc.ShutdownErrorCount)
		tc.Req.Equal(1, firstShutdown)
		tc.Req.Equal(1, secondShutdown)
	})
	t.Run("shutdown with two shutdown func, and first fails fatally", func(t *testing.T) {
		tc := newCase(t)
		e1 := "one err"
		e2 := "two err"
		e3 := "final err"
		tc.Shutdowns.OnShutdownError = func(w string, err error) {
			tc.ShutdownErrorCount++
			msg := err.Error()
			tc.Req.Contains(msg, e1)
			tc.Req.Contains(msg, e2)
			tc.Req.Contains(msg, e3)
		}
		firstShutdown := 0
		tc.Shutdowns.AddShutdown(func(ctx context.Context) error {
			firstShutdown++
			return fmt.Errorf(e1)
		})
		secondShutdown := 0
		tc.Shutdowns.AddShutdown(func(ctx context.Context) error {
			secondShutdown++
			return fmt.Errorf(e2)
		})
		tc.Shutdowns.ShutdownAll(testPhase, fmt.Errorf(e3))
		tc.Req.Equal(1, tc.ShutdownErrorCount)
		tc.Req.Equal(1, firstShutdown)
		tc.Req.Equal(1, secondShutdown)
	})
}
