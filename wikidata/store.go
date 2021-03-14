package wikidata

type Politician struct {
	Name string

	WikiPageTitle  string
	WikiArticleURL string
}

type PoliStore struct {
	politicians map[string]Politician
}

// Get returns, if possible, a wikipedia article with the given title is in this store
func (s *PoliStore) Get(pageTitle string) (p Politician, ok bool) {
	p, ok = s.politicians[pageTitle]
	return
}

func (s *PoliStore) Contains(pageTitle string) (ok bool) {
	_, ok = s.politicians[pageTitle]
	return
}
