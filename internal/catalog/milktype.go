package catalog

import "fmt"

type MilkType int

const (
	MilkNone MilkType = iota
	MilkWhole
	MilkAlmond
	MilkOat
)

func (m MilkType) String() string {
	switch m {
	case MilkNone:
		return "none"
	case MilkWhole:
		return "whole"
	case MilkAlmond:
		return "almond"
	case MilkOat:
		return "oat"
	default:
		return fmt.Sprintf("MilkType(%d)", int(m))
	}
}
