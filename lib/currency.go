package lib

var coinExchangeRates = []struct {
	value int
	name  string
}{
	{10_000_000_000, "Mithril"},
	{100_000_000, "Iridium"},
	{1_000_000, "Platinum"},
	{10_000, "Gold"},
	{100, "Silver"},
	{1, "Bronze"},
}

type Denomination struct {
	Value int
	Name  string
}

type Tender struct {
	amount int
}

func NewTender(v int) Tender {
	return Tender{v}
}

func (t *Tender) Add(v Tender) {
	if v.amount < 0 {
		panic("currency::adding negative values is not allowed")
	}
	t.amount += v.amount
}

func (t *Tender) Subtract(v Tender) {
	if v.amount > t.amount {
		panic("currency::subtracting more than balance not allowed")
	}
	t.amount -= v.amount
}

func (t Tender) Denominations() []Denomination {
	if t.amount == 0 {
		return []Denomination{}
	}

	coins := make([]Denomination, len(coinExchangeRates))
	remaining := t.amount

	for i, rate := range coinExchangeRates {
		if t.amount >= rate.value {
			coins[i] = Denomination{remaining / rate.value, rate.name}
			remaining %= rate.value
		} else {
			coins[i] = Denomination{0, rate.name}
		}
	}

	return coins
}
