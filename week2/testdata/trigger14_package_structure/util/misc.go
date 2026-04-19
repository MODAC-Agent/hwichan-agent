package util

func Zero() int { return 0 }

func CoalesceString(a, b string) string {
	if a == "" {
		return b
	}
	return a
}
