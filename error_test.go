package backend

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_that_Error_can_be_converted_to_string(t *testing.T) {
	err := Error("test")
	require.Equal(t, "test", err.Error())
}
