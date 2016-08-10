package envstruct

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
)

var ReportWriter io.Writer = os.Stdout

// WriteReport will take a struct that is setup for envstruct and print
// out a report containing the struct field name, field type, environment
// variable for that field, whether or not the field is required and
// the value of that field. The report is written to `ReportWriter`
// which defaults to `os.StdOut`
func WriteReport(t interface{}) error {
	w := tabwriter.NewWriter(ReportWriter, 0, 8, 2, ' ', 0)

	fmt.Fprintln(w, "FIELD NAME:\tTYPE:\tENV:\tREQUIRED:\tVALUE:")

	val := reflect.ValueOf(t).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag

		tagProperties := extractSliceInputs(tag.Get("env"))
		envVar := strings.ToUpper(tagProperties[indexEnvVar])

		var isRequired bool
		if len(tagProperties) >= 2 {
			isRequired = tagProperties[indexRequired] == "required"
		}

		fmt.Fprintln(w, fmt.Sprintf(
			"%v\t%v\t%v\t%t\t%v",
			typeField.Name,
			valueField.Type(),
			envVar,
			isRequired,
			valueField))
	}

	return w.Flush()
}
