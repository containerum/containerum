package renderer

import (
	"context"

	"github.com/containerum/containerum/embark/pkg/kube"
	"golang.org/x/sync/errgroup"
)

type RenderedComponent struct {
	name    string
	objects []kube.Object
}

func NewRenderedObject(name string, objects ...kube.Object) RenderedComponent {
	SortObjects(objects)
	return RenderedComponent{
		name:    name,
		objects: objects,
	}
}

func (component RenderedComponent) Objects() []kube.Object {
	return append([]kube.Object{}, component.objects...)
}

func (component RenderedComponent) Name() string {
	return component.name
}

type Op = func(obj kube.Object) error

func (component RenderedComponent) ForEachObject(op Op) error {
	for _, obj := range component.objects {
		if err := op(obj); err != nil {
			return err
		}
	}
	return nil
}

type AsyncOp = func(ctx context.Context, obj kube.Object) error

func (component RenderedComponent) ForEachObjectGo(op AsyncOp) error {
	for _, batch := range ObjectsToBatches(component.objects) {
		var group, ctx = errgroup.WithContext(context.Background())
		for _, obj := range batch {
			obj := obj
			group.Go(func() error {
				return op(ctx, obj)
			})
		}
		if err := group.Wait(); err != nil {
			return err
		}
	}
	return nil
}
