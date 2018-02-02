package common

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

/*
CheckError is used to verify the return of function
Params: err error
Return:
*/
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

/*
IsPathOrFileExists is used to check if a File or Directory exists
Params: path string
Return: bool, error
*/
func IsPathOrFileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	// Unexpected result
	return false, err
}
