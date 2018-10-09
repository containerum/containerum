package renderer

import (
	"context"
	"fmt"
	"testing"

	"github.com/containerum/containerum/embark/pkg/kube"
	"github.com/containerum/containerum/embark/pkg/utils/why"
	"github.com/stretchr/testify/assert"
)

func TestRenderedComponent_ForEachObjectGo(test *testing.T) {
	var kinds = append([]string{
		"anything",
		"weird object",
	}, DefaultOrder()...)
	var objects = mocks(kinds)
	var component = NewRenderedObject("test component", nil, objects...)
	var results = make(chan string)

	go func() {
		if err := component.ForEachObjectGo(func(_ context.Context, obj kube.Object) error {
			results <- obj.Kind()
			return nil
		}); err != nil {
			test.Fatal(err)
		}
		close(results)
	}()
	var aggregated = aggregate(results)
	assert.ElementsMatch(test, aggregated, kinds)
	if len(aggregated) == 0 {
		test.Fatalf("there can't be empty aggregated result!")
	}
	if testing.Verbose() {
		why.PrintFromIter("aggregated results", len(aggregated), func(i int) (string, error) {
			var kind = aggregated[i]
			return fmt.Sprintf("%-16s %2d", kind, ObjectPriority(kind)), nil
		})
	}
	var prevPriority = ObjectPriority(aggregated[0])
	var prevKind = aggregated[0]
	for _, kind := range aggregated {
		var priority = ObjectPriority(kind)
		if prevPriority > priority {
			test.Fatalf("unconsistent object ordering: %q(%d) followed by %q(%d)", kind, priority, prevKind, prevPriority)
		}
		prevPriority = priority
		prevKind = kind
	}
}

func aggregate(results <-chan string) []string {
	var sl []string
	for result := range results {
		sl = append(sl, result)
	}
	return sl
}
