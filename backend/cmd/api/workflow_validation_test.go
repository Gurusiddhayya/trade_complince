package main

import (
	"fmt"
	"strings"
	"testing"
)

// Gate 6 validates the customer journey configuration for all 13 products.
// This deliberately tests record-keeping and banking-assistance surfaces rather
// than treating regulatory UI as the primary workflow.
func TestAll13ProductCustomerWorkflows(t *testing.T) {
	for _, p := range products {
		t.Run(p.Code, func(t *testing.T) {
			tx := &Transaction{
				ID:            "TEST-" + p.Code,
				ReferenceNo:   fmt.Sprintf("TEST-%s-000001", p.Code),
				ProductCode:   p.Code,
				ProductType:   p.Name,
				Currency:      "USD",
				DeclaredValue: 1000,
				InvoiceValue:  1000,
				Counterparty:  "Test Counterparty",
				Description:   "Gate 6 customer workflow test",
			}

			if got := workflowFor(tx); len(got) < 4 {
				t.Fatalf("workflow too short: %v", got)
			}
			if got := productData(p.Code); len(got) == 0 {
				t.Fatal("missing product data configuration")
			}
			if got := productDocuments(p.Code); len(got) == 0 {
				t.Fatal("missing product document configuration")
			}
			if got := productRequirements(p.Code); len(got) == 0 {
				t.Fatal("missing product requirement configuration")
			}
			if got := productGuidance(p.Code); len(got) == 0 {
				t.Fatal("missing product guidance")
			}
			if got := productHelp(p.Code); len(got) == 0 {
				t.Fatal("missing product help")
			}

			reqs := requirementsFor(tx, nil)
			if len(reqs) == 0 {
				t.Fatal("transaction did not produce contextual requirements")
			}
			for _, r := range reqs {
				if r.ID == "" || r.TransactionID != tx.ID || r.Title == "" || r.Status == "" {
					t.Fatalf("invalid contextual requirement: %+v", r)
				}
			}
			summary := requirementSummary(reqs)
			if summary["total"] != len(reqs) || summary["open"] == 0 {
				t.Fatalf("unexpected requirement summary: %+v", summary)
			}
		})
	}
}

func TestTransaction360CustomerSurfaces(t *testing.T) {
	for _, p := range products {
		tx := &Transaction{ID: "TEST", ProductCode: p.Code, ProductType: p.Name, Currency: "USD", DeclaredValue: 1000, InvoiceValue: 1000}
		wf := strings.Join(workflowFor(tx), " | ")
		if wf == "" {
			t.Fatalf("%s has no workflow", p.Code)
		}
		// Every product must expose at least one plain-language record/help surface.
		if len(productHelp(p.Code)) == 0 || len(productDocuments(p.Code)) == 0 {
			t.Fatalf("%s missing help/document surfaces", p.Code)
		}
	}
}

func TestRequirementSuggestionsHaveActionableIDs(t *testing.T) {
	tx := &Transaction{ID: "TX-123", ProductCode: "EXP_GOODS", ProductType: "Export of Goods", InvoiceValue: 1000}
	reqs := requirementsFor(tx, nil)
	seen := map[string]bool{}
	for _, r := range reqs {
		if seen[r.ID] {
			t.Fatalf("duplicate requirement ID %q", r.ID)
		}
		seen[r.ID] = true
		if r.ID == "SUGGESTED" || !strings.HasPrefix(r.ID, "SUGGESTED-") {
			t.Fatalf("requirement suggestion has non-actionable ID: %q", r.ID)
		}
	}
}
