package internal

const (
	// この環境におけるintのサイズ 32 or 64
	IntSize = 32 << (^uint(0) >> 63)

	// この環境におけるintが表現できる最大値 math.MaxInt32 or math.MaxInt64
	MaxInt = (1 << (IntSize - 1)) - 1
)
