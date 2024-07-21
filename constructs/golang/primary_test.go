package golang

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLiteral_ToCode(t *testing.T) {
	// Test case: int literal
	literal := &Literal{
		Type:  "int",
		Value: "123",
	}
	expected := "123"
	if result := literal.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	arrayIndexLiteral := &Literal{
		Value: "c",
		Indexes: &Literal{
			Value: 1,
		},
	}
	expected = "c[1]"
	if result := arrayIndexLiteral.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	arrayIndexLiteral2 := &Literal{
		Value:   "c",
		Indexes: 2,
	}
	expected = `c[2]`
	if result := arrayIndexLiteral2.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	mapIndexLiteral := &Literal{
		Value: "c",
		Indexes: &Literal{
			Value: "a",
		},
	}
	expected = `c["a"]`
	if result := mapIndexLiteral.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	mapIndexLiteral2 := &Literal{
		Value:   "c",
		Indexes: "a",
	}
	expected = `c["a"]`
	if result := mapIndexLiteral2.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	nestedMapLiteral := &Literal{
		Value: "c",
		Indexes: &Literal{
			Value: "a",
			Indexes: &Literal{
				Value: "b",
			},
		},
	}

	expected = `c[a["b"]]`
	if result := nestedMapLiteral.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	nestedArrayLiteral := &Literal{
		Value: "c",
		Indexes: &Literal{
			Value: "a",
			Indexes: &Literal{
				Value: 2,
			},
		},
	}

	expected = `c[a[2]]`
	if result := nestedArrayLiteral.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	array3DIndexLiteral := &Literal{
		Value: "c",
		Indexes: []*Literal{
			{
				Value: 1,
			},
			{
				Value: 2,
			},
			{
				Value: 3,
			},
		},
	}

	expected = `c[1][2][3]`
	if result := array3DIndexLiteral.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	map3DIndexLiteral := &Literal{
		Value: "c",
		Indexes: []*Literal{
			{
				Value: "a",
			},
			{
				Value: "b",
			},
			{
				Value: "c",
			},
		},
	}

	expected = `c["a"]["b"]["c"]`
	if result := map3DIndexLiteral.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	mapIndexSimple := &Literal{
		Value:   "c",
		Indexes: "a",
	}
	expected = `c["a"]`
	if result := mapIndexSimple.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}
}

func TestFunctionCall_ToCode(t *testing.T) {
	// Test case: Simple function call
	funcCall := &FunctionCall{
		Receiver: "foo",
		Function: "FnName",
		Args:     []string{"a", "b"},
		Output:   []interface{}{&Literal{Value: "c", Indexes: "one"}, "err"},
	}
	expected := `c["one"], err = foo.FnName(a, b)`
	if result := funcCall.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}
}

func TestCodeElement_ToCode(t *testing.T) {
	// Test case: Arithmetic operation
	ce := &CodeElement{
		Add: []string{"a", "b"},
	}
	expected := "(a + b)"
	if result := ce.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	// Test case: Logical operation
	ce = &CodeElement{
		And: []string{"true", "false"},
	}
	expected = "(true && false)"
	if result := ce.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	// Test case: Compare operation
	ce = &CodeElement{
		Equal: []string{"a", "b"},
	}
	expected = "(a == b)"
	if result := ce.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	// Test case: Bitwise operation
	ce = &CodeElement{
		BitwiseAnd: []string{"a", "b"},
	}
	expected = "(a & b)"
	if result := ce.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	// Test case: Assignment
	ce = &CodeElement{
		Assign: &Assignment{
			Left:  "a",
			Right: "b",
		},
	}
	expected = "a = b"
	if result := ce.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	// Test case: New Assignment
	ce = &CodeElement{
		NewAssign: &NewAssignment{
			Left:  []string{"a", "b"},
			Right: []string{"1", "2"},
		},
	}
	expected = "a, b := 1, 2"
	if result := ce.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	// Test case: If statement
	ce = &CodeElement{
		If: &IfElement{
			Condition: &CodeElement{GreaterThan: []string{"a", "b"}},
			Then: []*CodeElement{{Assign: &Assignment{
				Left:  "a",
				Right: &CodeElement{Add: []string{"a", "b"}},
			}}},
		},
	}
	expected = "if (a > b) {\n\ta = (a + b)\n}"
	if result := ce.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}

	// Test case: For statement
	ce = &CodeElement{
		RepeatCond: &RepeatByCondition{
			Condition: &CodeElement{GreaterThan: []string{"a", "b"}},
			Body: []*CodeElement{{Assign: &Assignment{
				Left:  "a",
				Right: &CodeElement{Add: []string{"a", "b"}},
			}}},
		},
	}
	expected = "for (a > b) {\n\ta = (a + b)\n}"
	if result := ce.ToCode(); result != expected {
		t.Errorf("ToCode() = %v, want %v", result, expected)
	}
}

