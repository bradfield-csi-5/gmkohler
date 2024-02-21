package iter

// limitIterator limits the number of output rows.
type limitIterator struct {
	source Iterator
	limit  int
	sent   int
}

func NewLimitIterator(source Iterator, limit int) Iterator {
	// TODO: Consider (Iterator, err) return signature, and fail when limit < 0
	return &limitIterator{
		source: source,
		limit:  limit,
		sent:   0,
	}
}

func (l *limitIterator) Init() {
	l.source.Init()
}

func (l *limitIterator) Next() *Tuple {
	if l.sent >= l.limit {
		return nil
	}

	var next = l.source.Next()
	if next == nil {
		return nil
	}

	l.sent++
	return next
}

func (l *limitIterator) Close() {
	l.source.Close()
}
