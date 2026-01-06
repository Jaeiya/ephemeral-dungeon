package lib

type CoinRate struct {
	value int
	name  string
}

type CoinAmount struct {
	Value int
	Name  string
}

var coinExchangeRates = []CoinRate{
	{10_000_000_000, "Mithril"},
	{100_000_000, "Iridium"},
	{1_000_000, "Platinum"},
	{10_000, "Gold"},
	{100, "Silver"},
	{1, "Bronze"},
}

var coinMap = map[int]string{}

type Currency struct {
	balance int
}

func NewCurrency(v int) Currency {
	return Currency{v}
}

func (c *Currency) Add(v int) {
	if v < 0 {
		panic("currency::adding negative values is not allowed")
	}
	c.balance += v
}

func (c *Currency) Subtract(v int) {
	if v > c.balance {
		panic("currency::subtracting more than balance not allowed")
	}
	c.balance -= v
}

func (c Currency) Balance() int {
	return c.balance
}

func (c Currency) Amount() []CoinAmount {
	if c.balance == 0 {
		return []CoinAmount{}
	}

	coins := make([]CoinAmount, len(coinExchangeRates))
	remaining := c.balance

	for i, rate := range coinExchangeRates {
		if c.balance >= rate.value {
			coins[i] = CoinAmount{remaining / rate.value, rate.name}
			remaining %= rate.value
		} else {
			coins[i] = CoinAmount{0, rate.name}
		}
	}

	return coins
}
