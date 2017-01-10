package main

import (
	"testing"
	. "github.com/aandryashin/matchers"
)

func TestIsValidEmail(t *testing.T) {
	AssertThat(t, isValidEmail(""), Is{false})
	AssertThat(t, isValidEmail("test@example.com"), Is{true})
	AssertThat(t, isValidEmail("_valid@example.com"), Is{true})
	AssertThat(t, isValidEmail("-@-"), Is{false})
}
