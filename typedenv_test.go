package typedenv

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Example() {
	type config struct {
		Host     string        `env:"HOST"`
		Port     int           `env:"PORT"`
		Timeout  time.Duration `env:"TIMEOUT"`
		LogLevel slog.Level    `env:"LOG_LEVEL"`
	}

	os.Setenv("HOST", "localhost")
	os.Setenv("PORT", "8080")
	os.Setenv("TIMEOUT", "1s")
	os.Setenv("LOG_LEVEL", "debug")

	cfg, err := Load[config]()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	fmt.Printf("%#v\n", cfg)
	// Output: typedenv.config{Host:"localhost", Port:8080, Timeout:1000000000, LogLevel:-4}
}

func runDecodeValueCases(t *testing.T, tests map[string]struct {
	raw        string
	value      func() reflect.Value
	wantErr    error
	wantErrMsg string
	check      func(t *testing.T, val reflect.Value)
},
) {
	t.Helper()
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			val := tc.value()
			err := decodeValue(tc.raw, val)
			switch {
			case tc.wantErr != nil:
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("decodeValue() error = %v, want errors.Is(%v)", err, tc.wantErr)
				}
			case tc.wantErrMsg != "":
				if err == nil || !strings.Contains(err.Error(), tc.wantErrMsg) {
					t.Errorf("decodeValue() error = %v, want message containing %q", err, tc.wantErrMsg)
				}
			case err != nil:
				t.Errorf("decodeValue() unexpected error: %v", err)
			}
			if err == nil && tc.check != nil {
				tc.check(t, val)
			}
		})
	}
}

func TestDecodeValue(t *testing.T) {
	runDecodeValueCases(t, map[string]struct {
		raw        string
		value      func() reflect.Value
		wantErr    error
		wantErrMsg string
		check      func(t *testing.T, val reflect.Value)
	}{
		"non-settable value returns error": {
			raw:        "value",
			value:      func() reflect.Value { return reflect.ValueOf("immutable") },
			wantErrMsg: "field not settable",
		},
		"unsupported type returns error": {
			raw:     "value",
			value:   func() reflect.Value { var c complex64; return reflect.ValueOf(&c).Elem() },
			wantErr: ErrUnsupportedType,
		},
	})
}

func TestDecodeValue_String(t *testing.T) {
	runDecodeValueCases(t, map[string]struct {
		raw        string
		value      func() reflect.Value
		wantErr    error
		wantErrMsg string
		check      func(t *testing.T, val reflect.Value)
	}{
		"string value is set from raw": {
			raw:   "hello",
			value: func() reflect.Value { var s string; return reflect.ValueOf(&s).Elem() },
			check: func(t *testing.T, val reflect.Value) {
				if got := val.String(); got != "hello" {
					t.Errorf("got %q, want %q", got, "hello")
				}
			},
		},
	})
}

func TestDecodeValue_Bool(t *testing.T) {
	runDecodeValueCases(t, map[string]struct {
		raw        string
		value      func() reflect.Value
		wantErr    error
		wantErrMsg string
		check      func(t *testing.T, val reflect.Value)
	}{
		"bool value is set from valid raw": {
			raw:   "true",
			value: func() reflect.Value { var b bool; return reflect.ValueOf(&b).Elem() },
			check: func(t *testing.T, val reflect.Value) {
				if !val.Bool() {
					t.Error("got false, want true")
				}
			},
		},
		"bool value with invalid raw returns error": {
			raw:     "notabool",
			wantErr: ErrParse,
			value:   func() reflect.Value { var b bool; return reflect.ValueOf(&b).Elem() },
		},
	})
}

