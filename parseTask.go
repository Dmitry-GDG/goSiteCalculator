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

	dijkstraSlice := make([]dijkstraData, 0)      // срез данных Дейкстры (чисел и знаков)
	dijkstraSliceSigns := make([]dijkstraData, 0) // стек операторов
	lastPriority1 := 0                            // для обработки отрицательных чисел и последовательности знаков + и -, кол-во последовательных операторов + и -
	lastPriority23 := 0                           // кол-во последовательных операторов *, / и операторов возведения в степень

	for i := 0; i < len(tokens); i++ {
		v, err := strconv.Atoi(tokens[i])
		if err == nil { // если это число
			// fmt.Println("HERE")
			// fmt.Printf("tokens[%d] = %s\n", i, tokens[i])
			if lastPriority1 > 0 {
				// fmt.Println("lastPriority1 = ", lastPriority1, "len(dijkstraSliceSigns) = ", len(dijkstraSliceSigns))
				for len(dijkstraSliceSigns) > 0 && lastPriority1 > 1 {
					top := dijkstraSliceSigns[len(dijkstraSliceSigns)-1]
					if top.valSign == "-" { // встретился знак минус, меняем значение числа
						v *= -1
					}

					dijkstraSliceSigns = dijkstraSliceSigns[:len(dijkstraSliceSigns)-1] // удаляем знак из стека

					lastPriority1-- // уменьшаем количество последовательных операторов с приоритетом 1
					// i++
				}
				// dijkstraSlice = append(dijkstraSlice, dijkstraData{val: v})
				// top := dijkstraSliceSigns[len(dijkstraSliceSigns)-1]
				// dijkstraSliceSigns = dijkstraSliceSigns[:len(dijkstraSliceSigns)-1]
				// dijkstraSlice = append(dijkstraSlice, top)
				lastPriority1--
				if lastPriority1 > 1 {
					return []dijkstraData{}, fmt.Errorf("error: expression is not valid: ", tokens[i])
				}
				// i++
			} else {
			}
			dijkstraSlice = append(dijkstraSlice, dijkstraData{val: v})
			lastPriority1, lastPriority23 = 0, 0 // обнуляем маркеры

		} else { // это не число, проверяем дальше
			if tokens[i] == "(" { // открывающую скобку переносим в спомогательный срез (стек) операторов
				dijkstraSliceSigns = append(dijkstraSliceSigns, dijkstraData{valSign: tokens[i], isSign: true})
				lastPriority1, lastPriority23 = 0, 0 // обнуляем маркеры
			} else if tokens[i] == ")" {
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
					// если не было открывающейся скобки - ошибка
					return []dijkstraData{}, errors.New("error: expression is not valid, mismatched parentheses found")
				}
				lastPriority1, lastPriority23 = 0, 0 // обнуляем маркеры
			} else {
				// в остальных случаях проверяем на приоритет оператора
				priority, ok := dictSigns[tokens[i]] // это приоритет данного оператора
				if !ok {
					return []dijkstraData{}, fmt.Errorf("error: expression is not valid, unknown operator: %v", tokens[i])
				}
				if priority > 1 && (lastPriority1 > 0 || lastPriority23 > 0) { // после ряда операторов с приоритетом 1 не может идти оператор с другим приоритетом
					return []dijkstraData{}, fmt.Errorf("error: expression is not valid, not correct order of operators: %v", tokens[i])
				}
				if priority == 1 {
					lastPriority1++
					lastPriority23 = 0
				}
				if priority > 1 {
					lastPriority23++
					lastPriority1 = 0
				}

				if (priority != 1) || (priority == 1 && lastPriority1 == 1) {
					rightAssociative := tmpDictSigns[tokens[i]] // существовал? - if not set, tmpDictSigns will be false(left-associative)
					for len(dijkstraSliceSigns) > 0 {           // перебираем все элементы стека операторов, пока не попадётся открывающаяс скобка или стек не закончится
						top := dijkstraSliceSigns[len(dijkstraSliceSigns)-1]

						if top.valSign == "(" { // открвающая скобка нас не интересует, блок закончился, прекращаем
							break
						}

						prevPriority := dictSigns[top.valSign] // приоритет верхнего оператора стека операторов

						if (rightAssociative && priority < prevPriority) || (!rightAssociative && priority <= prevPriority) {
							// перекидываем в основной срез верхний оператор из стека операторов
							dijkstraSliceSigns = dijkstraSliceSigns[:len(dijkstraSliceSigns)-1]
							dijkstraSlice = append(dijkstraSlice, top)
							if prevPriority == 1 && lastPriority1 > 0 {
								lastPriority1--
							} else if lastPriority23 > 0 {
								lastPriority23--
							}
						} else {
							break
						}
					}
				} // end of for len(dijkstraSliceSigns) > 0
				// добавляем данный оператор в стек операторов
				// }
				dijkstraSliceSigns = append(dijkstraSliceSigns, dijkstraData{valSign: tokens[i], isSign: true})
			} // end of if token == "("
		}

	} // end of for _, token := range tokens

	// for _, token := range tokens {
	// 	v, err := strconv.Atoi(token)
	// 	if err == nil { // если это число
	// 		for lastMunus > 0 {
	// 			v *= -1
	// 			lastMunus--
	// 		}

	// 		// проверяем, есть ли знаки минус перед числом
	// 		for len(dijkstraSliceSigns) > 0 { // перебираем все элементы стека операторов, пока попадается знак минус или стек не закончится
	// 			top := dijkstraSliceSigns[len(dijkstraSliceSigns)-1]
	// 			if top.valSign != "-" { // не встретился знак минус, прекращаем
	// 				break
	// 			}
	// 			// предыдущий
	// 			v *= -1
	// 			prevPriority := dictSigns[top.valSign] // приоритет верхнего оператора стека операторов

	// 			if (rightAssociative && priority < prevPriority) || (!rightAssociative && priority <= prevPriority) {
	// 				// перекидываем в основной срез верхний оператор из стека операторов
	// 				dijkstraSliceSigns = dijkstraSliceSigns[:len(dijkstraSliceSigns)-1]
	// 				dijkstraSlice = append(dijkstraSlice, top)
	// 			} else {
	// 				break
	// 			}
	// 		}

	// 		dijkstraSlice = append(dijkstraSlice, dijkstraData{val: v})
	// 	} else { // это не число, проверяем дальше
	// 		if token == "(" { // открывающую скобку переносим в спомогательный срез (стек) операторов
	// 			dijkstraSliceSigns = append(dijkstraSliceSigns, dijkstraData{valSign: token, isSign: true})
	// 		} else if token == ")" {
	// 			// если попалась закрывающая скобка - переносим из стека все операторы, пока не попадётся открывающаяся скобка
	// 			found := false
	// 			for len(dijkstraSliceSigns) > 0 {
	// 				newMember := dijkstraSliceSigns[len(dijkstraSliceSigns)-1]
	// 				dijkstraSliceSigns = dijkstraSliceSigns[:len(dijkstraSliceSigns)-1]
	// 				if newMember.valSign == "(" {
	// 					found = true
	// 					break
	// 				} else {
	// 					dijkstraSlice = append(dijkstraSlice, newMember)
	// 				}
	// 			}
	// 			if !found {
	// 				// если не было открывающейся скобки - ошибка
	// 				return []dijkstraData{}, errors.New("error: expression is not valid, mismatched parentheses found")
	// 			}
	// 		} else {
	// 			// в остальных случаях проверяем на знак
	// 			priority, ok := dictSigns[token] // это приоритет данного оператора
	// 			if !ok {
	// 				return []dijkstraData{}, fmt.Errorf("error: expression is not valid, unknown operator: %v", token)
	// 			}

	// 			rightAssociative := tmpDictSigns[token] // существует?
	// 			for len(dijkstraSliceSigns) > 0 {       // перебираем все элементы стека операторов, пока не попадётся открывающаяс скобка или стек не закончится
	// 				top := dijkstraSliceSigns[len(dijkstraSliceSigns)-1]

	// 				if top.valSign == "(" { // открвающая скобка нас не интересует, блок закончился, прекращаем
	// 					break
	// 				}

	// 				prevPriority := dictSigns[top.valSign] // приоритет верхнего оператора стека операторов

	// 				if (rightAssociative && priority < prevPriority) || (!rightAssociative && priority <= prevPriority) {
	// 					// перекидываем в основной срез верхний оператор из стека операторов
	// 					dijkstraSliceSigns = dijkstraSliceSigns[:len(dijkstraSliceSigns)-1]
	// 					dijkstraSlice = append(dijkstraSlice, top)
	// 				} else {
	// 					break
	// 				}
	// 			} // end of for len(dijkstraSliceSigns) > 0
	// 			// добавляем данный оператор в стек операторов
	// 			dijkstraSliceSigns = append(dijkstraSliceSigns, dijkstraData{valSign: token, isSign: true})
	// 		} // end of if token == "("
	// 	}

	// } // end of for _, token := range tokens

	// перекидываем все оставшиеся операторы по порядку из стека операторов в основной срез
	for len(dijkstraSliceSigns) > 0 {
		// pop
		newMember := dijkstraSliceSigns[len(dijkstraSliceSigns)-1]
		dijkstraSliceSigns = dijkstraSliceSigns[:len(dijkstraSliceSigns)-1]

		if newMember.valSign == "(" {
			// непарная скобка осталась - ошибка
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
	log.Println("в бд добавлена новая строка:", data)

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
	log.Println("в бд изменены данные:", nameArr[0]+"="+strconv.Itoa(newData))

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
	log.Println("в бд изменены данные:", nameArr[0]+"="+strconv.Itoa(newData))

	// обновляем данные о последнем выражении в файле последнего конфига
	savePresentConfigToFile()
	return nil
}
