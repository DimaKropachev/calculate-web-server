package calculate

import "testing"

type TestCase struct {
	name       string
	expression string
	err        bool
}

func TestCheckExpression(t *testing.T) {
	cases := []TestCase{
		{
			name:       "simple-1",
			expression: "1+1",
			err:        false,
		},
		{
			name:       "simple-2",
			expression: "3-6",
			err:        false,
		},
		{
			name:       "priority-1",
			expression: "2 + 2*2",
			err:        false,
		},
		{
			name:       "priority-2",
			expression: "(2+2) * 2",
			err:        false,
		},
		{
			name:       "divizion",
			expression: "1/2",
			err:        false,
		},
		{
			name:       "brackets-1",
			expression: "2*(3+5/10)",
			err:        false,
		},
		{
			name:       "combine-1",
			expression: "23*11-37",
			err:        false,
		},
		{
			name:       "brackets-2",
			expression: "(3-6*5)-(4/2+22)",
			err:        false,
		},
		{
			name:       "combine-2",
			expression: "2- 44 *5",
			err:        false,
		},
		{
			name:       "combine-3",
			expression: "45/9*23",
			err:        false,
		},
		{
			name:       "empty",
			expression: "",
			err:        true,
		},
		{
			name:       "first operation",
			expression: "-13+3",
			err:        true,
		},
		{
			name:       "last operation",
			expression: "23-4+",
			err:        true,
		},
		{
			name:       "empty brackets",
			expression: "2*()-3",
			err:        true,
		},
		{
			name:       "invalid token",
			expression: "a +3",
			err:        true,
		},
		{
			name:       "operations brackets",
			expression: "2(3-4+)-3",
			err:        true,
		},
		{
			name:       "invalid brackets",
			expression: "1)*34",
			err:        true,
		},
		{
			name:       "2 operations",
			expression: "54/2++5",
			err:        true,
		},
		{
			name:       "2 numbers",
			expression: "2 2-345",
			err:        true,
		},
		{
			name:       "combine",
			expression: "+23-(34/)-2-",
			err:        true,
		},
	}

	for ind, tc := range cases {
		err := CheckExpression(tc.expression)
		if err == nil && tc.err {
			t.Fatalf("TestCase[%d] %s: expected error", ind, tc.name)
		}
		if err != nil && !tc.err {
			t.Fatalf("TestCase[%d] %s: unexpected error", ind, tc.name)
		}
	}

}
