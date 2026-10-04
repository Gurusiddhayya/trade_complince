package main

import "testing"

func TestRequirementsForAddsCustomerFriendlySuggestions(t *testing.T) {
	tx := &Transaction{ID: "1", ProductCode: "EXP_GOODS", InvoiceValue: 100, DeclaredValue: 100}
	got := requirementsFor(tx, nil)
	if len(got) < 3 {
		t.Fatalf("expected record-keeping/bank/buyer suggestions, got %d", len(got))
	}
	sum := requirementSummary(got)
	if sum["open"] == 0 || sum["need_confirmation"] == 0 {
		t.Fatalf("unexpected summary: %#v", sum)
	}
}

func TestRequirementsForDoesNotDuplicateExistingTitle(t *testing.T) {
	tx := &Transaction{ID: "1", ProductCode: "IMP_GOODS", InvoiceValue: 100}
	existing := []Requirement{{ID: "REQ-1", TransactionID: "1", Category: "BANK", Title: "Bank request / checklist", Status: "MISSING"}}
	got := requirementsFor(tx, existing)
	count := 0
	for _, r := range got {
		if r.Title == "Bank request / checklist" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("expected one bank requirement, got %d", count)
	}
}
