package parser

type MetaType int

const (
	MT_ERR MetaType = iota
	MT_BOOL
	MT_NUM
	MT_STR
	MT_ARROFSTR
	MT_ARROFOBJ
)

func GetMeta(meta string) MetaType {
	switch meta {
	case "stringField":
		return MT_STR
	case "arrayOfStrings":
		return MT_ARROFSTR
	case "arrayOfObjects":
		return MT_ARROFOBJ
	default:
		return MT_ERR
	}
}

// Dispatcher
var Dispatch = make(map[MetaType]ChunkParser)

func RegisterChunkParser(meta MetaType, parser ChunkParser) {
	// Install new parser
	Dispatch[meta] = parser
}
