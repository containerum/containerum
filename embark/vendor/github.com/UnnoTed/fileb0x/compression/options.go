package compression

// Options for compression
type Options struct {
	// activates the compression
	// default: false
	Compress bool

	// valid values are:
	// -> "NoCompression"
	// -> "BestSpeed"
	// -> "BestCompression"
	// -> "DefaultCompression"
	//
	// default: "DefaultCompression" // when: Compress == true && Method == ""
	Method string

	// true = do it yourself (the file is written as gzip into the memory file system)
	// false = decompress at run time (while writing file into memory file system)
	// default: false
	Keep bool
}
