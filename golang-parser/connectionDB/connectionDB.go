package connectionDB


import (
	"io/ioutil"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
)



// Ссылки на страницы
type Links struct {
	Id   	int
	Path 	string 	// где проверять
	Flux 	string 	// uid потока
	Parent 	string 	// проверена ссылки или нет
	Checked bool // Нужно ли проверять
}

// Ссылк на фотографии
type PhotoLinks struct {
	Id   	int
	Path 	string 	// откуда качать
	Dir 	string 	// путь у нас в каталоге
	Flux 	string 	// uid потока
	Checked bool 	// загружена фотка или нет

	LinksId int 	// id ссылки, откуда скачал (используется только для получения инфы из БД)
	LinksURL string // ссылка для поиска id в БД

	Date string	 	// Длф формирования директорий  (используется только для получения инфы из БД)
}


var (
	db *sql.DB
	arLinks []Links
	arPhotoLinks []PhotoLinks

	arLinksUpdate []Links
	arPicturesLinksUpdate []PhotoLinks
)


const (
	dbLogin = "golang"
	dbPassword = "%54321GoLang12345%"
	dbName = "golang"

	structurePath = "./structure.sql"

	queryBuff = 10000 // Размер буфера для выгрузки в БД

	limitGetLink = 100 // Сколько достаём ссылок из БД

	limitGetPhicturesLink = 100 // Сколько достаём ссылок на фото из БД
)


func Init() (*sql.DB) {
	// Соединение с бд
	_db, err_db := sql.Open("mysql", dbLogin + ":" + dbPassword + "@/" + dbName + "?multiStatements=true")
	if err_db != nil {
		panic(err_db.Error())
	}
	db = _db

	return db;
}


// Проверка структуры БД
func CheckStruckture() bool {
	content, err := ioutil.ReadFile(structurePath)
	if err != nil {
		checkError(err)
		panic(err)
	}

	result, err := db.Exec( string(content) )
	if err != nil {
		checkError(err)
		panic(err)
	}
	defer db.Close()


	rows, err := result.RowsAffected()
	if err != nil {
		checkError(err)
		panic(err)
	}

	if rows < 0 {
		return false
	}
	return true

}



// Первая ссылка
func StartPoint(link string)  {
	var sLink Links
	sLink.Path = link
	sLink.Parent = "#"

	arLinks = append(arLinks, sLink)
	Flush(true)
}






/////////////// LINKS ///////////////

// Запись
func PrepareLinks (link Links) bool {
	arLinks = append(arLinks, link)
	_,tag := Flush(false)
	if  in_array(tag, "#Links") {
		return true
	}
	return false
}


// Получение ссылок
func GetLinks () ([]Links){
	// Скорее всег в буфере что-то есть
	checkedLinks(true)

	// UUID 4
	nsec := uuid.Must(uuid.NewV4())

	// Резервирование записей
	updatePrepare,err := db.Prepare( `
		UPDATE Links
		SET flux = ?
		WHERE id IN (
		    SELECT id FROM(
		    	SELECT id
		        FROM Links
				WHERE flux = "" OR checked = true
		        ORDER BY id ASC
		        LIMIT 0, ?
		    ) tmp
		)
	` )
	checkError(err)
	defer updatePrepare.Close()

	update, err := updatePrepare.Exec(nsec, limitGetLink)
	checkError(err)


	// Обновили ли что-то
	rows, err := update.RowsAffected()
	checkError(err)

	if rows <= 0 {
		return []Links{}
	}


	// Выборка данных
	result, err := db.Query("SELECT id, path FROM Links WHERE flux = ?", nsec)
	checkError(err)


	var hit Links
	var resLikns []Links

	for result.Next() {
		var id int
        var path string

		err := result.Scan(&id, &path)
		if err != nil {
            fmt.Println(err.Error())
        }

		hit.Id = id
		hit.Path = path

		resLikns = append(resLikns, hit)
	}
	defer result.Close()

	return resLikns
}


// Для перепроверки ссылки
func PrepareCheckedLinks(link Links) {
	arLinksUpdate = append(arLinksUpdate, link)
	Flush(false)
}


func checkedLinks(check bool) {
	if(len(arLinksUpdate) > queryBuff || check) {

		if len(arLinksUpdate) > 0 {
			update,err := db.Prepare(`
				UPDATE Links
				SET checked = ?
				WHERE path = ?
			`)
			checkError(err)
			defer update.Close()

			count := queryBuff
			if (len(arLinksUpdate) < queryBuff) {
				count = len(arLinksUpdate)
			}

			var sliceArLinks []Links
			sliceArLinks = arLinksUpdate[:count]
			arLinksUpdate = arLinksUpdate[count:]

			for _,itemLink := range sliceArLinks {
				_, err := update.Exec(itemLink.Checked, itemLink.Path)
				checkError(err)
			}

			sliceArLinks = []Links{}
		}
	}
}








/////////////// PICTURES ///////////////

// Добавление в буфер запрос на добавление изображений
func PreparePicturesLinks (photoLinks PhotoLinks) bool {
	arPhotoLinks = append(arPhotoLinks, photoLinks)
	_,tag := Flush(false)
	if in_array(tag, "#PhotoLinks") {
		return true
	}
	return false
}


