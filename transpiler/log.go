package transpiler

// Log a interface for a log.
type Log interface {
	Name() string
}

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

// Name Returns the instance type name.
func (LogInternal) Name() string {
	return "internal"
}

func newLogStdout(stdout string) LogStdout {
	return LogStdout{Stdout: stdout}
}

// Name Returns the instance type name.
func (LogStdout) Name() string {
	return "stdout"
}

func newLogStderr(stderr string) LogStderr {
	return LogStderr{Stderr: stderr}
}

// Name Returns the instance type name.
func (LogStderr) Name() string {
	return "stderr"
}