func TestYamlToCode(t *testing.T) {
	// Test case: Arithmetic operation
	addYaml := `left: a
right: b
out: c`

	add := &Add{}
	yaml.Unmarshal([]byte(addYaml), add)
	t.Log(add)
	expected := "c = (a + b)"
	if result := add.ToCode(); result != expected {
		t.Errorf("AddToCode() = %v, want %v", result, expected)
	}

	expectedFib := `a, b, c := 0, 1, 0
for i := 0; (i < 10); i++ {
	c = (a + b)
	a, b = b, c
}
return c`

	fibYaml := `
steps:
- new_assign:
    left: [a, b, c]
    right: [0, 1, 0]
- repeat:
    init: 
        - {new_assign: {left: i, right: 0}}
    cond:  {lt: [i, 10]}
    step: 
        - post_inc: i
    body:
        - add:
            out: c
            left: a
            right: b
        - assign:
            left: [a, b]
            right: [b, c]
- return: c`

	fibCE := &CodeElement{}
	err := yaml.Unmarshal([]byte(fibYaml), fibCE)
	if err != nil {
		t.Error(err)
	}
	t.Log(fibCE.Steps[0].Assign)
	t.Log(fibCE.Steps[1].RepeatLoop.Init, fibCE.Steps[1].RepeatLoop.Condition,
		fibCE.Steps[1].RepeatLoop.Step, fibCE.Steps[1].RepeatLoop.Body)
	t.Log(fibCE.Steps[2].Return)
	if result := fibCE.ToCode(); result != expectedFib {
		t.Errorf("FibCode() = %v, want %v", result, expectedFib)
	}

	// Test case: Matrix multiplication
	// expected input is m * n matrix, n * p matrix
	// expected output is m * p matrix
	expectedMatMul := `c := make([][]int, len(a))
for i := 0; (i < len(a)); i++ {
	c[i] = make([]int, len(b[0]))
	for j := 0; (j < len(b[0])); j++ {
		c[i][j] = 0
		for k := 0; (k < len(b)); k++ {
			c[i][j] = (c[i][j] + (a[i][k] * b[k][j]))
		}
	}
}
return c`

	matMulYaml := `
steps:
- new_assign:
    left: c
    right: make([][]int, len(a))
- repeat:
    init: 
    - {new_assign: {left: i, right: 0}}
    cond:  {lt: [i, len(a)]}
    step: 
    - post_inc: i
    body:
    - assign:
        left: c[i]
        right: make([]int, len(b[0]))
    - repeat:
        init: 
        - {new_assign: {left: j, right: 0}}
        cond:  {lt: [j, "len(b[0])"]}
        step:
        - post_inc: j
        body:
        - assign:
            left: c[i][j]
            right: 0
        - repeat:
            init: 
            - {new_assign: {left: k, right: 0}}
            cond:  {lt: [k, len(b)]}
            step:
            - post_inc: k
            body:
            - add:
                out: "c[i][j]"
                left: "c[i][j]"
                right: 
                    mul: ["a[i][k]",  "b[k][j]"]

- return: c`

	matMulCE := &CodeElement{}
	err = yaml.Unmarshal([]byte(matMulYaml), matMulCE)
	if err != nil {
		t.Error(err)
	}
	t.Logf("matMulCE: %v", matMulCE)
	if result := matMulCE.ToCode(); result != expectedMatMul {
		t.Errorf("MatMulCode() = %v, want %v", result, expectedMatMul)
	}

	// Test case: merge sort algorithm
	expectedMergeSort := `result := make([]int, high - low + 1)
mid := (low + high) / 2
left := mergeSort(arr, low, mid)
right := mergeSort(arr, mid + 1, high)
i, j, k := 0, 0, low
for ((i < len(left)) && (j < len(right))) {
	if (left[i] < right[j]) {
		result[k] = left[i]
		i++
	} else {
		result[k] = right[j]
		j++
	}
	k++
}
for (i < len(left)) {
	result[k] = left[i]
	i++
	k++
}
for (j < len(right)) {
	result[k] = right[j]
	j++
	k++
}
return result`

	mergeSortYaml := `
steps:
- call:
    nout: result
    func: make
    args: ["[]int", high - low + 1]
- new_assign:
    left: mid
    right: (low + high) / 2
- call:
    nout: left
    func: mergeSort
    args: [arr, low, mid]
- call:
    nout: right
    func: mergeSort
    args: [arr, mid + 1, high]
- new_assign: 
    left: [i, j, k]
    right: [0, 0, low]
- repeat_cond:
    cond: 
        and:
        - {lt: [i, len(left)]}
        - {lt: [j, len(right)]}
    body:
    - cases:
        - cond: {lt: ["left[i]", "right[j]"]}
          body:
          - {assign: {left: "result[k]", right: "left[i]"}}
          - {post_inc: i}
        - body:
          - {assign: {left: "result[k]", right: "right[j]"}}
          - {post_inc: j}
    - {post_inc: k}
- repeat_cond:
    cond: {lt: [i, len(left)]}
    body:
    - {assign: {left: "result[k]", right: "left[i]"}}
    - {post_inc: i}
    - {post_inc: k}
- repeat_cond:
    cond: {lt: [j, len(right)]}
    body:
    - {assign: {left: "result[k]", right: "right[j]"}}
    - {post_inc: j}
    - {post_inc: k}
- return: result`

	mergeSortCE := &CodeElement{}
	err = yaml.Unmarshal([]byte(mergeSortYaml), mergeSortCE)
	if err != nil {
		t.Error(err)
	}
	t.Logf("mergeSortCE: %v", mergeSortCE)
	if result := mergeSortCE.ToCode(); result != expectedMergeSort {
		t.Errorf("MergeSortCode() = %v, want %v", result, expectedMergeSort)
	}

	// Test case: DFS
	// assume we have `adj` contains map, it has indexes from 0 to n - 1
	// `a->b` means `a` has a edge to `b`, `b->a` means `b` has a edge to `a`
	dfsExpected := `result := make([]int, 0)
visited := make([]bool, n)
for i := 0; i < n; i++ {
	visited[i] = false
}
for i := 0; i < n; i++ {
	if !visited[i] {
		stack := make([]int, 0)
		stackIndices := make([]int, 0)
		visited[i] = true
		result = append(result, i)
		stack = append(stack, i)
		stackIndices = append(stackIndices, 0)
		for (len(stack) > 0) {
			v := stack[len(stack) - 1]
			s := stackIndices[len(stackIndices) - 1]
			noNewVisits := true
			for j := s; j < len(adj[v]); j++ {
				w := adj[v][j]
				stackIndices[len(stackIndices) - 1] = j + 1
				if !visited[w] {
					visited[w] = true
					result = append(result, w)
					stack = append(stack, w)
					stackIndices = append(stackIndices, 0)
					noNewVisits = false
					break
				}
			}
			if noNewVisits {
				stack = stack[:len(stack) - 1]
				stackIndices = stackIndices[:len(stackIndices) - 1]
			}
		}
	}
}
return result`

	dfsYaml := `
steps:
- call:
    nout: result
    func: make
    args: ["[]int", 0]
- call:
    nout: visited
    func: make
    args: ["[]bool", n]
- repeat_n:
    iter: i
    limit: n
    body:
    - assign:
        left: "visited[i]"
        right: false
- repeat_n:
    iter: i
    limit: n
    body:
    - cases:
          - cond: {not: "visited[i]"}
            body:
            - call:
                nout: stack
                func: make
                args: ["[]int", 0]
            - call:
                nout: stackIndices
                func: make
                args: ["[]int", 0]
            - assign:
                left: "visited[i]"
                right: true
            - call:
                out: result
                func: append
                args: [result, i]
            - call:
                out: stack
                func: append
                args: [stack, i]
            - call:
                out: stackIndices
                func: append
                args: [stackIndices, 0]
            - repeat_cond:
                cond: {gt: [len(stack), 0]}
                body:
                - new_assign:
                    left: v
                    right: "stack[len(stack) - 1]"
                - new_assign:
                    left: s
                    right: "stackIndices[len(stackIndices) - 1]"
                - new_assign:
                    left: noNewVisits
                    right: true
                - repeat_n:
                    iter: j
                    start: s
                    limit: "len(adj[v])"
                    body:
                    - new_assign:
                        left: w
                        right: "adj[v][j]"
                    - assign:
                        left: "stackIndices[len(stackIndices) - 1]"
                        right: "j + 1"
                    - cases:
                      - cond: {not: "visited[w]"}
                        break: true
                        body:
                        - assign:
                            left: "visited[w]"
                            right: true
                        - call:
                            out: result
                            func: append
                            args: [result, w]
                        - call:
                            out: stack
                            func: append
                            args: [stack, w]
                        - call:
                            out: stackIndices
                            func: append
                            args: [stackIndices, 0]
                        - assign:
                            left: noNewVisits
                            right: false
                - if:
                    cond: noNewVisits
                    then:
                    - assign:
                        left: "stack"
                        right: "stack[:len(stack) - 1]"
                    - assign:
                        left: "stackIndices"
                        right: "stackIndices[:len(stackIndices) - 1]"
- return: result`

	dfsCodeElem := CodeElement{}
	err = yaml.Unmarshal([]byte(dfsYaml), &dfsCodeElem)
	if err != nil {
		t.Error(err)
	}
	t.Log(dfsCodeElem)

	if result := dfsCodeElem.ToCode(); result != dfsExpected {
		t.Errorf("dfsCodeElem.ToCode() = %v, want %v", result, dfsExpected)
	}

	maxFuncExpected := `if (a > b) {
	return a
} else {
	return b
}`
	maxFuncYaml := `
if:
  cond: {gt: [a, b]}
  then:
  - return: a
  else:
  - return: b`

	maxFuncCodeElem := CodeElement{}
	err = yaml.Unmarshal([]byte(maxFuncYaml), &maxFuncCodeElem)
	if err != nil {
		t.Error(err)
	}
	t.Log(maxFuncCodeElem)

	if result := maxFuncCodeElem.ToCode(); result != maxFuncExpected {
		t.Errorf("maxFuncCodeElem.ToCode() = %v, want %v", result, maxFuncExpected)
	}
}

