package wikidata

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Select all politicians, aka people with a abgeordnetenwatch.de id (P5355)
// You can edit this using https://query.wikidata.org/
const (
	poliquery = `SELECT DISTINCT ?item ?page_title ?article_url ?name ?first_name ?last_name ?partyHashtag ?partyTwittername ?partyShortname WHERE {
  ?item wdt:P5355 ?value;
    wdt:P1559 ?name.
  ?article schema:about ?item;
    schema:isPartOf <https://de.wikipedia.org/>;
    schema:name ?page_title.
  ?article_url schema:about ?item;
    schema:isPartOf <https://de.wikipedia.org/>.
  OPTIONAL {
    ?item wdt:P735 ?fval.
    ?fval wdt:P1705 ?first_name.
  }
  OPTIONAL {
    ?item wdt:P734 ?lval.
    ?lval wdt:P1705 ?last_name.
  }
  OPTIONAL {
    ?item wdt:P102 ?pval.
    ?pval wdt:P2572 ?partyHashtag.
  }
  OPTIONAL {
    ?item wdt:P102 ?pval.
    ?pval wdt:P2002 ?partyTwittername.
  }
  OPTIONAL {
    ?item wdt:P102 ?pval.
    ?pval wdt:P1813 ?partyShortname.
  }
  SERVICE wikibase:label { bd:serviceParam wikibase:language "de". }
}`

	queryURLPrefix = "https://query.wikidata.org/sparql?format=json&query="
)

var c = http.Client{
	Timeout: 30 * time.Second,
}

// Politicians returns a politicians store that contains all politicians that have an abgeordnetenwatch.de ID assigned to them on WikiData
func Politicians() (store PoliticianStore, err error) {
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

	store = PoliticianStore{
		politicians: make(map[string]Politician, len(data.Results.Bindings)),
	}

	for _, result := range data.Results.Bindings {
		p := result.toPoli()

		currentPoli, ok := store.politicians[p.WikiPageTitle]
		if ok {
			// Some pages are in there twice because of translations of certain fields.
			// Some politicians have two names, which results in two rows in the data.
			// Some are also in there twice because they changed party. It seems like keeping the first record of the data
			// is the correct way to go, but this might not be true for every politicians.
			// This should be fixed in the SPARQL query rather than here, but I don't really know how

			// If we have a first name that is more accurate, we use that
			// E.g. Andreas Scheuer is first seen as Franz Scheuer, but this fixes it
			if p.FirstName != "" && strings.HasPrefix(currentPoli.WikiPageTitle, p.FirstName) {
				goto breakout
			}
			continue
		}
	breakout:
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
	FirstName struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"first_name"`
	LastName struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"last_name"`
	PartyHashtag struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"partyHashtag"`
	PartyTwittername struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"partyTwittername"`
	PartyShortname struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"partyShortname"`
}

func (i *info) toPoli() Politician {
	var p = Politician{
		Name:             i.Name.Value,
		WikiPageTitle:    i.PageTitle.Value,
		WikiArticleURL:   i.ArticleURL.Value,
		FirstName:        i.FirstName.Value,
		LastName:         i.LastName.Value,
		partyHashtag:     i.PartyHashtag.Value,
		partyShortname:   i.PartyShortname.Value,
		partyTwittername: i.PartyTwittername.Value,
	}

	// Sometimes the first name/last name is missing, so we try to infer it from other information
	if p.FirstName == "" && p.LastName != "" {
		p.FirstName = strings.TrimSpace(strings.TrimSuffix(p.Name, p.LastName))
	}
	if p.LastName == "" && p.FirstName != "" {
		p.LastName = strings.TrimSpace(strings.TrimPrefix(p.Name, p.FirstName))
	}
	if p.FirstName == "" && p.LastName == "" {
		if f := strings.Fields(p.Name); len(f) == 2 {
			p.FirstName, p.LastName = f[0], f[1]
		}
	}

	return p
}
