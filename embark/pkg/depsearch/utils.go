package depsearch

import (
	"os"
	"path"
	"path/filepath"

	"github.com/containerum/containerum/embark/pkg/static"
)

type IndexBuilder func() (ComponentIndex, error)

func buildIndexFromComponentPaths(componentPaths []string) ComponentIndex {
	var index = make(ComponentIndex)
	for _, componentPath := range componentPaths {
		var componentName = path.Base(componentPath)
		index[componentName] = append(index[componentName], componentPath)
	}
	return index
}

func FindAllComponentDirFS(root string) IndexBuilder {
	var componentPaths []string
	var err = filepath.Walk(root, func(filePath string, info os.FileInfo, err error) error {
		if path.Base(filePath) != "templates" || err != nil || info == nil {
			return err
		}
		var componentPath = path.Base(filePath)
		componentPaths = append(componentPaths, componentPath)
		return nil
	})
	return func() (ComponentIndex, error) {
		return buildIndexFromComponentPaths(componentPaths), err
	}
}

func FindAllComponentsDirStatic() IndexBuilder {
	var errFunc = func(err error) IndexBuilder {
		return func() (ComponentIndex, error) {
			return nil, err
		}
	}
	var filePaths, walkDirsErr = static.WalkDirs("", true)
	if walkDirsErr != nil {
		return errFunc(walkDirsErr)
	}
	var componentPaths = make([]string, 0, len(filePaths)/2)
	for _, filePath := range filePaths {
		var finfo, statErr = static.FS.Stat(static.CTX, filePath)
		if statErr != nil {
			return errFunc(statErr)
		}
		if !finfo.IsDir() || path.Base(filePath) != "templates" {
			continue
		}
		var componentPath = path.Dir(filePath)
		componentPaths = append(componentPaths, componentPath)
	}
	return func() (ComponentIndex, error) {
		return buildIndexFromComponentPaths(componentPaths), nil
	}
}

type AppVersion struct {
	AppVersion string `yaml:"appVersion"`
}

func DirExists(dirName string) bool {
	var dir, err = os.Stat(dirName)
	if err != nil {
		return false
	}
	return dir.IsDir()
}

func FileExists(filePath string) bool {
	var file, err = os.Stat(filePath)
	if err != nil {
		return false
	}
	return !file.IsDir()
}
