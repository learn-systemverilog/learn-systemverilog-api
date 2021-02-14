package transpiler

import "log"

// Log A log.
type Log struct{}

// LogInternal A internal log.
type LogInternal struct {
	Log

	Message  string
	Severity string
}

// LogStdout A log from the stdout.
type LogStdout struct {
	Log

	Stdout string
}

// LogStderr A log from the stderr.
type LogStderr struct {
	Log

	Stderr string
}

func logInternal(msg, workspace, severity string) LogInternal {
	log.Println(msg)

	return LogInternal{Message: msg, Severity: severity}
}

func newLogStdout(stdout string) LogStdout {
	return LogStdout{Stdout: stdout}
}

func newLogStderr(stderr string) LogStderr {
	return LogStderr{Stderr: stderr}
}
