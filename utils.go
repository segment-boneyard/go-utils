package utils

// StringsToMap returns map used as a set from the given string slice.
func StringsToMap(s []string) map[string]struct{} {
	m := make(map[string]struct{})

	for _, v := range s {
		m[v] = struct{}{}
	}

	return m
}

// BlacklistKeys keys in the given map, returning a new map.
func BlacklistKeys(m map[string]interface{}, keys map[string]struct{}) map[string]interface{} {
	ret := make(map[string]interface{})

	for k, v := range m {
		if _, ok := keys[k]; ok {
			continue
		} else {
			ret[k] = v
		}
	}

	return ret
}
