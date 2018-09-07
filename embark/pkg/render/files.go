package render

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/containerum/containerum/embark/pkg/emberr"
)

type FileSet map[string]os.FileInfo

func (fileset FileSet) Contains(name string) bool {
	var _, contains = fileset[name]
	return contains
}

func (fileset FileSet) Names() []string {
	var names = make([]string, 0, len(fileset))
	for name := range fileset {
		names = append(names, name)
	}
	return names
}

func FileSetFromDir(dir string) (FileSet, error) {
	var info, lstatErr = os.Lstat(dir)
	if lstatErr != nil {
		return nil, emberr.ErrInvalidTemplateDir{
			Reason: lstatErr,
		}
	}
	var fileSet = make(FileSet)
	if !info.IsDir() {
		return nil, emberr.ErrInvalidTemplateDir{
			Reason: fmt.Errorf("expect dir, got %q(%v)", info.Name(), info.Mode()),
		}
	}
	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			var name = info.Name()
			var nameWithoutExt = FileNameWithoutExt(name)
			fileSet[name] = info
			fileSet[nameWithoutExt] = info
		}
		return nil
	})
	return fileSet, nil
}
