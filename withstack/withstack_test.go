package withstack

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func xxx() error {
	return New("simple error")
}
func yyy() error {
	return xxx()
}

type customError struct {
	val int
}

func (e customError) Error() string {
	return fmt.Sprintf("customError %d", e.val)
}
func makeCustomError() error {
	return customError{val: 42}
}
func TestCustomError(t *testing.T) {
	msg := fmt.Sprintf("%s", fmt.Errorf("%w", makeCustomError()))
	require.Equal(t, "customError 42", msg)
}

func aaa() error {
	return Errorf("aaa got an error: %w", makeCustomError())
}
func bbb() error {
	return aaa()
}
func TestErrorf(t *testing.T) {
	msg := fmt.Sprintf("%s", fmt.Errorf("%w", bbb()))
	require := require.New(t)
	require.Contains(msg, "withstack.aaa")
	require.Contains(msg, "withstack.bbb")
	require.Contains(msg, "withstack.TestErrorf")
}
