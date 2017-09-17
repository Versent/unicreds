package unicreds

import (
	"bytes"
	"testing"
)

func TestFprintSecret(t *testing.T) {
	var w bytes.Buffer
	var printSecretTests = []struct {
		in string
	}{
		{""},
		{"foo"},
		{"%v"},
		{"%#v"},
		{"%T"},
		{"%%"},
		{"%t"},
		{"%b"},
		{"%c"},
		{"%d"},
		{"%o"},
		{"%q"},
		{"%x"},
		{"%X"},
		{"%U"},
		{"%b"},
		{"%e"},
		{"%E"},
		{"%f"},
		{"%F"},
		{"%g"},
		{"%G"},
		{"%s"},
		{"%q"},
		{"%x"},
		{"%X"},
		{"%p"},
	}

	for _, noline := range []bool{false, true} {
		for _, tt := range printSecretTests {
			FprintSecret(&w, tt.in, noline)

			actual := w.String()
			expected := tt.in
			if !noline {
				expected += "\n"
			}

			if actual != expected {
				t.Errorf("Expected: %s, Actual: %s", expected, actual)
			}

			w.Reset()
		}
	}
}
