/*
Copyright 2021 The Kube Bind Authors.

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

package bootstrap

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"io"
	"strings"
	"text/template"
	"time"

	extensionsapiserver "k8s.io/apiextensions-apiserver/pkg/apiserver"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	apimachineryerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"
	"k8s.io/klog/v2"
)

// TransformFileFunc transforms a resource file before being applied to the cluster.
type TransformFileFunc func(bs []byte) ([]byte, error)

// Option allows to customize the bootstrap process.
type Option struct {
	// TransformFileFunc is a function that transforms a resource file before being applied to the cluster.
	TransformFile TransformFileFunc
}

// ReplaceOption allows to customize the bootstrap process.
func ReplaceOption(pairs ...string) Option {
	return Option{
		TransformFile: func(bs []byte) ([]byte, error) {
			if len(pairs)%2 != 0 {
				return nil, fmt.Errorf("odd number of arguments: %v", pairs)
			}
			for i := 0; i < len(pairs); i += 2 {
				bs = bytes.ReplaceAll(bs, []byte(pairs[i]), []byte(pairs[i+1]))
			}
			return bs, nil
		},
	}
}

// Bootstrap creates resources in a package's fs by
// continuously retrying the list. This is blocking, i.e. it only returns (with error)
// when the context is closed or with nil when the bootstrapping is successfully completed.
func Bootstrap(ctx context.Context, discoveryClient discovery.DiscoveryInterface, dynamicClient dynamic.Interface, batteriesIncluded sets.Set[string], fs embed.FS, opts ...Option) error {
	cache := memory.NewMemCacheClient(discoveryClient)
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(cache)

	// bootstrap non-crd resources
	transformers := make([]TransformFileFunc, 0, len(opts))
	for _, opt := range opts {
		transformers = append(transformers, opt.TransformFile)
	}
	return wait.PollUntilContextCancel(ctx, time.Second, true, func(ctx context.Context) (bool, error) {
		if err := CreateResourcesFromFS(ctx, dynamicClient, mapper, batteriesIncluded, fs, transformers...); err != nil {
			klog.V(1).Infof("Failed to bootstrap resources, retrying: %v", err)
			// invalidate cache if resources not found
			// xref: https://github.com/kcp-dev/kcp/issues/655
			cache.Invalidate()
			return false, nil
		}
		return true, nil
	})
}

// CreateResourcesFromFS creates all resources from a filesystem.
func CreateResourcesFromFS(ctx context.Context, client dynamic.Interface, mapper meta.RESTMapper, batteriesIncluded sets.Set[string], fs embed.FS, transformers ...TransformFileFunc) error {
	files, err := fs.ReadDir(".")
	if err != nil {
		return err
	}

	var errs []error
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if err := CreateResourceFromFS(ctx, client, mapper, batteriesIncluded, f.Name(), fs, transformers...); err != nil {
			errs = append(errs, err)
		}
	}
	return apimachineryerrors.NewAggregate(errs)
}

// CreateResourceFromFS creates given resource file.
func CreateResourceFromFS(ctx context.Context, client dynamic.Interface, mapper meta.RESTMapper, batteriesIncluded sets.Set[string], filename string, fs embed.FS, transformers ...TransformFileFunc) error {
	raw, err := fs.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("could not read %s: %w", filename, err)
	}

	if len(raw) == 0 {
		return nil // ignore empty files
	}

	d := kubeyaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(raw)))
	var errs []error
	for i := 1; ; i++ {
		doc, err := d.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}
		if len(bytes.TrimSpace(doc)) == 0 {
			continue
		}

		for _, transformer := range transformers {
			doc, err = transformer(doc)
			if err != nil {
				return err
			}
		}

		if err := createResourceFromFS(ctx, client, mapper, doc, batteriesIncluded); err != nil {
			errs = append(errs, fmt.Errorf("failed to create resource %s doc %d: %w", filename, i, err))
		}
	}
	return apimachineryerrors.NewAggregate(errs)
}

const annotationCreateOnlyKey = "bootstrap.kube-bind.io/create-only"
const annotationBattery = "bootstrap.kube-bind.io/battery"

func createResourceFromFS(ctx context.Context, client dynamic.Interface, mapper meta.RESTMapper, raw []byte, batteriesIncluded sets.Set[string]) error {
	type Input struct {
		Batteries map[string]bool
	}
	input := Input{
		Batteries: map[string]bool{},
	}
	for b := range batteriesIncluded {
		input.Batteries[b] = true
	}
	tmpl, err := template.New("manifest").Parse(string(raw))
	if err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, input); err != nil {
		return fmt.Errorf("failed to execute manifest: %w", err)
	}

	obj, gvk, err := extensionsapiserver.Codecs.UniversalDeserializer().Decode(buf.Bytes(), nil, &unstructured.Unstructured{})
	if err != nil {
		return fmt.Errorf("could not decode raw: %w", err)
	}
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("decoded into incorrect type, got %T, wanted %T", obj, &unstructured.Unstructured{})
	}

	if v, found := u.GetAnnotations()[annotationBattery]; found {
		partOf := strings.Split(v, ",")
		included := false
		for _, p := range partOf {
			if batteriesIncluded.Has(strings.TrimSpace(p)) {
				included = true
				break
			}
		}
		if !included {
			klog.V(4).Infof("Skipping %s because %s is/are not among included batteries %s", u.GetName(), v, batteriesIncluded)
			return nil
		}
	}

	m, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return fmt.Errorf("could not get REST mapping for %s: %w", gvk, err)
	}

	upserted, err := client.Resource(m.Resource).Namespace(u.GetNamespace()).Create(ctx, u, metav1.CreateOptions{})
	if err != nil {
		if apierrors.IsAlreadyExists(err) {
			existing, err := client.Resource(m.Resource).Namespace(u.GetNamespace()).Get(ctx, u.GetName(), metav1.GetOptions{})
			if err != nil {
				return err
			}

			if _, exists := existing.GetAnnotations()[annotationCreateOnlyKey]; exists {
				klog.V(3).Infof(
					"Skipping update of %s %s because it has the create-only annotation",
					gvk,
					qualifiedObjectName(existing),
				)

				return nil
			}

			u.SetResourceVersion(existing.GetResourceVersion())
			if _, err = client.Resource(m.Resource).Namespace(u.GetNamespace()).Update(ctx, u, metav1.UpdateOptions{}); err != nil {
				return fmt.Errorf("could not update %s %s: %w", gvk.Kind, qualifiedObjectName(existing), err)
			} else {
				klog.V(1).Infof("Updated %s %s", gvk, qualifiedObjectName(existing))
				return nil
			}
		}
		return err
	}

	klog.V(1).Infof("Bootstrapped %s %s", gvk.Kind, qualifiedObjectName(upserted))

	return nil
}

func qualifiedObjectName(obj metav1.Object) string {
	if len(obj.GetNamespace()) > 0 {
		return fmt.Sprintf("%s/%s", obj.GetNamespace(), obj.GetName())
	}
	return obj.GetName()
}
