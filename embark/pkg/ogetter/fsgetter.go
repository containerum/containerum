package ogetter

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/containerum/containerum/embark/pkg/emberr"
)

var (
	_ ObjectGetter = FSObjectGetter{}
)

// Loads file data from directory by name
// Example:
// dir/
//	|- deploy.yaml
//  |- svc.yml
// NewFSObjectGetter("dir").Object("svc", os.Stdout) will print "svc.yml" content to stdout
// Basically FSObjectGetter works as simple kv block storage
type FSObjectGetter struct {
	dir string

	cached nameToPath
}

func NewFSObjectGetter(dir string) FSObjectGetter {
	var expandedDir = os.ExpandEnv(dir)
	return FSObjectGetter{
		dir: path.Clean(expandedDir),
	}
}

func (getter FSObjectGetter) ObjectNames() []string {
	var o, readDirErr = getter.readDir()
	if readDirErr != nil {
		panic(readDirErr)
	}
	getter.cached = o
	return o.Names()
}

func (getter FSObjectGetter) templatesDir() string {
	return getter.dir
}

func (getter FSObjectGetter) readDir() (nameToPath, error) {
	if getter.cached != nil {
		return getter.cached.Copy(), nil
	}
	var objects = make(nameToPath)
	var readDirErr = filepath.Walk(
		getter.templatesDir(),
		func(objectPath string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			if IsIgnored(info.Name()) {
				return nil
			}
			var name = extractObjectNameFromFilename(info.Name())
			if name == "" || IsIgnored(info.Name()) {
				return nil
			}
			objects[name] = objectPath
			return nil
		})
	if readDirErr != nil {
		return nil, emberr.ErrReadDir{
			Dir:    getter.dir,
			Reason: readDirErr,
		}
	}
	getter.cached = objects.Copy()
	return objects, nil
}

func (getter FSObjectGetter) Object(name string, output io.Writer) error {
	var objects, getObjectsErr = getter.readDir()
	if getObjectsErr != nil {
		return getObjectsErr
	}
	var objectPath, exists = objects[name]
	if !exists {
		return emberr.ErrObjectNotFound{
			Name:              name,
			ObjectsWhichExist: objects.Names(),
		}
	}
	var objectFile, objectFileOpenErr = os.Open(path.Join(objectPath))
	if objectFileOpenErr != nil {
		return emberr.ErrUnableToOpenObjectFile{
			File:   objectPath,
			Reason: objectFileOpenErr,
		}
	}
	defer objectFile.Close()
	var _, writeObjectErr = io.Copy(output, objectFile)
	if writeObjectErr != nil {
		return emberr.ErrUnableToReadObjectFile{
			File:   objectPath,
			Reason: writeObjectErr,
		}
	}
	return writeObjectErr
}

type nameToPath map[string]string

func (names2path nameToPath) Names() []string {
	var names = make([]string, 0, len(names2path))
	for name := range names2path {
		names = append(names, name)
	}
	return names
}

func (names2path nameToPath) Copy() nameToPath {
	var cp = make(nameToPath, len(names2path))
	for k, v := range names2path {
		cp[k] = v
	}
	return cp
}
