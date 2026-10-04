package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Product struct{ Code, Name, Category string }

var products = []Product{
	{"EXP_GOODS", "Export of Goods", "Export"}, {"EXP_SERVICE", "Software / Service Export", "Export"}, {"IMP_GOODS", "Import of Goods", "Import"},
	{"HSS", "High Sea Sale", "Import"}, {"IMPORT_LC", "Import LC", "Trade Finance"}, {"BG", "Bank Guarantee", "Trade Finance"},
	{"EPC_PCFC", "EPC / PCFC", "Trade Finance"}, {"EXPORT_COLLECTION", "Export Bill Collection — LC & Non-LC", "Export"},
	{"TRADE_CREDIT", "Buyer Credit & Supplier Credit", "Trade Credit"}, {"MTT", "Merchanting Trade Transaction", "Special"},
	{"FDI", "FDI", "Investment"}, {"ODI", "ODI", "Investment"}, {"ECB", "ECB", "Borrowing"},
}

type Transaction struct {
	ID               string    `json:"id"`
	ReferenceNo      string    `json:"reference_no"`
	ProductType      string    `json:"product_type"`
	ProductCode      string    `json:"product_code"`
	Status           string    `json:"status"`
	Currency         string    `json:"currency"`
	DeclaredValue    float64   `json:"declared_value"`
	Counterparty     string    `json:"counterparty"`
	Country          string    `json:"country"`
	Description      string    `json:"description"`
	ContractValue    float64   `json:"contract_value"`
	InvoiceValue     float64   `json:"invoice_value"`
	RealizedValue    float64   `json:"realized_value"`
	OutstandingValue float64   `json:"outstanding_value"`
	Difference       float64   `json:"difference"`
	CreatedAt        time.Time `json:"created_at"`
}
type Payment struct {
	ID, TransactionID, Reference, PaymentType, Currency, PaymentDate, PurposeCode, BankReference, Direction string
	Amount                                                                                                  float64 `json:"amount"`
}
type PaymentWire struct {
	ID            string  `json:"id"`
	TransactionID string  `json:"transaction_id"`
	Reference     string  `json:"reference"`
	PaymentType   string  `json:"payment_type"`
	Currency      string  `json:"currency"`
	PaymentDate   string  `json:"payment_date"`
	PurposeCode   string  `json:"purpose_code"`
	BankReference string  `json:"bank_reference"`
	Direction     string  `json:"direction"`
	Amount        float64 `json:"amount"`
}
type Invoice struct {
	ID, TransactionID, InvoiceNumber, Currency, PaymentTerms, Description, Status string  `json:"-"`
	InvoiceValue                                                                  float64 `json:"invoice_value"`
}
type ShippingBill struct {
	ID, TransactionID, Number, Date, Port, ExportDate, Currency string  `json:"-"`
	ExportValue                                                 float64 `json:"export_value"`
}
type BillOfEntry struct {
	ID, TransactionID, Number, Date, Port, Currency string  `json:"-"`
	AssessableValue                                 float64 `json:"assessable_value"`
}
type Document struct {
	ID               string `json:"id"`
	TransactionID    string `json:"transaction_id"`
	Type             string `json:"document_type"`
	FileName         string `json:"file_name"`
	Status           string `json:"status"`
	Version          int    `json:"version"`
	StorageKey       string `json:"storage_key,omitempty"`
	MIMEType         string `json:"mime_type,omitempty"`
	FileSize         int64  `json:"file_size,omitempty"`
	SHA256           string `json:"sha256,omitempty"`
	UploadedBy       string `json:"uploaded_by,omitempty"`
	ParentDocumentID string `json:"parent_document_id,omitempty"`
}
type DocumentExtraction struct {
	ID          string  `json:"id"`
	DocumentID  string  `json:"document_id"`
	FieldName   string  `json:"field_name"`
	FieldValue  string  `json:"field_value,omitempty"`
	Confidence  float64 `json:"confidence"`
	Source      string  `json:"source"`
	Confirmed   bool    `json:"confirmed"`
	ConfirmedBy string  `json:"confirmed_by,omitempty"`
}
type DocumentFinding struct {
	Severity string   `json:"severity"`
	Code     string   `json:"code"`
	Message  string   `json:"message"`
	Fields   []string `json:"fields,omitempty"`
	Action   string   `json:"recommended_action,omitempty"`
}
type DocumentIntelligence struct {
	DocumentID               string               `json:"document_id"`
	Classification           string               `json:"classification"`
	ClassificationConfidence float64              `json:"classification_confidence"`
	ExtractionStatus         string               `json:"extraction_status"`
	ExtractedFields          []DocumentExtraction `json:"extracted_fields"`
	Findings                 []DocumentFinding    `json:"findings"`
}

type Action struct {
	Priority, Title, Reason, RecommendedAction, Status string
	DueDate                                            string `json:"due_date,omitempty"`
}

type Requirement struct {
	ID            string `json:"id"`
	TransactionID string `json:"transaction_id"`
	Category      string `json:"category"`
	Title         string `json:"title"`
	Description   string `json:"description,omitempty"`
	Status        string `json:"status"`
	Priority      string `json:"priority"`
	Source        string `json:"source,omitempty"`
	RequestedBy   string `json:"requested_by,omitempty"`
	AssignedTo    string `json:"assigned_to,omitempty"`
	DueDate       string `json:"due_date,omitempty"`
	Notes         string `json:"notes,omitempty"`
}

type Event struct {
	ID            string `json:"id"`
	TransactionID string `json:"transaction_id"`
	EventType     string `json:"event_type"`
	EventDate     string `json:"event_date"`
	Notes         string `json:"notes"`
}
type Declaration struct {
	ID            string            `json:"id"`
	TransactionID string            `json:"transaction_id"`
	Type          string            `json:"declaration_type"`
	Version       string            `json:"version"`
	Status        string            `json:"status"`
	Fields        map[string]string `json:"fields"`
}
type ReconciliationView struct {
	Mode          string   `json:"mode"`
	BaseValue     float64  `json:"base_value"`
	Currency      string   `json:"currency"`
	InwardValue   float64  `json:"inward_value"`
	OutwardValue  float64  `json:"outward_value"`
	RecordedValue float64  `json:"recorded_value"`
	Outstanding   float64  `json:"outstanding"`
	Difference    float64  `json:"difference"`
	Status        string   `json:"status"`
	Notes         []string `json:"notes"`
}

type Store struct {
	sync.RWMutex
	Transactions        map[string]*Transaction
	Payments            map[string][]Payment
	Documents           map[string][]Document
	Events              map[string][]Event
	Declarations        map[string][]Declaration
	DocumentExtractions map[string][]DocumentExtraction
	Invoices            map[string][]Invoice
	ShippingBills       map[string][]ShippingBill
	BillsOfEntry        map[string][]BillOfEntry
	Requirements        map[string][]Requirement
	Counter             int
}

