package typedenv

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

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
		if tagged {
			return fmt.Errorf("decodeStructField: can't decode into unexported field %v", t.Name)
		}

		return nil
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

	switch field.Kind() {
	case reflect.String:
		field.SetString(raw)

	case reflect.Bool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return err
		}

		field.SetBool(b)
	}

	return nil
}
