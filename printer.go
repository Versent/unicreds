package unicreds

import (
	"fmt"
	"io"

	"github.com/apex/log"
)

func FprintSecret(w io.Writer, secret string, noline bool) {
	log.WithField("noline", noline).Debug("print secret")
	if noline {
		fmt.Fprintf(w, "%s", secret)
	} else {
		fmt.Fprintln(w, secret)
	}
}
