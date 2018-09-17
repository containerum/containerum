package file

// File holds file's data
type File struct {
	OriginalPath string
	Name         string
	Path         string
	Data         string
	Bytes        []byte
	ReplacedText bool
	Tags         string
	Base         string
	Prefix       string
	Modified     string
}

// NewFile creates a new File
func NewFile() *File {
	f := new(File)
	return f
}
