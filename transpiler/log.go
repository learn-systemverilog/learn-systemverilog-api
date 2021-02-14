package transpiler

// LogInternal A internal log.
type LogInternal struct {
	Message  string
	Severity string
}

// LogStdout A log from the stdout.
type LogStdout struct {
	Stdout string
}

// LogStderr A log from the stderr.
type LogStderr struct {
	Stderr string
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
