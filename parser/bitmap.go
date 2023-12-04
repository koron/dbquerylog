package parser

type bitmap struct {
	b []byte
	m uint16
}

func (bm *bitmap) get(n uint16) bool {
	if n > bm.m {
		return false
	}
	i, j := n/8, n%8
	return bm.b[i]&(1<<j) != 0
}

// Unused private function.
//func (bm *bitmap) rank() uint16 {
//	c := bm.m
//	for i := uint16(0); i < bm.m; i++ {
//		if bm.get(i) {
//			c--
//		}
//	}
//	return c
//}
