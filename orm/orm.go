package orm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Actress struct {
	Id_actress   int    `gorm:"primaryKey;not null;autoIncrement"`
	Name         string `gorm:"size:191;uniqueIndex"`
	Url          string
	Idconsumidor int
	Idproductor  int
}

type Films struct {
	Id_film      int    `gorm:"primaryKey;not null;autoIncrement"`
	Name         string `gorm:"size:191;uniqueIndex"`
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

func Transaction_actress(db *gorm.DB, act string, link string, p int, c int) bool {
	var search_movies bool = true
	db.Transaction(func(tx *gorm.DB) error {
		To_insert := Actress{Name: act, Url: link, Idproductor: p, Idconsumidor: c}
		if err := tx.Table("actresses").Select("name", "url", "idproductor", "idconsumidor").Create(&To_insert).Error; err != nil {
			search_movies = false
			return err
		}
		return nil
	})
	return search_movies
}

func Transaction_movie(db *gorm.DB, c int, movie string, url string, actress string) {
	db.Transaction(func(tx *gorm.DB) error {
		To_insert := Films{Name: movie, Url: url, Idconsumidor: c}
		if err := tx.Table("films").Select("name", "url", "idconsumidor").Create(&To_insert).Error; err != nil {
			return err
		}
		return nil
	})
	db.Transaction(func(tx *gorm.DB) error {
		var id_film, id_actress int
		tx.Table("films").Select("id_film").Where("name = ?", movie).Scan(&id_film)
		tx.Table("actresses").Select("id_actress").Where("name = ?", actress).Scan(&id_actress)
		To_insert := Actress_Filmografy{Id_actress: id_actress, Id_film: id_film, Idconsumidor: c}
		if err := tx.Table("actress_filmografies").Create(&To_insert).Error; err != nil {
			return err
		}
		return nil
	})
}

func Connection_CreateTables() *gorm.DB {
	dsn := "root:helloworld@tcp(127.0.0.1:3308)/testapp?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	} else {
		db.Migrator().DropTable(&Actress_Filmografy{})
		db.Migrator().DropTable(&Actress{})
		db.Migrator().DropTable(&Films{})
		db.Migrator().CreateTable(&Actress{})
		db.Migrator().CreateTable(&Films{})
		db.Migrator().CreateTable(&Actress_Filmografy{})
	}
	return db
}
