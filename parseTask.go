package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// вспомогательная структура для парсинга выражений
type dijkstraData struct {
	val     int
	isSign  bool
	valSign string
}

// обрабатываемые программой арифметические знаки и скобки
func checkDictSigns(b byte) bool {
	if b == '(' || b == '+' || b == '-' || b == '*' || b == '/' || b == '^' || b == ')' {
		return true
	}
	return false
}

// разбиение на токены: числа, знаки, скобки
func tokenize(data string) ([]string, error) {
	outp := make([]string, 0)
	for i := 0; i < len(data); {
		tmpI := i
		tmp := ""
		for {
			if i < len(data) && data[i] >= '0' && data[i] <= '9' {
				tmp += string(data[i])
				i++
			} else {
				break
			}
		}
		if tmp != "" {
			outp = append(outp, tmp)
		}
		if i < len(data) && checkDictSigns(data[i]) {
			outp = append(outp, string(data[i]))
			i++
			if i == len(data) && data[i-1] != ')' {
				return []string{}, errors.New("error: expression is not valid")
			}
		}
		if i == tmpI {
			return []string{}, errors.New("error: expression is not valid")
		}
	}
	return outp, nil
}

// https://habr.com/ru/articles/596925/
// перевод выражение в форму Польской нотации
func dijkstra(tokens []string) ([]dijkstraData, error) {
	//	Список и приоритет операторов
	dictSigns := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
		"^": 3,
	}
	tmpDictSigns := make(map[string]bool, 0) // if not set, tmpDictSigns will be false(left-associative)

	dijkstraSlice := make([]dijkstraData, 0)      // срез данных Дейкстры
	dijkstraSliceSigns := make([]dijkstraData, 0) // стек операторов

	for _, token := range tokens {
		v, err := strconv.Atoi(token)
		if err == nil {
			dijkstraSlice = append(dijkstraSlice, dijkstraData{val: v})
		} else {
			if token == "(" {
				dijkstraSliceSigns = append(dijkstraSliceSigns, dijkstraData{valSign: token, isSign: true})
			} else if token == ")" {
				// если попалась закрывающая скобка - переносим из стека все операторы, пока не попадётся открывающаяся скобка
				found := false
				for len(dijkstraSliceSigns) > 0 {
					newMember := dijkstraSliceSigns[len(dijkstraSliceSigns)-1]
					dijkstraSliceSigns = dijkstraSliceSigns[:len(dijkstraSliceSigns)-1]
					if newMember.valSign == "(" {
						found = true
						break
					} else {
						dijkstraSlice = append(dijkstraSlice, newMember)
					}
				}
				if !found {
					// если не было открывающейся скобки
					return []dijkstraData{}, errors.New("error: expression is not valid, mismatched parentheses found")
				}
			} else {
				// log.Println("HERE")
				priority, ok := dictSigns[token]
				if !ok {
					return []dijkstraData{}, fmt.Errorf("error: expression is not valid, unknown operator: %v", token)
				}

				rightAssociative := tmpDictSigns[token]
				for len(dijkstraSliceSigns) > 0 {
					top := dijkstraSliceSigns[len(dijkstraSliceSigns)-1]

					if top.valSign == "(" {
						break
					}

					prevPriority := dictSigns[top.valSign]

					if (rightAssociative && priority < prevPriority) || (!rightAssociative && priority <= prevPriority) {
						// pop current operator
						dijkstraSliceSigns = dijkstraSliceSigns[:len(dijkstraSliceSigns)-1]
						dijkstraSlice = append(dijkstraSlice, top)
					} else {
						break
					}
				} // end of for len(dijkstraSliceSigns) > 0
				dijkstraSliceSigns = append(dijkstraSliceSigns, dijkstraData{valSign: token, isSign: true})
			} // end of if token == "("
		}

	} // end of for _, token := range tokens
	// log.Println("HERE")
	for len(dijkstraSliceSigns) > 0 {
		// pop
		newMember := dijkstraSliceSigns[len(dijkstraSliceSigns)-1]
		dijkstraSliceSigns = dijkstraSliceSigns[:len(dijkstraSliceSigns)-1]

		if newMember.valSign == "(" {
			return []dijkstraData{}, errors.New("error: expression is not valid, mismatched parentheses found")
		}
		dijkstraSlice = append(dijkstraSlice, newMember)
	}
	return dijkstraSlice, nil
}

// проверка задачи на уникальность (если задано такое условие)
func checkUnicTask(data string) (bool, error) {
	muExpressions.Lock()
	defer muExpressions.Unlock()

	f, err := os.Open(config_main.fileExpressions)
	if err != nil {
		return false, err
	}
	defer f.Close()
	fileScanner := bufio.NewScanner(f)
	for fileScanner.Scan() {
		exp := strings.Split(fileScanner.Text(), ":")
		// fmt.Println("data: ", data, "from db: ", exp[0])
		if len(exp) > 1 {
			if data == exp[1] {
				return false, nil
			}
		}
	}

	return true, nil
}

