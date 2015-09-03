package pigae

import (
	"fmt"
	"net/http"
	"strings"

	. "github.com/Deleplace/programming-idioms/pig"

	"github.com/gorilla/mux"

	"appengine"
)

//
// This file is about full text search of idioms and implementations.
//

// SearchResultsFacade is the facade for the Search Results page.
type SearchResultsFacade struct {
	PageMeta    PageMeta
	UserProfile UserProfile
	Q           string
	Results     []*Idiom
}

// This is a "word by word" search, not a rdbms "like" filter
func search(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	userProfile := readUserProfile(r)

	c := appengine.NewContext(r)
	q := vars["q"]
	words := SplitForSearching(q, true)

	numberMaxResults := 20
	hits, err := dao.searchIdiomsByWordsWithFavorites(c, words, userProfile.FavoriteLanguages, userProfile.SeeNonFavorite, numberMaxResults)
	if err != nil {
		return err
	}
	normalizedQ := strings.Join(words, " ")

	for _, idiom := range hits {
		implFavoriteLanguagesFirstWithOrder(idiom, userProfile.FavoriteLanguages, "", userProfile.SeeNonFavorite)
	}

	return listResults(w, r, normalizedQ, hits)
}

func listResults(w http.ResponseWriter, r *http.Request, q string, idioms []*Idiom) error {
	data := &SearchResultsFacade{
		PageMeta: PageMeta{
			PageTitle:   "Idioms for \"" + q + "\"",
			Toggles:     toggles,
			SearchQuery: q,
		},
		UserProfile: readUserProfile(r),
		Q:           q,
		Results:     idioms,
	}

	return templates.ExecuteTemplate(w, "page-list-results", data)
}

func searchRedirect(w http.ResponseWriter, r *http.Request) error {
	q := r.FormValue("q")
	if q == "" {
		// Small hack for own convenience: empty search -> homepage
		http.Redirect(w, r, "/", 301)
	}
	q = strings.Replace(q, " ", "+", -1)

	http.Redirect(w, r, hostPrefix()+"/search/"+q, http.StatusMovedPermanently)
	return nil
}

func listByLanguage(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	userProfile := readUserProfile(r)

	c := appengine.NewContext(r)
	langsStr := vars["langs"]
	langs := strings.Split(langsStr, "_")
	langs = MapStrings(langs, normLang)
	langs = RemoveEmptyStrings(langs)

	numberMaxResults := 20
	hits, err := dao.searchIdiomsByLangs(c, langs, numberMaxResults)
	if err != nil {
		return err
	}

	for _, idiom := range hits {
		implFavoriteLanguagesFirstWithOrder(idiom, userProfile.FavoriteLanguages, "", userProfile.SeeNonFavorite)
	}

	niceLangs := MapStrings(langs, printNiceLang)
	return listResults(w, r, fmt.Sprintf("Language=%v", niceLangs), hits)
}
