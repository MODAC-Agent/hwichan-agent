// Rule 13 위반: "common" 역시 책임이 불명확한 모호 이름.
// util / common / shared / helper / misc 모두 같은 안티패턴.
package common

type Pair struct {
	Key   string
	Value string
}

func MakePair(k, v string) Pair {
	return Pair{Key: k, Value: v}
}
