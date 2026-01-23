package models

import (
	"testing"
)

func TestAPIResponse(t *testing.T) {
	response := APIResponse{
		Success: true,
		Message: "Operation successful",
		Data:    map[string]string{"key": "value"},
		Error:   "",
	}

	if !response.Success {
		t.Error("Success should be true")
	}
	if response.Message != "Operation successful" {
		t.Errorf("Message = %v, want Operation successful", response.Message)
	}
	if response.Data == nil {
		t.Error("Data should not be nil")
	}
}

func TestPaginationResponse(t *testing.T) {
	data := []string{"item1", "item2", "item3"}
	pagination := PaginationResponse{
		Data:       data,
		Page:       1,
		Limit:      10,
		Total:      3,
		TotalPages: 1,
	}

	if pagination.Page != 1 {
		t.Errorf("Page = %v, want 1", pagination.Page)
	}
	if pagination.Limit != 10 {
		t.Errorf("Limit = %v, want 10", pagination.Limit)
	}
	if pagination.Total != 3 {
		t.Errorf("Total = %v, want 3", pagination.Total)
	}
	if pagination.TotalPages != 1 {
		t.Errorf("TotalPages = %v, want 1", pagination.TotalPages)
	}
}
