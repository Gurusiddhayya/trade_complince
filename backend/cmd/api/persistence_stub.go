package main

import "context"

type RuntimeSnapshot struct {
	Transactions        map[string]*Transaction         `json:"transactions"`
	Payments            map[string][]Payment            `json:"payments"`
	Documents           map[string][]Document           `json:"documents"`
	DocumentExtractions map[string][]DocumentExtraction `json:"document_extractions"`
	Events              map[string][]Event              `json:"events"`
	Declarations        map[string][]Declaration        `json:"declarations"`
	Invoices            map[string][]Invoice            `json:"invoices"`
	ShippingBills       map[string][]ShippingBill       `json:"shipping_bills"`
	BillsOfEntry        map[string][]BillOfEntry        `json:"bills_of_entry"`
	Requirements        map[string][]Requirement        `json:"requirements"`
	Counter             int                             `json:"counter"`
}

type PostgresPersistence struct{}

func openPostgres(context.Context) (*PostgresPersistence, error)            { return nil, nil }
func (*PostgresPersistence) Close() error                                   { return nil }
func (*PostgresPersistence) Load(context.Context) (*RuntimeSnapshot, error) { return nil, nil }
func (*PostgresPersistence) Save(context.Context, *Store) error             { return nil }
func (*PostgresPersistence) SaveNormalized(context.Context, *Store) error   { return nil }

func (*PostgresPersistence) LoadNormalized(context.Context, *Store) (bool, error) { return false, nil }
