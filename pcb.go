package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

var producers int = 5
var consumers int = 0
var BufferSize int = 9000

func main() {
	wp := &sync.WaitGroup{}
	wc := &sync.WaitGroup{}
	wc.Add(consumers)
	wp.Add(producers)
	c := make(chan string, BufferSize)
	canal := make(chan string, BufferSize)
	a := colly.NewCollector() //Iniciamos una variable para hacer web scrapping
	actress(c, a)
	for i := 0; i < producers; i++ {
		go producer(c, canal, i, wp)
	}
	for r := 0; r < consumers; r++ {
		go consumer(canal, wc, r)
	}
	close(c)
	wp.Wait()
	close(canal)
	wc.Wait()
}
func actress(canal chan string, c *colly.Collector) {
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

func producer(canal <-chan string, c_p chan<- string, n int, wp *sync.WaitGroup) {
	defer wp.Done()
	for x := range canal {
		fmt.Println("Producer: " + strconv.Itoa(n) + " - " + "Actress: " + x)
		c_p <- x
	}
	return
}

func consumer(link <-chan string, wg *sync.WaitGroup, i int) {
	defer wg.Done()
	for links := range link {
		parts := strings.Split(links, " - ")
		fmt.Println("Actress: " + parts[0] + " - Consumer: " + strconv.Itoa(i))
		movies(parts[1], parts[0])
	}
}

func movies(link string, actress string) {
	istable := false
	c := colly.NewCollector()
	c.OnHTML(".mw-parser-output", func(e *colly.HTMLElement) {
		e.ForEach(".wikitable", func(_ int, k *colly.HTMLElement) {
			istable = true
			fmt.Print("Movies: ")
			fmt.Println(k.ChildTexts("i"))
		})
		if !istable {
			e.ForEach("ul", func(_ int, l *colly.HTMLElement) {
				if len(l.ChildText("i")) != 0 {
					fmt.Print("Movies: ")
					fmt.Println(l.ChildTexts("i"))
				}
			})
		}
		fmt.Println("")
	})
	c.Visit("https://en.wikipedia.org" + link)
}