// Package typedenv is a strongly typed environment variable manager.
package typedenv

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"time"
)

// Load returns a new instance of the given struct, or an error. It fills
// public fields with operating system environment variable values. The struct
// fields must be tagged with `env`, which specifies the environment variable
// key value to use eg. `env:"APP_ENV"`.
func Load[S any]() (S, error) {
	s, err := decodeStruct[S](os.LookupEnv)
	if err != nil {
		return s, fmt.Errorf("TypedEnv.Load[](): %w", err)
	}

	return s, nil
}

type source func(key string) (string, bool)

func decodeStruct[S any](src source) (S, error) {
	var s S
	v := reflect.ValueOf(&s).Elem()
	if v.Kind() != reflect.Struct {
		return s, fmt.Errorf("decode: expected struct, got %v", v.Kind())
	}

	var errs []error
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		if err := decodeStructField(t.Field(i), v.Field(i), src); err != nil {
			errs = append(errs, err)
		}
	}

	return s, errors.Join(errs...)
}

func decodeStructField(t reflect.StructField, v reflect.Value, src source) error {
	tag, tagged := t.Tag.Lookup("env")
	if !tagged {
		return nil
	}

	if !t.IsExported() {
		return fmt.Errorf("decodeStructField: can't decode into unexported field %v", t.Name)
	}

	raw, ok := src(tag)
	if !ok {
		return fmt.Errorf("decodeStructField: no environment value for key %v", tag)
	}

	if err := decodeField(raw, v); err != nil {
		return fmt.Errorf("%s: %w", tag, err)
	}

	return nil
}

func decodeField(raw string, field reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// This is why the return error is named
			err = fmt.Errorf("decodeField: recovered from panic %v", r)
		}
	}()

	if !field.CanSet() {
		return fmt.Errorf("decodeField: field %v not settable", field.Kind())
	}

	switch field.Type() {
	case reflect.TypeFor[time.Duration]():
		d, err := time.ParseDuration(raw)
		if err != nil {
			return fmt.Errorf("decodeField: invalid %v", field.Type())
		}

		field.SetInt(int64(d))

		return nil

	case reflect.TypeFor[url.URL]():
		u, err := url.Parse(raw)
		if err != nil {
			return fmt.Errorf("decodeField: invalid %v", field.Type())
		}

		field.Set(reflect.ValueOf(*u))

		return nil
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(raw)

		return nil

	case reflect.Bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return fmt.Errorf("decodeField: invalid %v: %w", field.Type(), errors.Unwrap(err))
		}

		field.SetBool(b)

		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(raw, 10, field.Type().Bits())
		if err != nil {
			return fmt.Errorf("decodeField: invalid %v: %w", field.Type(), errors.Unwrap(err))
		}

		field.SetInt(i)

		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(raw, 10, field.Type().Bits())
		if err != nil {
			return fmt.Errorf("decodeField: invalid %v: %w", field.Type(), errors.Unwrap(err))
		}

		field.SetUint(u)

		return nil

	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(raw, field.Type().Bits())
		if err != nil {
			return fmt.Errorf("decodeField: invalid %v: %w", field.Type(), errors.Unwrap(err))
		}

		field.SetFloat(f)

		return nil
	}

	return fmt.Errorf("decodeField: unsupported type %v", field.Type())
}
