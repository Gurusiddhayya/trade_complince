package main

import (
	"fmt"
	"strings"
	"time"
)

type DeclarationField struct {
	Name              string `json:"name"`
	Value             string `json:"value,omitempty"`
	SourceType        string `json:"source_type"`
	SourceID          string `json:"source_id,omitempty"`
	CustomerConfirmed bool   `json:"customer_confirmed"`
}

type DeclarationPreparation struct {
	TransactionID      string             `json:"transaction_id"`
	DeclarationType    string             `json:"declaration_type"`
	Applicability      string             `json:"applicability"`
	RuleVersion        string             `json:"rule_version"`
	EffectiveFrom      string             `json:"effective_from"`
	Status             string             `json:"status"`
	Fields             []DeclarationField `json:"fields"`
	MissingFields      []string           `json:"missing_fields"`
	ReviewFindings     []DocumentFinding  `json:"review_findings"`
	SourceTraceability string             `json:"source_traceability"`
	Note               string             `json:"note"`
}

const declarationRuleVersion = "FEMA 23(R)/2026-RB"
const declarationEffectiveFrom = "2026-10-01"

func declarationTypeFor(t *Transaction) (string, string) {
	if t == nil {
		return "", "Not applicable"
	}
	switch t.ProductCode {
	case "EXP_GOODS", "EXP_SERVICE", "EXPORT_COLLECTION", "EPC_PCFC":
		return "EDF", "Potentially applicable export declaration; exact applicability depends on transaction/event facts and prescribed process."
	case "IMP_GOODS", "IMPORT_LC", "TRADE_CREDIT", "HSS", "MTT":
		return "IMPORT_DECLARATION", "Potential import declaration/reporting requirement; exact applicability depends on transaction/event facts and prescribed process."
	default:
		return "", "No generic trade declaration template selected by this engine."
	}
}

func declarationEventDate(t *Transaction) string {
	if t == nil {
		return ""
	}
	// Transaction date is a safe fallback only. Product workflows should provide
	// the specific shipment/invoice/remittance/BOE event date when available.
	if !t.CreatedAt.IsZero() {
		return t.CreatedAt.Format("2006-01-02")
	}
	return ""
}

func declarationApplicability(eventDate string) string {
	if eventDate == "" {
		return "REVIEW_REQUIRED"
	}
	d, err := time.Parse("2006-01-02", eventDate)
	if err != nil {
		return "REVIEW_REQUIRED"
	}
	effective, _ := time.Parse("2006-01-02", declarationEffectiveFrom)
	if d.Before(effective) {
		return "LEGACY_EVENT_REVIEW"
	}
	return "CURRENT_RULE_CANDIDATE"
}

func addDeclarationField(out *[]DeclarationField, name, value, sourceType, sourceID string, confirmed bool) {
	*out = append(*out, DeclarationField{Name: name, Value: value, SourceType: sourceType, SourceID: sourceID, CustomerConfirmed: confirmed})
}

func prepareDeclaration(t *Transaction, docs []Document, extractions map[string][]DocumentExtraction, invoices []Invoice, payments []Payment) DeclarationPreparation {
	typ, applicabilityText := declarationTypeFor(t)
	eventDate := declarationEventDate(t)
	applicability := declarationApplicability(eventDate)
	p := DeclarationPreparation{
		TransactionID: t.ID, DeclarationType: typ, Applicability: applicability,
		RuleVersion: declarationRuleVersion, EffectiveFrom: declarationEffectiveFrom,
		Status: "Customer review pending", SourceTraceability: "Every populated field identifies its source; extracted document values are eligible only after customer confirmation.",
		Note: applicabilityText + " The engine prepares data; it does not certify acceptance by an AD bank or authority.",
	}
	if typ == "" {
		p.Status = "Not applicable"
		return p
	}

	addDeclarationField(&p.Fields, "transaction_reference", t.ReferenceNo, "TRANSACTION", t.ID, true)
	addDeclarationField(&p.Fields, "product_code", t.ProductCode, "TRANSACTION", t.ID, true)
	addDeclarationField(&p.Fields, "counterparty", t.Counterparty, "TRANSACTION", t.ID, true)
	addDeclarationField(&p.Fields, "counterparty_country", t.Country, "TRANSACTION", t.ID, true)
	addDeclarationField(&p.Fields, "currency", t.Currency, "TRANSACTION", t.ID, true)
	addDeclarationField(&p.Fields, "description", t.Description, "TRANSACTION", t.ID, true)
	addDeclarationField(&p.Fields, "declared_value", fmt.Sprintf("%g", t.DeclaredValue), "TRANSACTION", t.ID, true)
	if t.ContractValue > 0 {
		addDeclarationField(&p.Fields, "contract_value", fmt.Sprintf("%g", t.ContractValue), "TRANSACTION", t.ID, true)
	}
	if t.InvoiceValue > 0 {
		addDeclarationField(&p.Fields, "invoice_value", fmt.Sprintf("%g", t.InvoiceValue), "TRANSACTION", t.ID, true)
	}

	for _, i := range invoices {
		addDeclarationField(&p.Fields, "invoice_number", i.InvoiceNumber, "INVOICE", i.ID, true)
		addDeclarationField(&p.Fields, "invoice_date", "", "INVOICE", i.ID, false)
		addDeclarationField(&p.Fields, "invoice_currency", i.Currency, "INVOICE", i.ID, true)
		addDeclarationField(&p.Fields, "invoice_payment_terms", i.PaymentTerms, "INVOICE", i.ID, i.PaymentTerms != "")
		addDeclarationField(&p.Fields, "invoice_description", i.Description, "INVOICE", i.ID, i.Description != "")
		break
	}
	for _, pay := range payments {
		if strings.EqualFold(pay.Direction, "INWARD") || strings.EqualFold(pay.Direction, "OUTWARD") {
			addDeclarationField(&p.Fields, "payment_reference", pay.Reference, "PAYMENT", pay.ID, true)
			addDeclarationField(&p.Fields, "payment_date", pay.PaymentDate, "PAYMENT", pay.ID, pay.PaymentDate != "")
			addDeclarationField(&p.Fields, "payment_currency", pay.Currency, "PAYMENT", pay.ID, pay.Currency != "")
			addDeclarationField(&p.Fields, "payment_amount", fmt.Sprintf("%g", pay.Amount), "PAYMENT", pay.ID, pay.Amount > 0)
			break
		}
	}

	// Promote confirmed document fields into the declaration source set without
	// overwriting transaction-level values. Duplicate field names are retained so
	// the customer can see source conflicts rather than having one silently win.
	for _, d := range docs {
		for _, x := range extractions[d.ID] {
			if !x.Confirmed || strings.TrimSpace(x.FieldValue) == "" {
				continue
			}
			addDeclarationField(&p.Fields, x.FieldName, x.FieldValue, "DOCUMENT_EXTRACTION", d.ID, true)
		}
	}

	for _, f := range p.Fields {
		if strings.TrimSpace(f.Value) == "" {
			p.MissingFields = append(p.MissingFields, f.Name)
		}
	}
	if len(p.MissingFields) > 0 {
		p.Status = "Customer review pending"
	}
	return p
}

func declarationFromPreparation(p DeclarationPreparation, version string) Declaration {
	fields := map[string]string{}
	for _, f := range p.Fields {
		if _, exists := fields[f.Name]; !exists && strings.TrimSpace(f.Value) != "" {
			fields[f.Name] = f.Value
		}
	}
	if version == "" {
		version = "2026.1"
	}
	return Declaration{TransactionID: p.TransactionID, Type: p.DeclarationType, Version: version, Status: p.Status, Fields: fields}
}
