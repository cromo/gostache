package main

import (
	"reflect"
	"testing"
)

func TestRender(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{`---
names: [ {name: chris}, {name: mark}, {name: scott} ]
---
{{#names}}
  Hi {{name}}!
{{/names}}`, []string{`  Hi chris!
  Hi mark!
  Hi scott!
`}}, {
			`---
name: chris
---
name: mark
---
name: scott
---
Hi {{name}}!
`, []string{"Hi chris!\n", "Hi mark!\n", "Hi scott!\n"}},
	}
	for _, c := range cases {
		template, data := splitTemplateAndData(c.in)
		got := renderTemplateWithData(template, data)
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("rendered %q, got %q, want %q", c.in, got, c.want)
		}
	}
}
