package envstruct

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	indexEnvVar int = iota
	indexRequired
	indexNoReport
)

// Load will use the `env` tags from a struct to populate the structs values and
// perform validations.
func Load(t interface{}) error {
	val := reflect.ValueOf(t).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		tagProperties := extractSliceInputs(tag.Get("env"))
		envVar := strings.ToUpper(tagProperties[indexEnvVar])
		envVal := os.Getenv(envVar)

		var required bool
		if len(tagProperties) >= 2 {
			required = tagProperties[indexRequired] == "required"
		}

		if isInvalid(envVal, required) {
			return fmt.Errorf("%s is required but was empty", envVar)
		}

		if envVal == "" {
			continue
		}

		err := setField(valueField, envVal)
		if err != nil {
			return err
		}
	}

	return nil
}

func setField(value reflect.Value, input string) error {
	switch value.Type() {
	case reflect.TypeOf(time.Second):
		return setDuration(value, input)
	case reflect.TypeOf(&url.URL{}):
		return setURL(value, input)
	}

	switch value.Kind() {
	case reflect.String:
		return setString(value, input)
	case reflect.Bool:
		return setBool(value, input)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return setInt(value, input)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return setUint(value, input)
	case reflect.Slice:
		return setSlice(value, input)
	}

	return nil
}

func extractSliceInputs(input string) []string {
	inputs := strings.Split(input, ",")

	for i, v := range inputs {
		inputs[i] = strings.TrimSpace(v)
	}

	return inputs
}

func isInvalid(input string, required bool) bool {
	return required && input == ""
}

func setDuration(value reflect.Value, input string) error {
	d, err := time.ParseDuration(input)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(d))

	return nil
}

func setURL(value reflect.Value, input string) error {
	u, err := url.Parse(input)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(u))

	return nil
}

func setString(value reflect.Value, input string) error {
	value.SetString(input)

	return nil
}

func setBool(value reflect.Value, input string) error {
	value.SetBool(input == "true" || input == "1")

	return nil
}

func setInt(value reflect.Value, input string) error {
	n, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return err
	}

	value.SetInt(int64(n))

	return nil
}

func setUint(value reflect.Value, input string) error {
	n, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return err
	}

	value.SetUint(uint64(n))

	return nil
}

func setSlice(value reflect.Value, input string) error {
	inputs := extractSliceInputs(input)

	rs := reflect.MakeSlice(value.Type(), len(inputs), len(inputs))
	for i, val := range inputs {
		err := setField(rs.Index(i), val)
		if err != nil {
			return err
		}
	}

	value.Set(rs)

	return nil
}
