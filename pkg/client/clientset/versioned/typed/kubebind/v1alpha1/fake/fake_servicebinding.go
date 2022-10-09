/*
Copyright The Kubectl Bind contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"

	v1alpha1 "github.com/kube-bind/kube-bind/pkg/apis/kubebind/v1alpha1"
)

// FakeServiceBindings implements ServiceBindingInterface
type FakeServiceBindings struct {
	Fake *FakeKubeBindV1alpha1
}

var servicebindingsResource = schema.GroupVersionResource{Group: "kube-bind.io", Version: "v1alpha1", Resource: "servicebindings"}

var servicebindingsKind = schema.GroupVersionKind{Group: "kube-bind.io", Version: "v1alpha1", Kind: "ServiceBinding"}

// Get takes name of the serviceBinding, and returns the corresponding serviceBinding object, and an error if there is any.
func (c *FakeServiceBindings) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ServiceBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(servicebindingsResource, name), &v1alpha1.ServiceBinding{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ServiceBinding), err
}

// List takes label and field selectors, and returns the list of ServiceBindings that match those selectors.
func (c *FakeServiceBindings) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ServiceBindingList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(servicebindingsResource, servicebindingsKind, opts), &v1alpha1.ServiceBindingList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ServiceBindingList{ListMeta: obj.(*v1alpha1.ServiceBindingList).ListMeta}
	for _, item := range obj.(*v1alpha1.ServiceBindingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested serviceBindings.
func (c *FakeServiceBindings) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(servicebindingsResource, opts))
}

// Create takes the representation of a serviceBinding and creates it.  Returns the server's representation of the serviceBinding, and an error, if there is any.
func (c *FakeServiceBindings) Create(ctx context.Context, serviceBinding *v1alpha1.ServiceBinding, opts v1.CreateOptions) (result *v1alpha1.ServiceBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(servicebindingsResource, serviceBinding), &v1alpha1.ServiceBinding{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ServiceBinding), err
}

// Update takes the representation of a serviceBinding and updates it. Returns the server's representation of the serviceBinding, and an error, if there is any.
func (c *FakeServiceBindings) Update(ctx context.Context, serviceBinding *v1alpha1.ServiceBinding, opts v1.UpdateOptions) (result *v1alpha1.ServiceBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(servicebindingsResource, serviceBinding), &v1alpha1.ServiceBinding{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ServiceBinding), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeServiceBindings) UpdateStatus(ctx context.Context, serviceBinding *v1alpha1.ServiceBinding, opts v1.UpdateOptions) (*v1alpha1.ServiceBinding, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(servicebindingsResource, "status", serviceBinding), &v1alpha1.ServiceBinding{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ServiceBinding), err
}

// Delete takes name of the serviceBinding and deletes it. Returns an error if one occurs.
func (c *FakeServiceBindings) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(servicebindingsResource, name, opts), &v1alpha1.ServiceBinding{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeServiceBindings) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(servicebindingsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ServiceBindingList{})
	return err
}

// Patch applies the patch and returns the patched serviceBinding.
func (c *FakeServiceBindings) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ServiceBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(servicebindingsResource, name, pt, data, subresources...), &v1alpha1.ServiceBinding{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ServiceBinding), err
}