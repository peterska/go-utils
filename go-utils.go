/*

Copyright (c) 2021 Peter Skarpetis

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

package goutils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"runtime"
	"strings"
	"syscall"
)

var (
	debuglevel = 0
	loglevel   = 0
	Log        = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
)

func shortFuncname(s string) string {
	f := strings.Split(s, "/")
	c := len(f)
	if c > 0 {
		return f[c-1]
	}
	return s
}

func init() {
	DisableProfiling()
}

func Callername() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	return shortFuncname(f.Name())
}

func Funcname() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return shortFuncname(f.Name())
}

func Trace() {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fmt.Printf("%s:%d %s\n", file, line, shortFuncname(f.Name()))
}

// return trun if running under kubernetes
func KubernetesPod() bool {
	_, exists := os.LookupEnv("KUBERNETES_PORT")
	return exists
}

func Setproduction() {
	SetDebuglevel(0)
}

func SetDebuglevel(level int) {
	debuglevel = level
	DisableProfiling()
}

func ProfilingEnabled() bool {
	return runtime.MemProfileRate != 0
}

func DisableProfiling() {
	runtime.MemProfileRate = 0
	runtime.SetBlockProfileRate(0)
	runtime.SetCPUProfileRate(0)
	runtime.SetMutexProfileFraction(0)
}

func EnableProfiling() {
	runtime.SetBlockProfileRate(1000)
	runtime.SetCPUProfileRate(1000)
	runtime.SetMutexProfileFraction(10)
	runtime.MemProfileRate = 10
}

func SetProfiling(on bool) {
	if on {
		EnableProfiling()
	} else {
		DisableProfiling()
	}
}

func Production() bool {
	return debuglevel == 0
}

func Development() bool {
	return debuglevel > 0
}

func Debuglevel() int {
	return debuglevel
}

func SetLoglevel(level int) {
	loglevel = level
}

func Loglevel() int {
	return loglevel
}

func PrintAsJSON(v interface{}) {
	json.NewEncoder(os.Stdout).Encode(v)
}

func ReadPassword(prompt string) ([]byte, error) {
	fmt.Fprintf(os.Stderr, "%v ", prompt)
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Fprintf(os.Stderr, "\n")
	return password, err
}

func ReadString(prompt string) ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(os.Stderr, "%v ", prompt)
	line, _ := reader.ReadString('\n')
	line = strings.TrimRight(line, " \n\r")
	return []byte(line), nil
}
