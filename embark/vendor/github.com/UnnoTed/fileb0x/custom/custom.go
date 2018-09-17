package custom

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/UnnoTed/fileb0x/compression"
	"github.com/UnnoTed/fileb0x/dir"
	"github.com/UnnoTed/fileb0x/file"
	"github.com/UnnoTed/fileb0x/updater"
	"github.com/UnnoTed/fileb0x/utils"
	"github.com/bmatcuk/doublestar"
	"github.com/karrick/godirwalk"
)

const hextable = "0123456789abcdef"

// SharedConfig holds needed data from config package
// without causing import cycle
type SharedConfig struct {
	Output      string
	Compression *compression.Gzip
	Updater     updater.Config
}

// Custom is a set of files with dedicaTed customization
type Custom struct {
	Files  []string
	Base   string
	Prefix string
	Tags   string

	Exclude []string
	Replace []Replacer
}

var (
	xx    = []byte(`\x`)
	start = []byte(`[]byte("`)
)

const lowerhex = "0123456789abcdef"

// Parse the files transforming them into a byte string and inserting the file
// into a map of files
func (c *Custom) Parse(files *map[string]*file.File, dirs **dir.Dir, config *SharedConfig) error {
	to := *files
	dirList := *dirs

	var newList []string
	for _, customFile := range c.Files {
		// get files from glob
		list, err := doublestar.Glob(customFile)
		if err != nil {
			return err
		}

		// insert files from glob into the new list
		newList = append(newList, list...)
	}

	// copy new list
	c.Files = newList

	// 0 files in the list
	if len(c.Files) == 0 {
		return errors.New("No files found")
	}

	// loop through files from glob
	for _, customFile := range c.Files {
		// gives error when file doesn't exist
		if !utils.Exists(customFile) {
			return fmt.Errorf("File [%s] doesn't exist", customFile)
		}

		cb := func(fpath string, d *godirwalk.Dirent) error {
			if config.Updater.Empty && !config.Updater.IsUpdating {
				log.Println("empty mode")
				return nil
			}

			// only files will be processed
			if d != nil && d.IsDir() {
				return nil
			}

			originalPath := fpath
			fpath = utils.FixPath(fpath)

			var fixedPath string
			if c.Prefix != "" || c.Base != "" {
				c.Base = strings.TrimPrefix(c.Base, "./")

				if strings.HasPrefix(fpath, c.Base) {
					fixedPath = c.Prefix + fpath[len(c.Base):]
				} else {
					if c.Base != "" {
						fixedPath = c.Prefix + fpath
					}
				}

				fixedPath = utils.FixPath(fixedPath)
			} else {
				fixedPath = utils.FixPath(fpath)
			}

			// check for excluded files
			for _, excludedFile := range c.Exclude {
				m, err := doublestar.Match(c.Prefix+excludedFile, fixedPath)
				if err != nil {
					return err
				}

				if m {
					return nil
				}
			}

			info, err := os.Stat(fpath)
			if err != nil {
				return err
			}

			if info.Name() == config.Output {
				return nil
			}

			// get file's content
			content, err := ioutil.ReadFile(fpath)
			if err != nil {
				return err
			}

			replaced := false

			// loop through replace list
			for _, r := range c.Replace {
				// check if path matches the pattern from property: file
				matched, err := doublestar.Match(c.Prefix+r.File, fixedPath)
				if err != nil {
					return err
				}

				if matched {
					for pattern, word := range r.Replace {
						content = []byte(strings.Replace(string(content), pattern, word, -1))
						replaced = true
					}
				}
			}

			// compress the content
			if config.Compression.Options != nil {
				content, err = config.Compression.Compress(content)
				if err != nil {
					return err
				}
			}

			dst := make([]byte, len(content)*4)
			for i := 0; i < len(content); i++ {
				dst[i*4] = byte('\\')
				dst[i*4+1] = byte('x')
				dst[i*4+2] = hextable[content[i]>>4]
				dst[i*4+3] = hextable[content[i]&0x0f]
			}

			f := file.NewFile()
			f.OriginalPath = originalPath
			f.ReplacedText = replaced
			f.Data = `[]byte("` + string(dst) + `")`
			f.Name = info.Name()
			f.Path = fixedPath
			f.Tags = c.Tags
			f.Base = c.Base
			f.Prefix = c.Prefix
			f.Modified = info.ModTime().String()

			if _, ok := to[fixedPath]; ok {
				f.Tags = to[fixedPath].Tags
			}

			// insert dir to dirlist so it can be created on b0x's init()
			dirList.Insert(path.Dir(fixedPath))

			// insert file into file list
			to[fixedPath] = f
			return nil
		}

		customFile = utils.FixPath(customFile)

		// unlike filepath.walk, godirwalk will only walk dirs
		f, err := os.Open(customFile)
		if err != nil {
			return err
		}

		defer f.Close()

		fs, err := f.Stat()
		if err != nil {
			return err
		}

		if fs.IsDir() {
			if err := godirwalk.Walk(customFile, &godirwalk.Options{
				Unsorted: true,
				Callback: cb,
			}); err != nil {
				return err
			}

		} else {
			if err := cb(customFile, nil); err != nil {
				return err
			}
		}
	}

	return nil
}
