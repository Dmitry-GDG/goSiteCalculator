package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Обработка обращения на главную страницу сайта
func getRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		return
	}
	if r.URL.Path != "/" && r.URL.Path != "/index" {
		http.NotFound(w, r)
		log.Println("error getting url", r.URL.Path, ": 404 page not found")
		return
	}
	log.Println("поступило обращение к главной странице")

	// fmt.Fprint(w, "Для использования калькулятора отправьте на сервер запрос типа:\n\tcurl http://localhost:8001/data/?data='2+6*7'")

	// path := filepath.Join("front", "index.html")
	// //создаем html-шаблон
	// tmpl, err := template.ParseFiles(path)
	// if err != nil {
	// 	http.Error(w, err.Error(), 400)
	// 	return
	// }
	// //выводим шаблон клиенту в браузер
	// err = tmpl.Execute(w, nil)
	// if err != nil {
	// 	http.Error(w, err.Error(), 400)
	// 	return
	// }
	//верстаем контент страницы в виде обычной строки

	content := ""
	f, err := os.Open("./front/up.html")
	fileScanner := bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	f, err = os.Open("./front/main.html")
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	f, err = os.Open("./front/middle.html")
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	content += `<p style="font-size: 1.4em;">В калькуляторе реализованы следующие функции: сложение, вычитание, умножение, деление, возведение в степень, скобки и приоритет операций.<br>Для использования калькулятора отправьте на сервер запрос типа:<br><br>curl http://localhost:` + config_main.port + `/data/?data="2+6*7"<br>curl http://localhost:` + config_main.port + `/data/?data="3+4*2/(1-5)^2^3"</p>`
	f, err = os.Open("./front/down.html")
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	//создаем html-шаблон
	tmpl, err := template.New("example").Parse(content)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	tmpl.Execute(w, content)
}

// Обработка обращения на страницу ввода данных
func getData(w http.ResponseWriter, r *http.Request) {
	// name := r.URL.Query().Get("data")
	name, err := url.PathUnescape(url.QueryEscape(r.URL.Query().Get("data")))
	if err != nil {
		// fmt.Errorf(err.Error())
		log.Println(err.Error())
	}
	if name != "" {
		log.Println("получено новое выражение: ", name)
		if err := parseAndSaveTask(name); err != nil {
			http.Error(w, err.Error(), 400)
			log.Println(err)
		}
	} else {
		log.Println("поступило обращение к странице ввода данных")
	}
	// fmt.Fprint(w, name)

	content := ""
	f, err := os.Open("./front/up.html")
	fileScanner := bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	f, err = os.Open("./front/data.html")
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	f, err = os.Open("./front/middle.html")
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()

	muConfigs.Lock()
	content += "Последнее полученное выражение: " + config_main.lastExpression + "<br><br>"
	content += `<p style="font-size: 1.4em;">Для использования калькулятора отправьте на сервер запрос типа:<br><br>curl http://localhost:` + config_main.port + `/data/?data="2+6*7"<br>curl http://localhost:` + config_main.port + `/data/?data="3+4*2/(1-5)^2^3"</p>`
	muConfigs.Unlock()

	f, err = os.Open("./front/down.html")
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	//создаем html-шаблон
	data := "getData"
	tmpl, err := template.New("data").Parse(content)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	tmpl.Execute(w, data)
}

// Обработка обращения на страницу получения информации о ранее отправлеенных задачах и статусе их решения
func getList(w http.ResponseWriter, r *http.Request) {
	log.Println("поступило обращение к странице со списком выражений")
	// fmt.Fprint(w, "Привет")

	content := ""
	f, err := os.Open("./front/up.html")
	fileScanner := bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	f, err = os.Open("./front/list.html")
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	f, err = os.Open("./front/middle.html")
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	muExpressions.Lock()
	f, err = os.Open(config_main.fileExpressions)
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		exp := strings.Split(fileScanner.Text(), ":")
		if len(exp) > 3 {
			if exp[3] == "0" {
				exp[3] = "status: waiting"
			} else if exp[3] == "1" {
				exp[3] = "status: in process"
			} else if exp[3] == "2" {
				exp[3] = "status: done"
			}
			content += "Номер задачи: " + exp[0] + "<br>&nbsp;" + exp[1] + "&nbsp;" + exp[3] + "<br>Выражение в виде, подготовленном для вычислений: " + exp[2]
			if exp[3] == "status: done" {
				content += "<br>&nbsp;Итог: " + exp[4]
			}
			content += "<br>-----------------<br>"
		}
	}
	f.Close()
	muExpressions.Unlock()

	// content += "-----------------<br>Выражения в виде, подготовленном для вычислений:<br>"

	// muDijkstra.Lock()

	// f, err = os.Open(config_main.fileDijkstra)

	// fileScanner = bufio.NewScanner(f)
	// for fileScanner.Scan() {
	// 	dataSlice := strings.Split(fileScanner.Text(), ":")
	// 	if len(dataSlice) > 0 {
	// 		content += dataSlice[0] + "<br>"
	// 	}
	// }
	// f.Close()
	// muDijkstra.Unlock()

	f, err = os.Open("./front/down.html")
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	//создаем html-шаблон
	tmpl, err := template.New("example").Parse(content)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	tmpl.Execute(w, content)
}

