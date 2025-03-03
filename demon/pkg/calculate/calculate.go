package calculate

import "fmt"

func Calc(num1, num2 float64, oper string) (float64, error) {
	switch oper {
	case "*":
		return num1 * num2, nil	
	case "/":
		return num1 / num2, nil
	case "+":
		return num1 + num2, nil
	case "-":
		return num1 - num2, nil
	default:
		return 0, fmt.Errorf("incorrect sign of a mathematical operation")
	}
}