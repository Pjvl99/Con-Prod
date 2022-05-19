package web_scraping

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

func Search_link(link string, actress string, movie []string, links []string) ([]string, []string) { //Revision de link
	resp, err := http.Get("https://en.wikipedia.org" + link)
	_ = err
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	doc.Find("div.hatnote.navigation-not-searchable > a").Each(func(_ int, sel *goquery.Selection) { //Este es el tag donde esta el link si existe
		href, ok := sel.Attr("href") //Lo va a buscar
		_ = ok
		act := strings.Split(actress, " ")[0]              //Lo separo
		if strings.Contains(href, act) && Contains(href) { //Ver si el link tiene un nombre de la actriz y contiene al menos una palabra clave
			x, y := Movies(href, actress, false) // x, y estan vacios
			movie = append(movie, x...)          // Los 3 puntos significan todo el array
			links = append(links, y...)
		}
	})
	return movie, links
}

func Contains(href string) bool {
	return strings.Contains(href, "film") || strings.Contains(href, "perf") || //Palabras clave
		strings.Contains(href, "screen") //Si contiene alguna de estas 3 es probable que existan peliculas en ese link
}

func Actress(canal chan string) {
	c := colly.NewCollector()                         //Iniciamos una variable para hacer web scrapping
	c.OnHTML(".div-col", func(e *colly.HTMLElement) { //Iteramos en cada tag por letra
		e.ForEach("li > a", func(_ int, k *colly.HTMLElement) { //Iteramos en cada tag donde se encuentra la informacion de actriz
			name := strings.Split(k.Attr("title"), " (")[0] //Retiramos los parentesis de mas dentro del nombre de la actriz y tomamos solo sus nombres
			link := k.Attr("href")                          //Tomamos los links respectivos de cada actriz
			value := name + " - " + link
			canal <- value // Lo subo al buffer (EL BUFFER PARA PRECARGAR)
		})
	})
	c.Visit("https://en.wikipedia.org/wiki/List_of_American_film_actresses") //El link donde buscamos a las actrice
}

func Movies(link string, actress string, visit bool) ([]string, []string) { // El Web Scrapping de las peliculas
	c := colly.NewCollector() // Iniciamos la variable para iniciar el web scrapping
	var movie, links []string
	istable := false
	islist := false
	c.OnHTML(".mw-parser-output", func(e *colly.HTMLElement) { // Es donde esta toda la informacion
		e.ForEach(".wikitable", func(_ int, k *colly.HTMLElement) { // Aqui es para buscar en caso sea tabla
			k.ForEach("td", func(i int, h *colly.HTMLElement) { // Aqui va por elemento dentro de la tabla
				if len(h.ChildText("i")) > 0 { // Ya va dentro del hijo
					x := "No link" //No hay link
					if len(h.ChildAttr("a", "href")) > 0 {
						x = h.ChildAttr("a", "href") // El link
					}
					istable = true
					movie = append(movie, h.ChildText("i")) // Nombre de la pelicula
					links = append(links, x)                //El link
				}
			})
		})
		e.ForEach("ul", func(_ int, k *colly.HTMLElement) { // Este es en caso este en lista
			k.ForEach("li", func(i int, h *colly.HTMLElement) { //Va por elemento de lista
				if len(h.ChildText("i")) > 0 && !istable { //Va por hijo
					x := "No link"
					if len(h.ChildAttr("a", "href")) > 0 { //Jala el link
						x = h.ChildAttr("a", "href")
					}
					islist = true
					movie = append(movie, h.ChildText("i")) //El nombre de la pelicula
					links = append(links, x)
				}
			})
		})
		e.ForEach(".toccolours", func(_ int, k *colly.HTMLElement) { //Tag raro
			k.ForEach("td", func(i int, h *colly.HTMLElement) { //Jala por lista
				if len(h.ChildText("i")) > 0 && !istable && !islist { //Va por hijo
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
	c.Visit("https://en.wikipedia.org" + link) //Link de wikipedia
	if visit {
		movie, links = Search_link(link, actress, movie, links) //Va a buscar si la actriz tiene filmografia en link externo
	}
	return movie, links
}
