package builder

import "io"

type Builder struct {
	Template io.Reader
	Values io.Reader
	Output io.Writer
}

func (builder Builder) Build() error{
	
}
