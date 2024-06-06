package cmd

import (
	"testing"

	"github.com/matryer/is"
)

func Test_InitialRoot_ShouldNotFail(t *testing.T) {
	assertIs := is.New(t)

	err := rootCmd.Execute()

	assertIs.NoErr(err)
}
