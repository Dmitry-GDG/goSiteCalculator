package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// заполнение текущего конфига новыми данными
func fillConfig(str string) {
	val := strings.Split(trimSpacesAndTabs(str), " ")
	muConfigs.Lock()
	defer muConfigs.Unlock()
	if val[0] == "port" {
		config_main.port = val[1]
	} else if val[0] == "checkUnicTask" {
		if val[1] == "1" {
			config_main.checkUnicTask = 1
			log.Println("в конфиге задано условие на проверку уникальности входных данных. Производим очистку БД.")
			clearDbFunc()
		} else if val[1] == "0" {
			config_main.checkUnicTask = 0
		}
	} else if val[0] == "conf" {
		config_main.conf = val[1]
	} else if val[0] == "fileExpressions" {
		config_main.fileExpressions = val[1]
	} else if val[0] == "lastExpression" {
		config_main.lastExpression = val[1]
	} else if val[0] == "oetDivide" {
		res, err := strconv.Atoi(val[1])
		if err == nil {
			if res < 0 {
				res = 0
			}
			config_main.oetDivide = res
		}
	} else if val[0] == "oetMinus" {
		res, err := strconv.Atoi(val[1])
		if err == nil {
			if res < 0 {
				res = 0
			}
			config_main.oetMinus = res
		}
	} else if val[0] == "oetMultiply" {
		res, err := strconv.Atoi(val[1])
		if err == nil {
			if res < 0 {
				res = 0
			}
			config_main.oetMultiply = res
		}
	} else if val[0] == "oetPlus" {
		res, err := strconv.Atoi(val[1])
		if err == nil {
			if res < 0 {
				res = 0
			}
			config_main.oetPlus = res
		}
	} else if val[0] == "oetPower" {
		res, err := strconv.Atoi(val[1])
		if err == nil {
			if res < 0 {
				res = 0
			}
			config_main.oetPower = res
		}
	} else if val[0] == "nextNumber" {
		res, err := strconv.Atoi(val[1])
		if err == nil {
			if res < 0 {
				res = 0
			}
			config_main.nextNumber = res
		}
	} else if val[0] == "qtyServers" {
		res, err := strconv.Atoi(val[1])
		if err == nil {
			if res < 1 {
				res = 1 // меньше 1 сервера нет смыслав устанавливать
			}
			config_main.qtyServers = res
		}
	}

}

// Сохранение текущего конфига в файл предыдущего конфига
func savePresentConfigToFile() {
	muConfigs.Lock()
	defer muConfigs.Unlock()

	configStr := ""

	configStr += "port " + config_main.port + "\n"
	configStr += "checkUnicTask " + strconv.Itoa(config_main.checkUnicTask) + "\n"
	configStr += "conf " + config_main.conf + "\n"
	configStr += "fileExpressions " + config_main.fileExpressions + "\n"
	configStr += "lastExpression " + config_main.lastExpression + "\n"
	configStr += "oetDivide " + strconv.Itoa(config_main.oetDivide) + "\n"
	configStr += "oetMinus " + strconv.Itoa(config_main.oetMinus) + "\n"
	configStr += "oetMultiply " + strconv.Itoa(config_main.oetMultiply) + "\n"
	configStr += "oetPlus " + strconv.Itoa(config_main.oetPlus) + "\n"
	configStr += "oetPower " + strconv.Itoa(config_main.oetPower) + "\n"
	configStr += "qtyServers " + strconv.Itoa(config_main.qtyServers) + "\n"
	configStr += "nextNumber " + strconv.Itoa(config_main.nextNumber)

	destFile, err := os.Create(config_main.conf)
	if err != nil {
		panic(err)
	}
	defer destFile.Close()

	_, err = io.WriteString(destFile, configStr)
	if err != nil {
		panic(err)
	}

}

// очистка файлов БД
func cleanPrevData() {
	muExpressions.Lock()
	defer muExpressions.Unlock()
	destFile, err := os.Create(config_main.fileExpressions)
	if err != nil {
		panic(err)
	}
	defer destFile.Close()

	_, err = io.WriteString(destFile, "")
	if err != nil {
		panic(err)
	}
}

// парсинг файла конфига
func parse(file string) {
	srcFile, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()
	fileScanner := bufio.NewScanner(srcFile)
	for fileScanner.Scan() {
		fillConfig(fileScanner.Text())
	}
}

// запрос и получение нового конфига
func getConfig() {
	var file string
	fmt.Print("Если требуется специальный конфиг, укажите название файла с конфигом; укажите 1, если желаете использовать конфиг последней загрузки: ")
	fmt.Scanln(&file)
	if file == "1" {
		parse(config_main.conf)
	} else {
		if file != "" && file != "0" {
			parse(file)
		}
		savePresentConfigToFile()
		cleanPrevData()
	}
}
