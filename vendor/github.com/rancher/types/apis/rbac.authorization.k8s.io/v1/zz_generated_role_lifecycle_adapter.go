package v1

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type RoleLifecycle interface {
	Create(obj *v1.Role) (*v1.Role, error)
	Remove(obj *v1.Role) (*v1.Role, error)
	Updated(obj *v1.Role) (*v1.Role, error)
}

type roleLifecycleAdapter struct {
	lifecycle RoleLifecycle
}

func (w *roleLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*v1.Role))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *roleLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*v1.Role))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *roleLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*v1.Role))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewRoleLifecycleAdapter(name string, client RoleInterface, l RoleLifecycle) RoleHandlerFunc {
	adapter := &roleLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, adapter, client.ObjectClient())
	return func(key string, obj *v1.Role) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
