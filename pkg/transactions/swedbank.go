package transactions

import (
	"fmt"
	"statements/pkg/ctime"
	"strconv"
	"strings"
	"time"
)

type SwedbankEntryType string

const (
	SwedbankEntryStartBalance    SwedbankEntryType = "10"
	SwedbankEntryTransaction     SwedbankEntryType = "20"
	SwedbankEntryTurnover        SwedbankEntryType = "82"
	SwedbankEntryEndBalance      SwedbankEntryType = "86"
	SwedbankEntryCurrentInterest SwedbankEntryType = "900"
)

type SwedbankTransactionType string

const (
	SwedbankTransactionStartBalance   SwedbankTransactionType = "AS"
	SwedbankTransactionToBank         SwedbankTransactionType = "INB"
	SwedbankTransactionToPrivate      SwedbankTransactionType = "PRV"
	SwedbankTransactionCapitalGains   SwedbankTransactionType = "AIA"
	SwedbankTransactionCurrentIntrest SwedbankTransactionType = "AI"
	SwedbankTransactionCommission     SwedbankTransactionType = "KOM"
	SwedbankTransactionTurnover       SwedbankTransactionType = "K2"
	SwedbankTransactionEndBalance     SwedbankTransactionType = "LS"
)

type SwedbankFlow string

const (
	SwedbankDebit  SwedbankFlow = "D"
	SwedbankCredit SwedbankFlow = "K"
)

type SwedbankTransaction struct {
	AccountNumber   string
	EntryType       SwedbankEntryType
	Date            time.Time
	AccountHolder   string
	Description     string
	Value           uint // Stored as decimal in CSV
	Currency        string
	Flow            SwedbankFlow
	ArchiveCode     string
	TransactionType SwedbankTransactionType
	ReferenceNumber string
	DocumentNumber  string
}

func NewSwedbankTransaction(row []string) (SwedbankTransaction, error) {
	var t SwedbankTransaction

	entryType, err := parseEntryType(row[1])
	if err != nil {
		return t, err
	}

	transactionType, err := parseTransactionType(row[9])
	if err != nil {
		return t, err
	}

	flow, err := parseFlow(row[7])
	if err != nil {
		return t, err
	}

	date, err := time.Parse(ctime.LittleEndianDateOnly, row[2])
	if err != nil {
		return t, err
	}

	vp := strings.Split(row[5], ",")
	w, err := strconv.ParseUint(vp[0], 10, 64)
	if err != nil {
		return t, err
	}
	d, err := strconv.ParseUint(vp[1], 10, 64)
	if err != nil {
		return t, err
	}
	value := w*100 + d

	t = SwedbankTransaction{
		AccountNumber:   row[0],
		EntryType:       entryType,
		Date:            date,
		AccountHolder:   row[3],
		Description:     row[4],
		Value:           uint(value),
		Currency:        row[6],
		Flow:            flow,
		ArchiveCode:     row[8],
		TransactionType: transactionType,
		ReferenceNumber: row[10],
		DocumentNumber:  row[11],
	}

	return t, nil
}

// Parses an entry type from a raw string into an enum.
func parseEntryType(t string) (SwedbankEntryType, error) {
	switch t {
	case "10":
		fallthrough
	case "20":
		fallthrough
	case "82":
		fallthrough
	case "86":
		fallthrough
	case "900":
		return SwedbankEntryType(t), nil
	default:
		return "", fmt.Errorf("invalid Swedbank entry type provided: %s", t)
	}
}

// Parses a transaction type from a raw string into an enum.
func parseTransactionType(t string) (SwedbankTransactionType, error) {
	switch t {
	case "AS":
		fallthrough
	case "INB":
		fallthrough
	case "PRV":
		fallthrough
	case "AIA":
		fallthrough
	case "AI":
		fallthrough
	case "KOM":
		fallthrough
	case "K2":
		fallthrough
	case "LS":
		return SwedbankTransactionType(t), nil
	default:
		return "", fmt.Errorf("invalid Swedbank transaction type provided: %s", t)
	}
}

// Parses cash flow from a raw string into an enum.
func parseFlow(t string) (SwedbankFlow, error) {
	switch t {
	case "D":
		fallthrough
	case "K":
		return SwedbankFlow(t), nil
	default:
		return "", fmt.Errorf("invalid Swedbank flow provided: %s", t)
	}
}
