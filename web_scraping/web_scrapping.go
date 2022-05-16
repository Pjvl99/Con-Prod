package web_scraping

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func Search_link(link string, actress string, movie []string, links []string) ([]string, []string) {
	resp, err := http.Get("https://en.wikipedia.org" + link)
	_ = err
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	doc.Find("div.hatnote.navigation-not-searchable > a").Each(func(_ int, sel *goquery.Selection) {
		href, ok := sel.Attr("href")
		_ = ok
		act := strings.Split(actress, " ")[0]
		if strings.Contains(href, act) && Contains(href) {
			x, y := Movies(href, actress, false)
			movie = append(movie, x...)
			links = append(links, y...)
		}
	})
	return movie, links
}

func Contains(href string) bool {
	return strings.Contains(href, "film") || strings.Contains(href, "perf") ||
		strings.Contains(href, "screen")
}

func Actress(canal chan string) {
	c := colly.NewCollector() //Iniciamos una variable para hacer web scrapping
	c.OnHTML(".div-col", func(e *colly.HTMLElement) { //Iteramos en cada tag por letra
		e.ForEach("li > a", func(_ int, k *colly.HTMLElement) { //Iteramos en cada tag donde se encuentra la informacion de actriz
			name := strings.Split(k.Attr("title"), " (")[0] //Retiramos los parentesis de mas dentro del nombre de la actriz y tomamos solo sus nombres
			link := k.Attr("href")                          //Tomamos los links respectivos de cada actriz
			value := name + " - " + link
			canal <- value
		})
	})
	c.Visit("https://en.wikipedia.org/wiki/List_of_American_film_actresses") //El link donde buscamos a las actrices
}

func Movies(link string, actress string, visit bool) ([]string, []string) {
	c := colly.NewCollector()
	var movie, links []string
	istable := false
	islist := false
	c.OnHTML(".mw-parser-output", func(e *colly.HTMLElement) {
		e.ForEach(".wikitable", func(_ int, k *colly.HTMLElement) {
			k.ForEach("td", func(i int, h *colly.HTMLElement) {
				if len(h.ChildText("i")) > 0 {
					x := "No link"
					if len(h.ChildAttr("a", "href")) > 0 {
						x = h.ChildAttr("a", "href")
					}
					istable = true
					movie = append(movie, h.ChildText("i"))
					links = append(links, x)
				}
			})
		})
		e.ForEach("ul", func(_ int, k *colly.HTMLElement) {
			k.ForEach("li", func(i int, h *colly.HTMLElement) {
				if len(h.ChildText("i")) > 0 && !istable {
					x := "No link"
					if len(h.ChildAttr("a", "href")) > 0 {
						x = h.ChildAttr("a", "href")
					}
					islist = true
					movie = append(movie, h.ChildText("i"))
					links = append(links, x)
				}
			})
		})
		e.ForEach(".toccolours", func(_ int, k *colly.HTMLElement) {
			k.ForEach("td", func(i int, h *colly.HTMLElement) {
				if len(h.ChildText("i")) > 0 && !istable && !islist {
					x := "No link"
					if len(h.ChildAttr("a", "href")) > 0 {
						x = h.ChildAttr("a", "href")
					}
					movie = append(movie, h.ChildText("i"))
					links = append(links, x)
				}
			})
		})
	})
	c.Visit("https://en.wikipedia.org" + link)
	if visit {
		movie, links = Search_link(link, actress, movie, links)
	}
	return movie, links
}
