/*
Copyright 2017 The Kubernetes Authors.

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

package explain

import (
	"bytes"
	"testing"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestModel(t *testing.T) {
	schema := resources.LookupResource(schema.GroupVersionKind{
		Group:   "",
		Version: "v1",
		Kind:    "OneKind",
	})
	if schema == nil {
		t.Fatal("Couldn't find schema v1.OneKind")
	}

	tests := []struct {
		path []string
		want string
	}{
		{
			want: `DESCRIPTION:
     OneKind has a short description

FIELDS:
   field1	<Object> -required-
     Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla ut lacus ac
     enim vulputate imperdiet ac accumsan risus. Integer vel accumsan lectus.
     Praesent tempus nulla id tortor luctus, quis varius nulla laoreet. Ut orci
     nisi, suscipit id velit sed, blandit eleifend turpis. Curabitur tempus ante
     at lectus viverra, a mattis augue euismod. Morbi quam ligula, porttitor sit
     amet lacus non, interdum pulvinar tortor. Praesent accumsan risus et ipsum
     dictum, vel ullamcorper lorem egestas.

   field2	<[]map[string]string>
     This is an array of object of PrimitiveDef

`,
			path: []string{},
		},
		{
			want: `RESOURCE: field1 <Object>

DESCRIPTION:
     Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla ut lacus ac
     enim vulputate imperdiet ac accumsan risus. Integer vel accumsan lectus.
     Praesent tempus nulla id tortor luctus, quis varius nulla laoreet. Ut orci
     nisi, suscipit id velit sed, blandit eleifend turpis. Curabitur tempus ante
     at lectus viverra, a mattis augue euismod. Morbi quam ligula, porttitor sit
     amet lacus non, interdum pulvinar tortor. Praesent accumsan risus et ipsum
     dictum, vel ullamcorper lorem egestas.

     This is another kind of Kind

FIELDS:
   array	<[]integer>
     This array must be an array of int

   int	<integer>
     This int must be an int

   object	<map[string]string>
     This is an object of string

   primitive	<string>

   string	<string> -required-
     This string must be a string

`,
			path: []string{"field1"},
		},
		{
			want: `FIELD: string <string>

DESCRIPTION:
     This string must be a string
`,
			path: []string{"field1", "string"},
		},
		{
			want: `FIELD: array <[]integer>

DESCRIPTION:
     This array must be an array of int

     This is an int in an array
`,
			path: []string{"field1", "array"},
		},
	}

	for _, test := range tests {
		buf := bytes.Buffer{}
		if err := PrintModelDescription(test.path, &buf, schema, false); err != nil {
			t.Fatalf("Failed to PrintModelDescription for path %v: %v", test.path, err)
		}
		got := buf.String()
		if got != test.want {
			t.Errorf("Got:\n%v\nWant:\n%v\n", buf.String(), test.want)
		}
	}
}
