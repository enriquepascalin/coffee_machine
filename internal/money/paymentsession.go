package money

import "fmt"

type PaymentSession struct {
	counts map[Denomination]int
}

func NewPaymentSession() *PaymentSession {
	return &PaymentSession{
		counts: make(map[Denomination]int),
	}
}

func (p *PaymentSession) Insert(denomination Denomination, count int) error {
	if !isAcceptedDenomination(denomination) {
		return fmt.Errorf("%w: %v", ErrInvalidDenomination, denomination.Cents())
	}

	if count <= 0 {
		return fmt.Errorf("invalid count: %d", count)
	}

	p.counts[denomination] += count
	return nil
}

func (p *PaymentSession) Total() Cents {
	var total Cents

	for denomination, count := range p.counts {
		total += denomination.Cents() * Cents(count)
	}

	return total
}

func (p *PaymentSession) Snapshot() map[Denomination]int {
	cloned := make(map[Denomination]int, len(p.counts))

	for denomination, count := range p.counts {
		cloned[denomination] = count
	}

	return cloned
}

func (p *PaymentSession) Commit(cashBox *CashBox) error {
	for denomination, count := range p.counts {
		if err := cashBox.Add(denomination, count); err != nil {
			return err
		}
	}

	p.Reset()
	return nil
}

func (p *PaymentSession) Reset() {
	p.counts = make(map[Denomination]int)
}

func isAcceptedDenomination(denomination Denomination) bool {
	for _, accepted := range AcceptedDenominationsDesc() {
		if denomination == accepted {
			return true
		}
	}

	return false
}
