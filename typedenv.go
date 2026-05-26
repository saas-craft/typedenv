// Package typedenv decodes OS environment variables into a struct
package typedenv

import (
	"encoding"
	"errors"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ErrNotStruct       = errors.New("expected struct")
	ErrUnexportedField = errors.New("unexported field")
	ErrNotFound        = errors.New("variable not found for key")
	ErrParse           = errors.New("invalid value")
	ErrUnsupportedType = errors.New("unsupported type")
	ErrInvalidDefault  = errors.New("invalid default value")
)

// Load reads operating system environment variables into a new instance of S,
// which must be a struct. Exported fields tagged with `env:"KEY"` are populated
// by looking up KEY in the environment; untagged fields are left at their zero
// value.
//
// Supported field types: string, bool, the int and uint families, the float
// family, time.Duration, and url.URL.
//
// Load returns an error if a tagged variable is missing from the environment,
// fails to parse, or targets an unexported field. Errors from multiple fields
// are joined.
func Load[S any]() (S, error) {
	s, err := decodeStruct[S](os.LookupEnv)
	if err != nil {
		return s, fmt.Errorf("typedenv: %w", err)
	}

	return s, nil
}

type source func(key string) (string, bool)

func decodeStruct[S any](lookup source) (S, error) {
	var s S
	structVal := reflect.ValueOf(&s).Elem()
	if structVal.Kind() != reflect.Struct {
		return s, fmt.Errorf("%w: got %v", ErrNotStruct, structVal.Kind())
	}

	var errs []error
	structType := structVal.Type()
	for i := range structType.NumField() {
		if err := decodeStructField(structType.Field(i), structVal.Field(i), lookup); err != nil {
			errs = append(errs, err)
		}
	}

	return s, errors.Join(errs...)
}

type fieldSpec struct {
	key        string
	defaultVal string
	hasDefault bool
}

func decodeStructField(field reflect.StructField, val reflect.Value, lookup source) error {
	tag, tagged := field.Tag.Lookup("env")
	if !tagged {
		return nil
	}

	if !field.IsExported() {
		return fmt.Errorf("%q: %w", field.Name, ErrUnexportedField)
	}

	spec := parseEnvTag(tag)
	raw, isDefault, err := resolveRaw(spec, lookup)
	if err != nil {
		return fmt.Errorf("%q: %w", spec.key, err)
	}

	if err := decodeValue(raw, val); err != nil {
		if isDefault {
			return fmt.Errorf("%q: %w: %v", spec.key, ErrInvalidDefault, err)
		}
		return fmt.Errorf("%q: %w", spec.key, err)
	}

	return nil
}

func parseEnvTag(tag string) fieldSpec {
	key, rest, hasOptions := strings.Cut(tag, ",")
	if !hasOptions {
		return fieldSpec{key: tag}
	}

	const prefix = "default="
	if strings.HasPrefix(rest, prefix) {
		return fieldSpec{key: key, defaultVal: rest[len(prefix):], hasDefault: true}
	}

	return fieldSpec{key: key}
}

func resolveRaw(spec fieldSpec, lookup source) (raw string, isDefault bool, err error) {
	if v, ok := lookup(spec.key); ok {
		return v, false, nil
	}

	if spec.hasDefault {
		return spec.defaultVal, true, nil
	}

	return "", false, ErrNotFound
}

func decodeValue(raw string, dest reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("internal panic: %v", r)
		}
	}()

	if !dest.CanSet() {
		return errors.New("field not settable")
	}

	if dest.Kind() == reflect.Pointer {
		elem := reflect.New(dest.Type().Elem())
		if err := decodeValue(raw, elem.Elem()); err != nil {
			return err
		}

		dest.Set(elem)

		return nil
	}

	if dest.CanAddr() {
		if u, ok := dest.Addr().Interface().(encoding.TextUnmarshaler); ok {
			if err := u.UnmarshalText([]byte(raw)); err != nil {
				return fmt.Errorf("%w: text unmarshaling failed", ErrParse)
			}

			return nil
		}
	}

	if dest.Type() == reflect.TypeFor[time.Duration]() {
		d, err := time.ParseDuration(raw)
		if err != nil {
			return fmt.Errorf("%w: invalid duration", ErrParse)
		}

		dest.SetInt(int64(d))

		return nil
	}

	if dest.Type().ConvertibleTo(reflect.TypeFor[url.URL]()) {
		u, err := url.Parse(raw)
		if err != nil {
			return fmt.Errorf("%w: invalid url", ErrParse)
		}

		dest.Set(reflect.ValueOf(*u).Convert(dest.Type()))

		return nil
	}

	switch dest.Kind() {
	case reflect.String:
		dest.SetString(raw)

		return nil

	case reflect.Bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return fmt.Errorf("%w: invalid bool", ErrParse)
		}

		dest.SetBool(b)

		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(raw, 10, dest.Type().Bits())
		if err != nil {
			return numericErr(err, dest.Type())
		}

		dest.SetInt(i)

		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(raw, 10, dest.Type().Bits())
		if err != nil {
			return numericErr(err, dest.Type())
		}

		dest.SetUint(u)

		return nil

	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(raw, dest.Type().Bits())
		if err != nil {
			return numericErr(err, dest.Type())
		}

		dest.SetFloat(f)

		return nil
	}

	return fmt.Errorf("%w: %v", ErrUnsupportedType, dest.Type())
}

func numericErr(err error, typ reflect.Type) error {
	if errors.Is(err, strconv.ErrRange) {
		return fmt.Errorf("%w: %v out of range", ErrParse, typ)
	}

	return fmt.Errorf("%w: invalid %v", ErrParse, typ)
}
