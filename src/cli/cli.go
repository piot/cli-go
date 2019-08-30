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

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"github.com/piot/log-go/src/clog"
)

func runWithLog(cliStructReference interface{}, log *clog.Log) error {
	ctx := kong.Parse(cliStructReference)
	err := ctx.Run(log)
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
}

func Run(cliStructReference interface{}, options RunOptions) {
	name := os.Args[0]
	color.New(color.FgCyan).Fprintf(os.Stderr, "%v %v\n", name, options.Version)

	log := clog.DefaultLog()
	err := runWithLog(cliStructReference, log)
	if err != nil {
		log.Err(err)
		os.Exit(-1)
	}
	os.Exit(0)
}
