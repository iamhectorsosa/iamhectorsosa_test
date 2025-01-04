package templates

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data any) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func New() *Templates {
	var files []string
	if err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".svg") {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return nil
	}

	if err := filepath.Walk("components", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			files = append(files, path)
		}
		return nil
	}); err != nil {
		return nil
	}

	funcMap := template.FuncMap{
		"formatItems": formatItems,
	}

	tmpl := template.New("").Funcs(funcMap)

	tmpl = template.Must(tmpl.ParseFiles(files...))

	return &Templates{
		templates: tmpl,
	}
}

type processedItem struct {
	Title       string
	Date        string
	Href        string
	Description string
}

func formatItems(items []SectionItem) []processedItem {
	sectionItems := []processedItem{}
	for _, item := range items {
		sectionItems = append(sectionItems, processedItem{
			Title:       item.Title,
			Date:        formatDates(item.StartDate, item.EndDate),
			Href:        item.Href,
			Description: item.Description,
		})
	}
	return sectionItems
}

func formatDates(startDate string, endDate *string) string {
	startYear, err := time.Parse("2006", startDate)
	if err != nil {
		return "Invalid start date format"
	}

	if endDate == nil {
		return fmt.Sprintf("%d", startYear.Year())
	}

	if *endDate == "Now" || *endDate == "now" {
		currentYear := time.Now().Year()
		return fmt.Sprintf("%d — Now · %d yrs", startYear.Year(), currentYear-startYear.Year())
	}

	endYear, err := time.Parse("2006", *endDate)
	if err != nil {
		return "Invalid end date format"
	}

	return fmt.Sprintf("%d — %d · %d yrs", startYear.Year(), endYear.Year(), endYear.Year()-startYear.Year())
}
