package depsearch

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/containerum/containerum/embark/pkg/static"

	"github.com/zyedidia/micro/tools/semver"
	"gopkg.in/yaml.v2"
)

type FileReader func(filePath string) (data []byte, err error)

type Searcher struct {
	ComponentIndex

	baseDir        string
	fileReader     FileReader
	versionedCache map[versionedComponent][]string
}

type versionedComponent struct {
	Name    string
	Version string
}

func FS(chartDir string) (Searcher, error) {
	var index, buildIndexErr = FindAllComponentDirFS(chartDir)()
	if buildIndexErr != nil {
		return Searcher{}, buildIndexErr
	}
	return Searcher{
		baseDir:        chartDir,
		fileReader:     ioutil.ReadFile,
		ComponentIndex: index,
	}, nil
}

func Static() Searcher {
	var index, buildIndexErr = FindAllComponentsDirStatic()()
	if buildIndexErr != nil {
		panic(buildIndexErr)
	}
	return Searcher{
		baseDir:        "static FS",
		fileReader:     static.ReadFile,
		ComponentIndex: index,
	}
}

func (searcher Searcher) Versions(name string) []string {
	var versions []string
	for _, chartPath := range searcher.ResolveNameToPath(name) {
		var manifestData, readManifestErr = searcher.fileReader(path.Join(chartPath, "Chart.yaml"))
		if readManifestErr != nil {
			manifestData, readManifestErr = searcher.fileReader(path.Join(chartPath, "Chart.yml"))
			if readManifestErr != nil {
				continue
			}
		}
		var vers AppVersion
		if err := yaml.Unmarshal(manifestData, &vers); err != nil {
			continue
		}
		versions = append(versions, vers.AppVersion)
	}
	return versions
}

func (searcher Searcher) ResolveVersion(name, version string) (chartPath string, err error) {
	var paths = searcher.ResolveNameToPath(name)
	if len(paths) == 0 {
		return "", fmt.Errorf("chart %q not found in %q", name, searcher.baseDir)
	}
	var isClaimedVersion = func(got string) bool {
		return got == version || version == ""
	}
	var claimed, parseUserVersionErr = semver.ParseTolerant(version)
	if parseUserVersionErr == nil {
		isClaimedVersion = func(got string) bool {
			if got == version {
				return true
			}
			var gotVers, parseGotVersErr = semver.Parse(got)
			if parseGotVersErr != nil {
				return false
			}
			return gotVers.Major == claimed.Major && gotVers.Minor <= claimed.Minor
		}
	}
	for _, chartPath := range paths {
		var manifestData, readManifestErr = searcher.fileReader(path.Join(chartPath, "Chart.yaml"))
		if readManifestErr != nil {
			manifestData, readManifestErr = searcher.fileReader(path.Join(chartPath, "Chart.yml"))
			if readManifestErr != nil {
				continue
			}
		}
		var vers AppVersion
		if err := yaml.Unmarshal(manifestData, &vers); err != nil {
			continue
		}
		if isClaimedVersion(vers.AppVersion) {
			return chartPath, nil
		}
	}
	return "", fmt.Errorf("unable to find satisfactory version for %q(%s)", name, version)
}

func (searcher Searcher) ResolveNameToPath(chartName string) (paths []string) {
	return searcher.ComponentIndex.ResolveChartNameToPaths(chartName)
}