func TestDecodeValue_Int(t *testing.T) {
	runDecodeValueCases(t, map[string]struct {
		raw        string
		value      func() reflect.Value
		wantErr    error
		wantErrMsg string
		check      func(t *testing.T, val reflect.Value)
	}{
		"int value is set from raw": {
			raw:   "10",
			value: func() reflect.Value { var i int; return reflect.ValueOf(&i).Elem() },
			check: func(t *testing.T, val reflect.Value) {
				if got := val.Int(); got != 10 {
					t.Errorf("got %d, want 10", got)
				}
			},
		},
		"int8 value is set from raw": {
			raw:   "8",
			value: func() reflect.Value { var i int8; return reflect.ValueOf(&i).Elem() },
			check: func(t *testing.T, val reflect.Value) {
				if got := val.Int(); got != 8 {
					t.Errorf("got %d, want 8", got)
				}
			},
		},
		"int16 value is set from raw": {
			raw:   "16",
			value: func() reflect.Value { var i int16; return reflect.ValueOf(&i).Elem() },
			check: func(t *testing.T, val reflect.Value) {
				if got := val.Int(); got != 16 {
					t.Errorf("got %d, want 16", got)
				}
			},
		},
		"int32 value is set from raw": {
			raw:   "32",
			value: func() reflect.Value { var i int32; return reflect.ValueOf(&i).Elem() },
			check: func(t *testing.T, val reflect.Value) {
				if got := val.Int(); got != 32 {
					t.Errorf("got %d, want 32", got)
				}
			},
		},
		"int64 value is set from raw": {
			raw:   "64",
			value: func() reflect.Value { var i int64; return reflect.ValueOf(&i).Elem() },
			check: func(t *testing.T, val reflect.Value) {
				if got := val.Int(); got != 64 {
					t.Errorf("got %d, want 64", got)
				}
			},
		},
		"int value with invalid raw returns error":     {raw: "notanumber", wantErr: ErrParse, value: func() reflect.Value { var i int; return reflect.ValueOf(&i).Elem() }},
		"int8 value with overflow raw returns error":   {raw: "128", wantErr: ErrParse, value: func() reflect.Value { var i int8; return reflect.ValueOf(&i).Elem() }},
		"int8 value with underflow raw returns error":  {raw: "-129", wantErr: ErrParse, value: func() reflect.Value { var i int8; return reflect.ValueOf(&i).Elem() }},
		"int16 value with overflow raw returns error":  {raw: "32768", wantErr: ErrParse, value: func() reflect.Value { var i int16; return reflect.ValueOf(&i).Elem() }},
		"int16 value with underflow raw returns error": {raw: "-32769", wantErr: ErrParse, value: func() reflect.Value { var i int16; return reflect.ValueOf(&i).Elem() }},
		"int32 value with overflow raw returns error":  {raw: "2147483648", wantErr: ErrParse, value: func() reflect.Value { var i int32; return reflect.ValueOf(&i).Elem() }},
		"int32 value with underflow raw returns error": {raw: "-2147483649", wantErr: ErrParse, value: func() reflect.Value { var i int32; return reflect.ValueOf(&i).Elem() }},
		"int64 value with overflow raw returns error":  {raw: "9223372036854775808", wantErr: ErrParse, value: func() reflect.Value { var i int64; return reflect.ValueOf(&i).Elem() }},
		"int64 value with underflow raw returns error": {raw: "-9223372036854775809", wantErr: ErrParse, value: func() reflect.Value { var i int64; return reflect.ValueOf(&i).Elem() }},
	})
}

func TestDecodeValue_Uint(t *testing.T) {
	runDecodeValueCases(t, map[string]struct {
		raw        string
		value      func() reflect.Value
		wantErr    error
		wantErrMsg string
		check      func(t *testing.T, val reflect.Value)
	}{
		"uint value is set from raw": {
			raw:   "10",
			value: func() reflect.Value { var u uint; return reflect.ValueOf(&u).Elem() },
			check: func(t *testing.T, val reflect.Value) {
				if got := val.Uint(); got != 10 {
					t.Errorf("got %d, want 10", got)
				}
			},
		},
		"uint8 value is set from raw": {
			raw:   "8",
			value: func() reflect.Value { var u uint8; return reflect.ValueOf(&u).Elem() },
			check: func(t *testing.T, val reflect.Value) {
				if got := val.Uint(); got != 8 {
					t.Errorf("got %d, want 8", got)
				}
			},
		},
		"uint16 value is set from raw": {
			raw:   "16",
			value: func() reflect.Value { var u uint16; return reflect.ValueOf(&u).Elem() },
			check: func(t *testing.T, val reflect.Value) {
				if got := val.Uint(); got != 16 {
					t.Errorf("got %d, want 16", got)
				}
			},
		},
		"uint32 value is set from raw": {
			raw:   "32",
			value: func() reflect.Value { var u uint32; return reflect.ValueOf(&u).Elem() },
			check: func(t *testing.T, val reflect.Value) {
				if got := val.Uint(); got != 32 {
					t.Errorf("got %d, want 32", got)
				}
			},
		},
		"uint64 value is set from raw": {
			raw:   "64",
			value: func() reflect.Value { var u uint64; return reflect.ValueOf(&u).Elem() },
			check: func(t *testing.T, val reflect.Value) {
				if got := val.Uint(); got != 64 {
					t.Errorf("got %d, want 64", got)
				}
			},
		},
		"uint value with invalid raw returns error":    {raw: "-1", wantErr: ErrParse, value: func() reflect.Value { var u uint; return reflect.ValueOf(&u).Elem() }},
		"uint8 value with overflow raw returns error":  {raw: "256", wantErr: ErrParse, value: func() reflect.Value { var u uint8; return reflect.ValueOf(&u).Elem() }},
		"uint16 value with overflow raw returns error": {raw: "65536", wantErr: ErrParse, value: func() reflect.Value { var u uint16; return reflect.ValueOf(&u).Elem() }},
		"uint32 value with overflow raw returns error": {raw: "4294967296", wantErr: ErrParse, value: func() reflect.Value { var u uint32; return reflect.ValueOf(&u).Elem() }},
		"uint64 value with overflow raw returns error": {raw: "18446744073709551616", wantErr: ErrParse, value: func() reflect.Value { var u uint64; return reflect.ValueOf(&u).Elem() }},
	})
}

