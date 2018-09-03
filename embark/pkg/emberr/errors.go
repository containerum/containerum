package emberr

import "fmt"

// ---------------------------------------------------------------------------------------------------------------------

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

type ErrUnableToInstallTiler struct {
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

type ErrUnsupportedKubeObjectType string

func (err ErrUnsupportedKubeObjectType) Error() string {
	return fmt.Sprintf("unsupported kube object %q", string(err))
}

func (err ErrUnsupportedKubeObjectType) Unwrap() error {
	return nil
}
