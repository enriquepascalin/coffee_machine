package money

type Denomination Cents

const (
	Denom5C  Denomination = 5
	Denom10C Denomination = 10
	Denom25C Denomination = 25
	Denom50C Denomination = 50
	Denom1D  Denomination = 100
	Denom2D  Denomination = 200
	Denom5D  Denomination = 500
	Denom10D Denomination = 1000
	Denom20D Denomination = 2000
)

func (d Denomination) Cents() Cents { return Cents(d) }

func AcceptedDenominationsDesc() []Denomination {
	return []Denomination{
		Denom20D, Denom10D, Denom5D, Denom2D, Denom1D,
		Denom50C, Denom25C, Denom10C, Denom5C,
	}
}
