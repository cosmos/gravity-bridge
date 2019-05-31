package log

// -----------------------------------------------------
//    Log
//
//		Intializes logging with formatted strings.
// -----------------------------------------------------

import (
	"os"

	logging "github.com/op/go-logging"
)

// Logger instance
var Log = logging.MustGetLogger("logger")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func init() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backendFormatter := logging.NewBackendFormatter(backend, format)

	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(logging.DEBUG, "")

	// Set the backends to be used.
	logging.SetBackend(backendLeveled)
}