// Получение фоток
func GetPhotoLinks () []PhotoLinks {
	// Скорее всег в буфере что-то есть
	checkedPicturesLinks(true)

	// UUID 4
	nsec := uuid.Must(uuid.NewV4())

	// Резервирование записей
	updatePrepare,err := db.Prepare( `
		UPDATE PhotoLinks
		SET flux = ?
		WHERE id IN (
		    SELECT id FROM(
		    	SELECT id
		        FROM PhotoLinks
				WHERE flux = "" OR checked = false
		        ORDER BY id ASC
		        LIMIT 0, ?
		    ) tmp
		)
	` )
	checkError(err)
	defer updatePrepare.Close()

	update, err := updatePrepare.Exec(nsec, limitGetPhicturesLink)
	checkError(err)

	// Обновили ли что-то
	rows, err := update.RowsAffected()
	checkError(err)

	if rows <= 0 {
		return []PhotoLinks{}
	}

	// Выборка данных
	result, err := db.Query("SELECT id, path, linksId, date FROM PhotoLinks WHERE flux = ? OR checked = false", nsec)
	checkError(err)

	var hit PhotoLinks
	var resLikns []PhotoLinks

	for result.Next() {
		var (
			id int
	        path string
			linksId int
			date string
		)

		err := result.Scan(&id, &path, &linksId, &date)
		if err != nil {
			fmt.Println("Panic", err)
			panic(err.Error())
		}

		hit.Id = id
		hit.Path = path
		hit.LinksId = linksId
		hit.Date = date

		resLikns = append(resLikns, hit)
	}
	defer result.Close()

	return resLikns
}







/////////////// BUFFER ///////////////

// Загрузить все из буфера
func Flush(check bool) (bool, []string) {
	result := []string{}

	if ( ( len(arLinks) + len(arPhotoLinks) ) >= queryBuff * 2 ) || check {

		// ССЫЛКИ | Links
		if len(arLinks) > 0 {

			insert, err := db.Prepare("INSERT IGNORE INTO Links(id, path, flux, parent, date) VALUES( NULL, ?, ?, ?, NOW() )")
			checkError(err)

			defer insert.Close()

			// max count
			count := queryBuff
			if (len(arLinks) < queryBuff) {
				count = len(arLinks)
			}

			var slicearLinks []Links
			slicearLinks = arLinks[:count]
			arLinks = arLinks[count:]

			for _,itemLink := range slicearLinks {
				_, err2 := insert.Exec(itemLink.Path, itemLink.Flux, itemLink.Parent)
				checkError(err2)

			}

			slicearLinks = []Links{}

			result = append(result, "#Links")
		}




		// ФОТКИ | PhotoLinks
		if len(arPhotoLinks) > 0 {
			insert2, err4 := db.Prepare(
				`INSERT IGNORE INTO PhotoLinks (path, linksId, dir, flux, checked, date)
				VALUES(
					?,
					(@LinksId :=(SELECT id FROM Links WHERE path = ?)),
					(REPLACE(?,'#LinksId#',@LinksId)),
					?,
					?,
					NOW()
				)`)
			checkError(err4)
			defer insert2.Close()

			//count
			count := queryBuff
			if (len(arPhotoLinks) < queryBuff) {
				count = len(arPhotoLinks)
			}

			var sliceArPhotoLinks []PhotoLinks
			sliceArPhotoLinks = arPhotoLinks[:count]
			arPhotoLinks = arPhotoLinks[count:]

			for _,itemPhotoLink := range sliceArPhotoLinks {
				_, err5 := insert2.Exec(itemPhotoLink.Path, itemPhotoLink.LinksURL, itemPhotoLink.Dir, itemPhotoLink.Flux, itemPhotoLink.Checked)
				checkError(err5)
			}

			sliceArPhotoLinks = []PhotoLinks{}

			result = append(result, "#PhotoLinks")
		}
	}


	go checkedLinks(check)
	go checkedPicturesLinks(check)

	return false, result
}








/////////////// PICTURES CHECKED ///////////////

// Добавление в буфер для пометки фоток в БД как "checked"
func PrepareCheckedPicturesLinks (photoLinks PhotoLinks) {
	arPicturesLinksUpdate = append(arPicturesLinksUpdate, photoLinks)
	Flush(false)
}


// Отправка буфера запросов в БД для отметки "checked" у фоток
func checkedPicturesLinks(check bool) {
	if(len(arPicturesLinksUpdate) > queryBuff || check) {

		if len(arPicturesLinksUpdate) > 0 {
			update,err := db.Prepare( `
				UPDATE PhotoLinks
				SET checked = ?
				WHERE id = ?
			` )
			checkError(err)
			defer update.Close()

			count := queryBuff
			if (len(arPicturesLinksUpdate) < queryBuff) {
				count = len(arPicturesLinksUpdate)
			}

			var sliceArPhotoLinks []PhotoLinks
			sliceArPhotoLinks = arPicturesLinksUpdate[:count]
			arPicturesLinksUpdate = arPicturesLinksUpdate[count:]

			for _,itemPhotoLink := range sliceArPhotoLinks {
				_, err := update.Exec(itemPhotoLink.Checked, itemPhotoLink.Id)
				checkError(err)
			}

			sliceArPhotoLinks = []PhotoLinks{}
		}
	}
}








/////////////// FUNCTIONS ///////////////

func checkError(err error) {
    if err != nil {
		fmt.Println("Panic", err)
        panic(err)
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
