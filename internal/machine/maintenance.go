package machine

type Maintenance struct {
	gaugePercent int
}

func NewMaintenance(gaugePercent int) Maintenance {
	if gaugePercent < 0 {
		gaugePercent = 0
	}

	return Maintenance{
		gaugePercent: gaugePercent,
	}
}

func (m *Maintenance) GaugePercent() int {
	return m.gaugePercent
}

func (m *Maintenance) Increase(delta int) {
	m.gaugePercent += delta
}

func (m *Maintenance) NeedsService() bool {
	return m.gaugePercent >= 100
}

func (m *Maintenance) Reset() {
	m.gaugePercent = 0
}
