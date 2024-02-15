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

// изменение статуса взятой в работу задачи (на "в работе")
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
				exp[3] = "1" // меняем статус на "в работе"
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

	tmpArr := make([]int, 0) // вспомогательный срез для хранения чисел
	findErr := false
	var result string
	var resNow int
	taskSplit := strings.Split(task[2], " ")
	for i := 0; i < len(taskSplit); i++ {
		nbr, err := strconv.Atoi(taskSplit[i])
		if err != nil && len(tmpArr) > 1 { // попался арифметический знак
			v2 := tmpArr[len(tmpArr)-1]
			v1 := tmpArr[len(tmpArr)-2]
			tmpArr = tmpArr[:len(tmpArr)-2]
			if taskSplit[i] == "/" {
				if v2 == 0 {
					findErr = true
					result = "nil"
					log.Printf("ошибка: деление на ноль не поддерживается")
					break
				}
				resNow = calcDivide(v1, v2)
			} else if taskSplit[i] == "*" {
				resNow = calcMultiply(v1, v2)
			} else if taskSplit[i] == "-" {
				if v2 > v1 {
					findErr = true
					result = "nil"
					log.Printf("ошибка: операции с отрицательнеыми числами не поддерживаются")
					break
				}
				resNow = calcMinus(v1, v2)
			} else if taskSplit[i] == "+" {
				resNow = calcPlus(v1, v2)
			} else if taskSplit[i] == "^" {
				resNow = calcPower(v1, v2)
			}
			tmpArr = append(tmpArr, resNow) // запишем результат вычислений во вспомогательный срез
		} else { // перекладываем число во вспомогательный срез
			tmpArr = append(tmpArr, nbr)
		}
		if findErr {
			break
		}
	}

	if result == "" {
		result = strconv.Itoa(resNow)
	}
	saveResultToFile(task[0], result) // переводим задачу в статус "выполнено" и сохраняем результат
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
