package installer

import "github.com/containerum/containerum/embark/pkg/ogetter"

type Installer struct {
	MetaProvider             ogetter.ObjectGetter
	ComponentObjectsProvider ogetter.ObjectGetter
	TempDir                  string
}
