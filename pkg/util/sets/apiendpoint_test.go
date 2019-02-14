/*
Copyright 2014 The Kubernetes Authors.

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

package sets

import (
	"testing"

	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

func TestAPIEndpointSet(t *testing.T) {
	s := APIEndpointSet{}
	if len(s) != 0 {
		t.Errorf("Expected len=0: %d", len(s))
	}
	s.Insert(clusterv1.APIEndpoint{Host: "a"}, clusterv1.APIEndpoint{Host: "b"})
	if len(s) != 2 {
		t.Errorf("Expected len=2: %d", len(s))
	}
	s.Insert(clusterv1.APIEndpoint{Host: "c"})
	if s.Has(clusterv1.APIEndpoint{Host: "d"}) {
		t.Errorf("Unexpected contents: %#v", s)
	}
	if !s.Has(clusterv1.APIEndpoint{Host: "a"}) {
		t.Errorf("Missing contents: %#v", s)
	}
	s.Delete(clusterv1.APIEndpoint{Host: "a"})
	if s.Has(clusterv1.APIEndpoint{Host: "a"}) {
		t.Errorf("Unexpected contents: %#v", s)
	}
}
