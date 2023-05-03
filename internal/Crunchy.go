package internal

import (
	"bufio"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
)

// Interfaces with crunchy-cli

const CRUNCHY_CLI_BIN = "/usr/bin/crunchy-cli"

var CRUNCHY_PARAMTERS = []string{"archive",
	"-s", "de-DE",
	"-o", "%v/{series_name}/S{season_number}E{relative_episode_number} {title}.mkv",
	"--skip-existing",
	"--default-subtitle", "de-DE",
	"%v"}

func execCrunchy(item QueueItem) error {

	eventHandler.handleEvent(Event{
		Id:      Process,
		Item:    item,
		Message: fmt.Sprintf("Execute crunchy-cli for Id:%v, Name:%v, Url:%v", item.Id, item.Name, item.CrunchyrollUrl),
	})
	defer eventHandler.handleEvent(Event{
		Id:      Processed,
		Item:    item,
		Message: fmt.Sprintf("Finished crunchy-cli for Id:%v, Name:%v, Url:%v", item.Id, item.Name, item.CrunchyrollUrl),
	})

	HandleErrorB(os.MkdirAll(configuration.LogDestination, 0664))

	preparedParameters := make([]string, len(CRUNCHY_PARAMTERS))
	copy(preparedParameters, CRUNCHY_PARAMTERS)
	preparedParameters[4] = fmt.Sprintf(CRUNCHY_PARAMTERS[4], configuration.MediaDestination)
	preparedParameters[8] = fmt.Sprintf(CRUNCHY_PARAMTERS[8], item.CrunchyrollUrl)

	command := exec.Command(CRUNCHY_CLI_BIN, preparedParameters...)

	stdoutFile, err := os.OpenFile(fmt.Sprintf("%v/%v-out.txt", configuration.LogDestination, item.Id), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		HandleErrorB(stdoutFile.Close())
		return err
	}

	stderrFile, err := os.OpenFile(fmt.Sprintf("%v/%v-err.txt", configuration.LogDestination, item.Id), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		HandleErrorB(stdoutFile.Close())
		HandleErrorB(stderrFile.Close())
		return err
	}
	stdout, err := command.StdoutPipe()
	if err != nil {
		HandleErrorB(stdoutFile.Close())
		HandleErrorB(stderrFile.Close())
		return err
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		HandleErrorB(stdoutFile.Close())
		HandleErrorB(stderrFile.Close())
		return err
	}
	var mw io.Writer
	var me io.Writer
	if configuration.Debug {
		mw = io.MultiWriter(os.Stdout, stdoutFile)
		me = io.MultiWriter(os.Stderr, stderrFile)
	} else {
		mw = stdoutFile
		me = stderrFile
	}
	log.Debugf("Start Process: %v", CRUNCHY_CLI_BIN)
	for i, parameter := range preparedParameters {
		log.Debugf("Paramter: %v", parameter)
		log.Debugf("Command String %v", command.Args[i])
	}
	err = command.Start()
	log.Debugf("Command String %v", command.String())

	// Print the output of the program to the console and file
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			HandleError(fmt.Fprintln(mw, scanner.Text()))
			eventHandler.handleEvent(Event{
				Id:   InfoLogUpdated,
				Item: item,
			})
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			HandleError(fmt.Fprintln(me, scanner.Text()))
			eventHandler.handleEvent(Event{
				Id:   ErrLogUpdated,
				Item: item,
			})
		}
	}()

	if err != nil {
		HandleErrorB(stdoutFile.Close())
		HandleErrorB(stderrFile.Close())
		return err
	}
	err = command.Wait()
	HandleErrorB(stdoutFile.Close())
	HandleErrorB(stderrFile.Close())
	return err
}

func crunchyValidation() error {
	fileInfo, err := os.Stat(CRUNCHY_CLI_BIN)
	if os.IsNotExist(err) {
		return err
	}
	if fileInfo.IsDir() {
		return errors.New(CRUNCHY_CLI_BIN + " is an Directory")
	}
	if !isExecAny(fileInfo.Mode()) {
		return errors.New(CRUNCHY_CLI_BIN + " is not executable")
	}

	return nil
}
