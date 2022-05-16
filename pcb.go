package main

import (
	"fmt"
	"os"
	"pcb/orm"
	"pcb/web_scraping"
	"strconv"
	"strings"
	"sync"

	"gorm.io/gorm"
)

func main() {
	wp, wc, c, canal, db, consumers, producers := create_variables()
	web_scraping.Actress(c)
	for i := 0; i < producers; i++ {
		go producer(c, canal, i, wp, db)
	}
	for r := 0; r < consumers; r++ {
		go consumer(canal, wc, r, db)
	}
	close(c)
	wp.Wait()
	close(canal)
	wc.Wait()
}

func create_variables() (*sync.WaitGroup, *sync.WaitGroup, chan string, chan string, *gorm.DB, int, int) {
	var consumers, producers int = 4, 4
	if len(os.Args) > 1 {
		producers, _ = strconv.Atoi(os.Args[1])
	}
	if len(os.Args) > 2 {
		consumers, _ = strconv.Atoi(os.Args[2])
	}
	var BufferSize int = 1000
	wp := &sync.WaitGroup{}
	wc := &sync.WaitGroup{}
	wc.Add(consumers)
	wp.Add(producers)
	c := make(chan string, 3000)
	canal := make(chan string, BufferSize)
	db := orm.Connection_CreateTables()
	return wp, wc, c, canal, db, consumers, producers
}

func producer(canal <-chan string, c_p chan<- string, n int, wp *sync.WaitGroup, db *gorm.DB) {
	defer wp.Done()
	for x := range canal {
		fmt.Println("Producer #" + strconv.Itoa(n) + " - " + "Actress: " + x)
		x += " - " + strconv.Itoa(n)
		select {
		case c_p <- x:
		default:
			fmt.Println("BUFFER IS FULL")
		}
	}
}

func consumer(link <-chan string, wg *sync.WaitGroup, r int, db *gorm.DB) {
	defer wg.Done()
	var films []string
	for links := range link {
		fmt.Println("Consumer #" + strconv.Itoa(r) + " " + links)
		films = append(films, links)
	}
	for _, links := range films {
		parts := strings.Split(links, " - ")
		p, err := strconv.Atoi(parts[2])
		_ = err
		if orm.Transaction_actress(db, parts[0], parts[1], p, r) {
			movie, links := web_scraping.Movies(parts[1], parts[0], true)
			unique(movie, links, parts[0], r, db)
		}
	}
}

func unique(s []string, l []string, actress string, c int, db *gorm.DB) {
	inResult := make(map[string]bool)
	for x, str := range s {
		str = strings.ToLower(str)
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			if strings.Contains(l[x], "/wiki/") || l[x] == "No link" {
				orm.Transaction_movie(db, c, strings.Title(str), l[x], actress)
			}
		}
	}
}
