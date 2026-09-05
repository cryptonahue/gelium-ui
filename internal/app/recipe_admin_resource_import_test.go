package app

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func postRecipeAdminImport(t *testing.T, csv string) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "projects.csv")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte(csv))
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/recipes/admin-resource/import", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res := httptest.NewRecorder()
	New().ServeHTTP(res, req)
	return res
}

func TestRecipeAdminResourceImportCSVCreatesRows(t *testing.T) {
	resetRecipeResourceStore()
	res := postRecipeAdminImport(t, "name,status,date,owner\nImported one,Active,2027-01-02,Ada\nImported two,Done,2027-01-03,Bob\n")
	if res.Code != http.StatusSeeOther {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusSeeOther)
	}
	if got := res.Header().Get("Location"); got != "/recipes/admin-resource" {
		t.Fatalf("Location = %q, want list", got)
	}
	if item, ok := resourceDemoStore.get("imported-one"); !ok || item.Status != "Active" {
		t.Fatalf("imported item = %+v, exists = %v", item, ok)
	}
	getReq := httptest.NewRequest(http.MethodGet, "/recipes/admin-resource", nil)
	getRes := httptest.NewRecorder()
	New().ServeHTTP(getRes, getReq)
	body := getRes.Body.String()
	if !bytes.Contains([]byte(body), []byte("Projects imported")) {
		t.Error("import should report a persistent success banner")
	}
}

func TestRecipeAdminResourceImportCSVRejectsInvalidRowsWithoutMutation(t *testing.T) {
	resetRecipeResourceStore()
	before := len(resourceDemoStore.snapshot())
	res := postRecipeAdminImport(t, "name,status,date,owner\nValid,Active,2027-01-02,Ada\nBroken,Unknown,not-a-date,Bob\n")
	if res.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", res.Code, http.StatusUnprocessableEntity)
	}
	if !bytes.Contains(res.Body.Bytes(), []byte("Row 3 is invalid")) {
		t.Error("invalid import should identify the failing row")
	}
	if got := len(resourceDemoStore.snapshot()); got != before {
		t.Errorf("invalid import changed store size to %d, want %d", got, before)
	}
}
