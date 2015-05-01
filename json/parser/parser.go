package parser

func Parse(meta, jsn string) string {
	return Dispatch[GetMeta(meta)].Parse(jsn)
}
