/*
Copyright 2016 The Kubernetes Authors.
Copyright 2020 Authors of Arktos - file modified.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package net

import (
	"io"
	"net"
	"net/url"
	"os"
	"reflect"
	"strings"
	"syscall"
)

// IPNetEqual checks if the two input IPNets are representing the same subnet.
// For example,
//	10.0.0.1/24 and 10.0.0.0/24 are the same subnet.
//	10.0.0.1/24 and 10.0.0.0/25 are not the same subnet.
func IPNetEqual(ipnet1, ipnet2 *net.IPNet) bool {
	if ipnet1 == nil || ipnet2 == nil {
		return false
	}
	if reflect.DeepEqual(ipnet1.Mask, ipnet2.Mask) && ipnet1.Contains(ipnet2.IP) && ipnet2.Contains(ipnet1.IP) {
		return true
	}
	return false
}

// Returns if the given err is "connection reset by peer" error.
func IsConnectionReset(err error) bool {
	if urlErr, ok := err.(*url.Error); ok {
		err = urlErr.Err
	}
	if opErr, ok := err.(*net.OpError); ok {
		err = opErr.Err
	}
	if osErr, ok := err.(*os.SyscallError); ok {
		err = osErr.Err
	}
	if errno, ok := err.(syscall.Errno); ok && errno == syscall.ECONNRESET {
		return true
	}
	return false
}

// Returns if the given err is "connection refused" error
func IsConnectionRefused(err error) bool {
	if urlErr, ok := err.(*url.Error); ok {
		err = urlErr.Err
	}
	if opErr, ok := err.(*net.OpError); ok {
		err = opErr.Err
	}
	if osErr, ok := err.(*os.SyscallError); ok {
		err = osErr.Err
	}
	if errno, ok := err.(syscall.Errno); ok && errno == syscall.ECONNREFUSED {
		return true
	}
	return false
}

// IsTimeout returns true if the given error is a network timeout error
func IsTimeout(err error) bool {
	neterr, ok := err.(net.Error)
	return ok && neterr != nil && neterr.Timeout()
}

// IsProbableEOF returns true if the given error resembles a connection termination
// scenario that would justify assuming that the watch is empty.
// These errors are what the Go http stack returns back to us which are general
// connection closure errors (strongly correlated) and callers that need to
// differentiate probable errors in connection behavior between normal "this is
// disconnected" should use the method.
func IsProbableEOF(err error) bool {
	if err == nil {
		return false
	}
	if uerr, ok := err.(*url.Error); ok {
		err = uerr.Err
	}
	msg := err.Error()
	switch {
	case err == io.EOF:
		return true
	case msg == "http: can't write HTTP request on broken connection":
		return true
	case strings.Contains(msg, "http2: server sent GOAWAY and closed the connection"):
		return true
	case strings.Contains(msg, "connection reset by peer"):
		return true
	case strings.Contains(strings.ToLower(msg), "use of closed network connection"):
		return true
	}
	return false
}
