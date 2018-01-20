package common

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

/*
Check is used to verify the return of function
Params: err error
Return:
*/
func Check(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
