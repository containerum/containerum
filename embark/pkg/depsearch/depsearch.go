package depsearch

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/zyedidia/micro/tools/semver"
	"gopkg.in/yaml.v2"
)

type ChartIndex map[string][]string

func (index ChartIndex) ResolveChartNameToPaths(chartName string) []string {
	return append([]string{}, index[chartName]...)
}

func (index ChartIndex) Len() int {
	return len(index)
}

func (index ChartIndex) Contains(chartName string) bool {
	var _, contains = index[chartName]
	return contains
}

func (index ChartIndex) Names() []string {
	var names = make([]string, 0, index.Len())
	for name := range index {
		names = append(names, name)
	}
	return names
}

type Searcher struct {
	baseDir string
	ChartIndex
}

func (searcher Searcher) ChartNames() []string {
	return searcher.ChartNames()
}

func (searcher Searcher) Versions(name string) []string {
	var versions []string
	for _, chartPath := range searcher.ResolveNameToPath(name) {
		var manifestData, readManifestErr = ioutil.ReadFile(path.Join(chartPath, "Chart.yaml"))
		if readManifestErr != nil {
			manifestData, readManifestErr = ioutil.ReadFile(path.Join(chartPath, "Chart.yml"))
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
		return got == version
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
		var manifestData, readManifestErr = ioutil.ReadFile(path.Join(chartPath, "Chart.yaml"))
		if readManifestErr != nil {
			manifestData, readManifestErr = ioutil.ReadFile(path.Join(chartPath, "Chart.yml"))
			if readManifestErr != nil {
				continue
			}
		}
		var vers AppVersion
		if err := yaml.Unmarshal(manifestData, &vers); err != nil {
			continue
		}
		fmt.Println(vers)
		if isClaimedVersion(vers.AppVersion) {
			return chartPath, nil
		}
	}
	return "", fmt.Errorf("unable to find satisfactory version for %q(%s)", name, version)
}

func (searcher Searcher) ResolveNameToPath(chartName string) (paths []string) {
	return searcher.ChartIndex.ResolveChartNameToPaths(chartName)
}

func NewSearcher(chartDir string) (Searcher, error) {
	var searcher = Searcher{
		baseDir:    chartDir,
		ChartIndex: make(ChartIndex),
	}
	return searcher, searchCharts(chartDir, searcher.ChartIndex)
}

func searchCharts(chartDir string, index ChartIndex) error {
	var subcharts = path.Join(chartDir, "charts")
	if !DirExists(chartDir) {
		return nil
	}
	return filepath.Walk(subcharts, func(currentPath string, info os.FileInfo, err error) error {
		var dir, _ = path.Split(currentPath)
		if path.Base(dir) != "charts" {
			return nil
		}
		if info != nil && info.IsDir() {
			index[info.Name()] = append(index[info.Name()], currentPath)
		}
		return nil
	})
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
