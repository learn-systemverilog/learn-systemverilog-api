package api

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/learn-systemverilog/learn-systemverilog-api/transpiler"
)

// Transpile Begin transpilation and send logs in realtime.
func Transpile(c *gin.Context) {
	logs := make(chan interface{})
	outputChan := make(chan string)

	go func() {
		output, err := transpiler.Transpile(transpiler.DummyWorkingCode, logs)

		close(logs)

		if err == nil {
			outputChan <- output
		}

		close(outputChan)
	}()

	_ = c.Stream(func(w io.Writer) bool {
		if log, ok := <-logs; ok {
			name := getLogName(log)
			c.SSEvent(name, log)

			return true
		}

		if output, ok := <-outputChan; ok {
			c.SSEvent("output", output)
		}

		return false
	})

	// Cleaning the channels so the transpiler can continue.
	for range logs {
	}
	for range outputChan {
	}
}

func getLogName(log interface{}) string {
	if _, ok := log.(transpiler.LogInternal); ok {
		return "internal"
	}

	if _, ok := log.(transpiler.LogStdout); ok {
		return "stdout"
	}

	if _, ok := log.(transpiler.LogStderr); ok {
		return "stderr"
	}

	return "unknown"
}
