package internal

func Min(n ...uint16) uint16 {
	var m uint16 = 0xffff
	for _, o := range n {
		if o < m {
			m = o
		}
	}
	return m
}

func Max(n ...uint16) uint16 {
	var m uint16
	for _, o := range n {
		if o > m {
			m = o
		}
	}
	return m
}

func Mid(n ...uint16) uint16 {
	return (Min(n...) + Max(n...)) >> 1
}
