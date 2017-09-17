package unicreds

import (
	"fmt"

	"github.com/apex/log"
)

func PrintSecret(secret string, noline bool) {
	log.WithField("noline", noline).Debug("print secret")
	if noline {
		fmt.Printf("%s", secret)
	} else {
		fmt.Println(secret)
	}
}
