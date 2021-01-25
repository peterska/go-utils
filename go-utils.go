/*
 *
 * Copyright (c) 2021 Peter Skarpetis
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
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

	// logging object
	Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
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

// returns caller function name
func Callername() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	return shortFuncname(f.Name())
}

// returns function name
func Funcname() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return shortFuncname(f.Name())
}

// return the function name and line number of the caller
func Trace() {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fmt.Printf("%s:%d %s\n", file, line, shortFuncname(f.Name()))
}

// return true if running under kubernetes
func KubernetesPod() bool {
	_, exists := os.LookupEnv("KUBERNETES_PORT")
	return exists
}

// set production mode that can be queried by other functions
func Setproduction() {
	SetDebuglevel(0)
}

// set the debug level that can be queried by other functions
func SetDebuglevel(level int) {
	debuglevel = level
	DisableProfiling()
}

// check whether go profiling is enabled
func ProfilingEnabled() bool {
	return runtime.MemProfileRate != 0
}

// disables go profiling
func DisableProfiling() {
	runtime.MemProfileRate = 0
	runtime.SetBlockProfileRate(0)
	runtime.SetCPUProfileRate(0)
	runtime.SetMutexProfileFraction(0)
}

// enables go profiling
func EnableProfiling() {
	runtime.SetBlockProfileRate(1000)
	runtime.SetCPUProfileRate(1000)
	runtime.SetMutexProfileFraction(10)
	runtime.MemProfileRate = 10
}

// Set the state of go profiling
func SetProfiling(on bool) {
	if on {
		EnableProfiling()
	} else {
		DisableProfiling()
	}
}

// check whether we are running in production mode
func Production() bool {
	return debuglevel == 0
}

// check whether we are running in development mode
func Development() bool {
	return debuglevel > 0
}

// return the current debug level
func Debuglevel() int {
	return debuglevel
}

// set the loglevel or verbosity
func SetLoglevel(level int) {
	loglevel = level
}

// check the loglevel
func Loglevel() int {
	return loglevel
}

// dump the passed interface as JSON
func PrintAsJSON(v interface{}) {
	json.NewEncoder(os.Stdout).Encode(v)
}

// read a password from the terminal, characters are not echoed
func ReadPassword(prompt string) ([]byte, error) {
	fmt.Fprintf(os.Stderr, "%v ", prompt)
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Fprintf(os.Stderr, "\n")
	return password, err
}

// read a string from the terminal
func ReadString(prompt string) ([]byte, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(os.Stderr, "%v ", prompt)
	line, _ := reader.ReadString('\n')
	line = strings.TrimRight(line, " \n\r")
	return []byte(line), nil
}
