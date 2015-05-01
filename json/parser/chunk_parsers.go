package parser

import (
	"encoding/json"
	"log"
	"sort"
	"strings"
)

type ChunkParser interface {
	Parse(string) string
}

// Chunk Parsers
type errorParser struct{}

func (errorParser) Parse(string) string {
	return "BADPARSER"
}

type simpleStringFieldParser struct{}

func (simpleStringFieldParser) Parse(jsn string) string {
	var data string
	if err := json.Unmarshal([]byte(strings.SplitN(jsn, ":", 2)[1]), &data); err != nil {
		log.Println("Parse string field error:", err)
		return ""
	}
	return data
}

type arrOfStringParser struct{}

func (arrOfStringParser) Parse(jsn string) string {
	arr := strings.SplitN(jsn, ":", 2)[1] // array part
	var data []interface{}
	if err := json.Unmarshal([]byte(arr), &data); err != nil {
		log.Println("Parse array of strings field error:", err)
		return ""
	}

	// Convert to string
	var res string
	for i, v := range data {
		if i > 0 {
			res += ","
		}
		res += v.(string)
	}
	return res
}

type arrOfObjectsParser struct{}

func (arrOfObjectsParser) Parse(jsn string) string {
	arr := strings.SplitN(jsn, ":", 2)[1] // array part
	var data []interface{}
	if err := json.Unmarshal([]byte(arr), &data); err != nil {
		log.Println("Parse array of objects field error:", err.Error()+" in `"+arr+"`")
		return ""
	}

	var res string
	for i, val := range data {
		if i > 0 {
			res += ","
		}

		// type switch
		switch val := val.(type) {
		case map[string]interface{}: // JSON object
			res += "{"

			tmp := make([]string, 0, 10)
			for k, v := range val {
				tmp = append(tmp, k+":"+v.(string))
			}

			// Make order stable
			sort.Strings(tmp)

			res += strings.Join(tmp, ",")
			res += "}"
		default:
			log.Println("Unexpected field type:", val)
		}
	}
	return res
}

// Register our chunk parser
func init() {
	// Register our chunk parsers
	RegisterChunkParser(MT_ERR, new(errorParser))
	RegisterChunkParser(MT_STR, new(simpleStringFieldParser))
	RegisterChunkParser(MT_ARROFSTR, new(arrOfStringParser))
	RegisterChunkParser(MT_ARROFOBJ, new(arrOfObjectsParser))
	// RegisterChunkParser(MT_ARROFOBJ, new(simpleNumFieldParser))
}