// Обработка обращения на страницу получения информации о настройках системы
func getSettings(w http.ResponseWriter, r *http.Request) {
	// log.Println("поступило обращение к странице со списком настроек системы")

	name, err := url.PathUnescape(url.QueryEscape(r.URL.Query().Get("settings")))
	if err != nil {
		log.Println(err.Error())
	}
	if name != "" {
		log.Println("получены новые настройки: ", name)
		if err := parseAndSaveSettings(name); err != nil {
			http.Error(w, err.Error(), 400)
			log.Println(err)
		}
	} else {
		log.Println("поступило обращение к странице со списком настроек системы")
	}
	// fmt.Fprint(w, name)

	content := ""
	f, err := os.Open("./front/up.html")
	fileScanner := bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	f, err = os.Open("./front/settings.html")
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	f, err = os.Open("./front/middle.html")
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()

	muConfigs.Lock()
	content += "Настройки серверов:<br>Время, требующееся для обработки сложения: " + strconv.Itoa(config_main.oetPlus) + " сек<br>Время, требующееся для обработки вычитания: " + strconv.Itoa(config_main.oetMinus) + " сек<br>Время, требующееся для обработки умножения: " + strconv.Itoa(config_main.oetMultiply) + " сек<br>Время, требующееся для обработки деления: " + strconv.Itoa(config_main.oetDivide) + " сек<br>Время, требующееся для обработки возведения в степень: " + strconv.Itoa(config_main.oetPower) + " сек"
	content += `<p style="font-size: 1.4em;">Для изменения настроек отправьте на сервер запрос типа:<br><br>curl http://localhost:` + config_main.port + `/settings/?settings="oetMinus=55"</p>`
	muConfigs.Unlock()

	f, err = os.Open("./front/down.html")
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	//создаем html-шаблон
	tmpl, err := template.New("example").Parse(content)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	tmpl.Execute(w, content)
}

// Обработка обращения на страницу получения информации о занятых ресурсах ситстемы
func getResources(w http.ResponseWriter, r *http.Request) {
	// log.Println("поступило обращение к странице со списком вычислительных мощностей")

	name, err := url.PathUnescape(url.QueryEscape(r.URL.Query().Get("resources")))
	if err != nil {
		log.Println(err.Error())
	}
	if name != "" {
		log.Println("получены новое количество серверов: ", name)
		if err := parseAndSaveServers(name); err != nil {
			http.Error(w, err.Error(), 400)
			log.Println(err)
		}
	} else {
		log.Println("поступило обращение к странице со списком вычислительных мощностей")
	}

	content := ""
	f, err := os.Open("./front/up.html")
	fileScanner := bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	f, err = os.Open("./front/resources.html")
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	f, err = os.Open("./front/middle.html")
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()

	muConfigs.Lock()
	content += "Общее количество серверов: " + strconv.Itoa(config_main.qtyServers) + "<br>в том числе количество занятых серверов: " + strconv.Itoa(config_main.qtyBusyServers)
	content += `<p style="font-size: 1.4em;">Для изменения настроек отправьте на сервер запрос типа:<br><br>curl http://localhost:` + config_main.port + `/resources/?resources="qtyServers=4"</p>`
	muConfigs.Unlock()

	f, err = os.Open("./front/down.html")
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	//создаем html-шаблон
	tmpl, err := template.New("example").Parse(content)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	tmpl.Execute(w, content)
}

// Обработка обращения на страницу очистки БД
func clearDbAttention(w http.ResponseWriter, r *http.Request) {
	log.Println("поступило обращение к странице очистки баз данных")

	content := ""
	f, err := os.Open("./front/up.html")
	fileScanner := bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	f, err = os.Open("./front/clearDbAttention.html")
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	f, err = os.Open("./front/middle.html")
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()

	content += `<p style="font-size: 1.4em;">Для очистки БД отправьте на сервер запрос:<br><br>curl http://localhost:` + config_main.port + `/clearDb</p>`
	// content += `<form><input class="styled" type="button" value="Очистить БД" /></form>`

	// 	content += `  <form action="">
	// <p>Для очистки БД введите yes и нажмите кнобку &laquo;Очистить БД&raquo;</p>
	// <p><input type="text" id="ddata"> <span id="status"></span></p>
	// <p><input type="button" value="Очистить БД" onclick="isOk()"></p>
	// </form>`

	f, err = os.Open("./front/down.html")
	fileScanner = bufio.NewScanner(f)
	for fileScanner.Scan() {
		content += fileScanner.Text()
	}
	f.Close()
	//создаем html-шаблон
	tmpl, err := template.New("example").Parse(content)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	tmpl.Execute(w, content)
}

// Очистка Баз Данных
func clearDbFunc() {
	muExpressions.Lock()
	defer muExpressions.Unlock()
	file, err := os.OpenFile(config_main.fileExpressions, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	file.Close()

	muDijkstra.Lock()
	defer muDijkstra.Unlock()
	file, err = os.OpenFile(config_main.fileDijkstra, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	file.Close()

	config_main.lastExpression = ""

	log.Println("Базы данных очищены")
}

func clearDb(w http.ResponseWriter, r *http.Request) {
	log.Println("поступило обращение на очистку баз данных")
	clearDbFunc()
	fmt.Fprint(w, "Базы данных очищены")
}
