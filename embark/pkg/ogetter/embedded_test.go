package ogetter

import (
	"bytes"
	"testing"
)

func TestStaticTemplatesObjectGetter(test *testing.T) {
	testStaticObjectGetter(test, StaticTemplates)
}

func TestStaticMetaObjectGetter(test *testing.T) {
	testStaticObjectGetter(test, StaticMeta)
}

func testStaticObjectGetter(test *testing.T, getter EmbeddedFSObjectGetter) {
	for _, objectName := range getter.ObjectNames() {
		var buf = &bytes.Buffer{}
		if err := getter.Object(objectName, buf); err != nil {
			test.Fatal(err)
		}
		test.Logf("%s:  %q %d bytes",
			objectName,
			func() string {
				var p, _ = getter.ObjectFilePath(objectName)
				return p
			}(),
			buf.Len())
	}
}
