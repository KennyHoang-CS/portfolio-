package engine

// Value represents any attribute value in a BUILD file.
// It can be a string, list of strings, bool, number, etc.
type Value interface{}

func AsString(v Value) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func AsStringList(v Value) []string {
	if list, ok := v.([]string); ok {
		return list
	}
	return nil
}