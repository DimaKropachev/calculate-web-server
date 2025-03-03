package orchestrator

import (
	"fmt"
	"strings"

	"github.com/DimaKropachev/calculate-web-server/server/internal/models"
	"github.com/DimaKropachev/calculate-web-server/server/pkg/calculate"
)

func Split(expression string, idExpr string) []*models.Task {
	exprs := make(map[string]string)
	tasks := []*models.Task{}

	miniExprId := 0
	miniExpr := ""
	stackExpr := []string{}
	stackOper := []string{}

	for ind, sym := range expression {

		if string(sym) == "+" || string(sym) == "-" {
			stackExpr = append(stackExpr, miniExpr)
			miniExpr = ""

			stackOper = append(stackOper, string(sym))
			continue
		}

		miniExpr += string(sym)

		if ind == len(expression)-1 {
			if len(miniExpr) != 0 {
				stackExpr = append(stackExpr, miniExpr)
				miniExpr = ""
			}
		}
	}

	for _, expr := range stackExpr {
		miniExprId++
		res := CreateTasks(expr, fmt.Sprintf("%s.%d", idExpr, miniExprId), exprs)
		tasks = append(tasks, res...)
	}

	miniExprId++
	taskId := 0
	for {
		if len(stackExpr) == 1 && len(stackOper) == 0 {
			break
		}

		taskId++

		expr1 := stackExpr[0]
		expr2 := stackExpr[1]
		oper := stackOper[0]

		bigExpr := expr1 + oper + expr2

		t1, ok1 := exprs[expr1]
		t2, ok2 := exprs[expr2]
		if ok1 {
			expr1 = ""
		}
		if ok2 {
			expr2 = ""
		}

		task := CreateTask(expr1, expr2, oper, fmt.Sprintf("%s.%d", idExpr, miniExprId), t1, t2, taskId)

		stackExpr[1] = bigExpr
		stackExpr = stackExpr[1:]
		stackOper = stackOper[1:]

		if _, ok := exprs[bigExpr]; !ok {
			exprs[bigExpr] = task.Id
		}

		tasks = append(tasks, task)
	}

	// Разкомментировав строку ниже в терминал будет выводится таблица с задачами, на которые было разбито выражение
	// PrintTasks(tasks)

	return tasks
}

func CreateTask(arg1, arg2, oper, idExpr, t1, t2 string, counter int) *models.Task {
	counter++

	task := &models.Task{
		Id:       fmt.Sprintf("%s.%d", idExpr, counter),
		Arg1:     arg1,
		Arg2:     arg2,
		Oper:     oper,
		Arg1Task: t1,
		Arg2Task: t2,
	}

	return task
}

func CreateTasks(expr string, id string, exprs map[string]string) []*models.Task {
	tasks := []*models.Task{}

	tokens := calculate.GetTokens(expr)

	stackNum := []string{}
	stackOper := []string{}

	for _, token := range tokens {
		if calculate.IsOperation(token) {
			stackOper = append(stackOper, token)
		} else if calculate.IsInteger(token) {
			stackNum = append(stackNum, token)
		}
	}

	tasksId := 0

	for {
		if len(stackNum) == 1 && len(stackOper) == 0 {
			break
		}

		tasksId++

		num1, num2 := stackNum[0], stackNum[1]
		oper := stackOper[0]

		expression := num1 + oper + num2

		t1, ok1 := exprs[num1]
		t2, ok2 := exprs[num2]
		if ok1 {
			num1 = ""
		}
		if ok2 {
			num2 = ""
		}

		task := CreateTask(num1, num2, oper, id, t1, t2, tasksId)

		stackNum[1] = expression
		stackNum = stackNum[1:]
		stackOper = stackOper[1:]

		if _, ok := exprs[expression]; !ok {
			exprs[expression] = task.Id
		}

		tasks = append(tasks, task)
	}

	return tasks
}

func PrintTasks(tasks []*models.Task) {
	fmt.Println("+------+------+------+------+------+------+")
	fmt.Println("|  id  | arg1 | arg2 |  t1  |  t2  | oper |")
	fmt.Println("+------+------+------+------+------+------+")
	for _, task := range tasks {
		arg1 := fmt.Sprintf("%v", task.Arg1) + strings.Repeat(" ", 6-len(fmt.Sprintf("%v", task.Arg1)))
		arg2 := fmt.Sprintf("%v", task.Arg2) + strings.Repeat(" ", 6-len(fmt.Sprintf("%v", task.Arg2)))
		t1 := task.Arg1Task + strings.Repeat(" ", 6-len(task.Arg1Task))
		t2 := task.Arg2Task + strings.Repeat(" ", 6-len(task.Arg2Task))

		fmt.Printf("|%s |%s|%s|%s|%s|  %s   |\n",
			task.Id,
			arg1,
			arg2,
			t1,
			t2,
			task.Oper,
		)
		fmt.Println("+------+------+------+------+------+------+")
	}
}
