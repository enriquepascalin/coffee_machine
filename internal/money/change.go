package money

import "fmt"

func MakeChange(amount Cents, available map[Denomination]int) (map[Denomination]int, error) {
	if amount < 0 {
		return nil, fmt.Errorf("negative change: %v", amount)
	}
	if amount == 0 {
		return map[Denomination]int{}, nil
	}

	denoms := AcceptedDenominationsDesc()
	best, ok := makeChangeRec(int(amount), denoms, available, 0)
	if !ok {
		return nil, ErrCannotMakeChange
	}
	return best, nil
}

func makeChangeRec(amount int, denoms []Denomination, available map[Denomination]int, idx int) (map[Denomination]int, bool) {
	if amount == 0 {
		return map[Denomination]int{}, true
	}
	if idx >= len(denoms) {
		return nil, false
	}

	d := denoms[idx]
	value := int(d)
	maxUse := available[d]
	if maxUse > amount/value {
		maxUse = amount / value
	}

	for use := maxUse; use >= 0; use-- {
		remaining := amount - use*value
		nextAvail := available
		if use > 0 {
			nextAvail = cloneCounts(available)
			nextAvail[d] -= use
		}

		sub, ok := makeChangeRec(remaining, denoms, nextAvail, idx+1)
		if !ok {
			continue
		}

		if use > 0 {
			sub[d] = use
		}
		return sub, true
	}

	return nil, false
}

func cloneCounts(m map[Denomination]int) map[Denomination]int {
	c := make(map[Denomination]int, len(m))
	for k, v := range m {
		c[k] = v
	}
	return c
}
