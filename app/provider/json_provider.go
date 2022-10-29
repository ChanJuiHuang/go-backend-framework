package provider

import jsoniter "github.com/json-iterator/go"

func provideJson() jsoniter.API {
	json := jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		CaseSensitive:          true,
	}.Froze()

	return json
}
