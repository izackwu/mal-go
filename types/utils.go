package types

// ToMalBool converts a Golang bool to mal bool
func ToMalBool(b bool) MalLiteral {
	if b {
		return MalTrue
	}
	return MalFalse
}

// NotMalBool is `!` for mal bool
func NotMalBool(mb MalLiteral) MalLiteral {
	switch mb {
	case MalFalse:
		return MalTrue
	case MalTrue:
		return MalFalse
	default:
		return mb
	}
}
