package v1alpha3

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"sigs.k8s.io/yaml"
	"testing"
)

func TestTemplateParameter_Required(t *testing.T) {
	type fields struct {
		paramYAML string
	}
	tests := []struct {
		name         string
		fields       fields
		wantRequired bool
	}{{
		name: "Parameter should be required if no default set",
		fields: fields{
			paramYAML: `
name: "fake-name"
`,
		},
		wantRequired: true,
	}, {
		name: "Parameter should not be required if default is empty string",
		fields: fields{
			paramYAML: `
name: "fake-name"
default: ""`,
		},
		wantRequired: false,
	}, {
		name: "Parameter should not be required if default is 0",
		fields: fields{
			paramYAML: `
name: "fake-name"
default: 0`,
		},
		wantRequired: false,
	}, {
		name: "Parameter should not be required if default is false",
		fields: fields{
			paramYAML: `
name: "fake-name"
default: false`,
		},
		wantRequired: false,
	}, {
		name: "Parameter should be required if default is nil",
		fields: fields{
			paramYAML: `
name: "fake-name"
default: null`,
		},
		wantRequired: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param := &TemplateParameter{}
			err := yaml.Unmarshal([]byte(tt.fields.paramYAML), param)
			assert.Nil(t, err)

			paramBytes, err := json.Marshal(param)
			assert.Nil(t, err)
			log.Println(string(paramBytes))

			assert.Equalf(t, tt.wantRequired, param.Required(), "Required()")
		})
	}
}
