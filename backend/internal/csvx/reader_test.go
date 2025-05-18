package csvx_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"net.bright-room.dev/calender-api/internal/csvx"
)

func TestReader_CSVFileReadingForVariousEncodings(t *testing.T) {
	type person struct {
		Name string `csv:"name"`
		Age  int    `csv:"age"`
	}

	tests := []struct {
		name     string
		filePath string
		encoding transform.Transformer
		useBOM   bool
		expected interface{}
	}{
		{
			name:     "UTF-8",
			filePath: "./testdata/utf8.csv",
			encoding: unicode.UTF8BOM.NewDecoder(),
			useBOM:   false,
			expected: []person{
				{Name: "Yamada taro", Age: 20},
				{Name: "Kojima naoki", Age: 30},
			},
		},
		{
			name:     "UTF-8 with BOM",
			filePath: "./testdata/utf8_bom.csv",
			encoding: unicode.UTF8BOM.NewDecoder(),
			useBOM:   true,
			expected: []person{
				{Name: "山田　太郎", Age: 20},
				{Name: "小島　直樹", Age: 30},
			},
		},
		{
			name:     "Shift-JIS",
			filePath: "./testdata/shift_jis.csv",
			encoding: japanese.ShiftJIS.NewDecoder(),
			useBOM:   false,
			expected: []person{
				{Name: "山田　太郎", Age: 20},
				{Name: "小島　直樹", Age: 30},
			},
		},
		{
			name:     "UTF-16BE with BOM",
			filePath: "./testdata/utf16be_bom.csv",
			encoding: unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewDecoder(),
			useBOM:   false,
			expected: []person{
				{Name: "山田　太郎", Age: 20},
				{Name: "小島　直樹", Age: 30},
			},
		},
		{
			name:     "UTF-16LE with BOM",
			filePath: "./testdata/utf16le_bom.csv",
			encoding: unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder(),
			useBOM:   false,
			expected: []person{
				{Name: "山田　太郎", Age: 20},
				{Name: "小島　直樹", Age: 30},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := csvx.Reader{
				Encoding:  tt.encoding,
				Delimiter: csvx.DelimiterComma,
				HasHeader: true,
			}

			file, _ := os.Open(tt.filePath)
			defer func(file *os.File) {
				_ = file.Close()
			}(file)

			var p []person
			_ = reader.Read(file, &p)

			assert.Equal(t, tt.expected, p)
		})
	}
}

func TestReader_ReadingWithUnsupportedEncoding(t *testing.T) {
	type person struct {
		Name string `csv:"name"`
		Age  int    `csv:"age"`
	}

	tests := []struct {
		name     string
		filePath string
		encoding transform.Transformer
		useBOM   bool
		expected interface{}
	}{
		{
			name:     "UTF-16BE",
			filePath: "./testdata/utf16be.csv",
			encoding: unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM).NewDecoder(),
			useBOM:   false,
			expected: []person{
				{Name: "Yamada taro", Age: 20},
				{Name: "Kojima naoki", Age: 30},
			},
		},
		{
			name:     "UTF-16LE",
			filePath: "./testdata/utf16le.csv",
			encoding: unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder(),
			useBOM:   false,
			expected: []person{
				{Name: "Yamada taro", Age: 20},
				{Name: "Kojima naoki", Age: 30},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := csvx.Reader{
				Encoding:  tt.encoding,
				Delimiter: csvx.DelimiterComma,
				HasHeader: true,
			}

			file, _ := os.Open(tt.filePath)
			defer func(file *os.File) {
				_ = file.Close()
			}(file)

			var p []person
			_ = reader.Read(file, &p)

			assert.NotEqual(t, tt.expected, p)
		})
	}
}

