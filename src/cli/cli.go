/*

MIT License

Copyright (c) 2018 Peter Bjorklund

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package cli

import (
	"os"
	"path/filepath"
	"reflect"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"github.com/piot/log-go/src/clog"
	"github.com/piot/log-go/src/clogint"
)

type LogConfigFile struct {
	ApplicationName string
	LogDirectory    string
}

type LogLevelString struct {
	Level string
}

type LogFileString struct {
	File string
}

type Options struct {
	LogLevel LogLevelString `help:"the log level"`
	LogFile  LogFileString  `help:"log output file" default:""`
}

func (o LogLevelString) Decode(ctx *kong.DecodeContext, target reflect.Value) error {
	stringField := target.FieldByName("Level")
	logLevelToken := ctx.Scan.Pop()
	stringField.SetString(logLevelToken.String())
	return nil
}

func (o LogFileString) Decode(ctx *kong.DecodeContext, target reflect.Value) error {
	stringField := target.FieldByName("File")
	logLevelToken := ctx.Scan.Pop()
	stringField.SetString(logLevelToken.String())
	return nil
}

func (o LogFileString) BeforeApply(log *clog.Log) error {
	fileName := o.File
	dir, file := filepath.Split(fileName)
	info := &LogConfigFile{ApplicationName: file, LogDirectory: dir}
	logFile, _ := createFileAndConsoleLog(info)
	log = logFile
	return nil
}

func (o LogLevelString) AfterApply(log *clog.Log) error {
	log.SetLogLevelUsingString(o.Level, clogint.Info)
	return nil
}

func runWithLog(cliStructReference interface{}, log *clog.Log, customArgs []interface{}) error {
	ctx := kong.Parse(cliStructReference, kong.Bind(log), kong.TypeMapper(reflect.TypeOf(LogLevelString{}), &LogLevelString{}), kong.TypeMapper(reflect.TypeOf(LogFileString{}), &LogFileString{}))
	logPlusRest := []interface{}{customArgs, log}
	err := ctx.Run(logPlusRest)
	return err
}

type ApplicationType = int

const (
	Daemon ApplicationType = iota
	Utility
)

type RunOptions struct {
	Version         string
	ApplicationType ApplicationType
	LogConfig       interface{}
}

func createFileAndConsoleLog(info *LogConfigFile) (*clog.Log, error) {
	loggers := make([]clog.Logger, 2)
	fileLog, fileLogErr := clog.NewFileLogger(info.LogDirectory, info.ApplicationName)
	if fileLogErr != nil {
		return nil, fileLogErr
	}
	loggers[0] = fileLog
	consoleLogger := clog.NewConsoleLogger()
	loggers[1] = consoleLogger

	multiLog, multiLogErr := clog.NewMultiLog(loggers...)
	if multiLogErr != nil {
		return nil, multiLogErr
	}
	return multiLog, nil
}

func internalRun(cliStructReference interface{}, options RunOptions, customArgs []interface{}) {
	name := filepath.Base(os.Args[0])
	color.New(color.FgCyan).Fprintf(os.Stderr, "%v %v\n", name, options.Version)

	var log *clog.Log
	if options.LogConfig != nil {
		fileLog, _ := options.LogConfig.(*LogConfigFile)
		if fileLog != nil {
			color.New(color.FgCyan).Fprintf(os.Stderr, "Multilog! %v %v\n", fileLog.ApplicationName, fileLog.LogDirectory)
			var multiLogErr error
			log, multiLogErr = createFileAndConsoleLog(fileLog)
			if multiLogErr != nil {
				panic(multiLogErr)
			}
		}
	} else {
		log = clog.DefaultLog()
	}
	log.SetLogLevel(clogint.Info)
	err := runWithLog(cliStructReference, log, customArgs)
	if err != nil {
		log.Err(err)
		os.Exit(-1)
	}
	os.Exit(0)
}

func RunWithArguments(cliStructReference interface{}, options RunOptions, customArgs ...interface{}) {
	internalRun(cliStructReference, options, customArgs)
}

func Run(cliStructReference interface{}, options RunOptions) {
	internalRun(cliStructReference, options, nil)
}
