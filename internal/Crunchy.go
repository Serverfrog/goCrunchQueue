package internal

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// Interfaces with crunchy-cli

const CRUNCHY_CLI_BIN = "/usr/bin/crunchy-cli"

var CRUNCHY_PARAMTERS = []string{"archive",
	"-s", "de-DE",
	"-o", "\"%v/{series_name}/S{season_number}E{relative_episode_number} {title}.mkv\"",
	"--skip-existing",
	"--default-subtitle", "de-DE",
	"%v"}

func execCrunchy(item QueueItem) {

	preparedParameters := CRUNCHY_PARAMTERS
	preparedParameters[4] = fmt.Sprintf(preparedParameters[4], configuration.MediaDestination)
	preparedParameters[8] = fmt.Sprintf(preparedParameters[8], item.CrunchyrollUrl)

	command := exec.Command(CRUNCHY_CLI_BIN, preparedParameters...)

	stdoutFile := HandleError(os.Create(fmt.Sprintf("%v/%v-out.txt", configuration.LogDestination, item.Id)))
	defer HandleErrorB(stdoutFile.Close())

	stderrFile := HandleError(os.Create(fmt.Sprintf("%v/%v-err.txt", configuration.LogDestination, item.Id)))
	defer HandleErrorB(stderrFile.Close())

	command.Stdout = stdoutFile
	command.Stderr = stderrFile
	HandleErrorB(command.Run())
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
