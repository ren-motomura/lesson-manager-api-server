package models

import (
	"time"

	"github.com/go-gorp/gorp"
)

type CashflowDirection int8

const (
	CashflowDirectionInbound CashflowDirection = iota
	CashflowDirectionOutbound
)

type CashflowLog struct {
	ID          int
	Amount      int
	CustomerID  int
	Direction   CashflowDirection
	CreatedAt   time.Time
	Description string
}

func registerCashflowLog(dbMap *gorp.DbMap) {
	t := dbMap.AddTableWithName(CashflowLog{}, "cashflow_logs").SetKeys(true, "ID")
	t.ColMap("ID").Rename("id")
	t.ColMap("Amount").Rename("amount")
	t.ColMap("CustomerID").Rename("customer_id")
	t.ColMap("Direction").Rename("direction")
	t.ColMap("CreatedAt").Rename("created_at")
	t.ColMap("Description").Rename("description")
}
