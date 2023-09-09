package main

import (
	"fmt" //вывод текста и его форматирование
	"io"
	"log"
	"math/rand"
	"net/http" //сетевой пакет
)

type ShortUrl struct { //объявляем структуру структуры
	ShortID   int
	Short_Url string
	LongUrl   string
}

// Write implements io.Writer.
func (ShortUrl) Write(p []byte) (n int, err error) {
	panic("unimplemented")
}

var SU ShortUrl //создаем структуру

func WebSepor(w http.ResponseWriter, r *http.Request) { //сепоратор типа запроса

	responseData, err := io.ReadAll(r.Body) //считывание тела пост запроса
	if err != nil {
		log.Fatal(err)
	}

	if r.Method == http.MethodPost && len(responseData) != 0 { //обработка пост запросов
		body := fmt.Sprintf("Method: %s 201 Created\r\n", r.Method)                   //создаем параметр body и выводим тип метода "Method: GET"
		body += fmt.Sprintf("Content-Type: %s\r\n", r.Header.Get("Content-Type"))     //читаем заголовок контент тайпа
		body += fmt.Sprintf("Content-Length: %s\r\n", r.Header.Get("Content-Length")) //читаем длину контента

		//заполняем данные структуры
		SU.ShortID = rand.Intn(100)
		SU.Short_Url = fmt.Sprintf("http://localhost:8080/%v", SU.ShortID)
		body += fmt.Sprintf(SU.Short_Url)

		SU.LongUrl = string(responseData) //запись параметра длинного урл

		w.Write([]byte(body)) //вывод боди

	} else if r.Method == http.MethodGet { //обработка гет зпросов
		rec := fmt.Sprintf(r.RequestURI)
		val := fmt.Sprint("/", SU.ShortID)

		if rec == val {
			body := fmt.Sprintf("Method: %s 307 Temporary Redirect\r\n", r.Method) //создаем параметр body и выводим тип метода "Method: POST"
			body += fmt.Sprintf("Location: %s\r\n", SU.LongUrl)
			w.Write([]byte(body))
		} else {
			body := fmt.Sprint("400\r\n")
			w.Write([]byte(body))
		}

	} else {
		body := fmt.Sprint("400\r\n")
		w.Write([]byte(body))
	}
}

// Тело программы
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, WebSepor) // Ловим запросы и отправляем в обработчик
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
