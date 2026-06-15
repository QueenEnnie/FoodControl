package httpapi

import (
	"net/http/httptest"
	"testing"
)

func TestParseID(t *testing.T) {
	tests := []struct {
		name    string
		pathID  string
		wantID  int64
		wantErr bool
	}{
		{
			name:   "valid id",
			pathID: "42",
			wantID: 42,
		},
		{
			name:    "non numeric id",
			pathID:  "abc",
			wantErr: true,
		},
		{
			name:    "zero id",
			pathID:  "0",
			wantErr: true,
		},
		{
			name:    "negative id",
			pathID:  "-1",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/products/"+tt.pathID, nil)
			req.SetPathValue("id", tt.pathID)

			gotID, err := parseID(req)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if gotID != tt.wantID {
				t.Fatalf("expected id %d, got %d", tt.wantID, gotID)
			}
		})
	}
}
