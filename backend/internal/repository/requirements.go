package repository

import "context"

type Requirement struct {
	ID, TransactionID, CompanyID, Category, Title, Description, Status, Priority, Source, RequestedBy, AssignedTo, DueDate, Notes string
}

type RequirementRepository interface {
	Create(context.Context, Requirement) (string, error)
	ListByTransaction(context.Context, string) ([]Requirement, error)
	UpdateStatus(context.Context, string, string, string) error
}
