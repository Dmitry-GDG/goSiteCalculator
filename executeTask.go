package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

// изменение статуса взятой в работу задачи (на "в работе" или "завершена")
func upperStatusTaskToFile(taskNbr string) {
	muExpressions.Lock()
	defer muExpressions.Unlock()
	f, _ := os.Open(config_main.fileExpressions)
	fileScanner := bufio.NewScanner(f)
	dataArr := make([]string, 0)
	for fileScanner.Scan() {
		data := fileScanner.Text()
		data = removeR(data)
		dataArr = append(dataArr, data)
	}
	dataArrNew := make([]string, 0)
	for _, v := range dataArr {
		exp := strings.Split(v, ":")
		if len(exp) > 4 {
			if exp[0] == taskNbr {
				// if exp[3] == "0" {
				exp[3] = "1" // меняем статус на "в работе"
				// } else if exp[3] == "1" {
				// 	exp[3] = "2" // меняем статус на "выполнено"
				// }
			}
			dataArrNew = append(dataArrNew, exp[0]+":"+exp[1]+":"+exp[2]+":"+exp[3]+":"+exp[4])
		}
	}
	f.Close()
	// запишем изменения в файл
	dataStr := ""
	for _, v := range dataArrNew {
		if v != "" {
			dataStr += v + "\n"
		}
	}
	f1, err := os.Create(config_main.fileExpressions)
	if err != nil {
		panic(err)
	}
	defer f1.Close()
	_, err = f1.WriteString(dataStr)
	if err != nil {
		panic(err)
	}
}

func calcPower(v1, v2 int) int {
	log.Printf("вычисляется возведение в степень %d ^ %d\n", v1, v2)
	muConfigs.Lock()
	ttime := config_main.oetPower
	muConfigs.Unlock()
	time.Sleep(time.Duration(ttime) * time.Second)
	res := int(math.Pow(float64(v1), float64(v2)))
	log.Printf("получен результат возведения в степень %d ^ %d = %d\n", v1, v2, res)
	return res
}

func calcMinus(v1, v2 int) int {
	log.Printf("вычисляется разность %d - %d\n", v1, v2)
	muConfigs.Lock()
	ttime := config_main.oetMinus
	muConfigs.Unlock()
	time.Sleep(time.Duration(ttime) * time.Second)
	res := v1 - v2
	log.Printf("получен результат разности %d - %d = %d\n", v1, v2, res)
	return res
}

func calcPlus(v1, v2 int) int {
	log.Printf("вычисляется сложение %d + %d\n", v1, v2)
	muConfigs.Lock()
	ttime := config_main.oetPlus
	muConfigs.Unlock()
	time.Sleep(time.Duration(ttime) * time.Second)
	res := v1 + v2
	log.Printf("получен результат сложения %d + %d = %d\n", v1, v2, res)
	return res
}

func calcMultiply(v1, v2 int) int {
	log.Printf("вычисляется произведение %d * %d\n", v1, v2)
	muConfigs.Lock()
	ttime := config_main.oetMultiply
	muConfigs.Unlock()
	time.Sleep(time.Duration(ttime) * time.Second)
	res := v1 * v2
	log.Printf("получен результат произведения %d * %d = %d\n", v1, v2, res)
	return res
}

func calcDivide(v1, v2 int) int {
	log.Printf("вычисляется деление %d / %d\n", v1, v2)
	muConfigs.Lock()
	ttime := config_main.oetDivide
	muConfigs.Unlock()
	time.Sleep(time.Duration(ttime) * time.Second)
	res := v1 / v2
	log.Printf("получен результат деления %d / %d = %d\n", v1, v2, res)
	return res
}

func saveResultToFile(taskNbr, result string) {
	muExpressions.Lock()
	defer muExpressions.Unlock()
	f, _ := os.Open(config_main.fileExpressions)
	fileScanner := bufio.NewScanner(f)
	dataArr := make([]string, 0)
	for fileScanner.Scan() {
		data := fileScanner.Text()
		data = removeR(data)
		dataArr = append(dataArr, data)
	}
	dataArrNew := make([]string, 0)
	for _, v := range dataArr {
		exp := strings.Split(v, ":")
		if len(exp) > 4 {
			if exp[0] == taskNbr {
				if exp[3] == "1" {
					exp[3] = "2" // меняем статус на "выполнено"
				}
				exp[4] = result
			}
			dataArrNew = append(dataArrNew, exp[0]+":"+exp[1]+":"+exp[2]+":"+exp[3]+":"+exp[4])
		}
	}
	f.Close()
	// запишем изменения в файл
	dataStr := ""
	for _, v := range dataArrNew {
		if v != "" {
			dataStr += v + "\n"
		}
	}
	f1, err := os.Create(config_main.fileExpressions)
	if err != nil {
		panic(err)
	}
	defer f1.Close()
	_, err = f1.WriteString(dataStr)
	if err != nil {
		panic(err)
	}
}

