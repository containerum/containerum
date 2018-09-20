package emberr

import (
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
	if err.Reason != nil {
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

var (
	_ Error = ErrUnmarshalYAML{}
)

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

var (
	_ Error = ErrInvalidTemplateDir{}
)

type ErrInvalidTemplateDir struct {
	defaultExitCoder
	Comment string
	Reason  error
}

func (err ErrInvalidTemplateDir) Unwrap() error { return err.Reason }

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

var (
	_ Error = ErrObjectNotFound{}
)

type ErrObjectNotFound struct {
	Name              string
	ObjectsWhichExist []string
}

func (err ErrObjectNotFound) Unwrap() error { return nil }

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

var (
	_ Error = ErrUnableToOpenObjectFile{}
)

type ErrUnableToOpenObjectFile struct {
	File   string
	Reason error
}

func (err ErrUnableToOpenObjectFile) Unwrap() error { return err.Reason }

func (err ErrUnableToOpenObjectFile) Error() string {
	return fmt.Sprintf("unable to open object file %q: %v", err.File, err.Reason)
}

// ---------------------------------------------------------------------------------------------------------------------

var (
	_ Error = ErrUnableToReadObjectFile{}
)

type ErrUnableToReadObjectFile struct {
	File   string
	Reason error
}

func (err ErrUnableToReadObjectFile) Unwrap() error { return err.Reason }

func (err ErrUnableToReadObjectFile) Error() string {
	return fmt.Sprintf("unable to read object file %q: %v", err.File, err.Reason)
}

// ---------------------------------------------------------------------------------------------------------------------

var (
	_ Error = ErrReadDir{}
)

type ErrReadDir struct {
	Dir    string
	Reason error
}

func (err ErrReadDir) Unwrap() error { return err.Reason }

func (err ErrReadDir) Error() string {
	return fmt.Sprintf("error while reading dir %q: %v", err.Dir, err.Reason)
}

// ---------------------------------------------------------------------------------------------------------------------

var (
	_ Error = ErrUnableToParseObject{}
)

type ErrUnableToParseObject struct {
	Name   string
	Reason error
}

func (err ErrUnableToParseObject) Unwrap() error { return err.Reason }

func (err ErrUnableToParseObject) Error() string {
	return fmt.Sprintf("unable to parse object %q as template: %v", err.Name, err.Reason)
}

// ---------------------------------------------------------------------------------------------------------------------

type ErrUnableToRenderObject struct {
	Name   string
	Reason error
}

func (err ErrUnableToRenderObject) Error() string {
	return fmt.Sprintf("unable to render object %q: %v", err.Name, err.Reason)
}

// ---------------------------------------------------------------------------------------------------------------------

var (
	_ Fatal = ErrUnableToCreateTempDir{}
)

type ErrUnableToCreateTempDir struct {
	defaultExitCoder
	Path   string
	Reason error
}

func (err ErrUnableToCreateTempDir) Unwrap() error { return err.Reason }

func (err ErrUnableToCreateTempDir) Error() string {
	return fmt.Sprintf("unable to create temp dir %q: %v", err.Path, err.Reason)
}

// ---------------------------------------------------------------------------------------------------------------------

var (
	_ Fatal = ErrUnableToLoadKubectlConfig{}
)

type ErrUnableToLoadKubectlConfig struct {
	defaultExitCoder
	Path   string
	Reason error
}

func (err ErrUnableToLoadKubectlConfig) Error() string {
	return fmt.Sprintf("unable to load kubectl config %q: %v", err.Path, err.Reason)
}

func (err ErrUnableToLoadKubectlConfig) Unwrap() error {
	return err.Reason
}
