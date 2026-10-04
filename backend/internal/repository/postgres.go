package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

type Postgres struct{ DB *sql.DB }

func NewPostgres(db *sql.DB) *Postgres { return &Postgres{DB: db} }

func (r *Postgres) CreateTransaction(ctx context.Context, t Transaction) (string, error) {
	var id string
	err := r.DB.QueryRowContext(ctx, `INSERT INTO transactions (company_id, reference_no, product_code, product_type, status, currency, declared_value, counterparty, country, description, contract_value, invoice_value) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING id::text`, t.CompanyID, t.ReferenceNo, t.ProductCode, t.ProductType, t.Status, t.Currency, t.DeclaredValue, t.Counterparty, t.Country, t.Description, t.ContractValue, t.InvoiceValue).Scan(&id)
	return id, err
}

func (r *Postgres) GetTransaction(ctx context.Context, id string) (Transaction, error) {
	var t Transaction
	err := r.DB.QueryRowContext(ctx, `SELECT id::text, company_id::text, reference_no, product_code, product_type, status, currency, counterparty, country, description, declared_value, contract_value, invoice_value FROM transactions WHERE id=$1`, id).Scan(&t.ID, &t.CompanyID, &t.ReferenceNo, &t.ProductCode, &t.ProductType, &t.Status, &t.Currency, &t.Counterparty, &t.Country, &t.Description, &t.DeclaredValue, &t.ContractValue, &t.InvoiceValue)
	if err != nil {
		return t, fmt.Errorf("get transaction: %w", err)
	}
	return t, nil
}

func (r *Postgres) ListTransactions(ctx context.Context, companyID string) ([]Transaction, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id::text, company_id::text, reference_no, product_code, product_type, status, currency, counterparty, country, description, declared_value, contract_value, invoice_value FROM transactions WHERE company_id=$1 ORDER BY created_at DESC`, companyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Transaction{}
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.CompanyID, &t.ReferenceNo, &t.ProductCode, &t.ProductType, &t.Status, &t.Currency, &t.Counterparty, &t.Country, &t.Description, &t.DeclaredValue, &t.ContractValue, &t.InvoiceValue); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *Postgres) CreateInvoice(ctx context.Context, i Invoice) (string, error) {
	var id string
	err := r.DB.QueryRowContext(ctx, `INSERT INTO invoices (transaction_id, invoice_number, invoice_date, currency, invoice_value, payment_terms, description, status) VALUES ($1,$2,NULLIF($3,'')::date,$4,$5,$6,$7,$8) RETURNING id::text`, i.TransactionID, i.InvoiceNumber, "", i.Currency, i.InvoiceValue, i.PaymentTerms, i.Description, i.Status).Scan(&id)
	return id, err
}

func (r *Postgres) ListInvoices(ctx context.Context, transactionID string) ([]Invoice, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id::text, transaction_id::text, invoice_number, COALESCE(currency,''), invoice_value, COALESCE(payment_terms,''), COALESCE(description,''), COALESCE(status,'') FROM invoices WHERE transaction_id=$1 ORDER BY created_at`, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Invoice{}
	for rows.Next() {
		var i Invoice
		if err := rows.Scan(&i.ID, &i.TransactionID, &i.InvoiceNumber, &i.Currency, &i.InvoiceValue, &i.PaymentTerms, &i.Description, &i.Status); err != nil {
			return nil, err
		}
		out = append(out, i)
	}
	return out, rows.Err()
}

