package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hoka-isdl/ISDL-BMS/backend/app/structure"
	// "github.com/hoka-isdl/ISDL-BMS/backend/app/structure"
)

var db *sql.DB

func Opendb() {

	db_name, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/bms_db?")

	if err != nil {
		panic(err.Error())
	}

	db = db_name
}
func CreateTask() {
	fmt.Println(1)
	Opendb()
	defer db.Close()

	insert, err := db.Prepare("INSERT INTO Students(id,name,password,email) VALUES(?, ?, ?, ?)")

	if err != nil {
		log.Fatal(err.Error())
	}
	// insert.Exec(2,"isbn","岡")
	insert.Exec("1116190059", "酒部健太郎", "sakabe", "sakabe.kentaro@mikilab.doshisha.ac.jp")
}

func CreateTask_tag() {
	fmt.Println(1)
	Opendb()
	defer db.Close()

	insert, err := db.Prepare("INSERT INTO Tags(id,tagname) VALUES(?, ?)")

	if err != nil {
		log.Fatal(err.Error())
	}
	// insert.Exec(2,"isbn","岡")
	insert.Exec("1", "機械学習")
	insert.Exec("2", "GAN")
	insert.Exec("3", "LSTM")
}

func CreateTask2() {
	fmt.Println(1)
	Opendb()
	defer db.Close()

	insert, err := db.Prepare("INSERT INTO Books(id,title,title_kana,tagid,ISBN,author,author_kana,publisher,item_caption,image_url) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		log.Fatal(err.Error())
	}
	// insert.Exec(2,"isbn","岡")
	insert.Exec(2, "実践力を身につけるPythonの教科書", "ジッセンリョクヲミニツケルパイソンノキョウカショ", "#1#", 9784839960247, "クジラ飛行机", "クジラヒコウヅクエ", "マイナビ出版", "基本文法から始めてアプリ開発までしっかり解説", "https://thumbnail.image.rakuten.co.jp/@0_mall/book/cabinet/0247/9784839960247.jpg?_ex=200x200")
	insert.Exec(4, "実践力を身につけるGANの教科書", "ジッセンリョクヲミニツケルギャンノキョウカショ", "#1##2#", 9784839960258, "クジラ飛行机", "クジラヒコウヅクエ", "マイナビ出版", "基本文法から始めてアプリ開発までしっかり解説", "https://thumbnail.image.rakuten.co.jp/@0_mall/book/cabinet/0247/9784839960247.jpg?_ex=200x200")
	insert.Exec(5, "実践力を身につけるLSTMの教科書", "ジッセンリョクヲミニツケルギャンノキョウカショ", "#1#", 9784839960259, "クジラ飛行机", "クジラヒコウヅクエ", "マイナビ出版", "基本文法から始めてアプリ開発までしっかり解説", "https://thumbnail.image.rakuten.co.jp/@0_mall/book/cabinet/0247/9784839960247.jpg?_ex=200x200")
}

func Researchbook(isbn string) string {

	// Opendb()
	// defer db.Close()

	// _ = db.QueryRow("SELECT name FROM User WHERE name = ? AND password = ?", name, password).Scan(&user.Name)

	// if user.Name == "" {
	// 	return false
	// }

	return isbn
}

func GetBookinfo(isbn string) (string, string, string) {
	var book structure.Books

	Opendb()
	defer db.Close()

	rows_title, err := db.Query("SELECT title, author, publisher FROM Books WHERE isbn = ?", isbn)

	if err != nil {
		log.Fatal(err.Error())
	}

	for rows_title.Next() {
		rows_title.Scan(&book.Book_title, &book.Book_author, &book.Book_publisher)
	}

	return book.Book_title, book.Book_author, book.Book_publisher
}

func Partitionid(num string) string {
	return "#" + num + "#"
}

func GetTagid(tagname string) string {
	var tag structure.Tags

	Opendb()
	defer db.Close()

	rows_title, err := db.Query("SELECT id, tagname FROM Tags WHERE tagname = ?", tagname)

	if err != nil {
		log.Fatal(err.Error())
	}

	for rows_title.Next() {
		rows_title.Scan(&tag.Id, &tag.TagName)
	}

	return tag.Id
}

func FilterBooks(tagid []string) [][]string {
	var book structure.Books
	var Filter_book_data [][]string
	var book_sql string
	var new_tagid []interface{}

	book_sql = "SELECT title, author, publisher from Books WHERE "
	for index, id := range tagid {
		if index == 0 {
			book_sql += book_sql + "tagid LIKE CONCAT(?, '%')"

		} else {
			book_sql += book_sql + " AND tagid LIKE CONCAT(?, '%')"
		}
		new_tagid = append(new_tagid, Partitionid(id))
	}
	// fmt.Print(sql)
	Opendb()
	defer db.Close()

	rows_all, err := db.Query(book_sql, new_tagid...)

	if err != nil {
		log.Fatal(err.Error())
	}

	for rows_all.Next() {
		rows_all.Scan(&book.Book_title, &book.Book_author, &book.Book_publisher)

		book_data := []string{book.Book_title, book.Book_author, book.Book_publisher}
		Filter_book_data = append(Filter_book_data, book_data)
	}

	return Filter_book_data
}