func newStore() *Store {
	s := &Store{Transactions: map[string]*Transaction{}, Payments: map[string][]Payment{}, Documents: map[string][]Document{}, DocumentExtractions: map[string][]DocumentExtraction{}, Events: map[string][]Event{}, Declarations: map[string][]Declaration{}, Invoices: map[string][]Invoice{}, ShippingBills: map[string][]ShippingBill{}, BillsOfEntry: map[string][]BillOfEntry{}, Requirements: map[string][]Requirement{}, Counter: 13}
	now := time.Now()
	samples := []Transaction{
		{ID: "1", ReferenceNo: "EXP-2026-000001", ProductCode: "EXP_GOODS", ProductType: "Export of Goods", Status: "In Progress", Currency: "USD", DeclaredValue: 100000, Counterparty: "ABC Trading Ltd", Country: "UAE", Description: "Export of machinery", ContractValue: 100000, InvoiceValue: 100000, CreatedAt: now},
		{ID: "2", ReferenceNo: "EXP-2026-000002", ProductCode: "EXP_SERVICE", ProductType: "Software / Service Export", Status: "In Progress", Currency: "USD", DeclaredValue: 50000, Counterparty: "Global Tech Inc", Country: "USA", Description: "Software implementation services", ContractValue: 50000, InvoiceValue: 50000, CreatedAt: now},
		{ID: "3", ReferenceNo: "IMP-2026-000001", ProductCode: "IMP_GOODS", ProductType: "Import of Goods", Status: "Attention Required", Currency: "USD", DeclaredValue: 100000, Counterparty: "Global Supplier Pte Ltd", Country: "Singapore", Description: "Import of machinery", ContractValue: 100000, InvoiceValue: 100000, CreatedAt: now},
		{ID: "4", ReferenceNo: "HSS-2026-000001", ProductCode: "HSS", ProductType: "High Sea Sale", Status: "In Progress", Currency: "USD", DeclaredValue: 75000, Counterparty: "Overseas Supplier / Indian Buyer", Country: "Singapore", Description: "High Sea Sale before clearance", ContractValue: 75000, InvoiceValue: 75000, CreatedAt: now},
		{ID: "5", ReferenceNo: "ILC-2026-000001", ProductCode: "IMPORT_LC", ProductType: "Import LC", Status: "In Progress", Currency: "USD", DeclaredValue: 200000, Counterparty: "Global Machinery Ltd", Country: "Germany", Description: "Import LC issuance and bill tracking", ContractValue: 200000, InvoiceValue: 200000, CreatedAt: now},
		{ID: "6", ReferenceNo: "BG-2026-000001", ProductCode: "BG", ProductType: "Bank Guarantee", Status: "In Progress", Currency: "USD", DeclaredValue: 100000, Counterparty: "Overseas Beneficiary", Country: "UAE", Description: "Performance guarantee", ContractValue: 100000, InvoiceValue: 0, CreatedAt: now},
		{ID: "7", ReferenceNo: "FIN-2026-000001", ProductCode: "EPC_PCFC", ProductType: "EPC / PCFC", Status: "In Progress", Currency: "USD", DeclaredValue: 80000, Counterparty: "ABC Trading Ltd", Country: "UAE", Description: "Pre-shipment export finance", ContractValue: 80000, InvoiceValue: 80000, CreatedAt: now},
		{ID: "8", ReferenceNo: "COL-2026-000001", ProductCode: "EXPORT_COLLECTION", ProductType: "Export Bill Collection — LC & Non-LC", Status: "In Progress", Currency: "USD", DeclaredValue: 120000, Counterparty: "Buyer Corp", Country: "UK", Description: "Export bill collection", ContractValue: 120000, InvoiceValue: 120000, CreatedAt: now},
		{ID: "9", ReferenceNo: "TC-2026-000001", ProductCode: "TRADE_CREDIT", ProductType: "Buyer Credit & Supplier Credit", Status: "In Progress", Currency: "USD", DeclaredValue: 300000, Counterparty: "Trade Finance Counterparty", Country: "USA", Description: "Trade credit monitoring", ContractValue: 300000, InvoiceValue: 300000, CreatedAt: now},
		{ID: "10", ReferenceNo: "MTT-2026-000001", ProductCode: "MTT", ProductType: "Merchanting Trade Transaction", Status: "Attention Required", Currency: "USD", DeclaredValue: 150000, Counterparty: "Supplier / Buyer Pair", Country: "UAE", Description: "Purchase leg, outward payment, sale leg and inward receipt", ContractValue: 150000, InvoiceValue: 150000, CreatedAt: now},
		{ID: "11", ReferenceNo: "FDI-2026-000001", ProductCode: "FDI", ProductType: "FDI", Status: "In Progress", Currency: "INR", DeclaredValue: 50000000, Counterparty: "Foreign Investor", Country: "USA", Description: "Foreign investment record and reporting preparation", ContractValue: 50000000, InvoiceValue: 0, CreatedAt: now},
		{ID: "12", ReferenceNo: "ODI-2026-000001", ProductCode: "ODI", ProductType: "ODI", Status: "In Progress", Currency: "USD", DeclaredValue: 100000, Counterparty: "Overseas JV", Country: "Singapore", Description: "Overseas investment monitoring", ContractValue: 100000, InvoiceValue: 0, CreatedAt: now},
		{ID: "13", ReferenceNo: "ECB-2026-000001", ProductCode: "ECB", ProductType: "ECB", Status: "In Progress", Currency: "USD", DeclaredValue: 1000000, Counterparty: "Overseas Lender", Country: "UK", Description: "External commercial borrowing record and monitoring", ContractValue: 1000000, InvoiceValue: 0, CreatedAt: now},
	}
	for i := range samples {
		s.Transactions[samples[i].ID] = &samples[i]
	}
	s.Payments["1"] = []Payment{{ID: "PAY-1", TransactionID: "1", Reference: "IRM-REC-001", PaymentType: "Advance / Export realization", Amount: 60000, Currency: "USD", PaymentDate: "2026-09-01", PurposeCode: "Customer confirmation required", Direction: "INWARD"}}
	s.Payments["2"] = []Payment{{ID: "PAY-2", TransactionID: "2", Reference: "IRM-REC-002", PaymentType: "Service realization", Amount: 30000, Currency: "USD", PaymentDate: "2026-09-03", PurposeCode: "Customer confirmation required", Direction: "INWARD"}}
	s.Payments["3"] = []Payment{{ID: "PAY-3", TransactionID: "3", Reference: "ORM-SETTLE-001", PaymentType: "Import settlement", Amount: 40000, Currency: "USD", PaymentDate: "2026-09-02", PurposeCode: "Customer confirmation required", Direction: "OUTWARD"}}
	s.Payments["10"] = []Payment{{ID: "PAY-10-OUT", TransactionID: "10", Reference: "MTT-OUT-001", PaymentType: "MTT outward purchase payment", Amount: 100000, Currency: "USD", PaymentDate: "2026-09-02", PurposeCode: "Customer confirmation required", Direction: "OUTWARD"}, {ID: "PAY-10-IN", TransactionID: "10", Reference: "MTT-IN-001", PaymentType: "MTT inward sale receipt", Amount: 90000, Currency: "USD", PaymentDate: "2026-09-10", PurposeCode: "Customer confirmation required", Direction: "INWARD"}}
	s.recalculateAll()
	return s
}
func (s *Store) recalculateAll() {
	for id := range s.Transactions {
		s.recalculate(id)
	}
}
func (s *Store) recalculate(id string) {
	t := s.Transactions[id]
	if t == nil {
		return
	}
	var p float64
	for _, x := range s.Payments[id] {
		p += x.Amount
	}
	t.RealizedValue = p
	base := t.InvoiceValue
	if base <= 0 {
		base = t.DeclaredValue
	}
	t.OutstandingValue = math.Max(base-p, 0)
	t.Difference = math.Max(p-base, 0)
	if base > 0 && t.OutstandingValue == 0 && t.Difference == 0 {
		t.Status = "Reconciled"
	} else if t.OutstandingValue > 0 {
		if t.Status == "Reconciled" {
			t.Status = "In Progress"
		}
	}
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(204)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func productByCode(code string) Product {
	for _, p := range products {
		if p.Code == code {
			return p
		}
	}
	return Product{Code: code, Name: code}
}
func actionsFor(t *Transaction, payments []Payment) []Action {
	a := []Action{}
	v := reconciliationFor(t, payments)
	if v.Outstanding > 0 && v.BaseValue > 0 {
		a = append(a, Action{"attention", "Review outstanding / reconciliation", fmt.Sprintf("%s %.2f remains unmatched.", t.Currency, v.Outstanding), "Review invoices, payments, allocations and supporting documents; contact the AD bank if clarification is needed.", "Open", ""})
	}
	if t.ProductCode == "MTT" && v.OutwardValue > 0 && v.InwardValue > 0 && v.OutwardValue != v.InwardValue {
		a = append(a, Action{"attention", "Review MTT two-leg difference", fmt.Sprintf("Purchase outlay %.2f and sale receipt %.2f are different.", v.OutwardValue, v.InwardValue), "Review purchase and sale legs separately and confirm the applicable monitoring position with the AD bank.", "Open", ""})
	}
	a = append(a, Action{"info", "Review regulatory position", "The app stores a customer-side position; it is not a live bank/EDPMS/IDPMS status.", "Confirm actual bank/authority status where applicable.", "Open", ""})
	return a
}
func main() {
	s := newStore()
	ctx := context.Background()
	persist, err := openPostgres(ctx)
	if err != nil {
		log.Fatalf("postgres initialization failed: %v", err)
	}
	defer persist.Close()
	if snap, err := persist.Load(ctx); err != nil {
		log.Fatalf("postgres load failed: %v", err)
	} else if snap != nil {
		s.Transactions, s.Payments, s.Documents, s.DocumentExtractions, s.Events, s.Declarations, s.Invoices, s.ShippingBills, s.BillsOfEntry, s.Requirements, s.Counter = snap.Transactions, snap.Payments, snap.Documents, snap.DocumentExtractions, snap.Events, snap.Declarations, snap.Invoices, snap.ShippingBills, snap.BillsOfEntry, snap.Requirements, snap.Counter
		s.recalculateAll()
	} else {
		if loaded, err := persist.LoadNormalized(ctx, s); err != nil {
			log.Fatalf("postgres normalized load failed: %v", err)
		} else if loaded {
			s.recalculateAll()
		}
	}
	persistStore := func() {
		if persist == nil {
			return
		}
		if err := persist.SaveNormalized(ctx, s); err != nil {
			log.Printf("postgres normalized persistence warning: %v", err)
		}
		if err := persist.Save(ctx, s); err != nil {
			log.Printf("postgres snapshot persistence warning: %v", err)
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		status := "in-memory"
		if persist != nil {
			status = "postgres"
		}
		writeJSON(w, 200, map[string]any{"status": "ok", "version": "2.17-customer-workflow-validation", "persistence": status, "security": securityConfigSummary()})
	})
	mux.HandleFunc("/api/v1/requirements", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			txID := r.URL.Query().Get("transaction_id")
			s.RLock()
			defer s.RUnlock()
			if txID != "" {
				writeJSON(w, 200, requirementsFor(s.Transactions[txID], s.Requirements[txID]))
				return
			}
			all := []Requirement{}
			for _, xs := range s.Requirements {
				all = append(all, xs...)
			}
			writeJSON(w, 200, all)
			return
		}
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		var in Requirement
		if json.NewDecoder(r.Body).Decode(&in) != nil || strings.TrimSpace(in.TransactionID) == "" || strings.TrimSpace(in.Title) == "" {
			writeJSON(w, 400, map[string]string{"error": "transaction_id and title are required"})
			return
		}
		s.Lock()
		defer s.Unlock()
		if s.Transactions[in.TransactionID] == nil {
			writeJSON(w, 404, map[string]string{"error": "transaction not found"})
			return
		}
		if in.Category == "" {
			in.Category = "BANK"
		}
		if in.Status == "" {
			in.Status = "MISSING"
		}
		if in.Priority == "" {
			in.Priority = "NORMAL"
		}
		in.ID = fmt.Sprintf("REQ-%d-%d", s.Counter, time.Now().UnixNano())
		s.Requirements[in.TransactionID] = append(s.Requirements[in.TransactionID], in)
		persistStore()
		writeJSON(w, 201, in)
	})
	mux.HandleFunc("/api/v1/requirements/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 4 {
			http.NotFound(w, r)
			return
		}
		id := parts[3]
		if r.Method != "PATCH" {
			http.NotFound(w, r)
			return
		}
		var in struct {
			Status string `json:"status"`
			Notes  string `json:"notes"`
		}
		if json.NewDecoder(r.Body).Decode(&in) != nil || in.Status == "" {
			writeJSON(w, 400, map[string]string{"error": "status is required"})
			return
		}
		s.Lock()
		defer s.Unlock()
		for txID, xs := range s.Requirements {
			for i := range xs {
				if xs[i].ID == id {
					xs[i].Status = in.Status
					xs[i].Notes = in.Notes
					s.Requirements[txID] = xs
					persistStore()
					writeJSON(w, 200, xs[i])
					return
				}
			}
		}
		writeJSON(w, 404, map[string]string{"error": "requirement not found"})
	})
	mux.HandleFunc("/api/v1/products", func(w http.ResponseWriter, r *http.Request) { writeJSON(w, 200, products) })
	mux.HandleFunc("/api/v1/dashboard", func(w http.ResponseWriter, r *http.Request) {
		s.RLock()
		defer s.RUnlock()
		var exp, imp, fin float64
		attention := 0
		for _, t := range s.Transactions {
			if strings.HasPrefix(t.ProductCode, "EXP") || t.ProductCode == "EXPORT_COLLECTION" {
				exp += t.DeclaredValue
			}
			if t.ProductCode == "IMP_GOODS" || t.ProductCode == "HSS" || t.ProductCode == "IMPORT_LC" {
				imp += t.DeclaredValue
			}
			if t.ProductCode == "EPC_PCFC" || t.ProductCode == "BG" || t.ProductCode == "TRADE_CREDIT" {
				fin += t.DeclaredValue
			}
			if t.Status == "Attention Required" || t.OutstandingValue > 0 {
				attention++
			}
		}
		writeJSON(w, 200, map[string]any{"trade_health": map[string]string{"status": "Attention required", "note": "Customer-side position; confirm with AD bank/authority where applicable."}, "transactions": len(s.Transactions), "products": len(products), "export_value": exp, "import_value": imp, "trade_finance_value": fin, "attention": attention, "requirements": requirementDashboard(s), "edpms": "Customer-recorded position", "idpms": "Customer-recorded position", "declarations": "Preparation / review", "ebrc": "Readiness tracking", "rules_version": "FEMA 23(R)/2026-RB + versioned rules; applicability determined by transaction/event date"})
	})
	mux.HandleFunc("/api/v1/transactions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			s.RLock()
			defer s.RUnlock()
			out := []*Transaction{}
			for _, t := range s.Transactions {
				out = append(out, t)
			}
			writeJSON(w, 200, out)
			return
		}
		if r.Method == "POST" {
			var in Transaction
			if json.NewDecoder(r.Body).Decode(&in) != nil {
				writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
				return
			}
			s.Lock()
			defer s.Unlock()
			s.Counter++
			in.ID = fmt.Sprint(s.Counter)
			p := productByCode(in.ProductCode)
			if p.Code == "" {
				p = productByCode("EXP_GOODS")
			}
			in.ProductType = p.Name
			if in.ProductCode == "" {
				in.ProductCode = p.Code
			}
			prefix := "TRD"
			if p.Category == "Export" {
				prefix = "EXP"
			} else if p.Category == "Import" {
				prefix = "IMP"
			}
			in.ReferenceNo = fmt.Sprintf("%s-2026-%06d", prefix, s.Counter)
			in.Status = "In Progress"
			in.CreatedAt = time.Now()
			s.Transactions[in.ID] = &in
			persistStore()
			writeJSON(w, 201, in)
			return
		}
		http.NotFound(w, r)
	})
	mux.HandleFunc("/api/v1/transactions/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 4 {
			http.NotFound(w, r)
			return
		}
		id := parts[3]
		s.RLock()
		defer s.RUnlock()
		t := s.Transactions[id]
		if t == nil {
			http.NotFound(w, r)
			return
		}
		if len(parts) == 5 && parts[4] == "360" {
			writeJSON(w, 200, map[string]any{"transaction": t, "payments": s.Payments[id], "documents": s.Documents[id], "document_intelligence": transactionDocumentIntelligence(id, s.Documents[id], s.DocumentExtractions, t), "reconciliation": reconciliationFor(t, s.Payments[id]), "regulatory": regulatoryFor(t), "actions": actionsFor(t, s.Payments[id]), "workflow": workflowFor(t), "help_me": helpFor(t), "events": s.Events[id], "declarations": s.Declarations[id], "requirements": requirementsFor(t, s.Requirements[id]), "requirement_summary": requirementSummary(requirementsFor(t, s.Requirements[id]))})
			return
		}
		http.NotFound(w, r)
	})
	mux.HandleFunc("/api/v1/payments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		var p PaymentWire
		if json.NewDecoder(r.Body).Decode(&p) != nil {
			writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
			return
		}
		s.Lock()
		defer s.Unlock()
		if s.Transactions[p.TransactionID] == nil {
			writeJSON(w, 404, map[string]string{"error": "transaction not found"})
			return
		}
		p.ID = fmt.Sprintf("PAY-%d-%d", len(s.Payments[p.TransactionID])+1, time.Now().UnixNano())
		if p.Currency == "" {
			p.Currency = s.Transactions[p.TransactionID].Currency
		}
		direction := strings.ToUpper(p.Direction)
		if direction != "INWARD" && direction != "OUTWARD" {
			pt := strings.ToLower(p.PaymentType)
			if strings.Contains(pt, "outward") || strings.Contains(pt, "payment") || strings.Contains(pt, "settlement") {
				direction = "OUTWARD"
			} else {
				direction = "INWARD"
			}
		}
		s.Payments[p.TransactionID] = append(s.Payments[p.TransactionID], Payment{ID: p.ID, TransactionID: p.TransactionID, Reference: p.Reference, PaymentType: p.PaymentType, Amount: p.Amount, Currency: p.Currency, PaymentDate: p.PaymentDate, PurposeCode: p.PurposeCode, BankReference: p.BankReference, Direction: direction})
		s.recalculate(p.TransactionID)
		persistStore()
		writeJSON(w, 201, p)
	})
	mux.HandleFunc("/api/v1/products/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 4 {
			http.NotFound(w, r)
			return
		}
		p := productByCode(parts[3])
		if p.Code == "" {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, 200, map[string]any{"product": p, "workflow": workflowFor(&Transaction{ProductCode: p.Code}), "guidance": productGuidance(p.Code), "required_data": productData(p.Code), "requirements": productRequirements(p.Code), "document_types": productDocuments(p.Code), "help": productHelp(p.Code)})
	})
	mux.HandleFunc("/api/v1/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		var e Event
		if json.NewDecoder(r.Body).Decode(&e) != nil {
			writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
			return
		}
		s.Lock()
		defer s.Unlock()
		if s.Transactions[e.TransactionID] == nil {
			writeJSON(w, 404, map[string]string{"error": "transaction not found"})
			return
		}
		e.ID = fmt.Sprintf("EVT-%d-%d", s.Counter, time.Now().UnixNano())
		if e.EventDate == "" {
			e.EventDate = time.Now().Format("2006-01-02")
		}
		s.Events[e.TransactionID] = append(s.Events[e.TransactionID], e)
		persistStore()
		writeJSON(w, 201, e)
	})
	mux.HandleFunc("/api/v1/documents/extractions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		var in struct {
			DocumentID, FieldName, FieldValue, Source, ConfirmedBy string
			Confidence                                             float64
			Confirmed                                              bool
		}
		if json.NewDecoder(r.Body).Decode(&in) != nil || in.DocumentID == "" || in.FieldName == "" {
			writeJSON(w, 400, map[string]string{"error": "document_id and field_name are required"})
			return
		}
		s.RLock()
		var found *Document
		for _, docs := range s.Documents {
			for i := range docs {
				if docs[i].ID == in.DocumentID {
					found = &docs[i]
					break
				}
			}
			if found != nil {
				break
			}
		}
		s.RUnlock()
		if found == nil {
			writeJSON(w, 404, map[string]string{"error": "document not found"})
			return
		}
		x := DocumentExtraction{ID: fmt.Sprintf("EXT-%d-%d", s.Counter, time.Now().UnixNano()), DocumentID: in.DocumentID, FieldName: strings.TrimSpace(in.FieldName), FieldValue: in.FieldValue, Confidence: clampConfidence(in.Confidence), Source: defaultSource(in.Source), Confirmed: in.Confirmed, ConfirmedBy: in.ConfirmedBy}
		s.Lock()
		s.DocumentExtractions[in.DocumentID] = upsertExtraction(s.DocumentExtractions[in.DocumentID], x)
		s.Unlock()
		persistStore()
		writeJSON(w, 201, x)
	})

	mux.HandleFunc("/api/v1/documents/intelligence/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 5 {
			http.NotFound(w, r)
			return
		}
		docID := parts[4]
		s.RLock()
		defer s.RUnlock()
		var doc *Document
		for _, docs := range s.Documents {
			for i := range docs {
				if docs[i].ID == docID {
					doc = &docs[i]
					break
				}
			}
			if doc != nil {
				break
			}
		}
		if doc == nil {
			writeJSON(w, 404, map[string]string{"error": "document not found"})
			return
		}
		fields := append([]DocumentExtraction(nil), s.DocumentExtractions[docID]...)
		info := buildDocumentIntelligence(*doc, fields)
		writeJSON(w, 200, info)
	})
	mux.HandleFunc("/api/v1/documents/analyze/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 5 {
			http.NotFound(w, r)
			return
		}
		docID := parts[4]
		s.RLock()
		var doc *Document
		for _, docs := range s.Documents {
			for i := range docs {
				if docs[i].ID == docID {
					d := docs[i]
					doc = &d
					break
				}
			}
			if doc != nil {
				break
			}
		}
		fields := append([]DocumentExtraction(nil), s.DocumentExtractions[docID]...)
		s.RUnlock()
		if doc == nil {
			writeJSON(w, 404, map[string]string{"error": "document not found"})
			return
		}
		info := buildDocumentIntelligence(*doc, fields)
		writeJSON(w, 200, info)
	})
	mux.HandleFunc("/api/v1/transactions/document-intelligence/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 5 {
			http.NotFound(w, r)
			return
		}
		txID := parts[4]
		s.RLock()
		defer s.RUnlock()
		if s.Transactions[txID] == nil {
			writeJSON(w, 404, map[string]string{"error": "transaction not found"})
			return
		}
		out := transactionDocumentIntelligence(txID, s.Documents[txID], s.DocumentExtractions, s.Transactions[txID])
		writeJSON(w, 200, out)
	})

	mux.HandleFunc("/api/v1/declarations/prepare/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 5 {
			http.NotFound(w, r)
			return
		}
		txID := parts[4]
		s.RLock()
		t := s.Transactions[txID]
		if t == nil {
			s.RUnlock()
			writeJSON(w, 404, map[string]string{"error": "transaction not found"})
			return
		}
		docs := append([]Document(nil), s.Documents[txID]...)
		invoices := append([]Invoice(nil), s.Invoices[txID]...)
		payments := append([]Payment(nil), s.Payments[txID]...)
		extractions := make(map[string][]DocumentExtraction, len(docs))
		for _, d := range docs {
			extractions[d.ID] = append([]DocumentExtraction(nil), s.DocumentExtractions[d.ID]...)
		}
		s.RUnlock()
		prep := prepareDeclaration(t, docs, extractions, invoices, payments)
		if r.Method == "GET" {
			writeJSON(w, 200, prep)
			return
		}
		if prep.DeclarationType == "" {
			writeJSON(w, 422, prep)
			return
		}
		s.Lock()
		declaration := declarationFromPreparation(prep, "2026.1")
		declaration.ID = fmt.Sprintf("DEC-%d-%d", s.Counter, time.Now().UnixNano())
		s.Declarations[txID] = append(s.Declarations[txID], declaration)
		s.Unlock()
		persistStore()
		writeJSON(w, 201, map[string]any{"preparation": prep, "declaration": declaration})
	})

	mux.HandleFunc("/api/v1/declarations", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			s.RLock()
			defer s.RUnlock()
			writeJSON(w, 200, s.Declarations)
			return
		}
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		var d Declaration
		if json.NewDecoder(r.Body).Decode(&d) != nil {
			writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
			return
		}
		s.Lock()
		defer s.Unlock()
		if s.Transactions[d.TransactionID] == nil {
			writeJSON(w, 404, map[string]string{"error": "transaction not found"})
			return
		}
		d.ID = fmt.Sprintf("DEC-%d-%d", s.Counter, time.Now().UnixNano())
		if d.Version == "" {
			d.Version = "2026.1"
		}
		if d.Status == "" {
			d.Status = "Customer review pending"
		}
		if d.Fields == nil {
			d.Fields = map[string]string{}
		}
		s.Declarations[d.TransactionID] = append(s.Declarations[d.TransactionID], d)
		persistStore()
		writeJSON(w, 201, d)
	})
	mux.HandleFunc("/api/v1/documents/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		if err := r.ParseMultipartForm(25 << 20); err != nil {
			writeJSON(w, 400, map[string]string{"error": "invalid multipart form"})
			return
		}
		txID := r.FormValue("transaction_id")
		docType := strings.TrimSpace(r.FormValue("document_type"))
		if txID == "" || docType == "" {
			writeJSON(w, 400, map[string]string{"error": "transaction_id and document_type are required"})
			return
		}
		s.RLock()
		_, ok := s.Transactions[txID]
		s.RUnlock()
		if !ok {
			writeJSON(w, 404, map[string]string{"error": "transaction not found"})
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": "file is required"})
			return
		}
		defer file.Close()
		if header.Size > 25<<20 {
			writeJSON(w, 413, map[string]string{"error": "file exceeds 25 MB limit"})
			return
		}
		data, err := io.ReadAll(http.MaxBytesReader(w, file, 25<<20))
		if err != nil {
			writeJSON(w, 400, map[string]string{"error": "unable to read file"})
			return
		}
		sum := sha256.Sum256(data)
		docID := fmt.Sprintf("DOC-%d-%d", s.Counter, time.Now().UnixNano())
		base := filepath.Base(header.Filename)
		storageDir := os.Getenv("TCC_DOCUMENT_STORAGE")
		if storageDir == "" {
			storageDir = "./storage/documents"
		}
		if err := os.MkdirAll(filepath.Join(storageDir, txID), 0750); err != nil {
			writeJSON(w, 500, map[string]string{"error": "unable to prepare document storage"})
			return
		}
		storageName := docID + "-" + base
		storagePath := filepath.Join(storageDir, txID, storageName)
		if err := os.WriteFile(storagePath, data, 0600); err != nil {
			writeJSON(w, 500, map[string]string{"error": "unable to store document"})
			return
		}
		d := Document{ID: docID, TransactionID: txID, Type: docType, FileName: base, Status: "Review Required", Version: 1, StorageKey: filepath.ToSlash(filepath.Join(txID, storageName)), MIMEType: header.Header.Get("Content-Type"), FileSize: int64(len(data)), SHA256: hex.EncodeToString(sum[:]), UploadedBy: r.FormValue("uploaded_by")}
		s.Lock()
		s.Documents[txID] = append(s.Documents[txID], d)
		s.Unlock()
		persistStore()
		writeJSON(w, 201, d)
	})

	mux.HandleFunc("/api/v1/documents", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.NotFound(w, r)
			return
		}
		var d Document
		if json.NewDecoder(r.Body).Decode(&d) != nil {
			writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
			return
		}
		s.Lock()
		defer s.Unlock()
		if s.Transactions[d.TransactionID] == nil {
			writeJSON(w, 404, map[string]string{"error": "transaction not found"})
			return
		}
		d.ID = fmt.Sprintf("DOC-%d-%d", s.Counter, time.Now().UnixNano())
		if d.Version == 0 {
			d.Version = 1
		}
		if d.Status == "" {
			d.Status = "Review Required"
		}
		s.Documents[d.TransactionID] = append(s.Documents[d.TransactionID], d)
		persistStore()
		writeJSON(w, 201, d)
	})
	mux.HandleFunc("/api/v1/reconciliation/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 4 {
			http.NotFound(w, r)
			return
		}
		id := parts[3]
		s.RLock()
		defer s.RUnlock()
		t := s.Transactions[id]
		if t == nil {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, 200, reconciliationFor(t, s.Payments[id]))
	})
	mux.HandleFunc("/api/v1/actions", func(w http.ResponseWriter, r *http.Request) {
		s.RLock()
		defer s.RUnlock()
		var out []map[string]any
		for _, t := range s.Transactions {
			for _, a := range actionsFor(t, s.Payments[t.ID]) {
				out = append(out, map[string]any{"transaction": t.ReferenceNo, "product": t.ProductType, "action": a})
			}
		}
		writeJSON(w, 200, out)
	})
	mux.HandleFunc("/api/v1/edpms", func(w http.ResponseWriter, r *http.Request) {
		s.RLock()
		defer s.RUnlock()
		var out []*Transaction
		for _, t := range s.Transactions {
			if t.ProductCode == "EXP_GOODS" || t.ProductCode == "EXP_SERVICE" || t.ProductCode == "EXPORT_COLLECTION" {
				out = append(out, t)
			}
		}
		writeJSON(w, 200, map[string]any{"position": "Customer-recorded position", "transactions": out})
	})
	mux.HandleFunc("/api/v1/idpms", func(w http.ResponseWriter, r *http.Request) {
		s.RLock()
		defer s.RUnlock()
		var out []*Transaction
		for _, t := range s.Transactions {
			if t.ProductCode == "IMP_GOODS" || t.ProductCode == "HSS" || t.ProductCode == "IMPORT_LC" {
				out = append(out, t)
			}
		}
		writeJSON(w, 200, map[string]any{"position": "Customer-recorded position", "transactions": out})
	})
	mux.HandleFunc("/api/v1/regulatory/updates", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}
		writeJSON(w, 200, map[string]any{"as_of": "2026-09-19", "updates": regulatoryUpdates2026})
	})

	mux.HandleFunc("/api/v1/knowledge", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, []map[string]string{{"title": "RBI / FEMA", "source": "RBI", "note": "Rules are versioned and event-date aware."}, {"title": "eBRC", "source": "DGFT", "note": "DGFT's eBRC system uses bank-transmitted IRMs and exporter self-certification."}, {"title": "ICC trade rules", "source": "ICC", "note": "UCP/ISBP/URDG/ISP apply only where relevant and incorporated."}, {"title": "Purpose codes", "source": "RBI / bank data", "note": "Use nature-driven suggestions; customer confirms with AD bank."}})
	})
	log.Println("Trade Compliance Cockpit V2 API listening on :8080")
	limit := 120
	if raw := os.Getenv("TCC_RATE_LIMIT_PER_MINUTE"); raw != "" {
		if n, e := strconv.Atoi(raw); e == nil && n > 0 {
			limit = n
		}
	}
	rate := newRateLimiter(limit)
	handler := securityHeaders(requestSecurity(rate, authenticate(secureCORS(mux))))
	log.Fatal(http.ListenAndServe(":8080", handler))
}
func reconciliationFor(t *Transaction, payments []Payment) ReconciliationView {
	v := ReconciliationView{Mode: "document/payment reconciliation", BaseValue: maxBase(t), Currency: t.Currency, Notes: []string{}}
	if t.ProductCode == "MTT" {
		v.Mode = "two-leg merchanting reconciliation"
		for _, p := range payments {
			if strings.Contains(strings.ToLower(p.PaymentType), "outward") {
				v.OutwardValue += p.Amount
			}
			if strings.Contains(strings.ToLower(p.PaymentType), "inward") {
				v.InwardValue += p.Amount
			}
		}
		v.RecordedValue = v.InwardValue
		v.Outstanding = math.Max(v.BaseValue-v.InwardValue, 0)
		v.Difference = math.Max(v.InwardValue-v.BaseValue, 0)
		v.Status = "Open"
		if v.OutwardValue > 0 && v.InwardValue > 0 && v.OutwardValue == v.InwardValue && v.Outstanding == 0 {
			v.Status = "Reconciled"
		}
		if v.OutwardValue > 0 && v.InwardValue > 0 && v.OutwardValue != v.InwardValue {
			v.Notes = append(v.Notes, "Purchase-leg outlay and sale-leg receipt differ; review the two legs separately.")
		}
		return v
	}
	if t.ProductCode == "FDI" || t.ProductCode == "ODI" || t.ProductCode == "ECB" || t.ProductCode == "BG" || t.ProductCode == "EPC_PCFC" || t.ProductCode == "TRADE_CREDIT" {
		v.Mode = "product exposure / monitoring"
		var recorded float64
		for _, p := range payments {
			recorded += p.Amount
		}
		v.RecordedValue = recorded
		v.Outstanding = 0
		v.Difference = 0
		v.Status = "Monitoring"
		v.Notes = append(v.Notes, "Recorded cash movements are shown separately from the product exposure; this is not an automatic regulatory closure conclusion.")
		return v
	}
	var inward, outward float64
	for _, p := range payments {
		pt := strings.ToLower(p.PaymentType)
		if strings.Contains(pt, "inward") || strings.Contains(pt, "receipt") || strings.Contains(pt, "realization") || strings.Contains(pt, "advance / export") || strings.Contains(pt, "service realization") {
			inward += p.Amount
		} else {
			outward += p.Amount
		}
	}
	v.InwardValue, v.OutwardValue = inward, outward
	v.RecordedValue = inward
	v.Outstanding = math.Max(v.BaseValue-inward, 0)
	v.Difference = math.Max(inward-v.BaseValue, 0)
	v.Status = "Open"
	if v.Outstanding == 0 && v.Difference == 0 && v.BaseValue > 0 {
		v.Status = "Reconciled"
	}
	if v.Difference > 0 {
		v.Notes = append(v.Notes, "Difference identified. Review allocation/supporting records with the AD bank where required.")
	}
	return v
}

