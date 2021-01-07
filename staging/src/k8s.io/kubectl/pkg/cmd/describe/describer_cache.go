package describe

import (
	"github.com/davecgh/go-spew/spew"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/kubectl/pkg/describe"
)

// describerCache can map from a meta.RESTMapping to a describe.ResourceDescriber.
//
// As of right now, this does not need to be synced since it is function-local, but
// it certainly could be updated to be thread-safe in the future.
type describerCache struct {
	cache map[string]describe.ResourceDescriber
}

func newDescriberCache() *describerCache {
	return &describerCache{cache: make(map[string]describe.ResourceDescriber)}
}

func (d *describerCache) key(mapping *meta.RESTMapping) string {
	key := struct {
		gvr   schema.GroupVersionResource
		gvk   schema.GroupVersionKind
		scope string
	}{
		gvr:   mapping.Resource,
		gvk:   mapping.GroupVersionKind,
		scope: string(mapping.Scope.Name()),
	}
	return (&spew.ConfigState{DisableMethods: true, Indent: " "}).Sprint(key)
}

func (d *describerCache) get(mapping *meta.RESTMapping) describe.ResourceDescriber {
	return d.cache[d.key(mapping)]
}

func (d *describerCache) put(mapping *meta.RESTMapping, describer describe.ResourceDescriber) {
	d.cache[d.key(mapping)] = describer
}
