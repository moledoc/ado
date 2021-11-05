package main_test

import (
	"testing"

	ado "gitlab.com/utt_meelis/ado"
)

// Test if flag count was updated in the help function
func TestHelpFlagCount(t *testing.T) {
	ado.Help(true)
}