func TestDecodeValue_Float(t *testing.T) {
	runDecodeValueCases(t, map[string]struct {
		raw        string
		value      func() reflect.Value
		wantErr    error
		wantErrMsg string
		check      func(t *testing.T, val reflect.Value)
	}{
		"float32 value is set from raw": {
			raw:   "1.5",
			value: func() reflect.Value { var f float32; return reflect.ValueOf(&f).Elem() },
			check: func(t *testing.T, val reflect.Value) {
				if got := val.Float(); got != 1.5 {
					t.Errorf("got %v, want 1.5", got)
				}
			},
		},
		"float64 value is set from raw": {
			raw:   "1.5",
			value: func() reflect.Value { var f float64; return reflect.ValueOf(&f).Elem() },
			check: func(t *testing.T, val reflect.Value) {
				if got := val.Float(); got != 1.5 {
					t.Errorf("got %v, want 1.5", got)
				}
			},
		},
		"float value with invalid raw returns error":    {raw: "notanumber", wantErr: ErrParse, value: func() reflect.Value { var f float64; return reflect.ValueOf(&f).Elem() }},
		"float32 value with overflow raw returns error": {raw: "3.5e38", wantErr: ErrParse, value: func() reflect.Value { var f float32; return reflect.ValueOf(&f).Elem() }},
		"float64 value with overflow raw returns error": {raw: "1e309", wantErr: ErrParse, value: func() reflect.Value { var f float64; return reflect.ValueOf(&f).Elem() }},
	})
}

func TestDecodeValue_Duration(t *testing.T) {
	runDecodeValueCases(t, map[string]struct {
		raw        string
		value      func() reflect.Value
		wantErr    error
		wantErrMsg string
		check      func(t *testing.T, val reflect.Value)
	}{
		"duration value is set from valid raw": {
			raw:   "1h30m",
			value: func() reflect.Value { var d time.Duration; return reflect.ValueOf(&d).Elem() },
			check: func(t *testing.T, val reflect.Value) {
				if got := time.Duration(val.Int()); got != 90*time.Minute {
					t.Errorf("got %v, want %v", got, 90*time.Minute)
				}
			},
		},
		"duration value with invalid raw returns error": {
			raw:     "notaduration",
			wantErr: ErrParse,
			value:   func() reflect.Value { var d time.Duration; return reflect.ValueOf(&d).Elem() },
		},
		"duration value with plain integer raw returns error": {
			raw:     "42",
			wantErr: ErrParse,
			value:   func() reflect.Value { var d time.Duration; return reflect.ValueOf(&d).Elem() },
		},
	})
}

