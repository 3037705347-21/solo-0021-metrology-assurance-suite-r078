package domain

import (
	"fmt"
	"strings"
	"time"
)

type EvidenceSeal struct {
	ID            string    `json:"id"`
	CaseID        string    `json:"case_id"`
	DeviceID      string    `json:"device_id"`
	MeasurementID string    `json:"measurement_id"`
	Reviewer      string    `json:"reviewer"`
	Summary       string    `json:"summary"`
	SealedAt      time.Time `json:"sealed_at"`
}

func NewEvidenceSeal(id, caseID, deviceID, measurementID, reviewer, summary string, at time.Time) (EvidenceSeal, error) {
	item := EvidenceSeal{
		ID:            strings.TrimSpace(id),
		CaseID:        strings.TrimSpace(caseID),
		DeviceID:      strings.TrimSpace(deviceID),
		MeasurementID: strings.TrimSpace(measurementID),
		Reviewer:      strings.TrimSpace(reviewer),
		Summary:       strings.TrimSpace(summary),
		SealedAt:      at.UTC(),
	}
	if err := item.Validate(); err != nil {
		return EvidenceSeal{}, err
	}
	return item, nil
}

func (s EvidenceSeal) Validate() error {
	if s.ID == "" || s.CaseID == "" || s.DeviceID == "" || s.MeasurementID == "" || s.Reviewer == "" || s.Summary == "" {
		return fmt.Errorf("%w: seal id, case, device, measurement, reviewer, and summary are required", ErrInvalidInput)
	}
	if s.SealedAt.IsZero() {
		return fmt.Errorf("%w: seal time is required", ErrInvalidInput)
	}
	return nil
}

type CustodyTransfer struct {
	ID        string    `json:"id"`
	CaseID    string    `json:"case_id"`
	From      string    `json:"from,omitempty"`
	To        string    `json:"to"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

func NewCustodyTransfer(id, caseID, from, to, reason string, at time.Time) (CustodyTransfer, error) {
	transfer := CustodyTransfer{
		ID:        strings.TrimSpace(id),
		CaseID:    strings.TrimSpace(caseID),
		From:      strings.TrimSpace(from),
		To:        strings.TrimSpace(to),
		Reason:    strings.TrimSpace(reason),
		CreatedAt: at.UTC(),
	}
	if err := transfer.Validate(); err != nil {
		return CustodyTransfer{}, err
	}
	return transfer, nil
}

func (c CustodyTransfer) Validate() error {
	if c.ID == "" || c.CaseID == "" || c.To == "" || c.Reason == "" {
		return fmt.Errorf("%w: custody id, case, receiver, and reason are required", ErrInvalidInput)
	}
	if c.From == c.To {
		return fmt.Errorf("%w: custody cannot stay with the same receiver", ErrInvalidInput)
	}
	if c.CreatedAt.IsZero() {
		return fmt.Errorf("%w: custody time is required", ErrInvalidInput)
	}
	return nil
}

func (c CustodyTransfer) IsOrigin() bool {
	return c.From == ""
}

// CurrentCustodian reports the receiver of the latest transfer in the chain.
func CurrentCustodian(chain []CustodyTransfer) (string, bool) {
	for index := len(chain) - 1; index >= 0; index-- {
		if chain[index].To != "" {
			return chain[index].To, true
		}
	}
	return "", false
}
