package main

import (
	"fmt"
	"net/url"
	"os"
	"database/sql"
	"gopkg.in/headzoo/surf.v1"
	"github.com/headzoo/surf/browser"

	CDB "parser/connectionDB"
)


var (
	db *sql.DB
)

func main() {
	if len(os.Args) < 2 {
		panic("Error argc")
	}

	if len(os.Args[1]) <= 0 {
		panic("Введите ссылку (Пример: http://site.net)")
	}

	link := os.Args[1];


	// Получаем домен
	u, err := url.Parse(link)
	if err != nil {
		panic(err)
	}
	if len(u.Hostname()) <= 0 {
		panic("Неверная ссылка (Пример: http://site.net)")
	}


	// БД
	db = CDB.Connect();
	// Проверка и создание таблиц БД
	// result := CDB.CheckStruckture()
	// if result == false {
	// 	panic("Проблема с таблицами")
	// }

	CDB.StartPoint(link)


    // Эмулятор
	bow := surf.NewBrowser()
	bow.AddRequestHeader("Accept", "text/html")
	bow.AddRequestHeader("Accept-Charset", "utf8")



	for {
		r := Loop(bow)

		if r == false {
			break
		}
	}

}


// Крутимся
func Loop(bow *browser.Browser) bool {
	fmt.Println("---Start Loop---")

	// Получаем ссылки для прохода
	arLinks := CDB.GetLinks()
	fmt.Println("Search links:", len(arLinks))

	if len(arLinks) <= 0 {
		fmt.Println("Flushing")
		CDB.Flush(true)

		arLinks = CDB.GetLinks()
		fmt.Println("Search links:", len(arLinks))

		if len(arLinks) <= 0 {
			return false
		}
	}



	for _,item := range arLinks {
		// Точка входа
		err := bow.Open(item.Path)
		if err != nil {
			panic(err)
		}

		// Собираем ссылки
		GetCollectionLinks(bow)

		// Собираем фотки
		GetCollectionPictures(bow)
	}

	return true
}


// Собираем ссылки
func GetCollectionLinks(boww *browser.Browser)  {
	links := boww.Links()
	fmt.Println("Find links:", len(links))
	// Сохраняем в базу ссылки
	SaveLinks(links)
}

//  Сохраняем ссылки
func SaveLinks(arLinks []*browser.Link)  {
	for _,item := range arLinks {
		res := CDB.GetNewLinks()
		res.Path = item.URL.String()
		res.Checked = false

		CDB.PrepareLinks(res)
	}
}






// Собираем ссылки с фотками
func GetCollectionPictures(boww *browser.Browser)  {
	images := boww.Images()
	fmt.Println("Find picture:", len(links))
	// Сохраняем в базу ссылки
	SavePicturesLinks(images)
}

//  Сохраняем ссылки на картинки
func SavePicturesLinks(arLinks []*browser.Image)  {
	for _,item := range arImages {
		res := CDB.GetNewPhotoLinks()
		res.Path = item.URL.String()
		res.Checked = false

		CDB.PreparePicturesLinks(res)
	}
}



// Скачивание фоток (асинхронно)
func DownloadPicturesLinks()  {
	// Получаем картинки

	// Скачиванием

}
