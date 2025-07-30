package querymap

import (
	"github.com/mitchellh/mapstructure"
	"golang.org/x/exp/maps"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

// anyList - type for storing a slice of arbitrary elements.
type anyList []any

// QueryMap is a map storing key-values from query-parameters.
// The value can be of one of the following types: string, []string, QueryMap, anyList.
type QueryMap map[string]any

// newQueryMap creates and returns an empty QueryMap.
func newQueryMap() QueryMap {
	return make(QueryMap)
}

// set sets the `untypedValue` value in the map by key `key`.
// If the value by key already exists, the method correctly merges the
// new data with old data (string, []string, anyList, QueryMap).
func (q QueryMap) set(key string, untypedValue any) QueryMap {
	untypedEntry, ok := q[key]
	if !ok {
		q[key] = untypedValue
		return q
	}

	switch entry := untypedEntry.(type) {
	case string:
		switch value := untypedValue.(type) {
		case string: // string1 + string2 = []string{string1, string2}
			q[key] = []string{entry, value}
		case []string: // string1 + []string{string2, string3} = []string{string1, string2, string3}
			q[key] = append([]string{entry}, value...)
		case anyList: // string1 + []any{var1, var2} = []any{string1, var1, var2}
			q[key] = append(anyList{entry}, value...)
		case QueryMap: // string1 + {key1: val1, key2: val2} = []any{string1, {key1: val1, key2: val2}}
			q[key] = anyList{entry, value}
		}

	case []string:
		switch value := untypedValue.(type) {
		case string: // []string{string1, string2} + string3 = []string{string1, string2, string3}
			q[key] = append(entry, value)
		case []string: // []string{string1} + []string{string2} = []string{string1, string2}
			q[key] = append(entry, value...)
		case anyList: // []string{string1} + []any{var1, var2} = []any{string1, var1, var2}
			slc := anyList{}
			for _, s := range entry {
				slc = append(slc, s)
			}
			q[key] = append(slc, value...)
		case QueryMap: // []string{string1} + {key1: val1, key2: val2} = []any{string1, {key1: val1, key2: val2}}
			slc := anyList{}
			for _, s := range entry {
				slc = append(slc, s)
			}
			q[key] = append(slc, value)
		}

	case anyList:
		switch value := untypedValue.(type) {
		case string: // []any{var1, var2} + string3 = []any{var1, var2, string3}
			q[key] = append(entry, value)
		case []string: // []any{var1, var2} + []string{string3, string4} = []any{var1, var2, string3, string4}
			for _, s := range value {
				entry = append(entry, s)
			}
			q[key] = entry
		case anyList: // []any{var1, var2} + []any{var3, var4} = []any{var1, var2, var3, var4} or can merge(not needed)
			q[key] = append(entry, value...)
		case QueryMap: // []any{var1, var2} + {key1: val1} = []any{var1, var2, {key1: val1}} or can merge(not needed)
			q[key] = append(entry, value)
		}

	case QueryMap:
		switch value := untypedValue.(type) {
		case string: // {key1: val1} + string2 = []any{{key1: val1}, string2}
			q[key] = append(anyList{entry}, value)
		case []string: // {key1: val1} + []string{string2, string3} = []any{{key1: val1}, string2, string3}
			slc := anyList{entry}
			for _, s := range value {
				slc = append(slc, s)
			}
			q[key] = slc
		case anyList: // {key1: val1} + []any{var1, var2} = []any{{key1: val1}, var1, var2} or can merge(not needed)
			slc := anyList{entry}
			slc = append(slc, value...)
			q[key] = slc
		case QueryMap: // {key1: val1} + {key2: val2} = {key1: val1, key2: val2}
			for typedValueKey, typedValueValue := range value {
				entry.set(typedValueKey, typedValueValue)
			}
		}
	}

	return q
}

// nestedQuery - recursively parses the key of the form "key[a][b]" and forms nested structures.
func nestedQuery(data QueryMap, key string, value []string) QueryMap {
	nextStart := strings.IndexRune(key, '[')
	nextEnd := strings.IndexRune(key, ']')

	currentKey := key

	if nextStart == -1 && nextEnd != -1 && len(key) == nextEnd+1 { // key]
		currentKey = key[:nextEnd]
	} else if nextStart != -1 && nextEnd != -1 && nextStart+1 == nextEnd { // key[]
		currentKey = key[:nextStart]
	} else if nextEnd+1 == nextStart { // key][
		currentKey = key[:nextEnd]
	} else if nextStart != -1 && nextStart < nextEnd { // key[a] or key[]
		currentKey = key[:nextStart]
	}

	// If the format is "key[]" or "key][]"
	//  key[] or key][] and no any text after
	//  regex: \[\]$ OR regex: \]\[$
	if nextStart+1 == nextEnd && nextEnd+1 == len(key) || nextEnd != -1 && nextStart > nextEnd && key[nextStart:] == "[]" && nextStart+2 == len(key) {
		return data.set(currentKey, value)
	}

	if nextStart != -1 {
		nextKey := ""
		if nextEnd != -1 && nextEnd < nextStart { // b][a
			nextKey = key[nextEnd+2:]
		} else { // b[]
			nextKey = key[nextStart+1:]
		}

		return data.set(currentKey, nestedQuery(newQueryMap(), nextKey, value))
	}

	// If there is only one value, write it as string
	if len(value) == 1 {
		return data.set(currentKey, value[0])
	}

	// Otherwise, we save the slice
	return data.set(currentKey, value)
}

// FromURL parses the *url.URL object and returns a QueryMap representing
// all its query parameters as a nested structure.
func FromURL(URL *url.URL) QueryMap {
	return FromValues(URL.Query())
}

// FromValues parses the url.Values object and returns a QueryMap representing
// all its query parameters as a nested structure.
func FromValues(urlQuery url.Values) QueryMap {
	data := newQueryMap()

	urlQueryKeys := maps.Keys(urlQuery)
	slices.Sort(urlQueryKeys)

	// First sort the keys for a predictable order
	for _, key := range urlQueryKeys {
		value := urlQuery[key]

		nestedQuery(data, key, value)
	}

	// Normalize the values (converting a set of numeric keys to a slice)
	for k, v := range data {
		data[k] = NormalizeSlicesNumbersIndexes(v)
	}

	return data
}

// ToStruct converts QueryMap into a structure of type T using mapstructure.
// The fields of the structure are read by the `json` tag.
func ToStruct[T any](m QueryMap) (*T, error) {
	var result T

	config := &mapstructure.DecoderConfig{Metadata: nil, Result: &result, WeaklyTypedInput: true, TagName: "json"}
	decoder, _ := mapstructure.NewDecoder(config)
	if err := decoder.Decode(m); err != nil {
		return nil, err
	}

	return &result, nil
}

// FromURLToStruct is a convenient function that combines FromURL and ToStruct.
// Accepts *url.URL and tries to convert the query string into a T structure.
func FromURLToStruct[T any](URL *url.URL) (*T, error) {
	return ToStruct[T](FromURL(URL))
}

// FromValuesToStruct is a convenient function that combines FromValues and ToStruct.
// Accepts url.Values and tries to convert the query string into a T structure.
func FromValuesToStruct[T any](values url.Values) (*T, error) {
	return ToStruct[T](FromValues(values))
}

// FromURLStringToStruct is an additional wrapper that parses the URL string,
// and then calls FromURLToStruct.
func FromURLStringToStruct[T any](URL string) (*T, error) {
	parsedUrl, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}

	return FromURLToStruct[T](parsedUrl)
}

// NormalizeSlicesNumbersIndexes recursively checks whether the value is
// a set of numeric keys, and if so, converts it to a slice (anyList).
// For example, QueryMap{"0": "first", "1": "second"} => []any{"first", "second"}.
func NormalizeSlicesNumbersIndexes(v any) any {
	switch value := v.(type) {
	case string:
		return value
	case []string:
		return value
	case QueryMap:
		total := 0
		keyAreNumbers := 0

		for key := range value {
			total++
			if _, err := strconv.Atoi(key); err == nil {
				keyAreNumbers++
			}
		}

		// If all keys are numbers, sort and turn into a slice
		if total == keyAreNumbers {
			slc := anyList{}

			valueKeys := maps.Keys(value)
			slices.Sort(valueKeys)
			for _, valueKey := range valueKeys {
				slc = append(slc, value[valueKey])
			}

			return slc
		}

		// Otherwise recursively process nested values
		for k, v := range value {
			value[k] = NormalizeSlicesNumbersIndexes(v)
		}

		return value
	case anyList:
		entry := make(anyList, len(value))
		for i, v := range value {
			entry[i] = NormalizeSlicesNumbersIndexes(v)
		}
		return entry
	}

	// If not one of the above cases, return as is
	return v
}
