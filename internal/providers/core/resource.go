package core

import (
	"context"
	"fmt"
)

type ResourceContext interface {
	context.Context

	NewController(ctx context.Context, typ string) (ResourceController, error)
	GetResourceType(ctx context.Context, iri string) (typ string, err error)

	BeginCreate(ctx context.Context, typ string, spec string) (opID string, err error)
	EndCreate(ctx context.Context, opID string, iri string, err error) error

	BeginUpdate(ctx context.Context, iri string, patch string) (opID string, err error)
	EndUpdate(ctx context.Context, opID string, err error) error

	BeginDelete(ctx context.Context, iri string) (opID string, err error)
	EndDelete(ctx context.Context, opID string, err error) error
}

type Resource struct {
	ctrl ResourceController
	iri  string
}

type ResourceController interface {
	Create(ctx context.Context, spec string) (iri string, err error)
	LoadIRI(iri string) error
	Exists(ctx context.Context) (bool, error)
	Update(ctx context.Context, patch string) error
	Delete(ctx context.Context) error
}

func CreateResource(ctx ResourceContext, typ string, spec string) (res *Resource, err error) {
	ctrl, err := ctx.NewController(ctx, typ)
	if err != nil {
		return nil, fmt.Errorf("creating %q controller: %w", typ, err)
	}
	res = &Resource{
		ctrl: ctrl,
	}

	opID, beginErr := ctx.BeginCreate(ctx, typ, spec)
	if beginErr != nil {
		return nil, fmt.Errorf("beginning: %w", beginErr)
	}
	defer func() {
		if endErr := ctx.EndCreate(ctx, opID, res.iri, err); endErr != nil {
			if err == nil {
				err = fmt.Errorf("ending: %w", endErr)
			}
		}
	}()

	res.iri, err = res.ctrl.Create(ctx, spec)
	return
}

func GetResource(ctx ResourceContext, iri string) (*Resource, error) {
	typ, err := ctx.GetResourceType(ctx, iri)
	if err != nil {
		return nil, fmt.Errorf("getting resource type: %w", err)
	}

	ctrl, err := ctx.NewController(ctx, typ)
	if err != nil {
		return nil, fmt.Errorf("creating %q controller: %w", typ, err)
	}

	res := &Resource{
		iri:  iri,
		ctrl: ctrl,
	}
	return res, nil
}

func (res *Resource) Update(ctx ResourceContext, spec string) (err error) {
	opID, beginErr := ctx.BeginUpdate(ctx, res.iri, spec)
	if beginErr != nil {
		return fmt.Errorf("beginning: %w", beginErr)
	}
	defer func() {
		if endErr := ctx.EndUpdate(ctx, opID, err); endErr != nil {
			if err == nil {
				err = fmt.Errorf("ending: %w", endErr)
			}
		}
	}()

	err = res.ctrl.Update(ctx, spec)
	return
}

func (res *Resource) Delete(ctx ResourceContext) (err error) {
	opID, beginErr := ctx.BeginDelete(ctx, res.iri)
	if beginErr != nil {
		return fmt.Errorf("beginning: %w", beginErr)
	}
	defer func() {
		if endErr := ctx.EndDelete(ctx, opID, err); endErr != nil {
			if err == nil {
				err = fmt.Errorf("ending: %w", endErr)
			}
		}
	}()

	err = res.ctrl.Delete(ctx)
	return
}
