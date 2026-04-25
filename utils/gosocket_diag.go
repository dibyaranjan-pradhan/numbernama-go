package utils

// GoSocketDiag implements go-socket's Logger interface (Printf only).
// Use it exclusively when wiring gosocket.Config.Logger so STAG code does not call Printf on *Logger directly.
type GoSocketDiag struct {
	L *Logger
}

// Printf forwards diagnostic lines to the shared app logger at info level.
func (d GoSocketDiag) Printf(format string, v ...interface{}) {
	if d.L == nil {
		return
	}
	d.L.Infof(format, v...)
}
