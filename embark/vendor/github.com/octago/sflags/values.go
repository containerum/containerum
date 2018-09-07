package sflags

//go:generate go run ./cmd/genvalues/main.go

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Value is the interface to the dynamic value stored in v flag.
// (The default value is represented as v string.)
//
// If v Value has an IsBoolFlag() bool method returning true, the command-line
// parser makes --name equivalent to -name=true rather than using the next
// command-line argument, and adds v --no-name counterpart for negating the
// flag.
type Value interface {
	String() string
	Set(string) error

	// pflag.Flag require this
	Type() string
}

// Getter is an interface that allows the contents of v Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
type Getter interface {
	Value
	Get() interface{}
}

// BoolFlag is an optional interface to indicate boolean flags
// that don't accept v value, and implicitly have v --no-<x> negation counterpart.
type BoolFlag interface {
	Value
	IsBoolFlag() bool
}

// RepeatableFlag is an optional interface for flags that can be repeated.
// required by kingpin
type RepeatableFlag interface {
	Value
	IsCumulative() bool
}

// === Custom values

type validateValue struct {
	Value
	validateFunc func(val string) error
}

func (v *validateValue) IsBoolFlag() bool {
	if boolFlag, casted := v.Value.(BoolFlag); casted {
		return boolFlag.IsBoolFlag()
	}
	return false
}

func (v *validateValue) IsCumulative() bool {
	if cumulativeFlag, casted := v.Value.(RepeatableFlag); casted {
		return cumulativeFlag.IsCumulative()
	}
	return false
}

func (v *validateValue) String() string {
	if v == nil || v.Value == nil {
		return ""
	}
	return v.Value.String()
}

func (v *validateValue) Set(val string) error {
	if v.validateFunc != nil {
		err := v.validateFunc(val)
		if err != nil {
			return err
		}
	}
	return v.Value.Set(val)
}

// HexBytes might be used if you want to parse slice of bytes as hex string.
// Original `[]byte` or `[]uint8` parsed as a list of `uint8`.
type HexBytes []byte

// Counter type is useful if you want to save number
// by using flag multiple times in command line.
// It's a boolean type, so you can use it without value.
// If you use `struct{count Counter}
// and parse it with `-count=10 ... -count .. -count`,
// then final value of `count` will be 12.
// Implements Value, Getter, BoolFlag, RepeatableFlag interfaces
type Counter int

var _ RepeatableFlag = (*Counter)(nil)

// Set method parses string from command line.
func (v *Counter) Set(s string) error {
	// flag package pass true if BoolFlag doesn't have an argument.
	if s == "" || s == "true" {
		*v++
		return nil
	}
	parsed, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}
	// -1 means that no specific value was passed, so increment
	if parsed == -1 {
		*v++
	} else {
		*v = Counter(parsed)
	}
	return nil
}

// Get method returns inner value for Counter.
func (v *Counter) Get() interface{} { return (int)(*v) }

// IsBoolFlag returns true, because Counter might be used without value.
func (v *Counter) IsBoolFlag() bool { return true }

// String returns string representation of Counter
func (v *Counter) String() string { return fmt.Sprintf("%d", *v) }

// IsCumulative returns true, because Counter might be used multiple times.
func (v *Counter) IsCumulative() bool { return true }

// Type returns `count` for Counter, it's mostly for pflag compatibility.
func (v *Counter) Type() string { return "count" }

// === Some patches for generated flags

// IsBoolFlag returns true. boolValue implements BoolFlag interface.
func (v *boolValue) IsBoolFlag() bool { return true }

// === Custom parsers

func parseIP(s string) (net.IP, error) {
	ip := net.ParseIP(strings.TrimSpace(s))
	if ip == nil {
		return nil, fmt.Errorf("failed to parse IP: %q", s)
	}
	return ip, nil
}

func parseTCPAddr(s string) (net.TCPAddr, error) {
	tcpADDR, err := net.ResolveTCPAddr("tcp", strings.TrimSpace(s))
	if err != nil {
		return net.TCPAddr{}, fmt.Errorf("failed to parse TCPAddr: %q", s)
	}
	return *tcpADDR, nil
}

func parseIPNet(s string) (net.IPNet, error) {
	_, ipNet, err := net.ParseCIDR(s)
	if err != nil {
		return net.IPNet{}, err
	}
	return *ipNet, nil
}
