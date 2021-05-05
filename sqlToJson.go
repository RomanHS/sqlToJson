package main

import (
	"database/sql"
	"encoding/json"
	"strings"
)

func sqlToJson(res *sql.Rows) ([]byte, error) {
	// создаем массив названий столбцов таблицы
	arrNamesColumns, _ := res.Columns()

	// получаем количество столбцов
	kolColumns := len(arrNamesColumns)

	// создаем отображения которое по ключу (названию столбца) будет хранить срез всех записей данного столбца
	resMap := make(map[string][]interface{}, kolColumns)

	// создаем срез который будет хранить все поля текущий строки таблицы в массиве байт
	tempLineByte := make([][]byte, kolColumns)

	// создаем срез который будет хранить ссылки на поля предыдущего среза
	// это нужно для метода Scan который принимает interface по ссылке таким образом записывая в него данные
	pTempLineByte := make([]interface{}, kolColumns)
	for i := 0; i < kolColumns; i++ {
		pTempLineByte[i] = &tempLineByte[i]
	}

	// создаем срез который будет хранить все поля текущий строки в соответствующим типе данных
	tempLine := make([]interface{}, kolColumns)

	// перебираем все строки тоблицы
	for res.Next() {
		// метод Scan записывает в срез tempLineByte текущую строку через срез pTempLineByte который хранит ссылки на поля tempLineByte
		err := res.Scan(pTempLineByte...)
		if err != nil {
			return nil, err
		}

		// перебираем все поля текущей строки
		for i := 0; i < kolColumns; i++ {
			// преобразовываем байтовое значение в значение соответствующего формата
			json.Unmarshal(tempLineByte[i], &tempLine[i])

			// если записалось значение nil то в поле была строка
			if tempLine[i] == nil {
				// добавляем значения в соответствующий ключу срез отображения resMap преобразуя его в строку и убираем переносы
				resMap[arrNamesColumns[i]] = append(resMap[arrNamesColumns[i]], strings.Trim(string(tempLineByte[i]), "\n"))
			} else {
				// иначе добавлять в срез соответствующие значения tempLine
				resMap[arrNamesColumns[i]] = append(resMap[arrNamesColumns[i]], tempLine[i])
			}
		}
	}

	// преобразовываем отображения в json и возвращаем в массиве байт
	return json.Marshal(resMap)

	// p.s. чтобы получить привычный json нужно просто преобразовать этот массив байт в строку
}
