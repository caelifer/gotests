package parser

// Parse parses a chunk of JSON encoded data, based on the type information
// provided by the meta.
func Parse(meta, jsn string) string {
	// Dispatch to the proper custom chunkParser
	return dispatcher[getMetaType(meta)].parse(jsn)
}
