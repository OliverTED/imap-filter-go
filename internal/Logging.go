package internal

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var _log_stderr *log.Logger = log.New(os.Stderr, "", 0)
var _log_null *log.Logger = log.New(ioutil.Discard, "", 0)

var _log_info *log.Logger = _log_stderr
var _log_verbose *log.Logger = _log_null

func SetupLogging(verbose bool, logfile string) func() {
	var out io.Writer = os.Stderr
	var log_file *os.File = nil

	if logfile != "" {
		log_file, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Panicln("error opening log file: ", logfile, err)
		}

		out = io.MultiWriter(out, log_file)
	}
	if verbose {
		_log_verbose = log.New(out, "LOG ", 0)
	} else {
		_log_verbose = _log_null
	}
	_log_info = log.New(out, "", 0)

	return func() {
		if log_file != nil {
			err := log_file.Close()
			if err != nil {
				log.Panicln("error closing log file: ", logfile, err)
			}
		}
	}
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
