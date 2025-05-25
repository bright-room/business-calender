package csvx

import (
	"reflect"
	"strings"

	"golang.org/x/xerrors"
)

type fieldInfo struct {
	name         string       // Field name in the struct
	header       string       // CSV header name
	required     bool         // Whether the field is required
	ignored      bool         // Whether the field should be ignored
	defaultValue string       // defaultValue value for the field
	format       string       // format string for date/time fields
	fieldType    reflect.Type // The type of the field
}

func parseStructTags(t reflect.Type) ([]fieldInfo, error) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, xerrors.Errorf("provided value is not a struct")
	}

	numFields := t.NumField()
	fields := make([]fieldInfo, 0, numFields)

	for i := 0; i < numFields; i++ {
		field := t.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		info := fieldInfo{
			name:      field.Name,
			fieldType: field.Type,
		}

		// Parse csv tag
		csvTag := field.Tag.Get("csv")
		if csvTag == "-" {
			info.ignored = true
			continue
		}

		parts := strings.Split(csvTag, ",")
		if len(parts) > 0 {
			info.header = parts[0]

			// Check for required flag
			if len(parts) > 1 && parts[1] == "required" {
				info.required = true
			}
		}

		// If no header is specified, use the field name
		if info.header == "" {
			info.header = field.Name
		}

		// Parse default tag
		defaultTag := field.Tag.Get("default")
		if defaultTag != "" {
			info.defaultValue = defaultTag
		}

		// Parse format tag
		formatTag := field.Tag.Get("format")
		if formatTag != "" {
			info.format = formatTag
		}

		fields = append(fields, info)
	}

	return fields, nil
}

func getHeaders(fields []fieldInfo) []string {
	headers := make([]string, 0, len(fields))
	for _, field := range fields {
		if !field.ignored {
			headers = append(headers, field.header)
		}
	}
	return headers
}
