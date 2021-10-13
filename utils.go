package main

//basic min function; go doesn't use generics btw
func minInt(a int, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}

//minPositiveInt guards against returning min(a,b)<=-1 if possible
func minPositiveInt(a int, b int) int {
	if a <= -1 && b >= 0 {
		return b
	}
	if b <= -1 && a >= 0 {
		return a
	}
	return minInt(a, b)
}
