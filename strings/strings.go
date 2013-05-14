package strings

func Reverse(s string) string {
	r := []rune(s)

	l := len(r)
	for i := 0; i < l/2; i++ {
		j := l - i - 1
		if j > i {
			// Exchange
			r[i], r[j] = r[j], r[i]
		}
	}

	return string(r)
}
