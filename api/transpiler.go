package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/learn-systemverilog/learn-systemverilog-api/transpiler"
)

// Transpile Begin transpilation and send logs in realtime.
func Transpile(c *gin.Context) {
	code := c.Query("code")

	logs := make(chan interface{})
	outputChan := make(chan string)

	go func() {
		output, err := transpiler.Transpile(code, logs)
		if err == nil {
			outputChan <- output
		}

		close(outputChan)
	}()

	sseSetup(c)

	for log := range logs {
		j, err := json.Marshal(log)
		if err != nil {
			panic(err)
		}

		sseStep(c, getLogName(log), string(j))
	}

	for output := range outputChan {
		j, err := json.Marshal(output)
		if err != nil {
			panic(err)
		}

		sseStep(c, "output", string(j))
	}

	sseClose(c)

	// Cleaning the channels so the transpiler can continue.
	for range logs {
	}
	for range outputChan {
	}
}

func sseSetup(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")

	c.Writer.Flush()
}

func sseStep(c *gin.Context, name, data string) {
	c.String(http.StatusOK, "data: %s\n\n", data)

	c.Writer.Flush()
}

func sseClose(c *gin.Context) {
	c.Status(http.StatusNoContent)
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
