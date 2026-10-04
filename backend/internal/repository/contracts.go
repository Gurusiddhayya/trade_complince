package repository

import "context"

type Transaction struct {
	ID, CompanyID, ReferenceNo, ProductCode, ProductType, Status, Currency, Counterparty, Country, Description string
	DeclaredValue, ContractValue, InvoiceValue                                                                 float64
}
type Invoice struct {
	ID, TransactionID, InvoiceNumber, Currency, PaymentTerms, Description, Status string
	InvoiceValue                                                                  float64
}
type Payment struct {
	ID, TransactionID, Reference, PaymentType, Currency, PaymentDate, PurposeCode, BankReference, Direction string
	Amount                                                                                                  float64
}

type TransactionRepository interface {
	Get(ctx context.Context, id string) (Transaction, error)
	ListByCompany(ctx context.Context, companyID string) ([]Transaction, error)
	Create(ctx context.Context, t Transaction) (string, error)
	UpdateStatus(ctx context.Context, id, status string) error
}
type InvoiceRepository interface {
	Create(ctx context.Context, i Invoice) (string, error)
	ListByTransaction(ctx context.Context, transactionID string) ([]Invoice, error)
}
type PaymentRepository interface {
	Create(ctx context.Context, p Payment) (string, error)
	ListByTransaction(ctx context.Context, transactionID string) ([]Payment, error)
}

type ShippingBill struct {
	ID, TransactionID, Number, Date, Port, ExportDate, Currency string
	ExportValue                                                 float64
}
type BillOfEntry struct {
	ID, TransactionID, Number, Date, Port, Currency string
	AssessableValue                                 float64
}
type Document struct {
	ID, TransactionID, DocumentType, FileName, Status, Source  string
	StorageKey, MIMEType, SHA256, UploadedBy, ParentDocumentID string
	FileSize                                                   int64
	Version                                                    int
}
type Event struct{ ID, TransactionID, EventType, EventDate, Notes string }
type Declaration struct {
	ID, TransactionID, DeclarationType, Version, Status string
	Fields                                              map[string]string
}
type Action struct{ ID, TransactionID, Priority, Title, Reason, RecommendedAction, Status, DueDate string }
type ReconciliationPosition struct {
	TransactionID, Currency, Status                                              string
	BaseValue, InwardValue, OutwardValue, RecordedValue, Outstanding, Difference float64
	Notes                                                                        []string
}
type RegulatoryPosition struct{ ID, TransactionID, Registry, PositionStatus, SourceType, EffectiveRuleVersion, AsOfDate, Notes string }
type AuditEvent struct {
	ID, CompanyID, TransactionID, EventType, ActorType, ActorID string
	Payload                                                     map[string]any
}

type ShippingBillRepository interface {
	Create(context.Context, ShippingBill) (string, error)
	ListByTransaction(context.Context, string) ([]ShippingBill, error)
}
type BillOfEntryRepository interface {
	Create(context.Context, BillOfEntry) (string, error)
	ListByTransaction(context.Context, string) ([]BillOfEntry, error)
}
type DocumentExtraction struct {
	ID, DocumentID, FieldName, FieldValue, Source, ConfirmedBy string
	Confidence                                                 float64
	Confirmed                                                  bool
}
type DocumentRepository interface {
	Create(context.Context, Document) (string, error)
	ListByTransaction(context.Context, string) ([]Document, error)
	CreateExtraction(context.Context, DocumentExtraction) (string, error)
	ListExtractions(context.Context, string) ([]DocumentExtraction, error)
}
type EventRepository interface {
	Create(context.Context, Event) (string, error)
	ListByTransaction(context.Context, string) ([]Event, error)
}
type DeclarationRepository interface {
	Create(context.Context, Declaration) (string, error)
	ListByTransaction(context.Context, string) ([]Declaration, error)
}
type ActionRepository interface {
	Create(context.Context, Action) (string, error)
	ListByTransaction(context.Context, string) ([]Action, error)
}
type ReconciliationRepository interface {
	Upsert(context.Context, ReconciliationPosition) error
	GetByTransaction(context.Context, string) (ReconciliationPosition, error)
}
type RegulatoryPositionRepository interface {
	Upsert(context.Context, RegulatoryPosition) error
	ListByTransaction(context.Context, string) ([]RegulatoryPosition, error)
}
type AuditRepository interface {
	Create(context.Context, AuditEvent) (string, error)
	ListByTransaction(context.Context, string) ([]AuditEvent, error)
}
