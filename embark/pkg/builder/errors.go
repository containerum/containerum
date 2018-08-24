package builder

import "fmt"

type ErrUnableToFetchChart struct {
	Chart  string
	Reason error
}

func (err ErrUnableToFetchChart) Error() string {
	return fmt.Sprintf("unable to fetch chart %q: %v", err.Chart, err.Reason)
}
