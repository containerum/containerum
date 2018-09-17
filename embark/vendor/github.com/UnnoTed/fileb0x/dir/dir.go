package dir

import "strings"

// Dir holds directory information to insert into templates
type Dir struct {
	List      [][]string
	Blacklist []string
}

// Exists checks if a directory exists or not
func (d *Dir) Exists(newDir string) bool {
	for _, dir := range d.Blacklist {
		if dir == newDir {
			return true
		}
	}

	return false
}

// Parse a directory to build a list of directories to be made at b0x.go
func (d *Dir) Parse(newDir string) []string {
	list := strings.Split(newDir, "/")

	var dirWalk []string

	for indx := range list {
		dirList := ""
		for i := -1; i < indx; i++ {
			dirList += list[i+1] + "/"
		}

		if !d.Exists(dirList) {
			if strings.HasSuffix(dirList, "//") {
				dirList = dirList[:len(dirList)-1]
			}

			dirWalk = append(dirWalk, dirList)
			d.Blacklist = append(d.Blacklist, dirList)
		}
	}

	return dirWalk
}

// Insert a new folder to the list
func (d *Dir) Insert(newDir string) {
	if !d.Exists(newDir) {
		d.Blacklist = append(d.Blacklist, newDir)
		d.List = append(d.List, d.Parse(newDir))
	}
}

// Clean dupes
func (d *Dir) Clean() []string {
	var cleanList []string

	for _, dirs := range d.List {
		for _, dir := range dirs {
			if dir == "./" || dir == "/" || dir == "." || dir == "" {
				continue
			}

			cleanList = append(cleanList, dir)
		}
	}

	return cleanList
}
