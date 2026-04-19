// Rule 13 위반: "util"은 책임이 드러나지 않는 모호한 이름.
// 어떤 책임이든 들어갈 수 있어 응집도 저하의 씨앗이 됨.
// 권장: 책임 단위로 분리 (예: time / encoding / validation).
package util

func Zero() int { return 0 }

func CoalesceString(a, b string) string {
	if a == "" {
		return b
	}
	return a
}
