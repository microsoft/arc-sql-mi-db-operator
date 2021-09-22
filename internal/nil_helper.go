package internal

func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func SafeInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func SafeBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func SetString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
