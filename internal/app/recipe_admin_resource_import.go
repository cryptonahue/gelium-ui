package app

import (
	"encoding/csv"
	"fmt"
	"geliumui/lib"
	"io"
	"net/http"
	"strings"
)

const recipeAdminImportMaxRows = 50

type recipeAdminImportView struct {
	AssetsVersion string
	Meta          metaView
	ThemeClass    string
	DataTheme     string
	Error         string
	ImportMaxRows int
	ListHref      string
}

func (s *server) recipeAdminResourceImport(w http.ResponseWriter, r *http.Request) {
	view := recipeAdminImportView{
		AssetsVersion: lib.AssetsVersion,
		Meta:          recipeAdminMeta("Import projects · Admin Resource recipe", "Import a bounded CSV file into the Admin Resource demo.", "/recipes/admin-resource/import"),
		ThemeClass:    themeClass(""),
		ImportMaxRows: recipeAdminImportMaxRows,
		ListHref:      "/recipes/admin-resource",
	}
	applyRequestChrome(r, &view)
	s.renderRecipeTemplate(w, http.StatusOK, "recipe-admin-resource-import", view)
}

func (s *server) recipeAdminResourceImportCSV(w http.ResponseWriter, r *http.Request) {
	view := recipeAdminImportView{
		AssetsVersion: lib.AssetsVersion,
		Meta:          recipeAdminMeta("Import projects · Admin Resource recipe", "Import a bounded CSV file into the Admin Resource demo.", "/recipes/admin-resource/import"),
		ThemeClass:    themeClass(""),
		ImportMaxRows: recipeAdminImportMaxRows,
		ListHref:      "/recipes/admin-resource",
	}
	r.Body = http.MaxBytesReader(w, r.Body, 2<<20)
	if err := r.ParseMultipartForm(2 << 20); err != nil {
		view.Error = "Upload a CSV file smaller than 2 MB."
		s.renderRecipeAdminImportError(w, r, view)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		view.Error = "Choose a CSV file before importing."
		s.renderRecipeAdminImportError(w, r, view)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil || !sameStrings(header, []string{"name", "status", "date", "owner"}) {
		view.Error = "The CSV header must be exactly: name,status,date,owner."
		s.renderRecipeAdminImportError(w, r, view)
		return
	}

	type importRow struct{ name, status, date, owner string }
	rows := make([]importRow, 0, recipeAdminImportMaxRows)
	for rowNumber := 2; ; rowNumber++ {
		record, readErr := reader.Read()
		if readErr == io.EOF {
			break
		}
		if readErr != nil || len(record) != 4 {
			view.Error = fmt.Sprintf("Row %d must contain exactly four CSV fields.", rowNumber)
			s.renderRecipeAdminImportError(w, r, view)
			return
		}
		name, status, date, owner := strings.TrimSpace(record[0]), strings.TrimSpace(record[1]), strings.TrimSpace(record[2]), strings.TrimSpace(record[3])
		if errs := validateRecipeResourceForm(name, status, date); len(errs) > 0 {
			view.Error = fmt.Sprintf("Row %d is invalid: %s", rowNumber, errs[0].Message)
			s.renderRecipeAdminImportError(w, r, view)
			return
		}
		if len(rows) >= recipeAdminImportMaxRows {
			view.Error = fmt.Sprintf("Imports are limited to %d rows.", recipeAdminImportMaxRows)
			s.renderRecipeAdminImportError(w, r, view)
			return
		}
		rows = append(rows, importRow{name, status, date, owner})
	}
	if len(rows) == 0 {
		view.Error = "The CSV contains no project rows."
		s.renderRecipeAdminImportError(w, r, view)
		return
	}

	for _, row := range rows {
		resourceDemoStore.create(row.name, row.status, row.date, row.owner)
	}
	resourceDemoStore.setBanner(recipeAdminSuccessBanner("Projects imported", fmt.Sprintf("%d project(s) were added to the projects list.", len(rows))))
	http.Redirect(w, r, "/recipes/admin-resource", http.StatusSeeOther)
}

func (s *server) renderRecipeAdminImportError(w http.ResponseWriter, r *http.Request, view recipeAdminImportView) {
	applyRequestChrome(r, &view)
	s.renderRecipeTemplate(w, http.StatusUnprocessableEntity, "recipe-admin-resource-import", view)
}

func sameStrings(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if strings.TrimSpace(got[i]) != want[i] {
			return false
		}
	}
	return true
}