// добавление в глобальную переменную новой невыполненной задачи
func addNewUnDoneTask(newData []string) {
	muUnDone.Lock()
	unDone = append(unDone, newData)
	muUnDone.Unlock()
}

// добавление в БД новой задачи
func addNewData(data, dijkstra string) error {
	muExpressions.Lock()
	defer muExpressions.Unlock()
	file, err := os.OpenFile(config_main.fileExpressions, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	muConfigs.Lock()
	newData := strconv.Itoa(config_main.nextNumber) + ":" + data + ":" + dijkstra + ":0:0\n"
	config_main.nextNumber++
	muConfigs.Unlock()
	if _, err = file.WriteString(newData); err != nil {
		muConfigs.Lock()
		config_main.nextNumber--
		muConfigs.Unlock()
		file.Close()
		return err
	}
	file.Close()
	var newDataArr []string
	newDataArr = append(newDataArr, strconv.Itoa(config_main.nextNumber-1))
	newDataArr = append(newDataArr, data)
	newDataArr = append(newDataArr, dijkstra)
	newDataArr = append(newDataArr, "0")
	newDataArr = append(newDataArr, "0")
	addNewUnDoneTask(newDataArr)
	savePresentConfigToFile()

	return nil
}

// основная функция парсинга и сохранения нового выражения
func parseAndSaveTask(data string) error {
	data = strings.Replace(data, "**", "^", -1)

	// проверим новое выражение на уникальность в бд, в случае необходимости
	if config_main.checkUnicTask == 1 {
		unic, err := checkUnicTask(data)
		if err != nil {
			panic(err)
		}
		if !unic {
			log.Printf("выражение %s ранее было добавлено в БД, посмотрите результат на соответствующей странице", data)
			return nil
		}
	}

	// разобьём выражение на токены (числа и знаки)
	tokens, err := tokenize(data)
	if err != nil {
		return err
	}

	// перенесём эти токены в слайс согласно правилу Польской нотации
	dijkstraSlice, err := dijkstra(tokens)
	if err != nil {
		// log.Println("error: выражение невалидно")
		return err
	}

	// переведём этот слайс в строку
	dijkstra := ""
	for i := 0; i < len(dijkstraSlice); i++ {
		if !dijkstraSlice[i].isSign {
			dijkstra += strconv.Itoa(dijkstraSlice[i].val)
		} else {
			dijkstra += dijkstraSlice[i].valSign
		}
		if i < len(dijkstraSlice)-1 {
			dijkstra += " "
		}
	}

	// добавим новое выражение в бд
	if err := addNewData(data, dijkstra); err != nil {
		panic(err)
	}
	log.Println("В бд добавлена новая строка:", data)

	// обновляем данные о последнем выражении в переменной
	muConfigs.Lock()
	config_main.lastExpression = data
	muConfigs.Unlock()

	// обновляем данные о последнем выражении в файле последнего конфига
	savePresentConfigToFile()

	return nil
}

// поступили новые настройки на изменение времени исполнения арифметических операций
func parseAndSaveSettings(name string) error {
	nameArr := strings.Split(name, "=")
	if len(nameArr) < 2 {
		return errors.New("error: new settings is not valid (not enough data)")
	}
	newData, err := strconv.Atoi(nameArr[1])
	if err != nil {
		return errors.New("error: new settings is not valid (need integer)")
	}
	if newData < 0 {
		newData = 0
	}

	muConfigs.Lock()
	if nameArr[0] == "oetDivide" {
		config_main.oetDivide = newData
	} else if nameArr[0] == "oetMinus" {
		config_main.oetMinus = newData
	} else if nameArr[0] == "oetMultiply" {
		config_main.oetMultiply = newData
	} else if nameArr[0] == "oetPlus" {
		config_main.oetPlus = newData
	} else if nameArr[0] == "oetPower" {
		config_main.oetPower = newData
	} else {
		muConfigs.Unlock()
		return errors.New("error: new settings is not valid (unknown settings)")
	}
	muConfigs.Unlock()
	log.Println("В бд изменены данные:", nameArr[0]+"="+strconv.Itoa(newData))

	// обновляем данные о последнем выражении в файле последнего конфига
	savePresentConfigToFile()
	return nil
}

// поступили новые настройки на изменение количества серверов
func parseAndSaveServers(name string) error {
	nameArr := strings.Split(name, "=")
	if len(nameArr) < 2 {
		return errors.New("error: new resources is not valid (not enough data)")
	}
	newData, err := strconv.Atoi(nameArr[1])
	if err != nil {
		return errors.New("error: new resources is not valid (need integer)")
	}
	if newData < 1 {
		newData = 1 // менее 1 сервера нет смысла устанавливать
	}

	muConfigs.Lock()
	if nameArr[0] == "qtyServers" {
		config_main.qtyServers = newData
	} else {
		muConfigs.Unlock()
		return errors.New("error: new resources is not valid (unknown settings)")
	}
	muConfigs.Unlock()
	log.Println("В бд изменены данные:", nameArr[0]+"="+strconv.Itoa(newData))

	// обновляем данные о последнем выражении в файле последнего конфига
	savePresentConfigToFile()
	return nil
}
