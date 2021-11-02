package randstring

const (
	lowerCaseAlpha = "abcdefghijklmnopqrstuvwxyz"
	upperCaseAlpha = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numeric        = "0123456789"
	nonZeroNumeric = "123456789"
	special        = "!@#$%^&*"
)

// Charsets is the entrypoint to use baked-in charsets as Option.
var Charsets = charset{}

type charset struct{}

func (_ charset) LowerCaseAlpha(count uint) Option {
	return Option{
		Charset: lowerCaseAlpha,
		Count:   count,
	}
}

func (_ charset) UpperCaseAlpha(count uint) Option {
	return Option{
		Charset: upperCaseAlpha,
		Count:   count,
	}
}

func (_ charset) Alpha(count uint) Option {
	return Option{
		Charset: lowerCaseAlpha + upperCaseAlpha,
		Count:   count,
	}
}

func (_ charset) AlphaNumeric(count uint) Option {
	return Option{
		Charset: lowerCaseAlpha + upperCaseAlpha + numeric,
		Count:   count,
	}
}

func (_ charset) Numeric(count uint) Option {
	return Option{
		Charset: numeric,
		Count:   count,
	}
}

func (_ charset) NonZeroNumeric(count uint) Option {
	return Option{
		Charset: nonZeroNumeric,
		Count:   count,
	}
}

func (_ charset) Special(count uint) Option {
	return Option{
		Charset: special,
		Count:   count,
	}
}
