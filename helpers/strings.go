package helpers

func StringsInclude(set []string, check string) bool {
	for _, member := range set {
		if member == check {
			return true
		}
	}
	return false
}
