package parser

type metaType int

// metaType values
const (
	MetaError metaType = iota
	MetaBoolField
	MetaNumberField
	MetaStringField
	MetaArrayOfStrings
	MetaArrayOfObjects
)

func getMetaType(meta string) metaType {
	switch meta {
	case "stringField":
		return MetaStringField
	case "arrayOfStrings":
		return MetaArrayOfStrings
	case "arrayOfObjects":
		return MetaArrayOfObjects
	default:
		return MetaError
	}
}

// Dispatcher a package global
var dispatcher = make(map[metaType]chunkJSONParser)

func registerChunkParser(meta metaType, parser chunkJSONParser) {
	// Install new parser
	dispatcher[meta] = parser
}