func executeTask(task []string) {
	log.Println("отправлена на выполнение новая задача: ", task[2])
	upperStatusTaskToFile(task[0]) // переводим задачу в статус "в работе"

	tmpArr := make([]int, 0)
	findErr := false
	var result int
	taskSplit := strings.Split(task[2], " ")
	for i := 0; i < len(taskSplit); i++ {
		nbr, err := strconv.Atoi(taskSplit[i])
		if err != nil && len(tmpArr) > 1 {
			v2 := tmpArr[len(tmpArr)-1]
			v1 := tmpArr[len(tmpArr)-2]
			// fmt.Println("v1,v2  ", v1, v2, "task[i] ", taskSplit[i])
			tmpArr = tmpArr[:len(tmpArr)-2]
			if taskSplit[i] == "/" {
				if v2 == 0 {
					findErr = true
					break
				}
				result = calcDivide(v1, v2)
			} else if taskSplit[i] == "*" {
				result = calcMultiply(v1, v2)
			} else if taskSplit[i] == "-" {
				result = calcMinus(v1, v2)
			} else if taskSplit[i] == "+" {
				result = calcPlus(v1, v2)
			} else if taskSplit[i] == "^" {
				result = calcPower(v1, v2)
			}
			tmpArr = append(tmpArr, result)
		} else {
			tmpArr = append(tmpArr, nbr)
		}
		if findErr {
			break
		}
	}

	// fmt.Println("result: ", result)

	saveResultToFile(task[0], strconv.Itoa(result)) // переводим задачу в статус "выполнено" и сохраняем результат
	muConfigs.Lock()
	config_main.qtyBusyServers-- // освобождаем один сервер
	muConfigs.Unlock()
}

// аркестратор, запускает задачи
func launchTasks() {
	// бесконечный цикл с паузой 1 сек
	for {
		muConfigs.Lock()

		if config_main.qtyServers > config_main.qtyBusyServers {
			muUnDone.Lock()
			if len(unDone) > 0 {
				for len(unDone) > 0 && config_main.qtyServers > config_main.qtyBusyServers { // пока есть задачи и свободные серверы
					// popFront
					newTask := unDone[0]
					unDone = unDone[1:]          // удаляем задачу из очереди
					go executeTask(newTask)      // запускаем задачу в отдельной горутине
					config_main.qtyBusyServers++ // занятых серверов стало больше
				}
			}
			muUnDone.Unlock()
		}

		muConfigs.Unlock()
		time.Sleep(time.Second)
	}
}

// проверяем, есть ли нерешенные с прошлого запуска задачи
func checkUndoneJobs() {
	muExpressions.Lock()
	nextNumber := config_main.nextNumber
	muExpressions.Unlock()

	if nextNumber > 1 { // были присланы выражения при прошлой работе
		muExpressions.Lock()
		defer muExpressions.Unlock()
		f, _ := os.Open(config_main.fileExpressions)
		fileScanner := bufio.NewScanner(f)
		dataArr := make([]string, 0)
		for fileScanner.Scan() {
			data := fileScanner.Text()
			data = removeR(data)
			dataArr = append(dataArr, data)
		}
		found := false // если встретятся нерешённые задачи, надо будет изменить их статус на "ожидание"
		dataArrNew := make([]string, 0)
		for _, v := range dataArr {
			exp := strings.Split(v, ":")
			if len(exp) > 4 {
				if exp[3] == "0" || exp[3] == "1" { //  с прошлого запуска остались "status: waiting" или "status: in process"
					exp[3] = "0" // меняем статус на "ожидание"
					found = true
					f.Close()
					// запишем задачу в очередь
					addNewUnDoneTask(exp)
					// muUnDone.Lock()
					// unDone = append(unDone, exp)
					// // fmt.Println(unDone)
					// muUnDone.Unlock()
				}
				dataArrNew = append(dataArrNew, exp[0]+":"+exp[1]+":"+exp[2]+":"+exp[3]+":"+exp[4])
			}
		}
		if found { // запишем изменения в файл
			dataStr := ""
			for _, v := range dataArrNew {
				if v != "" {
					dataStr += v + "\n"
				}
			}
			f1, err := os.Create(config_main.fileExpressions)
			if err != nil {
				panic(err)
			}
			defer f1.Close()
			_, err = f1.WriteString(dataStr)
			if err != nil {
				panic(err)
			}
		} else {
			f.Close()
		}
	}
}

func readDijkstra() error {
	dijkstraSlice := make([]string, 0)
	muDijkstra.Lock()
	defer muDijkstra.Unlock()

	f, err := os.Open(config_main.fileDijkstra)
	if err != nil {
		return err
	}
	defer f.Close()
	fileScanner := bufio.NewScanner(f)
	for fileScanner.Scan() {
		dataSlice := strings.Split(fileScanner.Text(), ":")
		if len(dataSlice) > 0 {
			dijkstraSlice = append(dijkstraSlice, dataSlice[0])
		}
	}

	log.Println("Прочитано из файла dijkstra:", dijkstraSlice)
	return nil
}

// func savePreparedData(dijkstraSlice []dijkstraData) {
// 	muDijkstra.Lock()
// 	defer muDijkstra.Unlock()
// 	file, err := os.OpenFile("./db/dijkstra.db", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer file.Close()
// 	if _, err = file.WriteString(dijkstraSlice + "\n"); err != nil {
// 		panic(err)
// 	}
// 	log.Println("В бд добавлена новая строка:", dijkstraSlice)
// 	savePreparedData(dijkstraSlice)

// 	return nil
// }