func TestReader_ReadingJapaneseHeaders(t *testing.T) {
	type person struct {
		Name string `csv:"名前"`
		Age  int    `csv:"年齢"`
	}

	tests := []struct {
		name     string
		filePath string
		expected interface{}
	}{
		{
			name:     "日本語ヘッダーの読み込み",
			filePath: "./testdata/japanese_header.csv",
			expected: []person{
				{Name: "山田　太郎", Age: 20},
				{Name: "小島　直樹", Age: 30},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := csvx.Reader{
				Encoding:  unicode.UTF8.NewDecoder(),
				Delimiter: csvx.DelimiterComma,
				HasHeader: true,
			}

			file, _ := os.Open(tt.filePath)
			defer func(file *os.File) {
				_ = file.Close()
			}(file)

			var p []person
			_ = reader.Read(file, &p)

			assert.Equal(t, tt.expected, p)
		})
	}
}

func TestReader_HeaderDoesNotExist(t *testing.T) {
	type person struct {
		Name string `csv:"name"`
		Age  int    `csv:"age"`
	}

	tests := []struct {
		name     string
		filePath string
		expected interface{}
	}{
		{
			name:     "ヘッダーが存在しない場合を考慮した読み込み",
			filePath: "./testdata/no_header.csv",
			expected: []person{
				{Name: "山田　太郎", Age: 20},
				{Name: "小島　直樹", Age: 30},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := csvx.Reader{
				Encoding:  unicode.UTF8.NewDecoder(),
				Delimiter: csvx.DelimiterComma,
				HasHeader: false,
			}

			file, _ := os.Open(tt.filePath)
			defer func(file *os.File) {
				_ = file.Close()
			}(file)

			var p []person
			_ = reader.Read(file, &p)

			assert.Equal(t, tt.expected, p)
		})
	}
}

func TestReader_CharactersOtherThanDelimitedColon(t *testing.T) {
	type person struct {
		Name string `csv:"name"`
		Age  int    `csv:"age"`
	}

	tests := []struct {
		name      string
		filePath  string
		delimiter csvx.Delimiter
		expected  interface{}
	}{
		{
			name:      "タブ",
			filePath:  "./testdata/delimiter_tab.csv",
			delimiter: csvx.DelimiterTab,
			expected: []person{
				{Name: "Yamada taro", Age: 20},
				{Name: "Kojima naoki", Age: 30},
			},
		},
		{
			name:      "セミコロン",
			filePath:  "./testdata/delimiter_semicolon.csv",
			delimiter: csvx.DelimiterSemicolon,
			expected: []person{
				{Name: "Yamada taro", Age: 20},
				{Name: "Kojima naoki", Age: 30},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := csvx.Reader{
				Encoding:  unicode.UTF8.NewDecoder(),
				Delimiter: tt.delimiter,
				HasHeader: true,
			}

			file, _ := os.Open(tt.filePath)
			defer func(file *os.File) {
				_ = file.Close()
			}(file)

			var p []person
			_ = reader.Read(file, &p)

			assert.Equal(t, tt.expected, p)
		})
	}
}

func TestReader_RequiredFieldErrors(t *testing.T) {
	type person struct {
		Name       string `csv:"name"`
		Age        int    `csv:"age"`
		Occupation string `csv:"occupation,required"`
	}

	tests := []struct {
		name     string
		filePath string
	}{
		{
			name:     "必須フィールドが存在しない場合のエラー",
			filePath: "./testdata/required_field_errors.csv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := csvx.Reader{
				Encoding:  unicode.UTF8.NewDecoder(),
				Delimiter: csvx.DelimiterComma,
				HasHeader: true,
			}

			file, _ := os.Open(tt.filePath)
			defer func(file *os.File) {
				_ = file.Close()
			}(file)

			var p []person
			err := reader.Read(file, &p)

			assert.Error(t, err)
		})
	}
}

