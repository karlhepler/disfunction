package log

type Logger interface {
	Debugf(string, ...any)
	Error(error)
}
