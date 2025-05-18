package csvx

import (
	"bytes"
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

// Writer provides functionality to write structs to CSV data
type Writer struct {
	Encoding  transform.Transformer // Character encoding transformer
	Delimiter Delimiter             // Field delimiter
	UseCRLF   bool                  // True to use \r\n as the line terminator
	HasHeader bool                  // Whether CSV has a header row
}

// NewDefaultWriter creates a new Writer with default configuration
func NewDefaultWriter() *Writer {
	return &Writer{
		Encoding:  unicode.UTF8.NewEncoder(),
		Delimiter: DelimiterComma,
		UseCRLF:   false,
	}
}

// Write writes a slice of structs to CSV format
func (w *Writer) Write(writer io.Writer, data interface{}) error {
	// Get the value of the data
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}

	// Ensure data is a slice
	if dataValue.Kind() != reflect.Slice {
		return fmt.Errorf("data must be a slice, got %T", data)
	}

	// If the slice is empty, return early
	if dataValue.Len() == 0 {
		return xerrors.Errorf("empty CSV data")
	}

	// Get the type of the slice elements
	elemType := dataValue.Type().Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	// Parse struct tags
	fields, err := parseStructTags(elemType)
	if err != nil {
		return err
	}

	// First write to a buffer so we can apply encoding transformation
	var buf bytes.Buffer

	// Create a CSV writer
	csvWriter := csv.NewWriter(&buf)
	csvWriter.Comma = rune(w.Delimiter)

	// Get headers
	headers := getHeaders(fields)

	// Write the header row if configured
	if w.HasHeader {
		if err := csvWriter.Write(headers); err != nil {
			return err
		}
	}

	// Write each row
	for i := 0; i < dataValue.Len(); i++ {
		rowValue := dataValue.Index(i)
		if rowValue.Kind() == reflect.Ptr {
			rowValue = rowValue.Elem()
		}

		// Create a row with values for each field
		row := make([]string, 0, len(headers))

		for _, field := range fields {
			if field.ignored {
				continue
			}

			// Get the field value
			fieldValue := rowValue.FieldByName(field.name)

			// Convert the field value to string
			strValue, err := getFieldStringValue(fieldValue, field.format)
			if err != nil {
				return fmt.Errorf("error getting string value for field %s: %w", field.name, err)
			}

			// If the field is empty and has a default value, use the default
			if strValue == "" && field.defaultValue != "" {
				strValue = field.defaultValue
			}

			// If the field is required and empty, return an error
			if field.required && strValue == "" {
				return xerrors.Errorf("required field is missing: %s", field.header)
			}

			row = append(row, strValue)
		}

		// Write the row
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}

	// Flush the writer
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return err
	}

	// Apply encoding transformation
	transformedWriter := transform.NewWriter(writer, w.Encoding)

	// Write the transformed data
	_, err = transformedWriter.Write(buf.Bytes())
	if err != nil {
		return err
	}

	// Close the transformer to flush any remaining data
	err = transformedWriter.Close()
	if err != nil {
		return err
	}

	return nil
}

// getFieldStringValue converts a field value to a string
func getFieldStringValue(field reflect.Value, format string) (string, error) {
	if !field.IsValid() {
		return "", nil
	}

	// Handle other types
	switch field.Kind() {
	case reflect.String:
		return field.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(field.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(field.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(field.Float(), 'f', -1, 64), nil
	case reflect.Bool:
		return strconv.FormatBool(field.Bool()), nil
	case reflect.Struct:
		// Handle time.Time type specially
		if field.Type() == reflect.TypeOf(time.Time{}) {
			t := field.Interface().(time.Time)

			// If the time is zero, return an empty string
			if t.IsZero() {
				return "", nil
			}

			// If no format is provided, use a default format
			if format == "" {
				format = time.RFC3339
			}

			// format the time using the provided format
			return t.Format(format), nil
		}
		fallthrough
	default:
		return "", xerrors.Errorf("invalid field type: %s", field.Kind().String())
	}
}

// WriteString writes a slice of structs to a CSV string
func (w *Writer) WriteString(data interface{}) (string, error) {
	var buf strings.Builder
	if err := w.Write(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
