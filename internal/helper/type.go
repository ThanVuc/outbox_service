package helper

import "github.com/thanvuc/go-core-lib/eventbus"

func IntPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func BoolPtr(b bool) *bool {
	if !b {
		return nil
	}
	return &b
}

func ExchangeNamePtr(exchange eventbus.ExchangeName) *eventbus.ExchangeName {
	if exchange == "" {
		return nil
	}
	return &exchange
}
