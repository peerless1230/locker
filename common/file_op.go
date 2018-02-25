package common

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

/*
MkdirAll is used to make dir by path
Params: path string, mode os.FileMode
Return: error
*/
func MkdirAll(path string, mode os.FileMode) error {
	err := os.MkdirAll(path, mode)
	if err != nil {
		log.Errorf("Mkdir: %s error: %v", path, err)
	} else {
		log.Debugf("Mkdir: %s successfully.", path)
	}

	return err
}

/*
RmdirAll is used to remove dir by path
Params: path string
Return: error
*/
func RmdirAll(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		log.Errorf("Remove dir: %s error: %v", path, err)
	} else {
		log.Debugf("Remove dir: %s successfully.", path)
	}

	return err
}
