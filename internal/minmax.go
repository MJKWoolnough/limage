package internal

func min(n ...uint16) uint16 {
	var m uint16 = 0xffff
	for _, o := range n {
		if o < m {
			m = o
		}
	}
	return m
}

func max(n ...uint16) uint16 {
	var m uint16
	for _, o := range n {
		if o > m {
			m = o
		}
	}
	return m
}

func mid(n ...uint16) uint16 {
	return (Min(n...) + Max(n...)) >> 1
}
