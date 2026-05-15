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
			raw: "value",
			field: func() reflect.Value {
				var c complex64
				return reflect.ValueOf(&c).Elem()
			},
		},
		"int field is set from raw": {
			raw: "10",
			field: func() reflect.Value {
				var i int
				return reflect.ValueOf(&i).Elem()
			},
			check: func(t *testing.T, field reflect.Value) {
				if got := field.Int(); got != 10 {
					t.Errorf("got %d, want 10", got)
				}
			},
		},
		"int8 field is set from raw": {
			raw: "8",
			field: func() reflect.Value {
				var i int8
				return reflect.ValueOf(&i).Elem()
			},
			check: func(t *testing.T, field reflect.Value) {
				if got := field.Int(); got != 8 {
					t.Errorf("got %d, want 8", got)
				}
			},
		},
		"int16 field is set from raw": {
			raw: "16",
			field: func() reflect.Value {
				var i int16
				return reflect.ValueOf(&i).Elem()
			},
			check: func(t *testing.T, field reflect.Value) {
				if got := field.Int(); got != 16 {
					t.Errorf("got %d, want 16", got)
				}
			},
		},
		"int32 field is set from raw": {
			raw: "32",
			field: func() reflect.Value {
				var i int32
				return reflect.ValueOf(&i).Elem()
			},
			check: func(t *testing.T, field reflect.Value) {
				if got := field.Int(); got != 32 {
					t.Errorf("got %d, want 32", got)
				}
			},
		},
		"int64 field is set from raw": {
			raw: "64",
			field: func() reflect.Value {
				var i int64
				return reflect.ValueOf(&i).Elem()
			},
			check: func(t *testing.T, field reflect.Value) {
				if got := field.Int(); got != 64 {
					t.Errorf("got %d, want 64", got)
				}
			},
		},
		"int field with invalid raw returns error": {
			raw:     "notanumber",
			wantErr: true,
			field: func() reflect.Value {
				var i int
				return reflect.ValueOf(&i).Elem()
			},
		},
		"uint field is set from raw": {
			raw: "10",
			field: func() reflect.Value {
				var u uint
				return reflect.ValueOf(&u).Elem()
			},
			check: func(t *testing.T, field reflect.Value) {
				if got := field.Uint(); got != 10 {
					t.Errorf("got %d, want 10", got)
				}
			},
		},
		"uint8 field is set from raw": {
			raw: "8",
			field: func() reflect.Value {
				var u uint8
				return reflect.ValueOf(&u).Elem()
			},
			check: func(t *testing.T, field reflect.Value) {
				if got := field.Uint(); got != 8 {
					t.Errorf("got %d, want 8", got)
				}
			},
		},
		"uint16 field is set from raw": {
			raw: "16",
			field: func() reflect.Value {
				var u uint16
				return reflect.ValueOf(&u).Elem()
			},
			check: func(t *testing.T, field reflect.Value) {
				if got := field.Uint(); got != 16 {
					t.Errorf("got %d, want 16", got)
				}
			},
		},
		"uint32 field is set from raw": {
			raw: "32",
			field: func() reflect.Value {
				var u uint32
				return reflect.ValueOf(&u).Elem()
			},
			check: func(t *testing.T, field reflect.Value) {
				if got := field.Uint(); got != 32 {
					t.Errorf("got %d, want 32", got)
				}
			},
		},
		"uint64 field is set from raw": {
			raw: "64",
			field: func() reflect.Value {
				var u uint64
				return reflect.ValueOf(&u).Elem()
			},
			check: func(t *testing.T, field reflect.Value) {
				if got := field.Uint(); got != 64 {
					t.Errorf("got %d, want 64", got)
				}
			},
		},
		"uint field with invalid raw returns error": {
			raw:     "-1",
			wantErr: true,
			field: func() reflect.Value {
				var u uint
				return reflect.ValueOf(&u).Elem()
			},
		},
		"float32 field is set from raw": {
			raw: "1.5",
			field: func() reflect.Value {
				var f float32
				return reflect.ValueOf(&f).Elem()
			},
			check: func(t *testing.T, field reflect.Value) {
				if got := field.Float(); got != 1.5 {
					t.Errorf("got %v, want 1.5", got)
				}
			},
		},
		"float64 field is set from raw": {
			raw: "1.5",
			field: func() reflect.Value {
				var f float64
				return reflect.ValueOf(&f).Elem()
			},
			check: func(t *testing.T, field reflect.Value) {
				if got := field.Float(); got != 1.5 {
					t.Errorf("got %v, want 1.5", got)
				}
			},
		},
		"float field with invalid raw returns error": {
			raw:     "notanumber",
			wantErr: true,
			field: func() reflect.Value {
				var f float64
				return reflect.ValueOf(&f).Elem()
			},
		},
		"int8 field with overflow raw returns error": {
			raw:     "128",
			wantErr: true,
			field:   func() reflect.Value { var i int8; return reflect.ValueOf(&i).Elem() },
		},
		"int8 field with underflow raw returns error": {
			raw:     "-129",
			wantErr: true,
			field:   func() reflect.Value { var i int8; return reflect.ValueOf(&i).Elem() },
		},
		"int16 field with overflow raw returns error": {
			raw:     "32768",
			wantErr: true,
			field:   func() reflect.Value { var i int16; return reflect.ValueOf(&i).Elem() },
		},
		"int16 field with underflow raw returns error": {
			raw:     "-32769",
			wantErr: true,
			field:   func() reflect.Value { var i int16; return reflect.ValueOf(&i).Elem() },
		},
		"int32 field with overflow raw returns error": {
			raw:     "2147483648",
			wantErr: true,
			field:   func() reflect.Value { var i int32; return reflect.ValueOf(&i).Elem() },
		},
		"int32 field with underflow raw returns error": {
			raw:     "-2147483649",
			wantErr: true,
			field:   func() reflect.Value { var i int32; return reflect.ValueOf(&i).Elem() },
		},
		"int64 field with overflow raw returns error": {
			raw:     "9223372036854775808",
			wantErr: true,
			field:   func() reflect.Value { var i int64; return reflect.ValueOf(&i).Elem() },
		},
		"int64 field with underflow raw returns error": {
			raw:     "-9223372036854775809",
			wantErr: true,
			field:   func() reflect.Value { var i int64; return reflect.ValueOf(&i).Elem() },
		},
		"uint8 field with overflow raw returns error": {
			raw:     "256",
			wantErr: true,
			field:   func() reflect.Value { var u uint8; return reflect.ValueOf(&u).Elem() },
		},
		"uint16 field with overflow raw returns error": {
			raw:     "65536",
			wantErr: true,
			field:   func() reflect.Value { var u uint16; return reflect.ValueOf(&u).Elem() },
		},
		"uint32 field with overflow raw returns error": {
			raw:     "4294967296",
			wantErr: true,
			field:   func() reflect.Value { var u uint32; return reflect.ValueOf(&u).Elem() },
		},
		"uint64 field with overflow raw returns error": {
			raw:     "18446744073709551616",
			wantErr: true,
			field:   func() reflect.Value { var u uint64; return reflect.ValueOf(&u).Elem() },
		},
		"float32 field with overflow raw returns error": {
			raw:     "3.5e38",
			wantErr: true,
			field:   func() reflect.Value { var f float32; return reflect.ValueOf(&f).Elem() },
		},
		"float64 field with overflow raw returns error": {
			raw:     "1e309",
			wantErr: true,
			field:   func() reflect.Value { var f float64; return reflect.ValueOf(&f).Elem() },
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
