package models

func reverse(heightByte *[]byte) {
	a := *heightByte
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	*heightByte = a
}