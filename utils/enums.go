package utils

import (
	"fmt"
)

type Type string

const (
	Goldankauf        Type = "goldankauf"
	Trauringe         Type = "trauringe"
	Verlobungsringe   Type = "verlobungsringe"
	Ohrlohstechen     Type = "ohrlohstechen"
	Sonstiges         Type = "sonstiges"
)

func StringToType(value string) (Type, error) {
	switch value {
	case string(Goldankauf):
		return Goldankauf, nil
	case string(Trauringe):
		return Trauringe, nil
	case string(Verlobungsringe):
		return Verlobungsringe, nil
	case string(Ohrlohstechen):
		return Ohrlohstechen, nil
	case string(Sonstiges):
		return Sonstiges, nil
	default:
		return "", fmt.Errorf("invalid type: %s", value)
	}
}
