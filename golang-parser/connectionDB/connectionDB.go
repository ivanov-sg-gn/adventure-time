package connectionDB


import (
	"io/ioutil"
	"database/sql"
	// "fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
)


var (
	db *sql.DB
	arLinks []Links
	arPhotoLinks []PhotoLinks
)

// Ссылки на страницы
type Links struct {
	Id   	int
	Path 	string 	// где проверять
	Flux 	string 	// uid потока
	Checked bool 	// проверена ссылки или нет
}

// Ссылк на фотографии
type PhotoLinks struct {
	Id   	int
	Path 	string 	// откуда качать
	Dir 	string 	// путь у нас в каталоге
	Flux 	string 	// uid потока
	Checked bool 	// загружена фотка или нет
	LinksId int 	// id ссылки, на которой данная картинка
	Date 	string
}

const (
	structurePath = "./structure.sql"
	queryBuff = 10000
	limitGetLink = 100
)



// Устанавливаем соединение
func Connect() (*sql.DB) {
	// Соединение с бд
	_db, err_db := sql.Open("mysql", "login:password@/golang?multiStatements=true")
	if err_db != nil {
		panic(err_db.Error())
	}
	db = _db

	return db;
}


// Проверяем таблицу
func CheckStruckture() bool {
	content, err := ioutil.ReadFile(structurePath)
	if err != nil {
		panic(err)
	}

	result, err2 := db.Exec( string(content) )
	if err2 != nil {
		panic(err2)
	}
	defer db.Close()


	rows, err3 := result.RowsAffected()
	if err3 != nil {
		panic(err3)
	}

	if rows < 0 {
		return false
	}
	return true

}



func StartPoint(link string)  {
	var sLink Links
	sLink.Path = link

	arLinks = append(arLinks, sLink)
	Flush(true)
}

///// Ссылки /////
func GetNewLinks() (Links) {
	var N Links
	return N;
}

// Запись
func PrepareLinks (link Links) {
	arLinks = append(arLinks, link)
	Flush(false)
}

// Получение ссылок
func GetLinks () ([]Links){
	// UUID 4
	nsec := uuid.Must(uuid.NewV4())

	// Резервирование записей
	updatePrepare,err := db.Prepare( `
		UPDATE Links
		SET flux = ?
		WHERE id IN (
		    SELECT id FROM(
		    	SELECT *
		        FROM Links
				WHERE flux = ""
		        ORDER BY id ASC
		        LIMIT 0, ?
		    ) tmp
		)
	` )
	if err != nil {
		panic(err)
	}
	defer updatePrepare.Close()

	update, err := updatePrepare.Exec(nsec, limitGetLink)
	if err != nil {
		panic(err)
	}

	// Обновили ли что-то
	rows, err := update.RowsAffected()
	if err != nil {
		panic(err)
	}
	if rows <= 0 {
		return arLinks
	}


	// Выборка данных
	result, err := db.Query("SELECT id, path FROM Links WHERE flux = ?", nsec)
	if err != nil {
		panic(err)
	}

	var hit Links
	var resLikns []Links

	for result.Next() {
		var id int
        var path string

		err := result.Scan(&id, &path)
		if err != nil {
            panic(err.Error())
        }

		hit.Id = id
		hit.Path = path

		resLikns = append(resLikns, hit)
	}
	defer result.Close()

	return resLikns
}






///// Фотки /////
func GetNewPhotoLinks () (PhotoLinks) {
	var N PhotoLinks
	return &N;
}

// Отложенный запрос на добавление изображений
func PreparePhotoLinks (photoLinks PhotoLinks) {
	arPhotoLinks = append(arPhotoLinks, photoLinks)
	Flush(false)
}

// Получение фоток
func GetPhotoLinks () {
	// Резервирование записей

	// Выборка данных
}







// Загрузить все отложенные запросы
func Flush(check bool) {

	if ( ( len(arLinks) + len(arPhotoLinks) ) >= queryBuff * 2 ) || check {

		// Links
		if len(arLinks) > 0 {
			insert, err := db.Prepare("INSERT IGNORE INTO Links(id, path, flux, checked, date) VALUES( NULL, ?, ?, ?, NOW() )")
			if err != nil {
				panic(err)
			}
			defer insert.Close()

			// count
			count := queryBuff
			if (len(arLinks) < queryBuff) {
				count = len(arLinks)
			}

			var slicearLinks []Links
			slicearLinks = arLinks[:count]
			arLinks = arLinks[count:]

			for _,itemLink := range slicearLinks {
				_, err2 := insert.Exec(itemLink.Path, itemLink.Flux, itemLink.Checked)
				if err2 != nil {
					panic(err2)
				}
			}

			slicearLinks = []Links
		}




		// PhotoLinks
		if len(arPhotoLinks) > 0 {
			insert2, err4 := db.Prepare("INSERT IGNORE INTO PhotoLinks(path, dir, flux, checked, linksId, date) VALUES( ?, ?, ?, ?, ?, NOW() )")
			if err4 != nil {
				panic(err4)
			}
			defer insert2.Close()

			// count
			count := queryBuff
			if (len(arPhotoLinks) < queryBuff) {
				count = len(arPhotoLinks)
			}

			var sliceArPhotoLinks []Links
			sliceArPhotoLinks = arPhotoLinks[:count]
			arPhotoLinks = arPhotoLinks[count:]


			for _,itemPhotoLink := range sliceArPhotoLinks {
				_, err5 := insert2.Exec(itemPhotoLink.Path, itemPhotoLink.Dir, itemPhotoLink.Flux, itemPhotoLink.Checked, itemPhotoLink.LinksId)
				if err5 != nil {
					panic(err5)
				}
			}

			slicearLinks = []PhotoLinks
		}

	}
}
