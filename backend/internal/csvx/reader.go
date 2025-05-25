package csvx

import (
	"encoding/csv"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/unicode"

	"golang.org/x/text/transform"
	"golang.org/x/xerrors"
)

// Reader provides functionality to read CSV data into structs
type Reader struct {
	Encoding  transform.Transformer // Character encoding transformer
	Delimiter Delimiter             // Field delimiter
	UseBOM    bool                  // UseBOM defines whether to use a BOM (Byte Order Mark) in the CSV encoding transformation.
	HasHeader bool                  // Whether CSV has a header row
}

// NewDefaultReader creates a new Reader with default configuration
func NewDefaultReader() *Reader {
	return &Reader{
		Encoding:  unicode.UTF8.NewDecoder(),
		Delimiter: DelimiterComma,
		UseBOM:    false,
		HasHeader: true,
	}
}

// Read reads CSV data from the given reader and maps it to a slice of the given struct type
func (r *Reader) Read(reader io.Reader, dest interface{}) error {
	// Get the type of the destination
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("destination must be a pointer to a slice, got %T", dest)
	}

	// Get the element type of the slice
	sliceValue := destValue.Elem()
	elemType := sliceValue.Type().Elem()

	// Parse struct tags
	fields, err := parseStructTags(elemType)
	if err != nil {
		return err
	}

	// Apply encoding transformation
	transformedReader := transform.NewReader(reader, r.Encoding)
	if r.UseBOM {
		transformedReader = transform.NewReader(transformedReader, unicode.BOMOverride(r.Encoding))
	}

	// Create a CSV reader
	csvReader := csv.NewReader(transformedReader)
	csvReader.Comma = rune(r.Delimiter)

	// Map of header indices
	headerIndices := make(map[string]int)

	// Handle header row if present
	if r.HasHeader {
		// Read the header row
		headers, err := csvReader.Read()
		if err != nil {
			return err
		}

		// Create a map of header indices
		for i, header := range headers {
			headerIndices[header] = i
		}

		// Check for required fields
		for _, field := range fields {
			if field.required {
				if _, ok := headerIndices[field.header]; !ok {
					return xerrors.Errorf("required field is missing: %s", field.header)
				}
			}
		}
	} else {
		// If no header, use field index as position
		for i, field := range fields {
			if !field.ignored {
				headerIndices[field.header] = i
			}
		}
	}

	// Read and process each row
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Create a new instance of the struct
		newElem := reflect.New(elemType).Elem()

		// Fill the struct fields
		for _, field := range fields {
			if field.ignored {
				continue
			}

			fieldValue := newElem.FieldByName(field.name)
			if !fieldValue.CanSet() {
				continue
			}

			// Get the value from the CSV record
			var strValue string
			if idx, ok := headerIndices[field.header]; ok && idx < len(record) {
				strValue = record[idx]
				// Apply the default value if the field is empty and has a default value
				if strValue == "" && field.defaultValue != "" {
					strValue = field.defaultValue
				}
			} else if field.defaultValue != "" {
				// Use default value if header not found but default is provided
				strValue = field.defaultValue
			} else if field.required {
				return xerrors.Errorf("required field is missing: %s", field.header)
			} else {
				// Skip this field
				continue
			}

			// Convert the string value to the appropriate type
			if err := setFieldValue(fieldValue, strValue, field.format); err != nil {
				return fmt.Errorf("error setting field %s: %w", field.name, err)
			}
		}

		// Append the new element to the slice
		sliceValue.Set(reflect.Append(sliceValue, newElem))
	}

	return nil
}

// setFieldValue converts a string value to the appropriate type and sets it on the given field
func setFieldValue(field reflect.Value, value string, format string) error {
	// Handle other types
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == "" {
			field.SetInt(0)
			return nil
		}
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == "" {
			field.SetUint(0)
			return nil
		}
		i, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(i)
	case reflect.Float32, reflect.Float64:
		if value == "" {
			field.SetFloat(0)
			return nil
		}
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	case reflect.Bool:
		if value == "" {
			field.SetBool(false)
			return nil
		}
		b, err := strconv.ParseBool(value)
		if err != nil {
			// Try to handle "0"/"1" as a bool
			switch value {
			case "0":
				field.SetBool(false)
				return nil
			case "1":
				field.SetBool(true)
				return nil
			}
		}
		field.SetBool(b)
	case reflect.Struct:
		// Handle time.Time type specially
		if field.Type() == reflect.TypeOf(time.Time{}) {
			if value == "" {
				// Set to zero time
				field.Set(reflect.ValueOf(time.Time{}))
				return nil
			}

			// If no format is provided, use a default format
			if format == "" {
				format = time.RFC3339
			}

			// Parse the time using the provided format
			t, err := time.Parse(format, value)
			if err != nil {
				return xerrors.Errorf("failed to parse time: %w", err)
			}

			// Set the time value
			field.Set(reflect.ValueOf(t))
			return nil
		}
		fallthrough
	default:
		return xerrors.Errorf("invalid field type: %s", field.Kind().String())
	}
	return nil
}

// ReadString reads CSV data from a string and maps it to a slice of the given struct type
func (r *Reader) ReadString(data string, dest interface{}) error {
	return r.Read(strings.NewReader(data), dest)
}