func TestDecodeValue_URL(t *testing.T) {
	runDecodeValueCases(t, map[string]struct {
		raw        string
		value      func() reflect.Value
		wantErr    error
		wantErrMsg string
		check      func(t *testing.T, val reflect.Value)
	}{
		"url value is set from valid raw": {
			raw:   "https://example.com/path?q=1",
			value: func() reflect.Value { var u url.URL; return reflect.ValueOf(&u).Elem() },
			check: func(t *testing.T, val reflect.Value) {
				got := val.Interface().(url.URL)
				if got.Scheme != "https" {
					t.Errorf("got scheme %q, want %q", got.Scheme, "https")
				}
				if got.Host != "example.com" {
					t.Errorf("got host %q, want %q", got.Host, "example.com")
				}
				if got.Path != "/path" {
					t.Errorf("got path %q, want %q", got.Path, "/path")
				}
			},
		},
		"url value with invalid raw returns error": {
			raw:     "://no-scheme",
			wantErr: ErrParse,
			value:   func() reflect.Value { var u url.URL; return reflect.ValueOf(&u).Elem() },
		},
	})
}

func TestDecodeValue_NamedURLType(t *testing.T) {
	type Endpoint url.URL
	runDecodeValueCases(t, map[string]struct {
		raw        string
		value      func() reflect.Value
		wantErr    error
		wantErrMsg string
		check      func(t *testing.T, val reflect.Value)
	}{
		"named type based on url.URL is decoded from valid url": {
			raw:   "https://example.com/path?q=1",
			value: func() reflect.Value { var u Endpoint; return reflect.ValueOf(&u).Elem() },
			check: func(t *testing.T, val reflect.Value) {
				if got := val.FieldByName("Scheme").String(); got != "https" {
					t.Errorf("got scheme %q, want %q", got, "https")
				}
				if got := val.FieldByName("Host").String(); got != "example.com" {
					t.Errorf("got host %q, want %q", got, "example.com")
				}
				if got := val.FieldByName("Path").String(); got != "/path" {
					t.Errorf("got path %q, want %q", got, "/path")
				}
			},
		},
		"named type based on url.URL with invalid raw returns error": {
			raw:     "://no-scheme",
			wantErr: ErrParse,
			value:   func() reflect.Value { var u Endpoint; return reflect.ValueOf(&u).Elem() },
		},
	})
}

type customText struct {
	val string
}

func (c *customText) UnmarshalText(b []byte) error {
	c.val = string(b)

	return nil
}

type failingText struct{}

func (f *failingText) UnmarshalText(b []byte) error {
	return errors.New("unmarshal error: " + string(b))
}

func TestDecodeValue_TextUnmarshaler(t *testing.T) {
	runDecodeValueCases(t, map[string]struct {
		raw        string
		value      func() reflect.Value
		wantErr    error
		wantErrMsg string
		check      func(t *testing.T, val reflect.Value)
	}{
		"text unmarshaler is called with raw value": {
			raw:   "hello",
			value: func() reflect.Value { var c customText; return reflect.ValueOf(&c).Elem() },
			check: func(t *testing.T, val reflect.Value) {
				got := val.FieldByName("val").String()
				if got != "hello" {
					t.Errorf("got %q, want %q", got, "hello")
				}
			},
		},
		"text unmarshaler error returns ErrParse": {
			raw:     "bad",
			value:   func() reflect.Value { var f failingText; return reflect.ValueOf(&f).Elem() },
			wantErr: ErrParse,
		},
	})
}

func TestDecodeValue_ErrorsOmitRawValue(t *testing.T) {
	const secret = "s3cr3t-v@lue"

	tests := map[string]func() reflect.Value{
		"bool":          func() reflect.Value { var b bool; return reflect.ValueOf(&b).Elem() },
		"int":           func() reflect.Value { var i int; return reflect.ValueOf(&i).Elem() },
		"uint":          func() reflect.Value { var u uint; return reflect.ValueOf(&u).Elem() },
		"float64":       func() reflect.Value { var f float64; return reflect.ValueOf(&f).Elem() },
		"time.Duration": func() reflect.Value { var d time.Duration; return reflect.ValueOf(&d).Elem() },
	}

	for name, valueFn := range tests {
		t.Run(name, func(t *testing.T) {
			err := decodeValue(secret, valueFn())
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if strings.Contains(err.Error(), secret) {
				t.Errorf("error exposes raw value: %v", err)
			}
		})
	}

	t.Run("TextUnmarshaler", func(t *testing.T) {
		var f failingText
		err := decodeValue(secret, reflect.ValueOf(&f).Elem())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if strings.Contains(err.Error(), secret) {
			t.Errorf("error exposes raw value: %v", err)
		}
	})

	// url.Parse accepts almost any string as a relative URL, so we prefix with
	// "://" to force a parse failure while still embedding the secret in the raw value.
	t.Run("url.URL", func(t *testing.T) {
		var u url.URL
		err := decodeValue("://"+secret, reflect.ValueOf(&u).Elem())
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if strings.Contains(err.Error(), secret) {
			t.Errorf("error exposes raw value: %v", err)
		}
	})
}