func TestAddToCode2(t *testing.T) {
	// Test case: Logical operation
	addYaml := `left: a
right: b
out: c`

	add := &Add{}
	yaml.Unmarshal([]byte(addYaml), add)
	t.Log(add)
	expected := "c = (a + b)"
	if result := AddToCode(add); result != expected {
		t.Errorf("AddToCode() = %v, want %v", result, expected)
	}

	addStrings := []string{"a", "b"}
	expected = "(a + b)"
	if result := AddToCode(addStrings); result != expected {
		t.Errorf("AddToCode() = %v, want %v", result, expected)
	}

}

// TestSQLYamlToCode description of the Go function.
// Create sqlYaml which contains yaml for SQL read through Go
// and compare the result with the expected
// t *testing.T.
func TestSQLYamlToCode(t *testing.T) {
	expectedFindGoCode := `stmt := db.preparedCache["{{ .Name }}"]
values, err := {{ .Name }}ParseParams(reqeust)
if err != nil {
	return nil, err
}
rows, err := stmt.Query(values...)
if err != nil {
	return nil, err
}
defer rows.Close()
var results []{{.ModelName}}
for rows.Next() {
	var item {{.ModelName}}
	err := rows.Scan({{ .ScanAttributes }})
	if err != nil {
		return nil, err
	}
	results = append(results, item)
}
return results, nil`

	findSQLYaml := `
steps:
    - lookup:
        nout: stmt
        obj: db
        name: preparedCache
        key: {val: "{{ .Name }}"}
    - call:
        nout: [values, err]
        func: "{{ .Name }}ParseParams"
        args: [reqeust]
        err_returns: [nil, err]
    - call:
        nout: [rows, err]
        obj: stmt
        func: Query
        args: ["values..."]
        err_returns: [nil, err]
        clean: 
            obj: rows
            func: Close
    - var:
        name: results
        type: "[]{{.ModelName}}"
    - repeat_cond:
        cond: {call: {func: Next, obj: rows}}
        body:
        - var:
            name: item
            type: "{{.ModelName}}"
        - call:
            nout: err
            obj: rows
            func: Scan
            args: "{{ .ScanAttributes }}"
            err_returns: [nil, err]
        - call:
            out: results
            func: append
            args: [results, item]
    - return: results, nil`

	sqlCodeElem := CodeElement{}
	err := yaml.Unmarshal([]byte(findSQLYaml), &sqlCodeElem)
	if err != nil {
		t.Error(err)
	}
	t.Log(sqlCodeElem)

	if result := sqlCodeElem.ToCode(); result != expectedFindGoCode {
		t.Errorf("sqlCodeElem.ToCode() = %v, want %v", result, expectedFindGoCode)
	}

	expectedUpdateCode := `stmt := db.preparedCache["UpdateByIdQuery"]
values := UpdateByIdParseParams(request)
result, err := stmt.Exec(values...)
if err != nil {
	return 0, err
}
rowsAffected, err := result.RowsAffected()
return rowsAffected, err`

	updateSQLYaml := `
steps:
    - lookup:
        nout: stmt
        obj: db
        name: preparedCache
        key: {val: "UpdateByIdQuery"}
    - call:
        nout: values
        func: UpdateByIdParseParams
        args: request
    - call:
        nout: [result, err]
        obj: stmt
        func: Exec
        args: "values..."
        err_returns: [0, err]
    - call:
        nout: [rowsAffected, err]
        obj: result
        func: RowsAffected
    - return: [rowsAffected, err]`

	updateCodeElem := CodeElement{}
	err = yaml.Unmarshal([]byte(updateSQLYaml), &updateCodeElem)
	if err != nil {
		t.Error(err)
	}
	t.Log(updateCodeElem)

	if result := updateCodeElem.ToCode(); result != expectedUpdateCode {
		t.Errorf("updateCodeElem.ToCode() = %v, want %v", result, expectedUpdateCode)
	}

}
