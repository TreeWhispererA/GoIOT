package tables

type ZeroHandling byte

const (
	ZH_SuppressZeroFeetAndInches               = 0
	ZH_ShowZeroFeetAndInches                   = 1
	ZH_ShowZeroFeetSuppressZeroInches          = 2
	ZH_SuppressZeroFeetShowZeroInches          = 3
	ZH_SuppressDecimalLeadingZeroes            = 4
	ZH_SuppressDecimalTrailingZeroes           = 8
	ZH_SuppressDecimalLeadingAndTrailingZeroes = 12
)
