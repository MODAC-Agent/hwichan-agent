package common

type Pair struct {
	Key   string
	Value string
}

func MakePair(k, v string) Pair {
	return Pair{Key: k, Value: v}
}
