package main

import (
	"fmt"
	"strconv"
	"strings"
)

func clampConfidence(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
func defaultSource(v string) string {
	if strings.TrimSpace(v) == "" {
		return "OCR"
	}
	return v
}
func upsertExtraction(xs []DocumentExtraction, x DocumentExtraction) []DocumentExtraction {
	for i := range xs {
		if strings.EqualFold(xs[i].FieldName, x.FieldName) {
			xs[i] = x
			return xs
		}
	}
	return append(xs, x)
}

type classificationRule struct {
	keys       []string
	name       string
	confidence float64
}

func classifyDocument(d Document) (string, float64) {
	text := strings.ToLower(d.Type + " " + d.FileName)
	rules := []classificationRule{
		{[]string{"invoice", "commercial invoice", "tax invoice"}, "Invoice", .98},
		{[]string{"purchase order", " po ", "po-"}, "Purchase Order", .95},
		{[]string{"contract", "agreement"}, "Contract / Agreement", .94},
		{[]string{"shipping bill", "shipping_bill"}, "Shipping Bill", .99},
		{[]string{"bill of entry", "boe"}, "Bill of Entry", .98},
		{[]string{"letter of credit", "lc"}, "Letter of Credit", .96},
		{[]string{"guarantee", "bg"}, "Bank Guarantee", .96},
		{[]string{"irm", "remittance advice", "firc", "bank advice"}, "Bank / Remittance Evidence", .92},
		{[]string{"orm", "outward remittance"}, "Outward Remittance Evidence", .94},
	}
	for _, r := range rules {
		for _, k := range r.keys {
			if strings.Contains(text, k) {
				return r.name, r.confidence
			}
		}
	}
	return d.Type, .60
}

func buildDocumentIntelligence(d Document, fields []DocumentExtraction) DocumentIntelligence {
	c, cc := classifyDocument(d)
	status := "Needs extraction / review"
	if len(fields) > 0 {
		status = "Extracted — customer confirmation required"
		all := true
		for _, f := range fields {
			if !f.Confirmed {
				all = false
				break
			}
		}
		if all {
			status = "Confirmed by customer"
		}
	}
	findings := []DocumentFinding{}
	for _, f := range fields {
		if f.Confidence > 0 && f.Confidence < .80 {
			findings = append(findings, DocumentFinding{Severity: "attention", Code: "LOW_CONFIDENCE", Message: fmt.Sprintf("%s was extracted with %.0f%% confidence.", f.FieldName, f.Confidence*100), Fields: []string{f.FieldName}, Action: "Review the source document and confirm or correct the value."})
		}
	}
	return DocumentIntelligence{DocumentID: d.ID, Classification: c, ClassificationConfidence: cc, ExtractionStatus: status, ExtractedFields: fields, Findings: findings}
}

func transactionDocumentIntelligence(txID string, docs []Document, all map[string][]DocumentExtraction, t *Transaction) map[string]any {
	insights := []DocumentIntelligence{}
	for _, d := range docs {
		insights = append(insights, buildDocumentIntelligence(d, all[d.ID]))
	}
	return map[string]any{"transaction_id": txID, "documents": insights, "cross_document_findings": compareDocumentFields(insights), "readiness": documentReadiness(t, docs, insights)}
}

func compareDocumentFields(docs []DocumentIntelligence) []DocumentFinding {
	type fieldValue struct{ doc, value string }
	fields := map[string][]fieldValue{}
	for _, d := range docs {
		for _, f := range d.ExtractedFields {
			if !f.Confirmed || strings.TrimSpace(f.FieldValue) == "" {
				continue
			}
			k := strings.ToLower(strings.TrimSpace(f.FieldName))
			fields[k] = append(fields[k], fieldValue{d.DocumentID, f.FieldValue})
		}
	}
	out := []DocumentFinding{}
	for field, vals := range fields {
		if len(vals) < 2 {
			continue
		}
		base := vals[0].value
		for _, v := range vals[1:] {
			if !sameNormalized(base, v.value) {
				out = append(out, DocumentFinding{Severity: "attention", Code: "FIELD_MISMATCH", Message: fmt.Sprintf("Confirmed %s differs across documents (%q vs %q).", field, base, v.value), Fields: []string{field}, Action: "Review the source documents and resolve the difference before using the value for a declaration or bank submission."})
				break
			}
		}
	}
	return out
}
func sameNormalized(a, b string) bool {
	a = strings.TrimSpace(strings.ToLower(a))
	b = strings.TrimSpace(strings.ToLower(b))
	if a == b {
		return true
	}
	af, ae := strconv.ParseFloat(strings.ReplaceAll(a, ",", ""), 64)
	bf, be := strconv.ParseFloat(strings.ReplaceAll(b, ",", ""), 64)
	return ae == nil && be == nil && af == bf
}
func documentReadiness(t *Transaction, docs []Document, insights []DocumentIntelligence) map[string]any {
	_ = t
	confirmed, review := 0, 0
	for _, d := range insights {
		if d.ExtractionStatus == "Confirmed by customer" {
			confirmed++
		} else if len(d.ExtractedFields) > 0 {
			review++
		}
	}
	status := "Information available for customer review"
	if len(docs) == 0 {
		status = "No documents uploaded"
	} else if review > 0 {
		status = "Review required"
	}
	return map[string]any{"status": status, "documents_uploaded": len(docs), "documents_with_confirmed_extractions": confirmed, "documents_needing_review": review, "note": "Document readiness is contextual. Absence of a document here is not, by itself, a regulatory violation or a bank checklist failure."}
}
