package ogetter

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/containerum/containerum/embark/pkg/emberr"
	"github.com/containerum/containerum/embark/pkg/static"
)

const StaticRootDir = ""

var (
	StaticTemplates = NewEmbeddedFSObjectGetter(path.Join(StaticRootDir, "templates"))
	StaticMeta      = NewEmbeddedFSObjectGetter(StaticRootDir)

	_ ObjectGetter = EmbeddedFSObjectGetter{}
)

// KV object store over static embedded filesystem
// Ignores files from subdirs
type EmbeddedFSObjectGetter struct {
	fnames map[string]string
}

// Panics if dir doesn't exist in static filesystem
// Ignored file defined in utils.go##Ignored()
// If
func NewEmbeddedFSObjectGetter(dir string) EmbeddedFSObjectGetter {
	var fnames, getAllFiles = static.WalkDirs(dir, false)
	if getAllFiles != nil {
		panic(fmt.Sprintf("[containerum/embark/pkg/ogetter.NewEmbeddedFSObjectGetter] unable to get static filenames: %v", getAllFiles))
	}
	var getter = EmbeddedFSObjectGetter{
		fnames: make(map[string]string, len(fnames)/3+2),
	}
	for _, fname := range fnames {
		{
			var fdir, _ = path.Split(fname)
			var fileInSubDir = dir == "" && fdir != ""
			var fileInAnotherSubDir = dir != "" && !strings.HasPrefix(fdir, dir)
			if fileInSubDir || fileInAnotherSubDir {
				continue
			}
		}
		var objectName = extractObjectNameFromFilename(fname)
		if objectName == "" || IsIgnored(fname) {
			continue
		}
		getter.fnames[objectName] = fname
	}
	return getter
}

// Returns object names
func (getter EmbeddedFSObjectGetter) ObjectNames() []string {
	var names = make([]string, 0, len(getter.fnames))
	for name := range getter.fnames {
		names = append(names, name)
	}
	return names
}

func (getter EmbeddedFSObjectGetter) ObjectFilePath(name string) (string, error) {
	var objectFilePath, objectExits = getter.fnames[name]
	if objectExits {
		return objectFilePath, nil
	}
	return "", emberr.ErrObjectNotFound{
		Name:              name,
		ObjectsWhichExist: getter.ObjectNames(),
	}
}

func (getter EmbeddedFSObjectGetter) Object(name string, output io.Writer) error {
	var fname, objectExists = getter.fnames[name]
	if objectExists {
		var data, readDataErr = static.ReadFile(fname)
		if readDataErr != nil {
			return emberr.ErrUnableToReadObjectFile{
				File:   fname,
				Reason: readDataErr,
			}
		}
		var _, writeErr = output.Write(data)
		return writeErr
	}
	return emberr.ErrObjectNotFound{
		Name:              name,
		ObjectsWhichExist: getter.ObjectNames(),
	}
}
