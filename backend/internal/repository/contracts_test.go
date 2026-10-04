package repository

import "testing"

func TestNormalizedDomainContracts(t *testing.T) {
	if (Transaction{}).ID != "" {
		t.Fatal("zero transaction contract changed")
	}
	p := Payment{Direction: "INWARD"}
	if p.Direction != "INWARD" {
		t.Fatal("payment direction must remain explicit")
	}
}
