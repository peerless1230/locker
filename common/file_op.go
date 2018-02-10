package common

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

/*
MkdirAll is used to make dir by path
Params: err error
Return:
*/
func MkdirAll(path string) error {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		log.Errorf("Mkdir in %s error: %v", path, err)
	}
	return err
}
