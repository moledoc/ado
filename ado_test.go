package main_test

import (
	"testing"

	ado "github.com/moledoc/ado"
)

// Test if flag count was updated in the help function
func TestHelpFlagCount(t *testing.T) {
	ado.Help(true)
}
