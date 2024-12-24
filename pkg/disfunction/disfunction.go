package disfunction

type Status int

const (
	StatusOK Status = iota
	StatusError
)

func DebugLog(msg string) string {
	return "[DEBUG] " + msg
}

func InfoLog(msg string) string {
	return "[INFO] " + msg
}

func ErrorLog(err error) string {
	return "[ERROR] " + err.Error()
}
