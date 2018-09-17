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

var StaticTemplates = NewEmbeddedFSObjectGetter(path.Join(StaticRootDir, "templates"))
var StaticMeta = NewEmbeddedFSObjectGetter(StaticRootDir)

type EmbeddedFSObjectGetter struct {
	fnames map[string]string
}

func NewEmbeddedFSObjectGetter(dir string) EmbeddedFSObjectGetter {
	var getter = EmbeddedFSObjectGetter{
		fnames: make(map[string]string),
	}
	var fnames, getAllFiles = static.WalkDirs(dir, false)
	if getAllFiles != nil {
		panic(fmt.Sprintf("[containerum/embark/pkg/ogetter.NewEmbeddedFSObjectGetter] unable to get static filenames: %v", getAllFiles))
	}
	for _, fname := range fnames {
		var objectName string
		if fdir, _ := path.Split(fname); fdir != dir {
			continue
		}
		{
			var ext = path.Ext(fname)
			objectName = strings.TrimPrefix(fname, dir+"/")
			objectName = strings.TrimSuffix(objectName, ext)
		}
		if objectName == "" || IsIgnored(fname) {
			continue
		}
		getter.fnames[objectName] = fname
	}
	return getter
}

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