func (r *Postgres) CreatePayment(ctx context.Context, p Payment) (string, error) {
	var id string
	err := r.DB.QueryRowContext(ctx, `INSERT INTO payments (transaction_id, reference, payment_type, currency, payment_date, purpose_code, bank_reference, amount, direction) VALUES ($1,$2,$3,$4,NULLIF($5,'')::date,$6,$7,$8,$9) RETURNING id::text`, p.TransactionID, p.Reference, p.PaymentType, p.Currency, p.PaymentDate, p.PurposeCode, p.BankReference, p.Amount, p.Direction).Scan(&id)
	return id, err
}
func (r *Postgres) ListPayments(ctx context.Context, transactionID string) ([]Payment, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id::text,transaction_id::text,reference,payment_type,COALESCE(currency,''),COALESCE(payment_date::text,''),COALESCE(purpose_code,''),COALESCE(bank_reference,''),amount,direction FROM payments WHERE transaction_id=$1 ORDER BY created_at`, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Payment{}
	for rows.Next() {
		var p Payment
		if err := rows.Scan(&p.ID, &p.TransactionID, &p.Reference, &p.PaymentType, &p.Currency, &p.PaymentDate, &p.PurposeCode, &p.BankReference, &p.Amount, &p.Direction); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *Postgres) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.DB.ExecContext(ctx, `UPDATE transactions SET status=$2, updated_at=now() WHERE id=$1`, id, status)
	return err
}

func (r *Postgres) EnsureCompany(ctx context.Context, companyID string) error {
	_, err := r.DB.ExecContext(ctx, `INSERT INTO companies(id,name) VALUES($1,$2) ON CONFLICT(id) DO NOTHING`, companyID, "Customer Company")
	return err
}

type TransactionStore struct{ *Postgres }

func (r TransactionStore) Get(ctx context.Context, id string) (Transaction, error) {
	return r.GetTransaction(ctx, id)
}
func (r TransactionStore) ListByCompany(ctx context.Context, companyID string) ([]Transaction, error) {
	return r.ListTransactions(ctx, companyID)
}
func (r TransactionStore) Create(ctx context.Context, t Transaction) (string, error) {
	return r.CreateTransaction(ctx, t)
}
func (r TransactionStore) UpdateStatus(ctx context.Context, id, status string) error {
	return r.Postgres.UpdateStatus(ctx, id, status)
}
func (r TransactionStore) EnsureCompany(ctx context.Context, id string) error {
	return r.Postgres.EnsureCompany(ctx, id)
}

type InvoiceStore struct{ *Postgres }

func (r InvoiceStore) Create(ctx context.Context, i Invoice) (string, error) {
	return r.CreateInvoice(ctx, i)
}
func (r InvoiceStore) ListByTransaction(ctx context.Context, id string) ([]Invoice, error) {
	return r.ListInvoices(ctx, id)
}

type PaymentStore struct{ *Postgres }

func (r PaymentStore) Create(ctx context.Context, p Payment) (string, error) {
	return r.CreatePayment(ctx, p)
}
func (r PaymentStore) ListByTransaction(ctx context.Context, id string) ([]Payment, error) {
	return r.ListPayments(ctx, id)
}

// Remaining Transaction 360 repositories.
func (r *Postgres) CreateShippingBill(ctx context.Context, s ShippingBill) (string, error) {
	var id string
	err := r.DB.QueryRowContext(ctx, `INSERT INTO shipping_bills(transaction_id,shipping_bill_number,shipping_bill_date,port,export_date,currency,export_value) VALUES($1,$2,NULLIF($3,'')::date,$4,NULLIF($5,'')::date,$6,$7) ON CONFLICT(transaction_id,shipping_bill_number) DO UPDATE SET shipping_bill_date=EXCLUDED.shipping_bill_date,port=EXCLUDED.port,export_date=EXCLUDED.export_date,currency=EXCLUDED.currency,export_value=EXCLUDED.export_value RETURNING id::text`, s.TransactionID, s.Number, s.Date, s.Port, s.ExportDate, s.Currency, s.ExportValue).Scan(&id)
	return id, err
}
func (r *Postgres) ListShippingBills(ctx context.Context, tx string) ([]ShippingBill, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id::text,transaction_id::text,shipping_bill_number,COALESCE(shipping_bill_date::text,''),COALESCE(port,''),COALESCE(export_date::text,''),COALESCE(currency,''),export_value FROM shipping_bills WHERE transaction_id=$1 ORDER BY created_at DESC`, tx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []ShippingBill{}
	for rows.Next() {
		var x ShippingBill
		if err := rows.Scan(&x.ID, &x.TransactionID, &x.Number, &x.Date, &x.Port, &x.ExportDate, &x.Currency, &x.ExportValue); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Postgres) CreateBOE(ctx context.Context, b BillOfEntry) (string, error) {
	var id string
	err := r.DB.QueryRowContext(ctx, `INSERT INTO bills_of_entry(transaction_id,boe_number,boe_date,port,currency,assessable_value) VALUES($1,$2,NULLIF($3,'')::date,$4,$5,$6) ON CONFLICT(transaction_id,boe_number) DO UPDATE SET boe_date=EXCLUDED.boe_date,port=EXCLUDED.port,currency=EXCLUDED.currency,assessable_value=EXCLUDED.assessable_value RETURNING id::text`, b.TransactionID, b.Number, b.Date, b.Port, b.Currency, b.AssessableValue).Scan(&id)
	return id, err
}
func (r *Postgres) ListBOE(ctx context.Context, tx string) ([]BillOfEntry, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id::text,transaction_id::text,boe_number,COALESCE(boe_date::text,''),COALESCE(port,''),COALESCE(currency,''),assessable_value FROM bills_of_entry WHERE transaction_id=$1 ORDER BY created_at DESC`, tx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []BillOfEntry{}
	for rows.Next() {
		var x BillOfEntry
		if err := rows.Scan(&x.ID, &x.TransactionID, &x.Number, &x.Date, &x.Port, &x.Currency, &x.AssessableValue); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Postgres) CreateDocument(ctx context.Context, d Document) (string, error) {
	var id string
	err := r.DB.QueryRowContext(ctx, `INSERT INTO documents(transaction_id,document_type,file_name,status,version,source,storage_key,mime_type,file_size,sha256,uploaded_by,parent_document_id) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NULLIF($12,'')::uuid) RETURNING id::text`, d.TransactionID, d.DocumentType, d.FileName, d.Status, d.Version, d.Source, d.StorageKey, d.MIMEType, d.FileSize, d.SHA256, d.UploadedBy, d.ParentDocumentID).Scan(&id)
	return id, err
}
func (r *Postgres) ListDocuments(ctx context.Context, tx string) ([]Document, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id::text,transaction_id::text,document_type,COALESCE(file_name,''),COALESCE(status,''),version,COALESCE(source,''),COALESCE(storage_key,''),COALESCE(mime_type,''),file_size,COALESCE(sha256,''),COALESCE(uploaded_by,''),COALESCE(parent_document_id::text,'') FROM documents WHERE transaction_id=$1 ORDER BY created_at DESC`, tx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Document{}
	for rows.Next() {
		var x Document
		if err := rows.Scan(&x.ID, &x.TransactionID, &x.DocumentType, &x.FileName, &x.Status, &x.Version, &x.Source, &x.StorageKey, &x.MIMEType, &x.FileSize, &x.SHA256, &x.UploadedBy, &x.ParentDocumentID); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Postgres) CreateDocumentExtraction(ctx context.Context, x DocumentExtraction) (string, error) {
	var id string
	err := r.DB.QueryRowContext(ctx, `INSERT INTO document_extractions(document_id,field_name,field_value,confidence,source,confirmed,confirmed_by,confirmed_at) VALUES($1,$2,$3,$4,$5,$6,$7,CASE WHEN $6 THEN now() ELSE NULL END) ON CONFLICT(document_id,field_name) DO UPDATE SET field_value=EXCLUDED.field_value,confidence=EXCLUDED.confidence,source=EXCLUDED.source,confirmed=EXCLUDED.confirmed,confirmed_by=EXCLUDED.confirmed_by,confirmed_at=EXCLUDED.confirmed_at RETURNING id::text`, x.DocumentID, x.FieldName, x.FieldValue, x.Confidence, x.Source, x.Confirmed, x.ConfirmedBy).Scan(&id)
	return id, err
}
func (r *Postgres) ListDocumentExtractions(ctx context.Context, documentID string) ([]DocumentExtraction, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id::text,document_id::text,field_name,COALESCE(field_value,''),COALESCE(confidence,0),COALESCE(source,''),confirmed,COALESCE(confirmed_by,'') FROM document_extractions WHERE document_id=$1 ORDER BY field_name`, documentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []DocumentExtraction{}
	for rows.Next() {
		var x DocumentExtraction
		if err := rows.Scan(&x.ID, &x.DocumentID, &x.FieldName, &x.FieldValue, &x.Confidence, &x.Source, &x.Confirmed, &x.ConfirmedBy); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}

func (r *Postgres) CreateEvent(ctx context.Context, e Event) (string, error) {
	var id string
	err := r.DB.QueryRowContext(ctx, `INSERT INTO transaction_events(transaction_id,event_type,event_date,notes) VALUES($1,NULLIF($2,'')::date,$3,$4) RETURNING id::text`, e.TransactionID, e.EventDate, e.EventType, e.Notes).Scan(&id)
	return id, err
}
func (r *Postgres) ListEvents(ctx context.Context, tx string) ([]Event, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id::text,transaction_id::text,event_type,COALESCE(event_date::text,''),COALESCE(notes,'') FROM transaction_events WHERE transaction_id=$1 ORDER BY event_date,created_at`, tx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Event{}
	for rows.Next() {
		var x Event
		if err := rows.Scan(&x.ID, &x.TransactionID, &x.EventType, &x.EventDate, &x.Notes); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Postgres) CreateDeclaration(ctx context.Context, d Declaration) (string, error) {
	b, _ := json.Marshal(d.Fields)
	var id string
	err := r.DB.QueryRowContext(ctx, `INSERT INTO declarations(transaction_id,declaration_type,version,status,fields) VALUES($1,$2,$3,$4,$5) RETURNING id::text`, d.TransactionID, d.DeclarationType, d.Version, d.Status, b).Scan(&id)
	return id, err
}
func (r *Postgres) ListDeclarations(ctx context.Context, tx string) ([]Declaration, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id::text,transaction_id::text,declaration_type,version,status,fields FROM declarations WHERE transaction_id=$1 ORDER BY updated_at DESC`, tx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Declaration{}
	for rows.Next() {
		var x Declaration
		var b []byte
		if err := rows.Scan(&x.ID, &x.TransactionID, &x.DeclarationType, &x.Version, &x.Status, &b); err != nil {
			return nil, err
		}
		x.Fields = map[string]string{}
		_ = json.Unmarshal(b, &x.Fields)
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Postgres) CreateAction(ctx context.Context, a Action) (string, error) {
	var id string
	err := r.DB.QueryRowContext(ctx, `INSERT INTO actions(transaction_id,priority,title,reason,recommended_action,status,due_date) VALUES($1,$2,$3,$4,$5,$6,NULLIF($7,'')::date) ON CONFLICT(transaction_id,title,status) DO UPDATE SET priority=EXCLUDED.priority,reason=EXCLUDED.reason,recommended_action=EXCLUDED.recommended_action,due_date=EXCLUDED.due_date,updated_at=now() RETURNING id::text`, a.TransactionID, a.Priority, a.Title, a.Reason, a.RecommendedAction, a.Status, a.DueDate).Scan(&id)
	return id, err
}
func (r *Postgres) ListActions(ctx context.Context, tx string) ([]Action, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id::text,transaction_id::text,priority,title,COALESCE(reason,''),COALESCE(recommended_action,''),status,COALESCE(due_date::text,'') FROM actions WHERE transaction_id=$1 ORDER BY due_date NULLS LAST,created_at DESC`, tx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Action{}
	for rows.Next() {
		var x Action
		if err := rows.Scan(&x.ID, &x.TransactionID, &x.Priority, &x.Title, &x.Reason, &x.RecommendedAction, &x.Status, &x.DueDate); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Postgres) UpsertReconciliation(ctx context.Context, x ReconciliationPosition) error {
	b, _ := json.Marshal(x.Notes)
	_, err := r.DB.ExecContext(ctx, `INSERT INTO reconciliation_positions(transaction_id,base_value,inward_value,outward_value,recorded_value,outstanding,difference,currency,status,notes,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,now()) ON CONFLICT(transaction_id) DO UPDATE SET base_value=EXCLUDED.base_value,inward_value=EXCLUDED.inward_value,outward_value=EXCLUDED.outward_value,recorded_value=EXCLUDED.recorded_value,outstanding=EXCLUDED.outstanding,difference=EXCLUDED.difference,currency=EXCLUDED.currency,status=EXCLUDED.status,notes=EXCLUDED.notes,updated_at=now()`, x.TransactionID, x.BaseValue, x.InwardValue, x.OutwardValue, x.RecordedValue, x.Outstanding, x.Difference, x.Currency, x.Status, b)
	return err
}
func (r *Postgres) GetReconciliation(ctx context.Context, tx string) (ReconciliationPosition, error) {
	var x ReconciliationPosition
	var b []byte
	err := r.DB.QueryRowContext(ctx, `SELECT transaction_id::text,currency,base_value,inward_value,outward_value,recorded_value,outstanding,difference,status,notes FROM reconciliation_positions WHERE transaction_id=$1`, tx).Scan(&x.TransactionID, &x.Currency, &x.BaseValue, &x.InwardValue, &x.OutwardValue, &x.RecordedValue, &x.Outstanding, &x.Difference, &x.Status, &b)
	if err != nil {
		return x, err
	}
	_ = json.Unmarshal(b, &x.Notes)
	return x, nil
}
func (r *Postgres) UpsertRegulatory(ctx context.Context, x RegulatoryPosition) error {
	_, err := r.DB.ExecContext(ctx, `INSERT INTO regulatory_positions(transaction_id,registry,position_status,source_type,effective_rule_version,as_of_date,notes,updated_at) VALUES($1,$2,$3,$4,$5,NULLIF($6,'')::date,$7,now()) ON CONFLICT(transaction_id,registry) DO UPDATE SET position_status=EXCLUDED.position_status,source_type=EXCLUDED.source_type,effective_rule_version=EXCLUDED.effective_rule_version,as_of_date=EXCLUDED.as_of_date,notes=EXCLUDED.notes,updated_at=now()`, x.TransactionID, x.Registry, x.PositionStatus, x.SourceType, x.EffectiveRuleVersion, x.AsOfDate, x.Notes)
	return err
}
func (r *Postgres) ListRegulatory(ctx context.Context, tx string) ([]RegulatoryPosition, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id::text,transaction_id::text,registry,position_status,source_type,COALESCE(effective_rule_version,''),COALESCE(as_of_date::text,''),COALESCE(notes,'') FROM regulatory_positions WHERE transaction_id=$1 ORDER BY registry`, tx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []RegulatoryPosition{}
	for rows.Next() {
		var x RegulatoryPosition
		if err := rows.Scan(&x.ID, &x.TransactionID, &x.Registry, &x.PositionStatus, &x.SourceType, &x.EffectiveRuleVersion, &x.AsOfDate, &x.Notes); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Postgres) CreateAudit(ctx context.Context, x AuditEvent) (string, error) {
	b, _ := json.Marshal(x.Payload)
	var id string
	err := r.DB.QueryRowContext(ctx, `INSERT INTO audit_events(company_id,transaction_id,event_type,actor_type,actor_id,payload) VALUES(NULLIF($1,'')::uuid,NULLIF($2,'')::uuid,$3,$4,$5,$6) RETURNING id::text`, x.CompanyID, x.TransactionID, x.EventType, x.ActorType, x.ActorID, b).Scan(&id)
	return id, err
}
func (r *Postgres) ListAudit(ctx context.Context, tx string) ([]AuditEvent, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id::text,COALESCE(company_id::text,''),COALESCE(transaction_id::text,''),event_type,actor_type,COALESCE(actor_id,''),payload FROM audit_events WHERE transaction_id=$1 ORDER BY occurred_at DESC`, tx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []AuditEvent{}
	for rows.Next() {
		var x AuditEvent
		var b []byte
		if err := rows.Scan(&x.ID, &x.CompanyID, &x.TransactionID, &x.EventType, &x.ActorType, &x.ActorID, &b); err != nil {
			return nil, err
		}
		x.Payload = map[string]any{}
		_ = json.Unmarshal(b, &x.Payload)
		out = append(out, x)
	}
	return out, rows.Err()
}

type ShippingBillStore struct{ *Postgres }

func (r ShippingBillStore) Create(c context.Context, x ShippingBill) (string, error) {
	return r.CreateShippingBill(c, x)
}
func (r ShippingBillStore) ListByTransaction(c context.Context, id string) ([]ShippingBill, error) {
	return r.ListShippingBills(c, id)
}

type BillOfEntryStore struct{ *Postgres }

func (r BillOfEntryStore) Create(c context.Context, x BillOfEntry) (string, error) {
	return r.CreateBOE(c, x)
}
func (r BillOfEntryStore) ListByTransaction(c context.Context, id string) ([]BillOfEntry, error) {
	return r.ListBOE(c, id)
}

type DocumentStore struct{ *Postgres }

func (r DocumentStore) Create(c context.Context, x Document) (string, error) {
	return r.CreateDocument(c, x)
}
func (r DocumentStore) ListByTransaction(c context.Context, id string) ([]Document, error) {
	return r.ListDocuments(c, id)
}
func (r DocumentStore) CreateExtraction(c context.Context, x DocumentExtraction) (string, error) {
	return r.CreateDocumentExtraction(c, x)
}
func (r DocumentStore) ListExtractions(c context.Context, id string) ([]DocumentExtraction, error) {
	return r.ListDocumentExtractions(c, id)
}

type EventStore struct{ *Postgres }

func (r EventStore) Create(c context.Context, x Event) (string, error) { return r.CreateEvent(c, x) }
func (r EventStore) ListByTransaction(c context.Context, id string) ([]Event, error) {
	return r.ListEvents(c, id)
}

type DeclarationStore struct{ *Postgres }

func (r DeclarationStore) Create(c context.Context, x Declaration) (string, error) {
	return r.CreateDeclaration(c, x)
}
func (r DeclarationStore) ListByTransaction(c context.Context, id string) ([]Declaration, error) {
	return r.ListDeclarations(c, id)
}

type ActionStore struct{ *Postgres }

func (r ActionStore) Create(c context.Context, x Action) (string, error) { return r.CreateAction(c, x) }
func (r ActionStore) ListByTransaction(c context.Context, id string) ([]Action, error) {
	return r.ListActions(c, id)
}

type ReconciliationStore struct{ *Postgres }

func (r ReconciliationStore) Upsert(c context.Context, x ReconciliationPosition) error {
	return r.UpsertReconciliation(c, x)
}
func (r ReconciliationStore) GetByTransaction(c context.Context, id string) (ReconciliationPosition, error) {
	return r.GetReconciliation(c, id)
}

type RegulatoryStore struct{ *Postgres }

func (r RegulatoryStore) Upsert(c context.Context, x RegulatoryPosition) error {
	return r.UpsertRegulatory(c, x)
}
func (r RegulatoryStore) ListByTransaction(c context.Context, id string) ([]RegulatoryPosition, error) {
	return r.ListRegulatory(c, id)
}

type AuditStore struct{ *Postgres }

func (r AuditStore) Create(c context.Context, x AuditEvent) (string, error) {
	return r.CreateAudit(c, x)
}
func (r AuditStore) ListByTransaction(c context.Context, id string) ([]AuditEvent, error) {
	return r.ListAudit(c, id)
}

type RequirementStore struct{ *Postgres }

func (r RequirementStore) Create(c context.Context, x Requirement) (string, error) {
	return r.CreateRequirement(c, x)
}
func (r RequirementStore) ListByTransaction(c context.Context, id string) ([]Requirement, error) {
	return r.ListRequirements(c, id)
}
func (r RequirementStore) UpdateStatus(c context.Context, id, status, notes string) error {
	return r.UpdateRequirementStatus(c, id, status, notes)
}

func (p *Postgres) CreateRequirement(c context.Context, x Requirement) (string, error) {
	q := `INSERT INTO requirements(id,transaction_id,company_id,category,title,description,status,priority,source,requested_by,assigned_to,due_date,notes) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NULLIF($12,''),$13) ON CONFLICT(id) DO NOTHING`
	_, err := p.DB.ExecContext(c, q, x.ID, x.TransactionID, x.CompanyID, x.Category, x.Title, x.Description, x.Status, x.Priority, x.Source, x.RequestedBy, x.AssignedTo, x.DueDate, x.Notes)
	return x.ID, err
}
func (p *Postgres) ListRequirements(c context.Context, transactionID string) ([]Requirement, error) {
	rows, err := p.DB.QueryContext(c, `SELECT id,transaction_id,company_id,category,title,description,status,priority,source,requested_by,assigned_to,COALESCE(to_char(due_date,'YYYY-MM-DD'),''),notes FROM requirements WHERE transaction_id=$1 ORDER BY COALESCE(due_date,'9999-12-31'), created_at`, transactionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Requirement
	for rows.Next() {
		var x Requirement
		if err := rows.Scan(&x.ID, &x.TransactionID, &x.CompanyID, &x.Category, &x.Title, &x.Description, &x.Status, &x.Priority, &x.Source, &x.RequestedBy, &x.AssignedTo, &x.DueDate, &x.Notes); err != nil {
			return nil, err
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (p *Postgres) UpdateRequirementStatus(c context.Context, id, status, notes string) error {
	_, err := p.DB.ExecContext(c, `UPDATE requirements SET status=$2,notes=$3,updated_at=now() WHERE id=$1`, id, status, notes)
	return err
}
