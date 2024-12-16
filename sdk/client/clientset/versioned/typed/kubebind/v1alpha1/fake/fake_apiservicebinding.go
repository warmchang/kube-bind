/*
Copyright 2024 The Kube Bind Authors.

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
	gentype "k8s.io/client-go/gentype"

	v1alpha1 "github.com/kube-bind/kube-bind/sdk/apis/kubebind/v1alpha1"
	kubebindv1alpha1 "github.com/kube-bind/kube-bind/sdk/client/clientset/versioned/typed/kubebind/v1alpha1"
)

// fakeAPIServiceBindings implements APIServiceBindingInterface
type fakeAPIServiceBindings struct {
	*gentype.FakeClientWithList[*v1alpha1.APIServiceBinding, *v1alpha1.APIServiceBindingList]
	Fake *FakeKubeBindV1alpha1
}

func newFakeAPIServiceBindings(fake *FakeKubeBindV1alpha1) kubebindv1alpha1.APIServiceBindingInterface {
	return &fakeAPIServiceBindings{
		gentype.NewFakeClientWithList[*v1alpha1.APIServiceBinding, *v1alpha1.APIServiceBindingList](
			fake.Fake,
			"",
			v1alpha1.SchemeGroupVersion.WithResource("apiservicebindings"),
			v1alpha1.SchemeGroupVersion.WithKind("APIServiceBinding"),
			func() *v1alpha1.APIServiceBinding { return &v1alpha1.APIServiceBinding{} },
			func() *v1alpha1.APIServiceBindingList { return &v1alpha1.APIServiceBindingList{} },
			func(dst, src *v1alpha1.APIServiceBindingList) { dst.ListMeta = src.ListMeta },
			func(list *v1alpha1.APIServiceBindingList) []*v1alpha1.APIServiceBinding {
				return gentype.ToPointerSlice(list.Items)
			},
			func(list *v1alpha1.APIServiceBindingList, items []*v1alpha1.APIServiceBinding) {
				list.Items = gentype.FromPointerSlice(items)
			},
		),
		fake,
	}
}