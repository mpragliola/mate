package mate

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	f := NewFrontMatterFromSource("")
	returnedType := fmt.Sprintf("%T", f)

	if returnedType != "*mate.FrontMatter" {
		t.Fatalf("wrong type \"%s\" returned", returnedType)
	}
}

func TestCreation(t *testing.T) {
	type scenario struct {
		source     string
		wantValues map[string]string
		wantBody   string
	}

	scenarios := map[string]scenario{
		"parses without front matter": {
			source:     `some text`,
			wantValues: nil,
			wantBody:   `some text`,
		},
		"parses with front matter": {
			source: `---
			foo: bar
			---
			body`,
			wantValues: map[string]string{
				"foo": "bar",
			},
			wantBody: `body`,
		},
		"does not parse incorrect front matter": {
			source: `
			---
			foo: bar
			---
			body`,
			wantValues: nil,
			wantBody: `
			---
			foo: bar
			---
			body`,
		},
	}

	for name, scenario := range scenarios {
		t.Run(
			name, func(t *testing.T) {
				fm := NewFrontMatterFromSource(scenario.source)

				if fm.GetBody() != scenario.wantBody {
					t.Fatalf("wrong body")
				}

				if len(fm.GetValues()) != len(scenario.wantValues) {
					t.Fatalf("wrong values")
				}

				for key, value := range scenario.wantValues {
					if val, ok := fm.GetValues()[key]; !ok || val != value {
						t.Fatalf("wrong value with key %s", key)
					}
				}
			},
		)
	}
}
