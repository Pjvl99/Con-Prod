package main

import (
	"fmt"
	"os"
	"pcb/orm"
	"pcb/web_scraping"
	"strconv"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

var BufferSize int = 4  //Tamaño del buffer

func main() {
	wp, wc, c, canal, db, consumers, producers, m, a := create_variables() //Aqui cree las variables
	web_scraping.Actress(c)                                                //AQUI COMIENZA EL WEB SCRAPPING SOLO DE ACTRICES
	for i := 0; i < producers; i++ {
		go producer(c, canal, i, wp, db) // Corren los productores, go = Paralelizacion
	}
	for r := 0; r < consumers; r++ {
		go consumer(canal, wc, r, db, &m, &a) // Corren los consumidores, go = Paralelizacion
	}
	close(c)     // SIgnifica que ya no va a seguir recibiendo informacion
	wp.Wait()    // Espera a que termine el productor
	close(canal) //Cerrar el buffer principal que comparten los productores y consumidores
	wc.Wait()    //Espera a que termine el consumidor

}

func create_variables() (*sync.WaitGroup, *sync.WaitGroup, chan string, chan string, *gorm.DB, int, int, sync.Mutex, sync.Mutex) {
	var consumers, producers int = 4, 4 // Valores por defualt del consumidor y el productor
	if len(os.Args) > 1 {               // Si el usuario envia solo 1 valor; Ej: pcb 4
		producers, _ = strconv.Atoi(os.Args[1]) //Primero recibe el productor
	}
	if len(os.Args) > 2 { // Si el usuario envia 2 valores
		consumers, _ = strconv.Atoi(os.Args[2]) // Ej; pcb 3 5 -> Consumidor serian 5
	}
	wp := &sync.WaitGroup{} //Un sincronizador para el productor -> wp
	wc := &sync.WaitGroup{} //Un sincronizador para el consumidor -> wc
	var m sync.Mutex
	var a sync.Mutex
	wc.Add(consumers)                                       //.Add -> Lo que hace es que es como un contador, La variable wc = # Consumidores
	wp.Add(producers)                                       //.Add -> Lo mismo con los productores -> wp = # Productores, 4 5 6
	c := make(chan string, 3000)                            //Esto solo es un buffer para precargar la informacion ESTE BUFFER NO ES EL QUE COMPARTE EL CONSUMIDOR Y PRODUCTOR
	canal := make(chan string, BufferSize)                  // ESTE BUFFER ES EL QUE COMPARTIRAN PRODUCTORES Y CONSUMIDORES
	db := orm.Connection_CreateTables()                     // AQUI CREO LA CONEXION A LA BASE DE DATOS Y LAS TABLAS
	return wp, wc, c, canal, db, consumers, producers, m, a // RETORNO LAS VARIABLES
}

func producer(canal chan string, c_p chan<- string, n int, wp *sync.WaitGroup, db *gorm.DB) {
	defer wp.Done()
	for x := range canal {
		i := 0
		x += " - " + strconv.Itoa(n)
		for { //Se itera 
			if len(c_p) != cap(c_p) { // Si el buffer aun tiene espacio se ingresa para insertar la info
				c_p <- x //Insertando en buffer
				fmt.Println("Actriz:" + x + " productor #" + strconv.Itoa(n))
				break
			} else if i == 200 {//Espera 20 segundos
				fmt.Println("DEJANDO PROCESO PUES EL TIEMPO DE ESPERA SE EXCEDIO, ATT: PRODUCTOR #" + strconv.Itoa(n)) // Luego abandona
				return
			} else {
				time.Sleep(100 * time.Millisecond) //Se duerme 100 milisegundos si el buffer esta lleno
				i += 1
			}
		}
	}
}

func consumer(link <-chan string, wg *sync.WaitGroup, r int, db *gorm.DB, m *sync.Mutex, a *sync.Mutex) {
	defer wg.Done()
	var films []string
	for links := range link { //Empieza a iterar sobre el buffer 
		fmt.Println("Consumer #" + strconv.Itoa(r) + " " + links)
		films = append(films, links) //Va guardando la informacion de las actrices en un array
	}
	for _, links := range films { //Se itera sobre el array de las informacion de la actriz
		parts := strings.Split(links, " - ") //Se separa para tener link y nombre por separaddo
		p, err := strconv.Atoi(parts[2]) //El [2] Es el # de Productor, se convierte a entero para insertarlo posteriormente en la base de datos
		_ = err
		if orm.Transaction_actress(db, parts[0], parts[1], p, r, a) { // Se inicia transaccion para insertar actriz
			movie, links := web_scraping.Movies(parts[1], parts[0], true) //Si no existe se empieza a buscar sus peliculas en wiki
			unique(movie, links, parts[0], r, db, m) //Remueve cualquier duplicado en el array
		}
	}
}

func unique(s []string, l []string, actress string, c int, db *gorm.DB, m *sync.Mutex) {
	inResult := make(map[string]bool) //Diccionario para evitar duplicados dentro del array
	for x, str := range s {
		str = strings.ToLower(str) //Se estandariza volviendo todas minusculas asi evitamos repeticiones
		if _, ok := inResult[str]; !ok { //Si no existe ingresa
			inResult[str] = true
			if strings.Contains(l[x], "/wiki/") || l[x] == "No link" { //Se limpian los links raros
				orm.Transaction_movie(db, c, strings.Title(str), l[x], actress, m) //Inicia transaccion por peliculas
			}
		}
	}
}