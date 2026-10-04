package main

import (
	"fmt"
	"strings"
)

func requirementsFor(t *Transaction, existing []Requirement) []Requirement {
	out := append([]Requirement(nil), existing...)
	if t == nil {
		return out
	}
	has := func(title string) bool {
		for _, r := range out {
			if strings.EqualFold(r.Title, title) {
				return true
			}
		}
		return false
	}
	add := func(cat, title, desc, priority, source string) {
		if !has(title) {
			// Suggestions need stable IDs so a customer can act on them from
			// Transaction 360 without confusing them with persisted records.
			id := fmt.Sprintf("SUGGESTED-%s-%s", t.ID, strings.ToUpper(strings.ReplaceAll(title, " ", "-")))
			out = append(out, Requirement{ID: id, TransactionID: t.ID, Category: cat, Title: title, Description: desc, Status: "NEED_CONFIRMATION", Priority: priority, Source: source})
		}
	}
	if t.InvoiceValue > 0 {
		add("INTERNAL", "Invoice record", "Keep the invoice number, date, currency, value and payment terms linked to this transaction.", "NORMAL", "Record-keeping")
	}
	add("BANK", "Bank request / checklist", "Record any document or information your AD bank has asked for. Upload the bank email/checklist if available.", "NORMAL", "Customer-entered")
	add("INTERNAL", "Transaction documents", "Keep the documents relevant to this transaction in the Document Vault.", "NORMAL", "Record-keeping")
	if strings.HasPrefix(t.ProductCode, "EXP") || t.ProductCode == "EXPORT_COLLECTION" {
		add("BUYER", "Buyer / client requirement", "Record any document, format, certificate or information requested by the overseas buyer/client.", "LOW", "Customer-entered")
	}
	if strings.HasPrefix(t.ProductCode, "IMP") || t.ProductCode == "HSS" || t.ProductCode == "IMPORT_LC" {
		add("SUPPLIER", "Supplier requirement", "Record any document or information requested from the overseas supplier.", "LOW", "Customer-entered")
	}
	return out
}

func requirementSummary(xs []Requirement) map[string]int {
	out := map[string]int{"total": len(xs), "open": 0, "missing": 0, "need_confirmation": 0, "submitted": 0, "completed": 0}
	for _, x := range xs {
		switch x.Status {
		case "MISSING":
			out["missing"]++
		case "NEED_CONFIRMATION":
			out["need_confirmation"]++
		case "SUBMITTED":
			out["submitted"]++
		case "COMPLETED", "ACCEPTED", "WAIVED":
			out["completed"]++
		}
		if x.Status != "COMPLETED" && x.Status != "ACCEPTED" && x.Status != "WAIVED" {
			out["open"]++
		}
	}
	return out
}

func requirementAction(t *Transaction, r Requirement) Action {
	return Action{Priority: "attention", Title: fmt.Sprintf("Requirement: %s", r.Title), Reason: r.Description, RecommendedAction: "Open the requirement, attach the relevant document or information, and update its status.", Status: "Open", DueDate: r.DueDate}
}

func requirementDashboard(s *Store) map[string]int {
	out := map[string]int{"total": 0, "open": 0, "missing": 0, "need_confirmation": 0, "submitted": 0, "completed": 0}
	for id, t := range s.Transactions {
		xs := requirementsFor(t, s.Requirements[id])
		sum := requirementSummary(xs)
		for k, v := range sum {
			out[k] += v
		}
	}
	return out
}