func TestReader_IfValueDoesNotExistDefaultValueIsAssigned(t *testing.T) {
	type person struct {
		Name       string `csv:"name"`
		Age        int    `csv:"age"`
		Occupation string `csv:"occupation" default:"無職"`
	}

	tests := []struct {
		name     string
		filePath string
		expected interface{}
	}{
		{
			name:     "値が存在しない場合デフォルト値が代入される",
			filePath: "./testdata/default_value_assigned.csv",
			expected: []person{
				{Name: "山田　太郎", Age: 20, Occupation: "会社員"},
				{Name: "小島　直樹", Age: 30, Occupation: "無職"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := csvx.Reader{
				Encoding:  unicode.UTF8.NewDecoder(),
				Delimiter: csvx.DelimiterComma,
				HasHeader: true,
			}

			file, _ := os.Open(tt.filePath)
			defer func(file *os.File) {
				_ = file.Close()
			}(file)

			var p []person
			_ = reader.Read(file, &p)

			assert.Equal(t, tt.expected, p)
		})
	}
}

func TestReader_NotParsedIfTagIsIgnore(t *testing.T) {
	type person struct {
		Name       string `csv:"name"`
		Age        int    `csv:"-"`
		Occupation string `csv:"occupation"`
	}

	tests := []struct {
		name     string
		filePath string
		expected interface{}
	}{
		{
			name:     "タグがignoreの場合パースされない",
			filePath: "./testdata/ignore.csv",
			expected: []person{
				{Name: "山田　太郎", Occupation: "会社員"},
				{Name: "小島　直樹", Occupation: "無職"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := csvx.Reader{
				Encoding:  unicode.UTF8.NewDecoder(),
				Delimiter: csvx.DelimiterComma,
				HasHeader: true,
			}

			file, _ := os.Open(tt.filePath)
			defer func(file *os.File) {
				_ = file.Close()
			}(file)

			var p []person
			_ = reader.Read(file, &p)

			assert.Equal(t, tt.expected, p)
		})
	}
}

func TestReader_ParsingToVariousDatetimeFormats(t *testing.T) {
	type timeParser struct {
		ISO8601DateTime time.Time `csv:"iso8601_datetime" format:"2006-01-02T15:04:05"`
		ISO8601Date     time.Time `csv:"iso8601_date"  format:"2006-01-02"`
		ISO8601Time     time.Time `csv:"iso8601_time"  format:"15:04:05"`
		RFCLikeDateTime time.Time `csv:"rfc_like_date" format:"2006-01-02 15:04:05"`
		JPDateTime      time.Time `csv:"jp_datetime" format:"2006/01/02 15:04:05"`
		JPDate          time.Time `csv:"jp_date" format:"2006/01/02"`
		CompactDateTime time.Time `csv:"compact_datetime" format:"20060102150405"`
		CompactDate     time.Time `csv:"compact_date" format:"20060102"`
		CompactTime     time.Time `csv:"compact_time" format:"150405"`
	}

	tests := []struct {
		name     string
		filePath string
		expected interface{}
	}{
		{
			name:     "様々な日時フォーマットへのパース",
			filePath: "./testdata/datetime.csv",
			expected: []timeParser{
				{
					ISO8601DateTime: time.Date(2025, 1, 1, 10, 50, 11, 0, time.UTC),
					ISO8601Date:     time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC),
					ISO8601Time:     time.Date(0, 1, 1, 9, 12, 50, 0, time.UTC),
					RFCLikeDateTime: time.Date(2025, 10, 5, 20, 26, 5, 0, time.UTC),
					JPDateTime:      time.Date(2024, 4, 25, 8, 0o0, 26, 0, time.UTC),
					JPDate:          time.Date(2024, 9, 25, 0, 0, 0, 0, time.UTC),
					CompactDateTime: time.Date(2025, 12, 31, 2, 10, 59, 0, time.UTC),
					CompactDate:     time.Date(2025, 2, 14, 0, 0, 0, 0, time.UTC),
					CompactTime:     time.Date(0, 1, 1, 1, 15, 0o0, 0, time.UTC),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := csvx.Reader{
				Encoding:  unicode.UTF8.NewDecoder(),
				Delimiter: csvx.DelimiterComma,
				HasHeader: true,
			}

			file, _ := os.Open(tt.filePath)
			defer func(file *os.File) {
				_ = file.Close()
			}(file)

			var p []timeParser
			_ = reader.Read(file, &p)

			assert.Equal(t, tt.expected, p)
		})
	}
}
