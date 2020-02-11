package main

import (
	"fmt"

	"time"

	"net/http"
	"net/url"

	"path/filepath"
	"strings"
	"strconv"

	"os"
	"io"

	"database/sql"
	"gopkg.in/headzoo/surf.v1"
	"github.com/headzoo/surf/browser"

	CDB "parser/connectionDB"
)


var (
	db 	*sql.DB
)

var	(
	typePhoto = []string{".png", ".jpg", ".jpeg"}

	replacer = strings.NewReplacer("://", "__", "/", "_", " ", "-")

	uploadDir = "./uploads"

	c chan string = make(chan string)
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
	checkError(err)
	if len(u.Hostname()) <= 0 {
		panic("Неверная ссылка (Пример: http://site.net)")
	}




	// БД
	db = CDB.Init()


	// Ссылка, откуда начнём искать
	CDB.StartPoint(link)


    // Эмулятор
	bow := surf.NewBrowser()
	bow.AddRequestHeader("Accept", "text/html")
	bow.AddRequestHeader("Accept-Charset", "utf8")



	// Лазание по страницам
	go linksF(bow, c)


	// Результат выполнения
	go func() {
	    for {
			fmt.Println(<- c)
		}
	}()




	var input string
    fmt.Scanln(&input)
}





/////////////// GO ///////////////

// Лазание по страницам
func linksF(bow *browser.Browser, c chan <- string)  {
	start := time.Now()
	for {
		r := linkWalking(bow)

		if r == false {
			c <- "Search links is finish"
			break
		} else {
			c <- "I find links"
		}
	}
	t := time.Now()
	elapsed := t.Sub(start).String()
	c <- "Search links time: " + elapsed

	// Выгружаем, если что-то осталось в буфере
	_, tags := CDB.Flush(true)
	//  Если все ссылки уже пройдены, можно проверить фотки
	if len(tags) <= 0 {
		c <- " -- Last images upload --"
		go photoF(c)
	} else {
		router(tags)
	}
}


// Скачка фоток (на конях ...)
func photoF(c chan <- string)  {
	start := time.Now()
	for {
		r := downloadPicturesLinks()

		if r == false {
			c <- "Search links pictures is finish"
			break
		} else {
			c <- "I find pictures links"
		}
	}
	t := time.Now()
	elapsed := t.Sub(start).String()
	c <- "Search links pictures time: " + elapsed

	// Выгружаем, если что-то осталось в буфере
	_, tags := CDB.Flush(true)
	router(tags)
}











// Лазаем
func linkWalking(bow *browser.Browser) bool {
	// Получаем ссылки для прохода
	arLinks := CDB.GetLinks()

	// Вдруг что-то в буфере валяется
	if len(arLinks) <= 0 {
		_,tags := CDB.Flush(true)
		router(tags)
		if in_array(tags, "#Links") {
			arLinks = CDB.GetLinks()
		}

		// Третьего шанса не даём, пошёл нах..
		if len(arLinks) <= 0 {
			return false
		}
	}


	for _,item := range arLinks {
		// Точка входа
		err := bow.Open(item.Path)
		checkErrorLog(err)
		if err != nil {
			continue
		}

		// Собираем ссылки
		getCollectionLinks(bow)

		// Собираем ссылки на фотки
		getCollectionPictures(bow)
	}

	return true
}




// Распределитель
func router(res []string)  {
	for _,tag := range res {
		switch tag {
			// case "#Links":
				// nothing =)
			case "#PhotoLinks":
				go photoF(c)
		}
	}
}





/////////////// LINKS ///////////////

// Собираем ссылки
func getCollectionLinks(boww *browser.Browser)  {
	links := boww.Links()
	// Сохраняем в базу ссылки
	go saveLinks(links, boww.Url().String())
}


//  Сохраняем ссылки в буфер
func saveLinks(arLinks []*browser.Link, url string)  {
	for _,item := range arLinks {
		CDB.PrepareLinks(CDB.Links{
			Path: item.URL.String(),
			Parent: url,
		})
	}

	CDB.PrepareCheckedLinks(CDB.Links{
		Path: url,
		Checked: false,
	})
}







/////////////// PICTURES ///////////////

// Собираем ссылки с фотками
func getCollectionPictures(boww *browser.Browser)  {
	images := boww.Images()
	// Сохраняем в базу ссылки
	go savePicturesLinks(images, boww.Url().String())
}


//  Сохраняем ссылки на картинки
func savePicturesLinks(arImages []*browser.Image, url string)  {
	needDownload := false

	for _,item := range arImages {
		// Фильтруем пока только по расширению
		if in_array( typePhoto, strings.ToLower( filepath.Ext(item.URL.String()) ) ) == false {
			continue
		}

		needDownload = CDB.PreparePicturesLinks(CDB.PhotoLinks{
			Path: item.URL.String(),
			Checked: false,
			Dir: strings.Join( []string{ uploadDir, "#LinksId#", replacer.Replace(item.URL.String()) } , "/"),
			LinksURL: url,
		})
	}

	// Если что-то добавили в БД, то качаем
	if needDownload == true {
		go photoF(c)
	}

}


// Скачивание фоток
func downloadPicturesLinks() bool {
	// Получаем фотки
	arPicturesLinks := CDB.GetPhotoLinks()

	// Вдруг что-то в буфере валяется
	if len(arPicturesLinks) <= 0 {
		_,tags := CDB.Flush(true)
		router(tags)
		if in_array(tags, "#PhotoLinks") {
			arPicturesLinks := CDB.GetPhotoLinks()
			// Третьего шанса не даём, пошёл нах..
			if len(arPicturesLinks) <= 0 {
				return false
			}
		}
		return false
	}


	for _,item := range arPicturesLinks {
		fileName := replacer.Replace(item.Path)
		filePuth := strings.Join( []string{ uploadDir, strconv.Itoa(item.LinksId) } , "/" )

		// Создаём директории, если их нет
		errr := os.MkdirAll(filePuth, os.ModePerm)
		checkErrorLog(errr)

		if errr == nil {
			if fileExists(filePuth + "/" + fileName) {
				// Помечаем, что данные файлы уже скачены
				if item.Checked == false {
					CDB.PrepareCheckedPicturesLinks(CDB.PhotoLinks{
						Id: item.Id,
						Checked: true,
					})
				}

				continue
			}

			file, err := os.Create( filePuth + "/" + fileName)
	    	checkErrorLog(err)

			if err == nil {
				go putFile(item, file)
			}
		}

	}

	return true
}


// Качаем к себе файл
func putFile(urlFile CDB.PhotoLinks, file *os.File) {
    resp,err := http.Get(urlFile.Path)
    checkErrorLog(err)
    defer resp.Body.Close()

    _,err2 := io.Copy(file, resp.Body)
    checkErrorLog(err2)
    defer file.Close()

	if err2 == nil {
		CDB.PrepareCheckedPicturesLinks(CDB.PhotoLinks{
			Id: urlFile.Id,
			Checked: true,
		})
	}
}







/////////////// FUNCTIONS ///////////////

func checkError(err error) {
    if err != nil {
        panic(err)
    }
}


func checkErrorLog(err error) {
    if err != nil {
        fmt.Println(err)
    }
}


func in_array(arr []string, str string) bool {
   for _, a := range arr {
      if a == str {
         return true
      }
   }
   return false
}


func fileExists(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}
