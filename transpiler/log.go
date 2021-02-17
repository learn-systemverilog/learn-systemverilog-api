package transpiler

// LogInternal A internal log.
type LogInternal struct {
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

// LogStdout A log from the stdout.
type LogStdout struct {
	Stdout string `json:"stdout"`
}

// LogStderr A log from the stderr.
type LogStderr struct {
	Stderr string `json:"stderr"`
}

func newLogInternal(msg, severity string) LogInternal {
	return LogInternal{Message: msg, Severity: severity}
}

func newLogStdout(stdout string) LogStdout {
	return LogStdout{Stdout: stdout}
}

func newLogStderr(stderr string) LogStderr {
	return LogStderr{Stderr: stderr}
}
