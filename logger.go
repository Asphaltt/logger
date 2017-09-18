package logger

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	debug   bool
	verbose bool

	logDir     = "log"
	logDay     = 0
	logLock    = sync.Mutex{}
	logDefault = os.Stdout
	logFile    *os.File
)

func IsDebug() bool {
	return debug
}

func offset() string {
	_, file, line, _ := runtime.Caller(2)
	return fileline(file, line)
}

func fileline(file string, line int) string {
	strs := strings.Split(file, "src/")
	file = strs[len(strs)-1]
	return fmt.Sprint(file, ":", line)
}

func check() {
	logLock.Lock()
	defer logLock.Unlock()

	now := time.Now().UTC()
	if logDay == now.Day() {
		return
	}

	logDay = now.Day()
	logProc := filepath.Base(os.Args[0])
	filename := filepath.Join(logDir,
		fmt.Sprintf("%s.%s.log", logProc, now.Format("2006-01-02")))
	newLog, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.SetOutput(logDefault)
		log.Println("cannot create the log file, log in STDOUT")
		return
	}
	logFile.Sync()
	logFile.Close()
	logFile = newLog
	log.SetOutput(logFile)
}

func Setup(isDebug, isVerbose bool, dir string, outputDefault *os.File) {
	debug = isDebug
	verbose = isVerbose
	logDir = dir
	logDefault = outputDefault
}

func Log(v ...interface{}) {
	check()
	log.Println(v...)
}

func Logf(fmt string, v ...interface{}) {
	check()
	log.Printf(offset()+" "+fmt, v...)
}

func Debug(v ...interface{}) {
	if !debug {
		return
	}
	check()
	l := append([]interface{}{offset(), "<DEBUG>"}, v...)
	log.Println(l...)
}

func Debugf(fmt string, v ...interface{}) {
	if !debug {
		return
	}
	check()
	log.Printf(offset()+" <DEBUG> "+fmt, v...)
}

func Info(v ...interface{}) {
	check()
	l := append([]interface{}{offset(), "<INFO>"}, v...)
	log.Println(l...)
}

func Infof(fmt string, v ...interface{}) {
	check()
	log.Printf(offset()+" <INFO> "+fmt, v...)
}

func Warning(v ...interface{}) {
	check()
	l := append([]interface{}{offset(), "<WARNING>"}, v...)
	log.Println(l...)
}

func Warningf(fmt string, v ...interface{}) {
	check()
	log.Printf(offset()+" <WARNING> "+fmt, v...)
}

func Error(v ...interface{}) {
	check()
	l := append([]interface{}{offset(), "<ERROR>"}, v...)
	log.Println(l...)
}

func Errorf(fmt string, v ...interface{}) {
	check()
	log.Printf(offset()+" <ERROR> "+fmt, v...)
}

func Trace(data []byte, v ...interface{}) {
	if debug {
		return
	}
	check()
	l := append([]interface{}{offset(), "<TRACE>"}, v...)
	log.Println(l...)
	if verbose {
		log.Println(hex.Dump(data))
	}
}

func Tracef(data []byte, fmt string, v ...interface{}) {
	if debug {
		return
	}
	check()
	log.Printf(offset()+" <TRACE> "+fmt, v...)
	if verbose {
		log.Println(hex.Dump(data))
	}
}

func Stderr(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}

func Stderrf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
}
