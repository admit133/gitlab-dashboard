package utils

func StringsContainString(strings []string, needle string) bool {
	for key := range strings {
		if needle == strings[key] {
			return true
		}
	}

	return false
}
