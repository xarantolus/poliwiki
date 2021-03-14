package wikidata

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

// Select all politicians, aka people with a abgeordnetenwatch.de id (P5355)
const (
	poliquery = `SELECT DISTINCT ?item ?page_title ?article_url ?name WHERE {
  ?item wdt:P5355 ?value;
    wdt:P1559 ?name.
  ?article schema:about ?item;
    schema:isPartOf <https://de.wikipedia.org/>;
    schema:name ?page_title.
  ?article_url schema:about ?item;
    schema:isPartOf <https://de.wikipedia.org/>.
  SERVICE wikibase:label { bd:serviceParam wikibase:language "de". }
}`

	queryURLPrefix = "https://query.wikidata.org/sparql?format=json&query="
)

var c = http.Client{
	Timeout: 30 * time.Second,
}

// Politicians returns a politicians store that contains all politicians that have an abgeordnetenwatch.de ID assigned to them on WikiData
func Politicians() (store PoliStore, err error) {
	var queryURL = queryURLPrefix + url.QueryEscape(poliquery)

	resp, err := c.Get(queryURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var data response

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return
	}

	store = PoliStore{
		politicians: make(map[string]Politician, len(data.Results.Bindings)),
	}

	for _, result := range data.Results.Bindings {
		p := result.toPoli()

		store.politicians[p.WikiPageTitle] = p
	}

	return
}

type response struct {
	Head struct {
		Vars []string `json:"vars"`
	} `json:"head"`
	Results result `json:"results"`
}

type result struct {
	Bindings []info `json:"bindings"`
}

type info struct {
	Item struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"item"`
	Name struct {
		XMLLang string `json:"xml:lang"`
		Type    string `json:"type"`
		Value   string `json:"value"`
	} `json:"name"`
	PageTitle struct {
		XMLLang string `json:"xml:lang"`
		Type    string `json:"type"`
		Value   string `json:"value"`
	} `json:"page_title"`
	ArticleURL struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"article_url"`
}

func (i *info) toPoli() Politician {
	return Politician{
		Name:           i.Name.Value,
		WikiPageTitle:  i.PageTitle.Value,
		WikiArticleURL: i.ArticleURL.Value,
	}
}
