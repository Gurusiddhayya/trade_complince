package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"trade-compliance-cockpit/backend/internal/repository"
)

type normalizedRepository struct {
	tx     repository.TransactionRepository
	inv    repository.InvoiceRepository
	pay    repository.PaymentRepository
	ship   repository.ShippingBillRepository
	boe    repository.BillOfEntryRepository
	doc    repository.DocumentRepository
	evt    repository.EventRepository
	dec    repository.DeclarationRepository
	action repository.ActionRepository
	recon  repository.ReconciliationRepository
	reg    repository.RegulatoryPositionRepository
	audit  repository.AuditRepository
	req    repository.RequirementRepository
}

func newNormalizedRepository(db *sql.DB) *normalizedRepository {
	if db == nil {
		return nil
	}
	p := repository.NewPostgres(db)
	return &normalizedRepository{
		tx: repository.TransactionStore{Postgres: p}, inv: repository.InvoiceStore{Postgres: p}, pay: repository.PaymentStore{Postgres: p},
		ship: repository.ShippingBillStore{Postgres: p}, boe: repository.BillOfEntryStore{Postgres: p}, doc: repository.DocumentStore{Postgres: p},
		evt: repository.EventStore{Postgres: p}, dec: repository.DeclarationStore{Postgres: p}, action: repository.ActionStore{Postgres: p},
		recon: repository.ReconciliationStore{Postgres: p}, reg: repository.RegulatoryStore{Postgres: p}, audit: repository.AuditStore{Postgres: p}, req: repository.RequirementStore{Postgres: p},
	}
}
func normalizedCompanyID() string {
	if v := os.Getenv("TCC_COMPANY_ID"); v != "" {
		return v
	}
	return "00000000-0000-0000-0000-000000000001"
}
func (n *normalizedRepository) SaveTransaction(ctx context.Context, t *Transaction) (string, error) {
	if n == nil || n.tx == nil {
		return "", nil
	}
	if db, ok := n.tx.(interface {
		EnsureCompany(context.Context, string) error
	}); ok {
		if err := db.EnsureCompany(ctx, normalizedCompanyID()); err != nil {
			return "", fmt.Errorf("ensure company: %w", err)
		}
	}
	id, err := n.tx.Create(ctx, repository.Transaction{ID: t.ID, CompanyID: normalizedCompanyID(), ReferenceNo: t.ReferenceNo, ProductCode: t.ProductCode, ProductType: t.ProductType, Status: t.Status, Currency: t.Currency, Counterparty: t.Counterparty, Country: t.Country, Description: t.Description, DeclaredValue: t.DeclaredValue, ContractValue: t.ContractValue, InvoiceValue: t.InvoiceValue})
	if err != nil {
		return "", fmt.Errorf("save transaction: %w", err)
	}
	return id, nil
}
func (n *normalizedRepository) SaveInvoice(c context.Context, i Invoice) error {
	if n == nil || n.inv == nil {
		return nil
	}
	_, e := n.inv.Create(c, repository.Invoice{ID: i.ID, TransactionID: i.TransactionID, InvoiceNumber: i.InvoiceNumber, Currency: i.Currency, InvoiceValue: i.InvoiceValue, PaymentTerms: i.PaymentTerms, Description: i.Description, Status: i.Status})
	return e
}
func (n *normalizedRepository) SavePayment(c context.Context, p Payment) error {
	if n == nil || n.pay == nil {
		return nil
	}
	_, e := n.pay.Create(c, repository.Payment{ID: p.ID, TransactionID: p.TransactionID, Reference: p.Reference, PaymentType: p.PaymentType, Currency: p.Currency, PaymentDate: p.PaymentDate, PurposeCode: p.PurposeCode, BankReference: p.BankReference, Amount: p.Amount, Direction: p.Direction})
	return e
}
func (n *normalizedRepository) SaveShippingBill(c context.Context, x ShippingBill) error {
	if n == nil || n.ship == nil {
		return nil
	}
	_, e := n.ship.Create(c, repository.ShippingBill{ID: x.ID, TransactionID: x.TransactionID, Number: x.Number, Date: x.Date, Port: x.Port, ExportDate: x.ExportDate, Currency: x.Currency, ExportValue: x.ExportValue})
	return e
}
func (n *normalizedRepository) SaveBOE(c context.Context, x BillOfEntry) error {
	if n == nil || n.boe == nil {
		return nil
	}
	_, e := n.boe.Create(c, repository.BillOfEntry{ID: x.ID, TransactionID: x.TransactionID, Number: x.Number, Date: x.Date, Port: x.Port, Currency: x.Currency, AssessableValue: x.AssessableValue})
	return e
}
func (n *normalizedRepository) SaveDocument(c context.Context, x Document) error {
	if n == nil || n.doc == nil {
		return nil
	}
	_, e := n.doc.Create(c, repository.Document{ID: x.ID, TransactionID: x.TransactionID, DocumentType: x.Type, FileName: x.FileName, Status: x.Status, Version: x.Version, StorageKey: x.StorageKey, MIMEType: x.MIMEType, FileSize: x.FileSize, SHA256: x.SHA256, UploadedBy: x.UploadedBy, ParentDocumentID: x.ParentDocumentID})
	return e
}
func (n *normalizedRepository) SaveDocumentExtraction(c context.Context, x DocumentExtraction) error {
	if n == nil || n.doc == nil {
		return nil
	}
	_, e := n.doc.CreateExtraction(c, repository.DocumentExtraction{ID: x.ID, DocumentID: x.DocumentID, FieldName: x.FieldName, FieldValue: x.FieldValue, Confidence: x.Confidence, Source: x.Source, Confirmed: x.Confirmed, ConfirmedBy: x.ConfirmedBy})
	return e
}
func (n *normalizedRepository) SaveEvent(c context.Context, x Event) error {
	if n == nil || n.evt == nil {
		return nil
	}
	_, e := n.evt.Create(c, repository.Event{ID: x.ID, TransactionID: x.TransactionID, EventType: x.EventType, EventDate: x.EventDate, Notes: x.Notes})
	return e
}
func (n *normalizedRepository) SaveDeclaration(c context.Context, x Declaration) error {
	if n == nil || n.dec == nil {
		return nil
	}
	_, e := n.dec.Create(c, repository.Declaration{ID: x.ID, TransactionID: x.TransactionID, DeclarationType: x.Type, Version: x.Version, Status: x.Status, Fields: x.Fields})
	return e
}
func (n *normalizedRepository) SaveAction(c context.Context, txID string, x Action) error {
	if n == nil || n.action == nil {
		return nil
	}
	_, e := n.action.Create(c, repository.Action{TransactionID: txID, Priority: x.Priority, Title: x.Title, Reason: x.Reason, RecommendedAction: x.RecommendedAction, Status: x.Status, DueDate: x.DueDate})
	return e
}
func (n *normalizedRepository) SaveReconciliation(c context.Context, txID string, x ReconciliationView) error {
	if n == nil || n.recon == nil {
		return nil
	}
	return n.recon.Upsert(c, repository.ReconciliationPosition{TransactionID: txID, Currency: x.Currency, BaseValue: x.BaseValue, InwardValue: x.InwardValue, OutwardValue: x.OutwardValue, RecordedValue: x.RecordedValue, Outstanding: x.Outstanding, Difference: x.Difference, Status: x.Status, Notes: x.Notes})
}
func (n *normalizedRepository) SaveRegulatory(c context.Context, txID, registry, status, source, version, asOf, notes string) error {
	if n == nil || n.reg == nil {
		return nil
	}
	return n.reg.Upsert(c, repository.RegulatoryPosition{TransactionID: txID, Registry: registry, PositionStatus: status, SourceType: source, EffectiveRuleVersion: version, AsOfDate: asOf, Notes: notes})
}
func (n *normalizedRepository) SaveAudit(c context.Context, txID, eventType string, payload map[string]any) error {
	if n == nil || n.audit == nil {
		return nil
	}
	_, e := n.audit.Create(c, repository.AuditEvent{CompanyID: normalizedCompanyID(), TransactionID: txID, EventType: eventType, ActorType: "system", ActorID: "system", Payload: payload})
	return e
}

func (n *normalizedRepository) SaveRequirement(c context.Context, x Requirement) error {
	if n == nil || n.req == nil {
		return nil
	}
	_, e := n.req.Create(c, repository.Requirement{ID: x.ID, TransactionID: x.TransactionID, CompanyID: normalizedCompanyID(), Category: x.Category, Title: x.Title, Description: x.Description, Status: x.Status, Priority: x.Priority, Source: x.Source, RequestedBy: x.RequestedBy, AssignedTo: x.AssignedTo, DueDate: x.DueDate, Notes: x.Notes})
	return e
}
