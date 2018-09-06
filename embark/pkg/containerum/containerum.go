package containerum

type Containerum struct {
	Component
	Components []Component
}

type Component struct {
	Name      string
	Version   string
	Resources []string
	Values    map[string]interface{}
}
