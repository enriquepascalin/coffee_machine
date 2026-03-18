package money

import "fmt"

type CashBox struct {
	counts map[Denomination]int
}

func NewCashBox(initial map[Denomination]int) CashBox {
	cloned := make(map[Denomination]int, len(initial))
	for k, v := range initial {
		cloned[k] = v
	}
	return CashBox{counts: cloned}
}

func (c CashBox) Count(d Denomination) int {
	return c.counts[d]
}

func (c CashBox) Total() Cents {
	var sum Cents
	for d, n := range c.counts {
		sum += d.Cents() * Cents(n)
	}
	return sum
}

func (c CashBox) Add(d Denomination, n int) error {
	if n < 0 {
		return fmt.Errorf("add negative count: %d", n)
	}
	c.counts[d] += n
	return nil
}

func (c CashBox) Remove(d Denomination, n int) error {
	if n < 0 {
		return fmt.Errorf("remove negative count: %d", n)
	}
	if c.counts[d] < n {
		return fmt.Errorf("%w: need %d of %v, have %d", ErrCannotMakeChange, n, d.Cents(), c.counts[d])
	}
	c.counts[d] -= n
	return nil
}

// Snapshot returns a copy for dry-run decision making.
func (c CashBox) Snapshot() map[Denomination]int {
	s := make(map[Denomination]int, len(c.counts))
	for k, v := range c.counts {
		s[k] = v
	}
	return s
}