func FilterBooks_ver2(tagid []string) [][]string {
	var book structure.Books
	var Filter_book_data [][]string

	Opendb()
	defer db.Close()

	rows_all, err := db.Query("SELECT title, author, publisher, id, tagid from Books")

	if err != nil {
		log.Fatal(err.Error())
	}

	for rows_all.Next() {
		rows_all.Scan(&book.Book_title, &book.Book_author, &book.Book_publisher, &book.Book_id, &book.Book_tagid)
		Filter := true
		for _, t := range tagid {
			fmt.Print(t)
			if !strings.Contains(book.Book_tagid, Partitionid(t)) {
				Filter = false
				break
			}
		}
		if Filter {
			book_data := []string{book.Book_title, book.Book_author, book.Book_publisher, GetRenterInfo(book.Book_id)}
			Filter_book_data = append(Filter_book_data, book_data)
		}
	}
	return Filter_book_data
}

func AddBookTag(tagid []string, isbn string){
	var tag_id string
	var alltag string

	Opendb()
	defer db.Close()

	rows_tagid, err := db.Query("SELECT tagid FROM Books WHERE isbn = ?", isbn)

	if err != nil {
		log.Fatal(err.Error())
	}

	for rows_tagid.Next() {
		rows_tagid.Scan(&tag_id)
	}
	fmt.Print(tag_id)

	for _, t := range tagid {
		alltag = alltag + Partitionid(t)
	}

	_ = db.QueryRow("UPDATE Books SET tagid = ? WHERE isbn = ?", alltag, isbn)
}

func GetUserinfo(id string, password string) (string, string) {
	var student structure.Students

	Opendb()
	defer db.Close()

	rows_title, err := db.Query("SELECT name, email FROM Students WHERE id = ? AND password = ?", id, password)

	if err != nil {
		log.Fatal(err.Error())
	}

	for rows_title.Next() {
		rows_title.Scan(&student.Name, &student.Email)
	}

	return student.Name, student.Email
}

func Rentbook(id string, isbn string) string {
	var book_id string
	var search_result string
	var res string

	search_result = "null"

	Opendb()
	defer db.Close()

	rows_bookid, err := db.Query("SELECT id FROM Books WHERE isbn = ?", isbn)

	if err != nil {
		log.Fatal(err.Error())
	}

	for rows_bookid.Next() {
		rows_bookid.Scan(&book_id)
	}
	fmt.Print(book_id)
	rows_histid, err := db.Query("SELECT id FROM Rent_hist WHERE bookid = ? AND returned=?", book_id, "1")

	if err != nil {
		log.Fatal(err.Error())
	}

	for rows_histid.Next() {
		rows_histid.Scan(&search_result)
	}
	// fmt.Print(search_result)
	if search_result == "null" {

		insert, err := db.Prepare("INSERT INTO Rent_hist(bookid,renterid,returned) VALUES(?, ?, ?)")
		if err != nil {
			log.Fatal(err.Error())
		}
		insert.Exec(book_id, id, "1")

		res = "貸出完了"
	} else {
		res = "貸出中です"
	}

	return res
}

func Returnbook(id string, isbn string) {
	var book_id string
	var search_result string

	Opendb()
	defer db.Close()

	rows_bookid, err := db.Query("SELECT id FROM Books WHERE isbn = ?", isbn)

	if err != nil {
		log.Fatal(err.Error())
	}

	for rows_bookid.Next() {
		rows_bookid.Scan(&book_id)
	}
	fmt.Print(book_id)

	rows_histid, err := db.Query("SELECT id FROM Rent_hist WHERE bookid = ? AND returned=?", book_id, "1")

	if err != nil {
		log.Fatal(err.Error())
	}

	for rows_histid.Next() {
		rows_histid.Scan(&search_result)
	}
	fmt.Print(search_result)

	if rows_histid != nil {
		_ = db.QueryRow("UPDATE Rent_hist SET returned = ? WHERE id = ? AND renterid = ?", "0", search_result, id)
	}
}

func GetAllBookData() [][]string {
	var all_book_data [][]string
	var book structure.Books

	Opendb()
	defer db.Close()

	rows_all, err := db.Query("SELECT title, isbn, author, publisher, id from Books")

	if err != nil {
		log.Fatal(err.Error())
	}
	for rows_all.Next() {
		rows_all.Scan(&book.Book_title, &book.Book_isbn, &book.Book_author, &book.Book_publisher, &book.Book_id)

		book_data := []string{book.Book_title, book.Book_isbn, book.Book_author, book.Book_publisher, book.Book_id}
		all_book_data = append(all_book_data, book_data)
	}

	return all_book_data
}

func GetRenterInfo(book_id string) string {
	var hist structure.Rent_hist
	var name string

	Opendb()
	defer db.Close()

	rows_renterid, err := db.Query("SELECT renterid FROM Rent_hist WHERE bookid = ? AND returned=?", book_id, "1")

	if err != nil {
		log.Fatal(err.Error())
	}

	for rows_renterid.Next() {
		rows_renterid.Scan(&hist.Hist_renterid)
	}

	name = "null"

	rows_name, err := db.Query("SELECT name FROM Students where id=?", hist.Hist_renterid)

	if err != nil {
		log.Fatal(err.Error())
	}

	for rows_name.Next() {
		rows_name.Scan(&name)
	}

	return name
}

func GetBookDetail(title string) (string, string, string, string) {
	var book structure.Books

	Opendb()
	defer db.Close()

	rows_detail, err := db.Query("SELECT author, publisher, item_caption, image_url from Books where title=?", title)

	if err != nil {
		log.Fatal(err.Error())
	}
	for rows_detail.Next() {
		rows_detail.Scan(&book.Book_author, &book.Book_publisher, &book.Book_caption, &book.Book_imageurl)
	}

	return book.Book_author, book.Book_publisher, book.Book_caption, book.Book_imageurl
}
