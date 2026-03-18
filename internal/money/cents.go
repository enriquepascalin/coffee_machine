package money

import "fmt"

type Cents int

func (c Cents) String() string {
	sign := ""
	v := c
	if v < 0 {
		sign = "-"
		v = -v
	}
	dollars := int(v) / 100
	cents := int(v) % 100
	return fmt.Sprintf("%s$%d.%02d", sign, dollars, cents)
}
