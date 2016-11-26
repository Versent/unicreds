package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/versent/unicreds"
	"testing"
)

func TestStripFix(t *testing.T) {

	testTable := []struct {
		Prefix    []string
		Suffix    []string
		Delimiter string
		Name      string
		Result    string
	}{
		{[]string{}, []string{}, ".", "NO_CHANGE", "NO_CHANGE"},
		{[]string{"PRE"}, []string{}, "_", "PRE_FIX", "FIX"},
		{[]string{"PRE"}, []string{"FIX"}, "_", "PRE_FIX", "FIX"},
		{[]string{""}, []string{"FIX"}, "_", "PRE_FIX", "PRE"},
		{[]string{"PRE", "FIX"}, []string{""}, "_", "PRE_FIX", "FIX"},
		{[]string{"PRE", "FIX"}, []string{"FIX"}, "_", "PRE_FIX_MORE", "FIX_MORE"},
		{[]string{"PRE", "FIX"}, []string{"FIX", "MORE"}, "_", "PRE_FIX_MORE", "FIX"},
		{[]string{"PRE_FIX"}, []string{"FIX"}, "_", "PRE_FIX_MORE", "MORE"},
		{[]string{"PRE"}, []string{}, "_", "PRE_", "PRE_"},
		{[]string{}, []string{"FIX"}, "_", "_FIX", "_FIX"},
	}

	for _, test := range testTable {
		c := &unicreds.DecryptedCredential{
			Credential: &unicreds.Credential{
				Name: test.Name,
			},
		}
		stripFix(c, test.Prefix, test.Suffix, test.Delimiter)
		assert.Equal(t, test.Result, c.Name)
	}

}
