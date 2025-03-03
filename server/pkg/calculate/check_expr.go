package calculate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	Err           error = errors.New("invalid expression")
	ErrStupidCalc error = errors.New("calculator is too stupid for that")
)

// Функция возвращает разбитую на числа и математические операции строку в виде массива
func GetTokens(expression string) []string {
	result := []string{}

	var num string
	// Перебираем строку по символам в цикле
	for _, r := range expression {
		sym := string(r)
		if sym == " " {
			// Если дошли до пробела, то если в переменной num есть значения, то добавляем его в массив
			if num != "" {
				result = append(result, num)
				// После добавления обновляем переменную num
				num = ""
			}
			continue
		}
		// Проверяем является ли символ строки знаком математического выражения
		if sym == "+" || sym == "-" || sym == "*" || sym == "/" || sym == "(" || sym == ")" {
			// Если да, то проверяем есть ли значение в переменной num. Если есть, то добавляем его в массив
			if num != "" {
				result = append(result, num)
				// После добавления обновляем переменную num
				num = ""
			}
			// Добавляем символ операции в массив
			result = append(result, sym)
			continue
		} else {
			// Если символ не является знаком математической операции, то добавляем его к переменной num (если этот символ не пробел)
			if sym != " " {
				num += sym
			}
		}
	}
	// После прохода по всем символам проверяем есть ли значения в переменной num
	if num != "" {
		// Если есть, то добавляем их в массив
		result = append(result, num)
	}

	return result
}

// Функция проверяет правильность написания выражения
func CheckExpression(expression string) error {
	tokens := GetTokens(expression)

	var (
		countBrackets int
		trueBrackets  bool
	)

	// Если выражение, пустое сразу возвращаем ошибку
	if len(tokens) == 0 {
		return Err
	}

	// Считаем скобки в выражении
	for _, token := range tokens {
		if token == "(" {
			countBrackets++
		} else if token == ")" {
			countBrackets--
		}
		// Если нарушен порядок скобок, то ...
		if countBrackets < 0 {
			return Err
		}
	}

	// Если скобки расставлены правильно, то проверяем содержимое внутри скобок
	if countBrackets == 0 {
		trueBrackets = true
		err := CheckBrackets(tokens)
		if err != nil {
			return Err
		}
	} else {
		trueBrackets = false
		return Err
	}

	// Создаем счетчики подряд идущих математических операций и чисел
	var (
		countNum, countOper int
	)

	// Перебираем токены выражения
	for i, token := range tokens {
		// Если token стоит на 1 месте и является знаком, то добавляем в resultError ошибку об этом
		if i == 0 {
			if IsOperation(token) {
				return Err
			}
		}
		// Если token стоит на последнем месте и является знаком, то добавляем в resultError ошибку об этом
		if i == len(tokens)-1 {
			if IsOperation(token) {
				return Err
			}
		}
		if token == "(" {
			// Если token - открывающаяся скобка
			// Обнуляем счетчики подряд идущих чисел и математических операций
			countNum = 0
			countOper = 0

			// Если скобки в выражении расставлены правильно, то проверяем стоящие перед и после скобки token-ы
			// Если token не соответствует правилам, то добавляем в map с ошибками индекс этого токена
			if trueBrackets {
				if tokens[i+1] == ")" {
					continue
				}
				if !IsInteger(tokens[i+1]) {
					return Err
				}
				if i > 0 {
					if !IsOperation(tokens[i-1]) {
						return Err
					}
				}
			}
		} else if token == ")" {
			// Если token - закрывающаяся скобка
			// Обнуляем счетчики подряд идущих чисел и математических операций
			countNum = 0
			countOper = 0

			// Если скобки в выражении расставлены правильно, то проверяем стоящие перед и после скобки token-ы
			// Если token не соответствует правилам, то добавляем в map с ошибками индекс этого токена
			if trueBrackets {
				if tokens[i-1] == "(" {
					continue
				}
				if !IsInteger(tokens[i-1]) {
					return Err
				}
				if i < len(tokens)-2 {
					if !IsOperation(tokens[i+1]) {
						return Err
					}
				}
			}
		} else {
			// Если token - не скобка
			if IsOperation(token) {
				// Если token является математической операцией, то счетчик подряд идущих математический операций увеличиваем, а счетчик подряд идущих чисел обнуляем
				countOper++
				countNum = 0
			} else if IsInteger(token) {
				// Если token является числом, то счетчик подряд идущих чисел увеличиваем, а счетчик подряд идущих математических операций обнуляем
				countNum++
				countOper = 0
			} else {
				// Если token не является ни числом, ни математической операцией, то добавляем в map с ошибками, ошибку о неизвестном token-е и индекс данного token-а
				// Также обнуляем счетчики подряд идущих чисел и математических операций
				countNum = 0
				countOper = 0
				return Err
			}

			// Если счетчик подряд идущих чисел равен 2, то добавлем в map с ошибками, ошибку об этом и индекс первого подряд идущего token-а
			if countNum == 2 {
				return Err
			}
			// Если счетчик подряд идущих математических операций равен 2, то добавлем в map с ошибками, ошибку об этом и индекс первого подряд идущего token-а
			if countOper == 2 {
				return Err
			}
		}
	}

	if strings.Count(expression, "(") > 0 {
		return ErrStupidCalc
	}

	return nil
}

// Функция проверяет наличие содержимого в скобках
func CheckBrackets(tokens []string) error {
	brackets := [][]int{}
	emptyBrackets := []string{}
	// Создаем map-ы для индексов открывающихся и закрывающихся скобок
	mapFirst := make(map[int]int)
	mapSecond := make(map[int]int)
	// Переменная для нумерации скобок
	cb := 0
	// Перебираем token-ы выражения
	for i, token := range tokens {
		// Если token - открывающаяся скобка, то добавляем его индекс в соответствующую map-у под соответствующим номером (cb)
		if token == "(" {
			cb++
			if _, ok := mapFirst[cb]; !ok {
				mapFirst[cb] = i
			}
		} else if token == ")" {
			// Если token - закрывающаяся скобка, то добавляем его индекс в соответствующую map-у под соответствующим номером (cb)
			if _, ok := mapSecond[cb]; !ok {
				mapSecond[cb] = i
			} else {
				for {
					cb--
					if _, ok := mapSecond[cb]; !ok {
						mapSecond[cb] = i
						break
					}
				}
			}
		}
	}

	// Перебираем map-ы и добавляем индексы token-ов с одинаковыми номерами в массив brackets
	for count := range mapFirst {
		brackets = append(brackets, []int{mapFirst[count], mapSecond[count]})
	}

	// Проверяем скобки на пустоту, и если скобка пустая то добавляем её в список emptyBrackets
	for _, bracket := range brackets {
		if len(tokens[bracket[0]+1:bracket[1]]) == 0 {
			emptyBrackets = append(emptyBrackets, fmt.Sprintf("[%d:%d]", bracket[0], bracket[1]))
		}
	}
	// Добавляем в resultError ошибки о пустых скобках
	if len(emptyBrackets) > 0 {
		return Err
	}

	return nil
}

// Функция проверяет является ли token числом
func IsInteger(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

// Функция проверяет является ли token допустимой математической операцией
func IsOperation(token string) bool {
	if token == "+" || token == "/" || token == "-" || token == "*" {
		return true
	}
	return false
}

// Функция проверяет является ли token скобкой
func IsBracket(token string) bool {
	if token == "(" || token == ")" {
		return true
	}
	return false
}
