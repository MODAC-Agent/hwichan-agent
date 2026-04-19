package pricing

type Price struct {
	Cents int64
}

func Sum(a, b Price) Price {
	return Price{Cents: a.Cents + b.Cents}
}
