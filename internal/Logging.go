package internal

import (
	"io/ioutil"
	"log"
	"os"
)

var _log_stderr *log.Logger = log.New(os.Stderr, "", 0)
var _log_null *log.Logger = log.New(ioutil.Discard, "", 0)

var _log_info *log.Logger = _log_stderr
var _log_verbose *log.Logger = _log_null

func SetupLogging(verbose bool) {
	if verbose {
		_log_verbose = log.New(os.Stderr, "LOG ", 0)
	} else {
		_log_verbose = _log_null
	}
}

func FinishLogging() {
}

func LError() *log.Logger {
	return _log_info
}

func LWarning() *log.Logger {
	return _log_info
}

func LInfo() *log.Logger {
	return _log_info
}

func LVerbose() *log.Logger {
	return _log_verbose
}
