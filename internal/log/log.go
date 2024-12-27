package log

type Logger interface {
	Debug(...any)
	Debugf(string, ...any)
	Error(error)
}
