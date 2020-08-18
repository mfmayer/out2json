package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"sync/atomic"
)

var (
	lineCount int64 = 0
)

const (
	lineIndexProperty string = "i"
	lineDataProoperty string = "d"
)

func printUsage() {
	fmt.Printf("%v: missing command\n", os.Args[0])
	fmt.Printf("Usage: %v COMMAND ARGUMENTS...\n", os.Args[0])
	fmt.Printf(" Try : %v ls -al\n", os.Args[0])
}

// write creates a line inder and wraps data into json format before writing it to out
func write(data interface{}, out io.Writer) {
	lineNo := strconv.Itoa(int(atomic.AddInt64(&lineCount, 1)))
	entry := map[string]interface{}{
		lineIndexProperty: lineNo,
		lineDataProoperty: data,
	}
	marshalledEntry, _ := json.Marshal(entry)
	out.Write(marshalledEntry)
	out.Write([]byte("\n"))
}

// processStream reads line by line from given reader and quits if reader is closed or EOF is reached
func processStream(r io.Reader, out io.Writer, wg *sync.WaitGroup) {
	defer wg.Done()
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadBytes('\n')
		if lineLen := len(line); lineLen > 0 {
			if line[lineLen-1] == '\n' {
				line = line[:lineLen-1]
			}
		}
		if len(line) > 0 {
			write(string(line), out)
		}
		if err != nil {
			break
		}
	}
}

func main() {
	if len(os.Args) <= 1 {
		printUsage()
		os.Exit(0)
	}

	// get parameters and prepare command
	params := []string{}
	if len(os.Args) > 1 {
		params = os.Args[2:]
	}
	cmd := exec.Command(os.Args[1], params...)

	// create pipes to read stderr and stdout
	stderrReader, stderrWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()

	cmd.Stdin = os.Stdin
	cmd.Stdout = stderrWriter
	cmd.Stderr = stdoutWriter

	wg := sync.WaitGroup{}
	wg.Add(2)
	go processStream(stderrReader, os.Stderr, &wg)
	go processStream(stdoutReader, os.Stdout, &wg)

	// start command
	if err := cmd.Start(); err != nil {
		write(err, os.Stderr)
		os.Exit(1)
	}

	// wait until process is finished
	exitStatus := 1
	if err := cmd.Wait(); err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if !ok {
			write(err, os.Stderr)
			exitStatus = 1
		}
		exitStatus = exitErr.ExitCode()
	}

	// close pipes
	stderrWriter.Close()
	stdoutWriter.Close()

	// wait for go routiens that process stderr & stdout
	wg.Wait()

	// exit with same exitstatus as started application
	os.Exit(exitStatus)
}
