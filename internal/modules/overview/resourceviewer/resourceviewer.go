/*
Copyright (c) 2019 VMware, Inc. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package resourceviewer

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/vmware/octant/internal/componentcache"
	"github.com/vmware/octant/internal/config"
	"github.com/vmware/octant/internal/modules/overview/objectvisitor"
	"github.com/vmware/octant/internal/queryer"
	"github.com/vmware/octant/pkg/store"
	"github.com/vmware/octant/pkg/view/component"
)

// ViewerOpt is an option for ResourceViewer.
type ViewerOpt func(*ResourceViewer) error

// WithDefaultQueryer configures ResourceViewer with the default visitor.
func WithDefaultQueryer(dashConfig config.Dash, q queryer.Queryer) ViewerOpt {
	return func(rv *ResourceViewer) error {
		visitor, err := objectvisitor.NewDefaultVisitor(dashConfig, q)
		if err != nil {
			return err
		}

		rv.visitor = visitor
		return nil
	}
}

// ResourceViewer visits an object and creates a view component.
type ResourceViewer struct {
	dashConfig config.Dash
	visitor    objectvisitor.Visitor
}

// New creates an instance of ResourceViewer.
func New(dashConfig config.Dash, opts ...ViewerOpt) (*ResourceViewer, error) {
	rv := &ResourceViewer{
		dashConfig: dashConfig,
	}

	for _, opt := range opts {
		if err := opt(rv); err != nil {
			return nil, errors.Wrap(err, "invalid resource viewer option")
		}
	}

	if rv.visitor == nil {
		return nil, errors.New("resource viewer visitor is nil")
	}

	return rv, nil
}

// Visit visits an object and creates a view component.

func (rv *ResourceViewer) Visit(ctx context.Context, object runtime.Object) (*component.ResourceViewer, error) {
	ctx, span := trace.StartSpan(ctx, "resourceViewer")
	defer span.End()

	handler, err := NewHandler(rv.dashConfig)
	if err != nil {
		return nil, errors.Wrap(err, "Create handler")
	}

	newCtx := context.Background()

	if err := rv.visitor.Visit(newCtx, object, handler); err != nil {
		return nil, err
	}

	accessor := meta.NewAccessor()
	uid, err := accessor.UID(object)
	if err != nil {
		return nil, err
	}

	return GenerateComponent(ctx, handler, uid)
}

// CachedResourceViewer returns a RV component from the component cache and starts a new visit.
func CachedResourceViewer(object runtime.Object, dashConfig config.Dash, q queryer.Queryer) componentcache.UpdateFn {
	return func(ctx context.Context, cacheChan chan componentcache.Event) (string, error) {
		var event componentcache.Event
		event.Name = "Resource Viewer"

		copyObject := object.DeepCopyObject()

		key, err := store.KeyFromObject(copyObject)
		if err != nil {
			return "", err
		}
		sKey := fmt.Sprintf("%s-%s", "resourceviewer", key.String())
		event.Key = sKey

		componentCache := dashConfig.ComponentCache()
		if _, ok := componentCache.Get(sKey); !ok {
			title := component.Title(component.NewText("Resource Viewer"))
			loading := component.NewLoading(title, "Resource Viewer")
			componentCache.Add(sKey, loading)
		}

		rv, err := New(dashConfig, WithDefaultQueryer(dashConfig, q))
		if err != nil {
			return sKey, err
		}

		go func() {
			c, err := rv.Visit(ctx, copyObject)
			event.Err = errors.WithMessage(err, "visiting object failed")
			event.Component = c
			cacheChan <- event
		}()

		return sKey, nil
	}
}
