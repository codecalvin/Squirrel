package utils

// Return keys of the given map with string as key type
func StringKeys(m map[string] interface{}) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func UInt8Keys(m map[uint8] *interface{}) (keys []uint8) {
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
