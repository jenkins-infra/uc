package update

func Contains(s []DepInfo, name string) bool {
	for _, a := range s {
		if a.Name == name {
			return true
		}
	}
	return false
}
