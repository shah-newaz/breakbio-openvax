package breakbio

import (
	"bytes"
	"context"
	"breakbio-openvax/pkg/breakbio/executor"
	"breakbio-openvax/pkg/breakbio/log"
	"fmt"
	"io"
	"io/ioutil"
	log2 "log"
	"net/http"
)

func Serve() {
	server := NewBreakBioServer(executor.NewExecutor())
	http.HandleFunc("/run", server.RunHandler)
	http.HandleFunc("/listen", server.Listener)
	log.Cyan("Starting BreakBio server")
	log2.Fatal(http.ListenAndServe(":8080", nil))
}

type CommandServerInterface interface {
	RunHandler(w http.ResponseWriter, r *http.Request)
}

type CommandServer struct {
	executorObj       executor.ExecutorInterface
	buffer            *bytes.Buffer
	execChan          chan int
	commandCompleted  bool
	flusher           http.Flusher
	responseTeeReader io.Reader
}

func NewBreakBioServer(executorObj executor.ExecutorInterface) *CommandServer {
	return &CommandServer{
		executorObj:      executorObj,
		commandCompleted: true,
	}
}

func (c *CommandServer) RunHandler(w http.ResponseWriter, r *http.Request) {
	log.Greenf("BreakBio received %v", r)
	if !c.commandCompleted {
		respondWithStatusCode(w, "BreakBio is already executing a command", http.StatusConflict)
		return
	}

	shell := "bash"
	command := r.FormValue("COMMAND")

	var initialValue []byte
	c.buffer = bytes.NewBuffer(initialValue)

	c.execChan = make(chan int, 1)

	go func() {
		c.commandCompleted = false
		exitCode := c.executorObj.Execute(c, shell, command)
		c.execChan <- exitCode
		log.Greenf("Command finished with exit code %d", exitCode)
		c.commandCompleted = true
	}()

	w.WriteHeader(http.StatusOK)
	respondWithMessage(w, "Command received")
	log.Green("Execution started")
}

func (c *CommandServer) WriteString(e string) (n int, err error) {
	n, err = c.buffer.WriteString(e)
	if err != nil {
		log.Redf("Failed to write the command output to the log file. Message: %+v", err)
		return n, err
	}

	if c.responseTeeReader != nil {
		_, err = ioutil.ReadAll(c.responseTeeReader)
		if err != nil {
			if err == io.EOF {
				return 0, nil
			}
			log.Redf("Error reading the logs. %v", err)
			return n, err
		}

		if c.flusher != nil {
			c.flusher.Flush()
		}
	}
	return n, nil
}

func (c *CommandServer) Listener(w http.ResponseWriter, r *http.Request) {
	log.Green("Client started listening log stream")

	f, ok := w.(http.Flusher)
	if !ok {
		panic("Writer does not implement flush")
	}
	c.flusher = f

	c.responseTeeReader = io.TeeReader(c.buffer, w)
	w.WriteHeader(http.StatusOK)

	c.flusher.Flush()
	c.handleInterruptions(r.Context(), w, c.buffer)
}

func respondWithMessage(w http.ResponseWriter, msg string) {
	log.Red(msg)
	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Redf("Error writing response %v", err)
	}
}

func respondWithStatusCode(w http.ResponseWriter, msg string, statusCode int) {
	log.Red(msg)
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(msg))
	if err != nil {
		log.Redf("Error writing response %v", err)
	}
}

func (c *CommandServer) handleInterruptions(ctx context.Context, writer http.ResponseWriter, reader io.Reader) {
	select {
	case <-ctx.Done():
		c.flusher = nil
		c.responseTeeReader = nil
		log.Red("Client stopped listening the logs")
		log.Red("Context stopped with " + ctx.Err().Error())
		return

	case exitCode := <-c.execChan:
		log.Green("Finalizing command execution")

		tail, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Redf("Error reading last part of log, Error: %v", err)
			return
		}
		_, err = writer.Write(tail)
		if err != nil {
			log.Redf("Error writing log, Error: %v", err)
			return
		}

		exitMessage := fmt.Sprintf("\n%s%d\n", "BREAKBIO", exitCode)
		_, err = writer.Write([]byte(exitMessage))
		if err != nil {
			log.Redf("Error writing exit message, Error: %v", err)
			return
		}

		if c.flusher != nil {
			c.flusher.Flush()
			c.flusher = nil
		}
		c.responseTeeReader = nil
		close(c.execChan)
		log.Green("Execution and processing completed")
		return

	}
}
