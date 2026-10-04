//go:build postgres

package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresPersistence struct {
	db         *sql.DB
	normalized *normalizedRepository
}

func openPostgres(ctx context.Context) (*PostgresPersistence, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, nil
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS runtime_snapshots (id TEXT PRIMARY KEY, payload JSONB NOT NULL, version TEXT NOT NULL, updated_at TIMESTAMPTZ NOT NULL DEFAULT now())`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("create runtime_snapshots: %w", err)
	}
	return &PostgresPersistence{db: db, normalized: newNormalizedRepository(db)}, nil
}
func (p *PostgresPersistence) Close() error {
	if p == nil || p.db == nil {
		return nil
	}
	return p.db.Close()
}
func (p *PostgresPersistence) Load(ctx context.Context) (*RuntimeSnapshot, error) {
	if p == nil {
		return nil, nil
	}
	var payload []byte
	err := p.db.QueryRowContext(ctx, `SELECT payload FROM runtime_snapshots WHERE id='primary'`).Scan(&payload)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load runtime snapshot: %w", err)
	}
	var s RuntimeSnapshot
	if err := json.Unmarshal(payload, &s); err != nil {
		return nil, fmt.Errorf("decode runtime snapshot: %w", err)
	}
	return &s, nil
}
func (p *PostgresPersistence) Save(ctx context.Context, s *Store) error {
	if p == nil {
		return nil
	}
	payload, err := json.Marshal(RuntimeSnapshot{Transactions: s.Transactions, Payments: s.Payments, Documents: s.Documents, DocumentExtractions: s.DocumentExtractions, Events: s.Events, Declarations: s.Declarations, Invoices: s.Invoices, ShippingBills: s.ShippingBills, BillsOfEntry: s.BillsOfEntry, Requirements: s.Requirements, Counter: s.Counter})
	if err != nil {
		return err
	}
	_, err = p.db.ExecContext(ctx, `INSERT INTO runtime_snapshots(id,payload,version,updated_at) VALUES('primary',$1,'2.6',now()) ON CONFLICT(id) DO UPDATE SET payload=EXCLUDED.payload,version=EXCLUDED.version,updated_at=now()`, payload)
	return err
}

func (p *PostgresPersistence) LoadNormalized(ctx context.Context, s *Store) (bool, error) {
	if p == nil || p.normalized == nil {
		return false, nil
	}
	txs, err := p.normalized.tx.ListByCompany(ctx, normalizedCompanyID())
	if err != nil {
		return false, fmt.Errorf("list normalized transactions: %w", err)
	}
	if len(txs) == 0 {
		return false, nil
	}
	s.Transactions = map[string]*Transaction{}
	s.Payments = map[string][]Payment{}
	s.Invoices = map[string][]Invoice{}
	for _, rt := range txs {
		t := &Transaction{ID: rt.ID, ReferenceNo: rt.ReferenceNo, ProductCode: rt.ProductCode, ProductType: rt.ProductType, Status: rt.Status, Currency: rt.Currency, DeclaredValue: rt.DeclaredValue, Counterparty: rt.Counterparty, Country: rt.Country, Description: rt.Description, ContractValue: rt.ContractValue, InvoiceValue: rt.InvoiceValue, CreatedAt: time.Now()}
		s.Transactions[t.ID] = t
		invs, err := p.normalized.inv.ListByTransaction(ctx, t.ID)
		if err != nil {
			return false, err
		}
		for _, i := range invs {
			s.Invoices[t.ID] = append(s.Invoices[t.ID], Invoice{ID: i.ID, TransactionID: i.TransactionID, InvoiceNumber: i.InvoiceNumber, Currency: i.Currency, InvoiceValue: i.InvoiceValue, PaymentTerms: i.PaymentTerms, Description: i.Description, Status: i.Status})
		}
		pays, err := p.normalized.pay.ListByTransaction(ctx, t.ID)
		if err != nil {
			return false, err
		}
		for _, pay := range pays {
			s.Payments[t.ID] = append(s.Payments[t.ID], Payment{ID: pay.ID, TransactionID: pay.TransactionID, Reference: pay.Reference, PaymentType: pay.PaymentType, Currency: pay.Currency, PaymentDate: pay.PaymentDate, PurposeCode: pay.PurposeCode, BankReference: pay.BankReference, Amount: pay.Amount, Direction: pay.Direction})
		}
		reqs, err := p.normalized.req.ListByTransaction(ctx, t.ID)
		if err != nil {
			return false, err
		}
		for _, rq := range reqs {
			s.Requirements[t.ID] = append(s.Requirements[t.ID], Requirement{ID: rq.ID, TransactionID: rq.TransactionID, Category: rq.Category, Title: rq.Title, Description: rq.Description, Status: rq.Status, Priority: rq.Priority, Source: rq.Source, RequestedBy: rq.RequestedBy, AssignedTo: rq.AssignedTo, DueDate: rq.DueDate, Notes: rq.Notes})
		}
	}
	return true, nil
}

func (p *PostgresPersistence) SaveNormalized(ctx context.Context, s *Store) error {
	if p == nil || p.normalized == nil {
		return nil
	}
	// First persist transactions. Database UUIDs become the durable IDs used by child rows.
	for key, t := range s.Transactions {
		id := t.ID
		if !isUUID(id) {
			newID, err := p.normalized.SaveTransaction(ctx, t)
			if err != nil {
				if strings.Contains(err.Error(), "duplicate key") {
					continue
				}
				return err
			}
			if newID != "" {
				delete(s.Transactions, key)
				t.ID = newID
				s.Transactions[newID] = t
				if pays, ok := s.Payments[key]; ok {
					for i := range pays {
						pays[i].TransactionID = newID
					}
					delete(s.Payments, key)
					s.Payments[newID] = pays
				}
				if invs, ok := s.Invoices[key]; ok {
					for i := range invs {
						invs[i].TransactionID = newID
					}
					delete(s.Invoices, key)
					s.Invoices[newID] = invs
				}
				id = newID
			}
		}
		for _, i := range s.Invoices[id] {
			if err := p.normalized.SaveInvoice(ctx, i); err != nil && !strings.Contains(err.Error(), "duplicate key") {
				return err
			}
		}
		for _, pay := range s.Payments[id] {
			if err := p.normalized.SavePayment(ctx, pay); err != nil && !strings.Contains(err.Error(), "duplicate key") {
				return err
			}
		}
		// Persist the remaining Transaction 360 domains.
		for _, x := range s.ShippingBills[id] {
			if err := p.normalized.SaveShippingBill(ctx, x); err != nil && !strings.Contains(err.Error(), "duplicate key") {
				return err
			}
		}
		for _, x := range s.BillsOfEntry[id] {
			if err := p.normalized.SaveBOE(ctx, x); err != nil && !strings.Contains(err.Error(), "duplicate key") {
				return err
			}
		}
		for _, x := range s.Documents[id] {
			if err := p.normalized.SaveDocument(ctx, x); err != nil && !strings.Contains(err.Error(), "duplicate key") {
				return err
			}
		}
		for _, x := range s.Documents[id] {
			for _, ex := range s.DocumentExtractions[x.ID] {
				if err := p.normalized.SaveDocumentExtraction(ctx, ex); err != nil && !strings.Contains(err.Error(), "duplicate key") {
					return err
				}
			}
		}
		for _, x := range s.Events[id] {
			if err := p.normalized.SaveEvent(ctx, x); err != nil && !strings.Contains(err.Error(), "duplicate key") {
				return err
			}
		}
		for _, x := range s.Declarations[id] {
			if err := p.normalized.SaveDeclaration(ctx, x); err != nil && !strings.Contains(err.Error(), "duplicate key") {
				return err
			}
		}
		for _, x := range s.Requirements[id] {
			if x.ID == "SUGGESTED" {
				continue
			}
			if err := p.normalized.SaveRequirement(ctx, x); err != nil && !strings.Contains(err.Error(), "duplicate key") {
				return err
			}
		}
		rv := reconciliationFor(t, s.Payments[id])
		if err := p.normalized.SaveReconciliation(ctx, id, rv); err != nil {
			return err
		}
		for _, a := range actionsFor(t, s.Payments[id]) {
			if err := p.normalized.SaveAction(ctx, id, a); err != nil && !strings.Contains(err.Error(), "duplicate key") {
				return err
			}
		}
		if err := p.normalized.SaveRegulatory(ctx, id, "CUSTOMER_POSITION", "REVIEW", "CUSTOMER_RECORDED", "2026.1", time.Now().Format("2006-01-02"), "Customer-side regulatory position; verify actual bank/authority status where applicable."); err != nil {
			return err
		}
		if err := p.normalized.SaveAudit(ctx, id, "PERSISTENCE_SYNC", map[string]any{"version": "2.10", "domains": "transaction-360"}); err != nil {
			return err
		}
	}
	return nil
}

func isUUID(v string) bool {
	if len(v) != 36 {
		return false
	}
	for i, c := range v {
		if i == 8 || i == 13 || i == 18 || i == 23 {
			if c != '-' {
				return false
			}
			continue
		}
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}
