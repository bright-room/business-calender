package csvx

// Delimiter represents the field delimiter for CSV files
type Delimiter rune

// Predefined delimiters
const (
	DelimiterComma     Delimiter = ','
	DelimiterTab       Delimiter = '\t'
	DelimiterSemicolon Delimiter = ';'
)
