package main

import (
	"log"
	"net/http"
	"sync"
)

// основной конфиг программы
type Config_main struct {
	port            string // порт, на котором запускается
	checkUnicTask   int    // проверка введённого выражения на уникальность (0 не проверять, 1 проверять)
	conf            string // файл для сохранения последнего конфига
	fileExpressions string // файл бд, в который будут записываться выражения
	fileDijkstra    string // файл для записи подготовленных выражений
	lastExpression  string // Последнее введённое выражение
	oetPlus         int    // operation execution time: + , seconds
	oetMinus        int    // operation execution time: - , seconds
	oetMultiply     int    // operation execution time: * , seconds
	oetDivide       int    // operation execution time: / , seconds
	oetPower        int    // operation execution time: ^ , seconds
	nextNumber      int    // номер, который будет присвоен следующему присланному выражению
	qtyServers      int    // количество возможных параллельных вычислений (серверов)
	qtyBusyServers  int    // количество занятых серверов вычислений
}

// type Task struct {
// 	v1, v2 int
// 	sign   string
// }

// type Expression struct {
// 	expression string
// 	isDone     bool
// }

var config_main Config_main
var muConfigs sync.Mutex
var muExpressions sync.Mutex
var muDijkstra sync.Mutex
var unDone [][]string // слайс для хранения невыполненных задач
var muUnDone sync.Mutex

// var muLastExpression sync.Mutex

func initConfig() {
	config_main.port = "8000"
	config_main.checkUnicTask = 0
	config_main.conf = "./configs/config_prev.conf"
	config_main.fileExpressions = "./db/expressions.db"
	config_main.fileDijkstra = "./db/dijkstra.db"
	config_main.lastExpression = ""
	config_main.oetDivide = 6
	config_main.oetMinus = 6
	config_main.oetMultiply = 6
	config_main.oetPlus = 6
	config_main.oetPower = 6
	config_main.nextNumber = 1
	config_main.qtyServers = 5
	config_main.qtyBusyServers = 0
}

func main() {
	// инициализируем конфиг начальными значениями
	initConfig()
	// получаем из файла новый конфиг
	getConfig()

	// проверяем, есть ли не завершённые с прошлого запуска программы работы (если при запуске указано сохранение настроек прошлого запуска)
	checkUndoneJobs()

	// запускаем вычисление выражений в отдельной горутине
	go launchTasks()

	// http server
	mux := http.NewServeMux()

	mux.HandleFunc("/data/", getData)
	mux.HandleFunc("/data", getData)
	mux.HandleFunc("/list/", getList)
	mux.HandleFunc("/list", getList)
	mux.HandleFunc("/settings/", getSettings)
	mux.HandleFunc("/settings", getSettings)
	mux.HandleFunc("/resources/", getResources)
	mux.HandleFunc("/resources", getResources)
	mux.HandleFunc("/clearDb/", clearDb)
	mux.HandleFunc("/clearDb", clearDb)
	mux.HandleFunc("/clearDbAttention/", clearDbAttention)
	mux.HandleFunc("/clearDbAttention", clearDbAttention)
	mux.HandleFunc("/", getRoot)

	log.Println("калькулятор запускается на порту: ", config_main.port)
	if err := http.ListenAndServe(":"+config_main.port, mux); err != nil {
		log.Fatal(err)
	}
}

// go run $(ls *.go)
// curl http://localhost:8080/metrics
// Ctrl + c
