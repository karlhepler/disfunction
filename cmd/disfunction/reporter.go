package main

type ErrorLogger interface {
	Error(err error)
}

type ConsoleReporter struct {
	Log ErrorLogger
}

func NewConsoleReporter(logger ErrorLogger) ConsoleReporter {
	return ConsoleReporter{Log: logger}
}

func (console ConsoleReporter) Report(err error) {
	console.Log.Error(err)
}
