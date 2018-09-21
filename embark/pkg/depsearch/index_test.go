package depsearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearcher_ResolveNameToPath(test *testing.T) {
	var searcher, err = NewSearcher("testdata/containerum")
	assert.Nil(test, err)
	test.Log(searcher.ResolveVersion("kube", "9.6."))
}
