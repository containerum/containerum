package depsearch

type ComponentIndex map[string][]string

func (index ComponentIndex) ResolveChartNameToPaths(componentName string) []string {
	return append([]string{}, index[componentName]...)
}

func (index ComponentIndex) Len() int {
	return len(index)
}

func (index ComponentIndex) Contains(chartName string) bool {
	var _, contains = index[chartName]
	return contains
}

func (index ComponentIndex) Names() []string {
	var names = make([]string, 0, index.Len())
	for name := range index {
		names = append(names, name)
	}
	return names
}
