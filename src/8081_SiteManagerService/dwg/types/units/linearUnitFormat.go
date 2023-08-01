package units

type LinearUnitFormat int16

const (
	LU_None           LinearUnitFormat = 0
	LU_Scientific     LinearUnitFormat = 1
	LU_Decimal        LinearUnitFormat = 2
	LU_Engineering    LinearUnitFormat = 3
	LU_Architectural  LinearUnitFormat = 4
	LU_Fractional     LinearUnitFormat = 5
	LU_WindowsDesktop LinearUnitFormat = 6
)