func maxBase(t *Transaction) float64 {
	if t.InvoiceValue > 0 {
		return t.InvoiceValue
	}
	return t.DeclaredValue
}
func regulatoryFor(t *Transaction) map[string]string {
	m := map[string]string{"position": "Customer-recorded", "authority_action": "Confirm with AD bank / competent authority", "rule_engine": "Versioned + effective-date aware"}
	if t.ProductCode == "EXP_GOODS" || t.ProductCode == "EXP_SERVICE" || t.ProductCode == "EXPORT_COLLECTION" {
		m["edpms"] = "Customer-recorded position"
	}
	if t.ProductCode == "IMP_GOODS" || t.ProductCode == "HSS" || t.ProductCode == "IMPORT_LC" {
		m["idpms"] = "Customer-recorded position"
	}
	if t.ProductCode == "EXP_SERVICE" {
		m["ebrc"] = "Readiness tracking; no live DGFT filing in this prototype"
	}
	return m
}
func workflowFor(t *Transaction) []string {
	m := map[string][]string{"EXP_GOODS": {"Contract", "Invoice / PI", "Advance / IRM", "Shipment", "Export documents", "Reconciliation", "EDPMS", "Declaration", "Closure"}, "EXP_SERVICE": {"Service contract", "Invoice", "Inward remittance / IRM", "Service/software declaration", "Reconciliation", "EDPMS", "eBRC readiness", "Closure"}, "IMP_GOODS": {"Contract", "Supplier invoice", "ORM / payment", "Import", "BOE / evidence", "Reconciliation", "IDPMS", "Declaration", "Closure"}, "HSS": {"Purchase leg", "HSS agreement", "Sale leg", "Documents", "Payment / settlement", "Reconciliation", "Regulatory review", "Closure"}, "IMPORT_LC": {"LC request", "Issuance", "Shipment", "Documents", "Bill", "Payment", "BOE", "IDPMS / closure"}, "BG": {"Request", "Guarantee issuance", "Margin / collateral", "Amendment / monitoring", "Claim / expiry", "Closure"}, "EPC_PCFC": {"Finance request", "Sanction / drawdown", "Export", "Realization", "Adjustment", "Closure"}, "EXPORT_COLLECTION": {"LC / Non-LC terms", "Bill lodgment", "Document dispatch", "Overseas bank", "Payment / realization", "Reconciliation", "Closure"}, "TRADE_CREDIT": {"Credit arrangement", "Shipment / utilization", "Payment", "Repayment", "Monitoring", "Closure"}, "MTT": {"Purchase contract", "Outward payment", "Sale contract", "Inward receipt", "Reconciliation", "Monitoring", "Closure"}, "FDI": {"Investment decision", "Inward funds", "Securities / investment", "Reporting preparation", "Records / closure"}, "ODI": {"Investment decision", "UIN / AD process", "Outward remittance", "Evidence / reporting", "APR / monitoring", "Exit / closure"}, "ECB": {"Borrowing arrangement", "Drawdown", "Utilization", "Reporting", "Repayment", "Closure"}}
	return m[t.ProductCode]
}
func productGuidance(code string) []string {
	g := map[string][]string{
		"EXP_GOODS":         {"Record contract/order, invoice/PI, shipment/export evidence, inward receipts and reconciliation. Track EDPMS position and applicable declarations; bank checklist is optional."},
		"EXP_SERVICE":       {"Record service/software contract, invoice, inward remittance/IRM and applicable declaration. Track EDPMS and eBRC readiness where relevant."},
		"IMP_GOODS":         {"Record contract, supplier invoice, ORM/payment, BOE/evidence of import and reconciliation. Track IDPMS customer position."},
		"HSS":               {"Record purchase leg, HSS arrangement, sale leg, document trail and fund movement. Keep HSS distinct from ordinary import."},
		"IMPORT_LC":         {"Record LC type/terms, issuance request, shipment, document presentation, bill, payment and BOE/evidence. Applicable ICC rules depend on the instrument and incorporation."},
		"BG":                {"Record guarantee type, beneficiary, amount, validity, margin/collateral, amendments, claims and expiry. Applicable guarantee rules depend on wording and incorporation."},
		"EPC_PCFC":          {"Record facility, drawdown, export linkage, utilization, realization and adjustment. Do not treat financing data as export realization automatically."},
		"EXPORT_COLLECTION": {"Record LC/non-LC terms, bill lodgment, document dispatch, overseas bank status, payment and realization."},
		"TRADE_CREDIT":      {"Record buyer/supplier credit arrangement, shipment, drawdown/payment, repayment schedule and regulatory monitoring."},
		"MTT":               {"Record purchase leg, outward payment, sale leg, inward receipt, reconciliation and monitoring period; fund outlay is tracked separately."},
		"FDI":               {"Record investment party and instrument, inward funds, securities/investment records and applicable reporting preparation. Counterparty information is customer-entered only."},
		"ODI":               {"Record investment, designated AD/UIN information where applicable, outward remittance, evidence, reporting and exit/monitoring."},
		"ECB":               {"Record borrowing arrangement, lender, drawdowns, end-use, reporting, repayments and closure; current rules must be evaluated by effective date."},
	}
	return g[code]
}
func productData(code string) []string {
	m := map[string][]string{
		"EXP_GOODS":         {"Counterparty", "Currency/value", "Contract/order", "Invoice/PI", "Shipment/export evidence", "Inward remittance/IRM", "Reconciliation"},
		"EXP_SERVICE":       {"Counterparty", "Currency/value", "Service/software contract", "Invoice", "Inward remittance/IRM", "Applicable declaration", "eBRC readiness"},
		"IMP_GOODS":         {"Supplier", "Currency/value", "Contract", "Supplier invoice", "ORM/payment", "BOE/evidence of import", "Reconciliation"},
		"HSS":               {"Original supplier", "Indian HSS buyer", "Purchase value", "Sale value", "HSS documents", "Payment trail"},
		"IMPORT_LC":         {"Applicant", "Beneficiary", "LC type", "Amount/currency", "Terms", "Shipment", "Documents", "Bill", "BOE"},
		"BG":                {"Applicant", "Beneficiary", "Guarantee type", "Amount", "Validity", "Margin/collateral", "Claim/expiry"},
		"EPC_PCFC":          {"Facility", "Currency", "Amount", "Export linkage", "Drawdown", "Realization", "Adjustment"},
		"EXPORT_COLLECTION": {"LC/non-LC", "Buyer", "Bill amount", "Documents", "Overseas bank", "Realization"},
		"TRADE_CREDIT":      {"Buyer/supplier credit", "Lender", "Amount", "Shipment", "Drawdown", "Repayment"},
		"MTT":               {"Overseas supplier", "Overseas buyer", "Purchase value", "Sale value", "Outward payment", "Inward receipt"},
		"FDI":               {"Foreign investor", "Instrument", "Amount/currency", "Investment date", "Securities", "Reporting records"},
		"ODI":               {"Overseas entity", "Investment amount", "UIN/AD details where applicable", "Remittance", "Evidence", "APR/reporting"},
		"ECB":               {"Lender", "Currency", "Amount", "Drawdown", "End-use", "Reporting", "Repayment"},
	}
	return m[code]
}
func productRequirements(code string) []map[string]any {
	m := map[string][]map[string]any{
		"EXP_GOODS":         {{"key": "counterparty", "label": "Overseas buyer", "required": true}, {"key": "currency", "label": "Currency", "required": true}, {"key": "declared_value", "label": "Transaction value", "required": true}, {"key": "contract_value", "label": "Contract / order value", "required": false}, {"key": "description", "label": "Goods / transaction description", "required": true}},
		"EXP_SERVICE":       {{"key": "counterparty", "label": "Overseas customer", "required": true}, {"key": "currency", "label": "Currency", "required": true}, {"key": "declared_value", "label": "Service value", "required": true}, {"key": "description", "label": "Service / software description", "required": true}, {"key": "contract_value", "label": "Contract value", "required": false}},
		"IMP_GOODS":         {{"key": "counterparty", "label": "Overseas supplier", "required": true}, {"key": "currency", "label": "Currency", "required": true}, {"key": "declared_value", "label": "Import value", "required": true}, {"key": "description", "label": "Goods description", "required": true}},
		"HSS":               {{"key": "counterparty", "label": "Relevant HSS counterparty", "required": true}, {"key": "currency", "label": "Currency", "required": true}, {"key": "declared_value", "label": "Transaction value", "required": true}, {"key": "description", "label": "HSS description", "required": true}},
		"IMPORT_LC":         {{"key": "counterparty", "label": "Beneficiary / supplier", "required": true}, {"key": "currency", "label": "LC currency", "required": true}, {"key": "declared_value", "label": "LC amount", "required": true}, {"key": "description", "label": "Goods / LC purpose", "required": true}},
		"BG":                {{"key": "counterparty", "label": "Beneficiary", "required": true}, {"key": "currency", "label": "Guarantee currency", "required": true}, {"key": "declared_value", "label": "Guarantee amount", "required": true}, {"key": "description", "label": "Guarantee purpose", "required": true}},
		"EPC_PCFC":          {{"key": "counterparty", "label": "Export counterparty / linkage", "required": false}, {"key": "currency", "label": "Facility currency", "required": true}, {"key": "declared_value", "label": "Facility amount", "required": true}, {"key": "description", "label": "Facility / export purpose", "required": true}},
		"EXPORT_COLLECTION": {{"key": "counterparty", "label": "Overseas buyer", "required": true}, {"key": "currency", "label": "Bill currency", "required": true}, {"key": "declared_value", "label": "Bill value", "required": true}, {"key": "description", "label": "Collection description", "required": true}},
		"TRADE_CREDIT":      {{"key": "counterparty", "label": "Lender / credit counterparty", "required": true}, {"key": "currency", "label": "Currency", "required": true}, {"key": "declared_value", "label": "Credit amount", "required": true}, {"key": "description", "label": "Credit purpose", "required": true}},
		"MTT":               {{"key": "counterparty", "label": "Primary counterparty", "required": true}, {"key": "currency", "label": "Currency", "required": true}, {"key": "declared_value", "label": "Transaction value", "required": true}, {"key": "description", "label": "Purchase/sale description", "required": true}},
		"FDI":               {{"key": "counterparty", "label": "Foreign investor", "required": true}, {"key": "currency", "label": "Investment currency", "required": true}, {"key": "declared_value", "label": "Investment value", "required": true}, {"key": "description", "label": "Investment description", "required": true}},
		"ODI":               {{"key": "counterparty", "label": "Overseas entity", "required": true}, {"key": "currency", "label": "Currency", "required": true}, {"key": "declared_value", "label": "Investment value", "required": true}, {"key": "description", "label": "Investment description", "required": true}},
		"ECB":               {{"key": "counterparty", "label": "Overseas lender", "required": true}, {"key": "currency", "label": "Borrowing currency", "required": true}, {"key": "declared_value", "label": "Borrowing amount", "required": true}, {"key": "description", "label": "Borrowing / end-use description", "required": true}},
	}
	return m[code]
}
func productDocuments(code string) []string {
	m := map[string][]string{
		"EXP_GOODS":         {"Contract / order", "Invoice / PI", "Shipping / export evidence", "Inward remittance / IRM"},
		"EXP_SERVICE":       {"Service / software contract", "Invoice", "Inward remittance / IRM", "Applicable declaration"},
		"IMP_GOODS":         {"Import contract", "Supplier invoice", "ORM / payment evidence", "BOE / evidence of import"},
		"HSS":               {"Purchase documents", "HSS agreement / sale document", "Sale invoice", "Payment trail"},
		"IMPORT_LC":         {"LC application / terms", "Commercial invoice", "Transport / shipment documents", "Bill / document presentation", "BOE / evidence"},
		"BG":                {"Guarantee request", "Underlying contract", "Guarantee text", "Margin / collateral records", "Amendment / claim records"},
		"EPC_PCFC":          {"Facility documents", "Sanction / drawdown", "Export linkage", "Realization evidence"},
		"EXPORT_COLLECTION": {"Collection instruction", "Invoice", "Transport documents", "Bill / collection documents", "Realization evidence"},
		"TRADE_CREDIT":      {"Credit agreement", "Shipment / invoice", "Drawdown / payment", "Repayment records"},
		"MTT":               {"Purchase contract", "Outward payment evidence", "Sale contract / invoice", "Inward receipt evidence"},
		"FDI":               {"Investment agreement / instrument documents", "Inward funds evidence", "Securities / allotment records", "Applicable reporting records"},
		"ODI":               {"Investment documents", "UIN / AD records where applicable", "Remittance evidence", "Evidence / reporting records"},
		"ECB":               {"Loan agreement", "Drawdown evidence", "End-use records", "Reporting records", "Repayment records"},
	}
	return m[code]
}
func productHelp(code string) []string {
	return []string{"Only information relevant to this transaction is requested.", "Bank checklists and bank-prescribed forms are optional uploads; if provided, the app can organize and pre-fill customer-review data.", "The app prepares and explains. The customer and AD bank / competent authority remain responsible for final decisions and submissions."}
}

func helpFor(t *Transaction) []string {
	return []string{"What needs attention? Review the Action Center.", "What is the regulatory status? This app records the customer position; it does not connect to or replace the bank/authority system.", "Which documents are needed? Requirements are contextual and may depend on transaction nature and the AD bank checklist."}
}
