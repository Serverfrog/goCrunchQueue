package internal

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func isExecAny(mode os.FileMode) bool {
	return mode&0111 != 0
}

func HandleError[T any](returnValue T, err error) T {
	if err != nil {
		log.Errorf("Error Occured: %v", err)
	}
	return returnValue
}

func HandleErrorB(err error) bool {
	if err != nil {
		log.Errorf("error Occured: %v", err)
	}
	return err != nil
}
func HandleErrorM(err error, msg string) bool {
	if err != nil {
		log.Errorf("error Occured ( %v ):  %v", msg, err)
	}
	return err != nil
}
func HandleCriticalError(err error) {
	if err != nil {
		log.Fatalf("error Occured: %v", err)
	}
}

func HandleFatalErrorf(err error, msg string) {
	if err != nil {
		log.Fatalf("Error Occured (%v): %v", msg, err)
	}
}
func HandleFatalErrorR[T any](returnValue T, err error) T {
	if err != nil {
		log.Fatalf("Error Occured : %v", err)
	}
	return returnValue
}
func DebugLog(msg string) {
	if !configuration.Debug {
		return
	}

	log.Debug(msg)
}
func DebugLogf(template string, args ...interface{}) {
	if !configuration.Debug {
		return
	}
	log.Debug(template, args)
}