func TestDecodeStructField(t *testing.T) {
	tests := map[string]struct {
		setup   func() (reflect.StructField, reflect.Value)
		src     source
		wantErr error
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
		"unexported field without env tag is skipped": {
			setup: func() (reflect.StructField, reflect.Value) {
				type s struct {
					name string //nolint:unused
				}
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
			wantErr: ErrUnexportedField,
		},
		"exported field with missing env key returns error": {
			setup: func() (reflect.StructField, reflect.Value) {
				type s struct {
					Name string `env:"NAME"`
				}
				v := reflect.ValueOf(&s{}).Elem()
				return reflect.TypeFor[s]().Field(0), v.Field(0)
			},
			src:     func(string) (string, bool) { return "", false },
			wantErr: ErrNotFound,
		},
		"exported field with env key is decoded": {
			setup: func() (reflect.StructField, reflect.Value) {
				type s struct {
					Name string `env:"NAME"`
				}
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
				type s struct {
					Active bool `env:"ACTIVE"`
				}
				v := reflect.ValueOf(&s{}).Elem()
				return reflect.TypeFor[s]().Field(0), v.Field(0)
			},
			src:     func(string) (string, bool) { return "notabool", true },
			wantErr: ErrParse,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tf, v := tc.setup()
			err := decodeStructField(tf, v, tc.src)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("decodeStructField() error = %v, want errors.Is(%v)", err, tc.wantErr)
				}
			} else if err != nil {
				t.Errorf("decodeStructField() unexpected error: %v", err)
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
			if !errors.Is(err, ErrNotStruct) {
				t.Errorf("got %v, want errors.Is(ErrNotStruct)", err)
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
			type config struct {
				Host string `env:"HOST"`
			}
			_, err := decodeStruct[config](func(string) (string, bool) { return "", false })
			if !errors.Is(err, ErrNotFound) {
				t.Errorf("got %v, want errors.Is(ErrNotFound)", err)
			}
		},
		"multiple field errors are joined": func(t *testing.T) {
			type config struct {
				Host string `env:"HOST"`
				Port string `env:"PORT"`
			}
			_, err := decodeStruct[config](func(string) (string, bool) { return "", false })
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("got %v, want errors.Is(ErrNotFound)", err)
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

func TestLoad(t *testing.T) {
	tests := map[string]func(t *testing.T){
		"empty struct returns zero value": func(t *testing.T) {
			type config struct{}
			got, err := Load[config]()
			if err != nil {
				t.Fatalf("got error %v, want nil", err)
			}
			if got != (config{}) {
				t.Errorf("got %v, want zero value", got)
			}
		},
		"env vars are loaded into struct fields": func(t *testing.T) {
			type config struct {
				Host string `env:"TYPEDENV_TEST_HOST"`
				Port int    `env:"TYPEDENV_TEST_PORT"`
			}
			t.Setenv("TYPEDENV_TEST_HOST", "localhost")
			t.Setenv("TYPEDENV_TEST_PORT", "8080")

			got, err := Load[config]()
			if err != nil {
				t.Fatalf("got error %v, want nil", err)
			}
			want := config{Host: "localhost", Port: 8080}
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		},
		"missing env key error is wrapped with function name": func(t *testing.T) {
			type config struct {
				Value string `env:"TYPEDENV_TEST_MISSING"`
			}
			_, err := Load[config]()
			if !errors.Is(err, ErrNotFound) {
				t.Fatalf("got %v, want errors.Is(ErrNotFound)", err)
			}
			if !strings.Contains(err.Error(), "typedenv:") {
				t.Errorf("got %q, want error containing \"typedenv:\"", err.Error())
			}
		},
	}

	for name, run := range tests {
		t.Run(name, run)
	}
}
