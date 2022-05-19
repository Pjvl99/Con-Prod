package orm

import (
	"errors"
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Actress struct {
	Id_actress   int `gorm:"primaryKey;not null;autoIncrement"`
	Name         string
	Url          string
	Idconsumidor int
	Idproductor  int
}

type Films struct {
	Id_film      int `gorm:"primaryKey;not null;autoIncrement"`
	Name         string
	Url          string
	Idconsumidor int
}

type Actress_Filmografy struct {
	Id_actress   int
	Id_film      int
	Idconsumidor int
	Films        Films   `gorm:"foreignKey:id_film"`
	Actress      Actress `gorm:"foreignKey:id_actress"`
}

func Transaction_actress(db *gorm.DB, act string, link string, p int, c int, a *sync.Mutex) bool { //La primera transaccion es para las actrices
	fmt.Println("Inicia transaccion de actriz: " + act)
	var search_movies bool = true
	var actress_result string
	a.Lock()                                 //Inicia transaccion, y bloquea
	db.Transaction(func(tx *gorm.DB) error { //Inicia la transaccion
		To_insert := Actress{Name: act, Url: link, Idproductor: p, Idconsumidor: c}
		tx.Table("actresses").Select("name").Where("name = ?", act).Scan(&actress_result)
		if len(actress_result) != 0 {
			return errors.New("ACTRIZ YA EXISTE")
		} else {
			if err := tx.Table("actresses").Select("name", "url", "idproductor", "idconsumidor").Create(&To_insert).Error; err != nil { //Es insertar
				search_movies = false //Si hay error ya no va a buscar peliculas de la actriz
				return err            // Hace rollback
			}
		}
		return nil //hace commit
	})
	a.Unlock()           //Termina la transaccion y desbloquea
	return search_movies //Retorna si hay que buscar o no las peliculas
}

func Transaction_movie(db *gorm.DB, c int, movie string, url string, actress string, m *sync.Mutex) { //Segunda transaccion tema de peliculas
	var movie_exist string
	m.Lock()                                 //Inicia transaccion y bloquea
	db.Transaction(func(tx *gorm.DB) error { //Inicia transaccion
		To_insert := Films{Name: movie, Url: url, Idconsumidor: c}
		tx.Table("films").Select("name").Where("name = ?", movie).Scan(&movie_exist)
		if len(movie_exist) != 0 {
			return errors.New("LA PELICULA YA EXISTE")
		} else {
			if err := tx.Table("films").Select("name", "url", "idconsumidor").Create(&To_insert).Error; err != nil { //Inserta en tabla
				return err //Si hay error hace rollback
			}
		}
		return nil //commit
	})
	db.Transaction(func(tx *gorm.DB) error { // La 3ra tabla donde estan actrices y peliculas juntos
		var id_film, id_actress int
		tx.Table("films").Select("id_film").Where("name = ?", movie).Scan(&id_film)                //Toma el id de la actriz
		tx.Table("actresses").Select("id_actress").Where("name = ?", actress).Scan(&id_actress)    //Id de la pelicula
		To_insert := Actress_Filmografy{Id_actress: id_actress, Id_film: id_film, Idconsumidor: c} //Crea el puntero
		if err := tx.Table("actress_filmografies").Create(&To_insert).Error; err != nil {          //Lo inserta
			return err //Rollback
		}
		return nil //Commit
	})
	m.Unlock() //Termina transaccion y desbloquea
}

func Connection_CreateTables() *gorm.DB {
	dsn := "root:helloworld@tcp(127.0.0.1:3308)/testapp?charset=utf8mb4&parseTime=True&loc=Local" //Connectarse a la base de datos
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})                                         //Defino que voy a usar mysql
	if err != nil {
		panic(err)
	} else {
		db.Migrator().DropTable(&Actress_Filmografy{}) //Drop table if exists
		db.Migrator().DropTable(&Actress{})
		db.Migrator().DropTable(&Films{})
		db.Migrator().CreateTable(&Actress{})
		db.Migrator().CreateTable(&Films{})
		db.Migrator().CreateTable(&Actress_Filmografy{}) //Create table
	}
	return db
}
