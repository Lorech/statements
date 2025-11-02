package transactions

import (
	"fmt"
	"statements/pkg/ctime"
	"time"
)

type Transaction struct {
	Date          time.Time
	AccountHolder string
	Description   string
	Value         int
	Currency      string
}

func (t Transaction) Csv() []string {
	return []string{
		t.Date.Format(ctime.LittleEndianDateOnly),
		t.AccountHolder,
		t.Description,
		fmt.Sprintf("%d,%d", t.Value/100, max(-t.Value%100, t.Value%100)),
		t.Currency,
	}
}
