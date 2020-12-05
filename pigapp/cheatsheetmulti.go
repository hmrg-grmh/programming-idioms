package main

import (
	"net/http"
	"sort"
	"strconv"

	. "github.com/Deleplace/programming-idioms/pig"
	"github.com/gorilla/mux"
	gaesearch "google.golang.org/appengine/search"
)

// CheatSheetMultipleFacade is the Facade for the Cheat Sheets with multiple languages.
type CheatSheetMultipleFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
	Langs       []string
	Lines       []cheatSheetLineMulti
}

type cheatSheetLineMulti struct {
	IdiomID            gaesearch.Atom // TODO string?
	IdiomTitle         gaesearch.Atom // TODO string?
	IdiomLeadParagraph gaesearch.Atom // TODO string?

	// Depth represents the max number of impl for one language in ByLanguage
	Depth []struct{}

	ByLanguage []cheatSheetLineDocs
}

func (chsh *cheatSheetLineMulti) hasCell(langIndex, alt int) bool {
	return alt < len(chsh.Depth) &&
		langIndex < len(chsh.ByLanguage) &&
		alt < len(chsh.ByLanguage[langIndex])
}

func (chsh *cheatSheetLineMulti) computeDepth() {
	n := 0
	for _, lineDocs := range chsh.ByLanguage {
		if len(lineDocs) > n {
			n = len(lineDocs)
		}
	}
	chsh.Depth = make([]struct{}, n)
}

func cheatsheetDouble(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	langs := []string{vars["lang1"], vars["lang2"]}
	ctx := r.Context()

	// Security belt. Might be changed if needed.
	limit := 1000

	var lines []cheatSheetLineMulti

	byIdiomID := map[int][]cheatSheetLineDocs{}
	idiomTitles := map[int]gaesearch.Atom{}
	idiomLeadParagraphs := map[int]gaesearch.Atom{}
	for langIndex, lang := range langs {
		cheatsheetLines, err := dao.getCheatSheet(ctx, lang, limit)
		if err != nil {
			return PiError{err.Error(), http.StatusInternalServerError}
		}
		for _, line := range cheatsheetLines {
			idiomID, err := strconv.Atoi(string(line.IdiomID))
			if err != nil {
				return err
			}
			if _, exists := byIdiomID[idiomID]; exists {
				byIdiomID[idiomID][langIndex] = append(byIdiomID[idiomID][langIndex], line)
			} else {
				byIdiomID[idiomID] = make([]cheatSheetLineDocs, len(langs))
				byIdiomID[idiomID][langIndex] = cheatSheetLineDocs{line}
				idiomTitles[idiomID] = (line.IdiomTitle)
				idiomLeadParagraphs[idiomID] = (line.IdiomLeadParagraph)
			}

		}
	}
	idiomIDs := make([]int, 0, len(byIdiomID))
	for id := range byIdiomID {
		idiomIDs = append(idiomIDs, id)
	}
	sort.Ints(idiomIDs)
	for _, idiomID := range idiomIDs {
		line := cheatSheetLineMulti{
			IdiomID:            gaesearch.Atom(strconv.Itoa(idiomID)),
			IdiomTitle:         idiomTitles[idiomID],
			IdiomLeadParagraph: idiomLeadParagraphs[idiomID],
			// Depth:              make([]struct{}, 2),
			ByLanguage: byIdiomID[idiomID],
		}
		line.computeDepth()
		lines = append(lines, line)
	}

	pageTitle := ""
	glue := ""
	for _, lang := range langs {
		pageTitle += glue + PrintNiceLang(lang)
		glue = ", "
	}
	pageTitle += " cheat sheet"
	data := CheatSheetMultipleFacade{
		PageMeta: PageMeta{
			PageTitle: pageTitle,
			Toggles:   toggles,
			ExtraCss:  []string{hostPrefix() + themeDirectory() + "/css/pages/cheatsheetmulti.css"},
		},
		UserProfile: readUserProfile(r),
		Langs:       langs,
		Lines:       lines,
	}

	if err := templates.ExecuteTemplate(w, "page-cheatsheet-multi", data); err != nil {
		return PiError{err.Error(), http.StatusInternalServerError}
	}
	return nil
}