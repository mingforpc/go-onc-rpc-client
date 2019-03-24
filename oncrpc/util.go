package oncrpc

func createFragmentHeader(size uint32, last bool) uint32 {
	header := size | uint32(0)
	if last {
		header = header | (1 << 31)
	}

	return header
}

func getFragment(header uint32) (uint32, bool) {
	size := header &^ (1 << 31)

	last := false
	if header&(1<<31) > 0 {
		last = true
	}

	return size, last
}
