package resolve

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	_URL_ALL_COUNTRIES = "https://univ.cc/world.php"

	_URL_COUNTRY_UNIVERSITY = "https://univ.cc/search.php?dom=%s&key=&start=%d"
)

type UnivClient struct {
}

func NewUnivClient() UnivClient {
	return UnivClient{}
}

// ResolveAllCountries resolves all countries
func (c UnivClient) ResolveAllCountries() ([]Country, error) {
	res, err := http.Get(_URL_ALL_COUNTRIES)
	if err != nil {
		return nil, fmt.Errorf("ResolveAllCountries err:%w", err)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("NewDocumentFromReader err:%w", err)
	}

	cs := make([]Country, 0)
	doc.Find("select[name='dom'] option").Each(func(idx int, s *goquery.Selection) {
		v, ok := s.Attr("value")
		if !ok {
			fmt.Printf("\t not found country value on resolve country info,idx:%d\n", idx)
			return
		}

		if v == "world" {
			return
		}
		name := strings.Split(s.Text(), " (")[0]

		if v == "hk" || v == "tw" { // Hong Kong and Taiwan are part of China
			name = "China " + name
		}

		cs = append(cs, Country{Name: name, Code: v})
	})

	return cs, nil
}

// ResolveUniversityInUnitedState resolves all university in United States
func (c UnivClient) ResolveUniversityInUnitedStates() ([]University, error) {
	cty := Country{Name: "United States", Code: "edu"}

	return c.ResolveUniversityByCountry(cty)
}

// ResolveUniversityByCountry resolves all university in the country
func (c UnivClient) ResolveUniversityByCountry(cty Country) ([]University, error) {
	start := 1
	us := make([]University, 0)
	for {
		fmt.Printf("\t will resolve university at start:%d\n", start)
		tmp, nextStart, err := c.resolveUniversityByPage(cty, start)
		if err != nil {
			return nil, fmt.Errorf("resolve university(%s) by page(%d) err:%w", cty.Name, nextStart, err)
		}

		us = append(us, tmp...)
		if start == nextStart {
			break
		}
		start = nextStart
	}

	return us, nil
}

func (c UnivClient) resolveUniversityByPage(cty Country, start int) ([]University, int, error) {
	nextStart := start
	res, err := http.Get(fmt.Sprintf(_URL_COUNTRY_UNIVERSITY, cty.Code, start))
	if err != nil {
		return nil, nextStart, fmt.Errorf("ResolveAllCountries err:%w", err)
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, nextStart, fmt.Errorf("NewDocumentFromReader err:%w", err)
	}

	us := make([]University, 0)
	doc.Find("table.fixedWidth tr td ol li a").Each(func(idx int, s *goquery.Selection) {
		v, ok := s.Attr("href")
		if !ok {
			fmt.Printf("\t not found university domain on resolve university info,country:%s,start:%d,idx:%d\n",
				cty.Name, start, idx)
			return
		}

		us = append(us, University{Country: cty.Name, Name: s.Text(), Domain: v})
	})
	hasNext := false
	doc.Find("nav.resultNavigation a").Each(func(idx int, s *goquery.Selection) {
		if !hasNext {
			hasNext = strings.Index(strings.ToLower(s.Text()), "next") > -1
		}
	})

	if hasNext {
		nextStart = start + 50 // 50 per page
	}

	return us, nextStart, nil
}
