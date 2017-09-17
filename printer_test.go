package unicreds

import (
	"os"
)

func ExamplePrintSecret() {
	w := os.Stdout
	var printSecretTests = []struct {
		in string
		// must manually specify output in Example tests
	}{
		{"foo"},
		{"foo."},
		{"Foo\\"},
		{"%"},
		{"%%"},
		{"%s"},
		{"%#v"},
	}

	for _, noline := range []bool{false, true} {
		for _, tt := range printSecretTests {
			PrintSecret(w, tt.in, noline)
		}
	}

	// Output:
	// foo
	// foo.
	// Foo\
	// %
	// %%
	// %s
	// %#v
	// foofoo.Foo\%%%%s%#v
}
