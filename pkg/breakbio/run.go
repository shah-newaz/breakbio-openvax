package breakbio

import (
	"bufio"
	"bytes"
	"breakbio-openvax/pkg/breakbio/executor"
	"breakbio-openvax/pkg/breakbio/log"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	StandardEOF     = "##BREAKBIO##="
	defaultExitCode = 500
)

type RunParams struct {
	Shell                 string
	Command               string
}

func (r *RunParams) getCommand() string {
	return r.Command
}

func RunCommand(params RunParams) {
	exitCode := doRun(executor.NewExecutor(), params)
	os.Exit(exitCode)
}

func doRun(localExecutor executor.ExecutorInterface, params RunParams) int {
	var initialValue []byte
	bufWriter := bytes.NewBuffer(initialValue)
	log.TeeGreen(bufWriter, fmt.Sprintf(`Running "%s"`, params.getCommand()))

	var exitCode int
	exitCode = localExecutor.Execute(bufWriter, params.Shell, params.getCommand())

	if exitCode == 0 {
		log.TeeGreen(bufWriter, "Exit code: 0")
	} else {
		log.TeeRed(bufWriter, fmt.Sprintf("Exit code: %d", exitCode))
	}
	fmt.Println()
	return exitCode
}

func listeningToTheLog(getUrl string, writer io.Writer) int {
	reconnects := 10
	if val, ok := os.LookupEnv("RECONNECT_ATTEMPTS"); !ok {
		if i, err := strconv.Atoi(val); err == nil {
			reconnects = i
		}
	}

	for i := 0; i < reconnects; i++ {
		resp, err := http.Get(getUrl)

		if err != nil {
			log.TeeRed(writer, fmt.Sprintf("Cannot connect log listener endpoint on %s, Error: %v", getUrl, err))
			continue
		}

		exitCode, err := streamLogsToStdOut(bufio.NewScanner(resp.Body), writer)

		if err == nil {
			if resp.Body != nil {
				resp.Body.Close()
			}
			return exitCode
		}

		log.TeeNoColor(writer, fmt.Sprintf("Reconnecing log listener. Attempt: %d, Reason: %v", i+1, err))
		if resp.Body != nil {
			resp.Body.Close()
		}
	}

	log.TeeRed(writer, "Failed to reconnect to a log listener endpoint.")
	return defaultExitCode
}

func streamLogsToStdOut(reader *bufio.Scanner, writer io.Writer) (int, error) {
	for reader.Scan() {
		line := reader.Text()
		if strings.Contains(line, StandardEOF) {
			exitCode, err := extractExitCode(line)
			if err != nil {
				log.TeeRed(writer, fmt.Sprintf("%v", err))
				return defaultExitCode, err
			}
			return exitCode, nil
		}
		log.TeeNoColor(writer, string(line))
	}
	if err := reader.Err(); err != nil {
		log.TeeRed(writer, fmt.Sprintf("reading error: %v", err))
		return defaultExitCode, err
	}

	return defaultExitCode, fmt.Errorf("Stream finished but exit code tag was missing")
}

func extractExitCode(line string) (int, error) {
	regex := regexp.MustCompile(StandardEOF + `(\d*)`)
	res := regex.FindStringSubmatch(line)

	if len(res) != 2 {
		return 0, fmt.Errorf("error parsing EOF token. Exit code not found")
	}

	exitCode, err := strconv.Atoi(res[1])
	return exitCode, err
}
