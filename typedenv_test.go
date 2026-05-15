package typedenv

import (
	"reflect"
	"testing"
)

func TestDecodeField(t *testing.T) {
	tests := map[string]struct {
		raw     string
		field   func() reflect.Value
		wantErr bool
		check   func(t *testing.T, field reflect.Value)
	}{
		"non-settable field returns error": {
			raw:     "value",
			field:   func() reflect.Value { return reflect.ValueOf("immutable") },
			wantErr: true,
		},
		"string field is set from raw": {
			raw: "hello",
			field: func() reflect.Value {
				var s string
				return reflect.ValueOf(&s).Elem()
			},
			check: func(t *testing.T, field reflect.Value) {
				if got := field.String(); got != "hello" {
					t.Errorf("got %q, want %q", got, "hello")
				}
			},
		},
		"bool field is set from valid raw": {
			raw: "true",
			field: func() reflect.Value {
				var b bool
				return reflect.ValueOf(&b).Elem()
			},
			check: func(t *testing.T, field reflect.Value) {
				if !field.Bool() {
					t.Error("got false, want true")
				}
			},
		},
		"bool field with invalid raw returns error": {
			raw:     "notabool",
			wantErr: true,
			field: func() reflect.Value {
				var b bool
				return reflect.ValueOf(&b).Elem()
			},
		},
		"unhandled kind returns no error and leaves field unchanged": {
			raw: "42",
			field: func() reflect.Value {
				var i int
				return reflect.ValueOf(&i).Elem()
			},
			check: func(t *testing.T, field reflect.Value) {
				if field.Int() != 0 {
					t.Errorf("got %d, want 0", field.Int())
				}
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			field := tc.field()
			err := decodeField(tc.raw, field)
			if (err != nil) != tc.wantErr {
				t.Errorf("decodeField() error = %v, wantErr %v", err, tc.wantErr)
			}
			if err == nil && tc.check != nil {
				tc.check(t, field)
			}
		})
	}
}

func TestDecodeStructField(t *testing.T) {
	tests := map[string]struct {
		setup   func() (reflect.StructField, reflect.Value)
		src     source
		wantErr bool
		check   func(t *testing.T, v reflect.Value)
	}{
		"field without env tag is skipped": {
			setup: func() (reflect.StructField, reflect.Value) {
				type s struct{ Name string }
				v := reflect.ValueOf(&s{}).Elem()
				return reflect.TypeFor[s]().Field(0), v.Field(0)
			},
			src: nil, // must not be called
		},
		"unexported field with env tag returns error": {
			setup: func() (reflect.StructField, reflect.Value) {
				type s struct {
					name string `env:"NAME"` //nolint:unused
				}
				v := reflect.ValueOf(&s{}).Elem()
				return reflect.TypeFor[s]().Field(0), v.Field(0)
			},
			src:     func(string) (string, bool) { return "value", true },
			wantErr: true,
		},
		"exported field with missing env key returns error": {
			setup: func() (reflect.StructField, reflect.Value) {
				type s struct{ Name string `env:"NAME"` }
				v := reflect.ValueOf(&s{}).Elem()
				return reflect.TypeFor[s]().Field(0), v.Field(0)
			},
			src:     func(string) (string, bool) { return "", false },
			wantErr: true,
		},
		"exported field with env key is decoded": {
			setup: func() (reflect.StructField, reflect.Value) {
				type s struct{ Name string `env:"NAME"` }
				v := reflect.ValueOf(&s{}).Elem()
				return reflect.TypeFor[s]().Field(0), v.Field(0)
			},
			src: func(string) (string, bool) { return "hello", true },
			check: func(t *testing.T, v reflect.Value) {
				if got := v.String(); got != "hello" {
					t.Errorf("got %q, want %q", got, "hello")
				}
			},
		},
		"decode failure returns wrapped error": {
			setup: func() (reflect.StructField, reflect.Value) {
				type s struct{ Active bool `env:"ACTIVE"` }
				v := reflect.ValueOf(&s{}).Elem()
				return reflect.TypeFor[s]().Field(0), v.Field(0)
			},
			src:     func(string) (string, bool) { return "notabool", true },
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tf, v := tc.setup()
			err := decodeStructField(tf, v, tc.src)
			if (err != nil) != tc.wantErr {
				t.Errorf("decodeStructField() error = %v, wantErr %v", err, tc.wantErr)
			}
			if err == nil && tc.check != nil {
				tc.check(t, v)
			}
		})
	}
}

func TestDecodeStruct(t *testing.T) {
	tests := map[string]func(t *testing.T){
		"non-struct type returns error": func(t *testing.T) {
			_, err := decodeStruct[string](func(string) (string, bool) { return "", false })
			if err == nil {
				t.Error("got nil, want error")
			}
		},
		"empty struct returns zero value": func(t *testing.T) {
			type config struct{}
			got, err := decodeStruct[config](func(string) (string, bool) { return "", false })
			if err != nil {
				t.Fatalf("got error %v, want nil", err)
			}
			if got != (config{}) {
				t.Errorf("got %v, want zero value", got)
			}
		},
		"struct with untagged fields returns zero value": func(t *testing.T) {
			type config struct{ Name string }
			got, err := decodeStruct[config](func(string) (string, bool) { return "", false })
			if err != nil {
				t.Fatalf("got error %v, want nil", err)
			}
			if got != (config{}) {
				t.Errorf("got %v, want zero value", got)
			}
		},
		"all tagged fields decoded successfully": func(t *testing.T) {
			type config struct {
				Host  string `env:"HOST"`
				Debug bool   `env:"DEBUG"`
			}
			vals := map[string]string{"HOST": "localhost", "DEBUG": "true"}
			src := func(key string) (string, bool) { v, ok := vals[key]; return v, ok }

			got, err := decodeStruct[config](src)
			if err != nil {
				t.Fatalf("got error %v, want nil", err)
			}
			want := config{Host: "localhost", Debug: true}
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		},
		"missing env key returns error": func(t *testing.T) {
			type config struct{ Host string `env:"HOST"` }
			_, err := decodeStruct[config](func(string) (string, bool) { return "", false })
			if err == nil {
				t.Error("got nil, want error")
			}
		},
		"multiple field errors are joined": func(t *testing.T) {
			type config struct {
				Host string `env:"HOST"`
				Port string `env:"PORT"`
			}
			_, err := decodeStruct[config](func(string) (string, bool) { return "", false })
			if err == nil {
				t.Fatal("got nil, want error")
			}
			joined, ok := err.(interface{ Unwrap() []error })
			if !ok {
				t.Fatalf("expected joined errors, got %T", err)
			}
			if got := len(joined.Unwrap()); got != 2 {
				t.Errorf("got %d errors, want 2", got)
			}
		},
	}

	for name, run := range tests {
		t.Run(name, run)
	}
}
