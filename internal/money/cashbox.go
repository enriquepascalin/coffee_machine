package money

import "fmt"

type CashBox struct {
	counts map[Denomination]int
}

func NewCashBox(initial map[Denomination]int) CashBox {
	cloned := make(map[Denomination]int, len(initial))
	for denomination, count := range initial {
		cloned[denomination] = count
	}

	return CashBox{
		counts: cloned,
	}
}

func (c *CashBox) Count(denomination Denomination) int {
	return c.counts[denomination]
}

func (c *CashBox) Total() Cents {
	var total Cents

	for denomination, count := range c.counts {
		total += denomination.Cents() * Cents(count)
	}

	return total
}

func (c *CashBox) Add(denomination Denomination, count int) error {
	if count < 0 {
		return fmt.Errorf("add negative count: %d", count)
	}

	c.counts[denomination] += count
	return nil
}

func (c *CashBox) Remove(denomination Denomination, count int) error {
	if count < 0 {
		return fmt.Errorf("remove negative count: %d", count)
	}

	if c.counts[denomination] < count {
		return fmt.Errorf(
			"%w: need %d of %v, have %d",
			ErrCannotMakeChange,
			count,
			denomination.Cents(),
			c.counts[denomination],
		)
	}

	c.counts[denomination] -= count
	return nil
}

func (c *CashBox) Snapshot() map[Denomination]int {
	snapshot := make(map[Denomination]int, len(c.counts))

	for denomination, count := range c.counts {
		snapshot[denomination] = count
	}

	return snapshot
}
