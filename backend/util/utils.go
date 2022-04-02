package util

func SlicePop[T interface{}](fuckinterIntermediateShit []T) T {
	if len(fuckinterIntermediateShit) == 0 {
		return *new(T)
	}
	return fuckinterIntermediateShit[len(fuckinterIntermediateShit)-1]
}
