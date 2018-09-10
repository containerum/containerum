package emberr

import (
	"bytes"
	"fmt"

	"github.com/agnivade/levenshtein"
)

type Error interface {
	error
	Unwrap() error
}

// ---------------------------------------------------------------------------------------------------------------------

var (
	_ Error = ErrUnableToFetchChart{}
)

type ErrUnableToFetchChart struct {
	Chart  string
	Reason error
}

func (err ErrUnableToFetchChart) Error() string {
	return fmt.Sprintf("unable to fetch chart %q: %v", err.Chart, err.Reason)
}

func (err ErrUnableToFetchChart) Unwrap() error {
	return err.Reason
}

// ---------------------------------------------------------------------------------------------------------------------

var (
	_ Error = ErrUnableToDownloadDependencies{}
)

type ErrUnableToDownloadDependencies struct {
	Reason error
}

func (err ErrUnableToDownloadDependencies) Error() string {
	if err.Reason == nil {
		return fmt.Sprintf("unable to download dependencies: %v", err.Reason)
	}
	return "unable to download dependencies"
}

func (err ErrUnableToDownloadDependencies) Unwrap() error {
	return err.Reason
}

// ---------------------------------------------------------------------------------------------------------------------

var (
	_ Error = ErrUnableToInstallChart{}
)

type ErrUnableToInstallChart struct {
	Prefix string
	Chart  string
	Reason error
}

func (err ErrUnableToInstallChart) Error() string {
	var ff = "unable to install chart"
	if err.Chart != "" {
		ff += fmt.Sprintf(" %q", err.Chart)
	}
	if err.Prefix != "" {
		ff += ": " + err.Prefix
	}
	if err.Reason != nil {
		ff += fmt.Sprintf(": %v", err.Reason)
	}
	return ff
}

func (err ErrUnableToInstallChart) Unwrap() error {
	return err.Reason
}

// ---------------------------------------------------------------------------------------------------------------------

var (
	_ Error = ErrUnableToLoadChart{}
)

type ErrUnableToLoadChart struct {
	Chart  string
	Reason error
}

func (err ErrUnableToLoadChart) Error() string {
	return fmt.Sprintf("unable to install chart %q: %v", err.Chart, err.Reason)
}

func (err ErrUnableToLoadChart) Unwrap() error {
	return err.Reason
}

// ---------------------------------------------------------------------------------------------------------------------

var (
	_ Fatal = ErrUnableToInstallTiler{}
)

type ErrUnableToInstallTiler struct {
	defaultExitCoder
	Prefix string
	Reason error
}

func (err ErrUnableToInstallTiler) Error() string {
	var ff = "unable to install tiller"
	if err.Prefix != "" {
		ff = ff + ": " + err.Prefix
	}
	return fmt.Sprintf(ff+": %v", err.Reason)
}

func (err ErrUnableToInstallTiler) Unwrap() error {
	return err.Reason
}

// ---------------------------------------------------------------------------------------------------------------------

var (
	_ Fatal = ErrUnsupportedKubeObjectType("")
)

type ErrUnsupportedKubeObjectType string

func (err ErrUnsupportedKubeObjectType) Error() string {
	return fmt.Sprintf("unsupported kube object %q", string(err))
}

func (ErrUnsupportedKubeObjectType) Unwrap() error {
	return nil
}

func (ErrUnsupportedKubeObjectType) ExitCode() int {
	return 1
}

// ---------------------------------------------------------------------------------------------------------------------

var (
	_ Fatal = ErrUnableToCreateKubeCLient{}
)

type ErrUnableToCreateKubeCLient struct {
	defaultExitCoder
	Comment string
	Reason  error
}

func (err ErrUnableToCreateKubeCLient) Error() string {
	var prefix = "unable to create kube client"
	if err.Comment != "" {
		prefix += " " + err.Comment
	}
	return fmt.Sprintf("%s: %v", prefix, err.Reason)
}

func (err ErrUnableToCreateKubeCLient) Unwrap() error {
	return err.Reason
}

// ---------------------------------------------------------------------------------------------------------------------

type Chain struct {
	head error
	tail []error
}

func NewChain(head error, tail ...error) Chain {
	return Chain{
		head: head,
		tail: tail,
	}
}

func (err Chain) Head() error {
	return err.head
}

func (err Chain) Error() string {
	var buf = bytes.NewBufferString(err.head.Error() + ":\n")
	for _, e := range err.tail {
		fmt.Fprintf(buf, "\t%s\n", e)
	}
	return buf.String()
}

func (err Chain) Unwrap() error {
	switch len(err.tail) {
	case 0:
		return err.head
	case 1:
		return Chain{
			head: err.tail[0],
		}
	default:
		return Chain{
			head: err.tail[0],
			tail: append([]error{}, err.tail[1:]...),
		}
	}
}

// ---------------------------------------------------------------------------------------------------------------------

type ErrUnmarshalYAML struct {
	Filename string
	Reason   error
}

func (err ErrUnmarshalYAML) Error() string {
	var prefix = "unable to unmarshal YAML"
	if err.Filename != "" {
		prefix += fmt.Sprintf(" file %q", err.Filename)
	}
	return fmt.Sprintf("%s: %v", prefix, err.Reason)
}

func (err ErrUnmarshalYAML) Unwrap() error {
	return err.Reason
}

// ---------------------------------------------------------------------------------------------------------------------

type ErrInvalidTemplateDir struct {
	defaultExitCoder
	Comment string
	Reason  error
}

func (err ErrInvalidTemplateDir) Error() string {
	var msg = "invalid template dir"
	if err.Comment != "" {
		msg = fmt.Sprintf("%s %s", msg, err.Comment)
	}
	if err.Reason != nil {
		msg = fmt.Sprintf("%s: %v", msg, err.Reason)
	}
	return msg
}

// ---------------------------------------------------------------------------------------------------------------------

type ErrObjectNotFound struct {
	Name              string
	ObjectsWhichExist []string
}

func (err ErrObjectNotFound) findNearest() string {
	var minDist = -1
	var nearest = err.Name
	for _, exists := range err.ObjectsWhichExist {
		var dist = levenshtein.ComputeDistance(err.Name, exists)
		if dist < minDist || minDist < 0 {
			minDist = dist
			nearest = exists
		}
	}
	return nearest
}

func (err ErrObjectNotFound) Error() string {
	return fmt.Sprintf("object %q not found, maybe you mean %q", err.Name, err.findNearest())
}

// ---------------------------------------------------------------------------------------------------------------------

type ErrUnableToOpenObjectFile struct {
	File   string
	Reason error
}

func (err ErrUnableToOpenObjectFile) Error() string {
	return fmt.Sprintf("unable to open object file %q: %v", err.File, err.Reason)
}

// ---------------------------------------------------------------------------------------------------------------------

type ErrUnableToReadObjectFile struct {
	File   string
	Reason error
}

func (err ErrUnableToReadObjectFile) Error() string {
	return fmt.Sprintf("unable to read object file %q: %v", err.File, err.Reason)
}

// ---------------------------------------------------------------------------------------------------------------------

type ErrReadDir struct {
	Dir    string
	Reason error
}

func (err ErrReadDir) Error() string {
	return fmt.Sprintf("error while reading dir %q: %v", err.Dir, err.Reason)
}
