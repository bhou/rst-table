package rsttable

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableWithoutGroup(t *testing.T) {
	table := NewTable()

	table.AddCol("name", DefaultColRender)
	table.AddCol("age", DefaultColRender)

	table.AddRow(map[string]any{"name": "John", "age": 20})
	table.AddRow(map[string]any{"name": "Jane", "age": 30})
	table.AddRow(map[string]any{"name": "Joe", "age": 40})
	table.AddRow(map[string]any{"name": "Jack", "age": 30})

	t1 := table.GenerateRstTable(nil)
	fmt.Println(t1)

	assert.Equal(t, strings.TrimSpace(t1), strings.TrimSpace(`
+------+------+
| name | age  |
+======+======+
| John | 20   |
+------+------+
| Jane | 30   |
+------+------+
| Joe  | 40   |
+------+------+
| Jack | 30   |
+------+------+`))
}

func TestTableWithGroup(t *testing.T) {
	table := NewTable()

	table.AddCol("name", DefaultColRender)
	table.AddCol("age", DefaultColRender)

	table.AddRow(map[string]any{"name": "John", "age": 20})
	table.AddRow(map[string]any{"name": "Jane", "age": 30})
	table.AddRow(map[string]any{"name": "Joe", "age": 40})
	table.AddRow(map[string]any{"name": "Jack", "age": 30})
	table.AddRow(map[string]any{"name": "Jill", "age": 30})

	t2 := table.GenerateRstTable([]string{"age"})
	fmt.Println(t2)

	assert.Equal(t, strings.TrimSpace(t2), strings.TrimSpace(`
+------+------+
| age  | name |
+======+======+
| 20   | John |
+------+------+
| 30   | Jane |
+      +------+
|      | Jack |
+      +------+
|      | Jill |
+------+------+
| 40   | Joe  |
+------+------+
`))
}

func TestTableWithMultipleColGroup(t *testing.T) {
	table := NewTable()

	table.AddCol("name", DefaultColRender)
	table.AddCol("age", DefaultColRender)
	table.AddCol("profession", DefaultColRender)

	table.AddRow(map[string]any{"name": "John", "age": 20, "profession": "student"})
	table.AddRow(map[string]any{"name": "Jane", "age": 30, "profession": "student"})
	table.AddRow(map[string]any{"name": "Joe", "age": 40, "profession": "teacher"})
	table.AddRow(map[string]any{"name": "Jack", "age": 30, "profession": "teacher"})
	table.AddRow(map[string]any{"name": "Jill", "age": 30, "profession": "teacher"})

	t2 := table.GenerateRstTable([]string{"profession", "age"})

	assert.Equal(t, strings.TrimSpace(t2), strings.TrimSpace(`
+------------+------------+------------+
| profession | age        | name       |
+============+============+============+
| student    | 20         | John       |
+            +------------+------------+
|            | 30         | Jane       |
+------------+------------+------------+
| teacher    | 30         | Jack       |
+            +            +------------+
|            |            | Jill       |
+            +------------+------------+
|            | 40         | Joe        |
+------------+------------+------------+
`))
	fmt.Println(t2)
}

func TestTableWithThreeColGroup(t *testing.T) {
	table := NewTable()

	table.AddCol("name", DefaultColRender)
	table.AddCol("age", DefaultColRender)
	table.AddCol("profession", DefaultColRender)
	table.AddCol("salary", DefaultColRender)

	table.AddRow(map[string]any{"name": "John", "age": 20, "profession": "student", "salary": 1000})
	table.AddRow(map[string]any{"name": "Jane", "age": 30, "profession": "student", "salary": 1000})
	table.AddRow(map[string]any{"name": "Joe", "age": 40, "profession": "teacher", "salary": 3000})
	table.AddRow(map[string]any{"name": "Jack", "age": 30, "profession": "teacher", "salary": 4000})
	table.AddRow(map[string]any{"name": "Jill", "age": 30, "profession": "teacher", "salary": 4000})

	t1 := table.GenerateRstTable([]string{"profession", "age", "salary"})
	fmt.Println(t1)
	assert.Equal(t, strings.TrimSpace(t1), strings.TrimSpace(`
+------------+------------+------------+------------+
| profession | age        | salary     | name       |
+============+============+============+============+
| student    | 20         | 1000       | John       |
+            +------------+------------+------------+
|            | 30         | 1000       | Jane       |
+------------+------------+------------+------------+
| teacher    | 30         | 4000       | Jack       |
+            +            +            +------------+
|            |            |            | Jill       |
+            +------------+------------+------------+
|            | 40         | 3000       | Joe        |
+------------+------------+------------+------------+
	`))

	t2 := table.GenerateRstTable([]string{"profession", "salary", "age"})
	fmt.Println(t2)
	assert.Equal(t, strings.TrimSpace(t2), strings.TrimSpace(`
+------------+------------+------------+------------+
| profession | salary     | age        | name       |
+============+============+============+============+
| student    | 1000       | 20         | John       |
+            +            +------------+------------+
|            |            | 30         | Jane       |
+------------+------------+------------+------------+
| teacher    | 3000       | 40         | Joe        |
+            +------------+------------+------------+
|            | 4000       | 30         | Jack       |
+            +            +            +------------+
|            |            |            | Jill       |
+------------+------------+------------+------------+
	`))
}

func TestTableWithGroupAndSort(t *testing.T) {
	table := NewTable()

	table.AddCol("name", DefaultColRender)
	table.AddCol("age", DefaultColRender)

	table.AddRow(map[string]any{"name": "John", "age": 20})
	table.AddRow(map[string]any{"name": "Jane", "age": 30})
	table.AddRow(map[string]any{"name": "Joe", "age": 40})
	table.AddRow(map[string]any{"name": "Jack", "age": 30})
	table.AddRow(map[string]any{"name": "Jill", "age": 30})

	t2 := table.GenerateRstTableWithCustomOrder([]string{"age"}, func(a, b Row) bool {
		if a.(map[string]any)["age"].(int) < b.(map[string]any)["age"].(int) {
			return false
		}
		return true
	})

	fmt.Println(t2)

	assert.Equal(t, strings.TrimSpace(t2), strings.TrimSpace(`
+------+------+
| age  | name |
+======+======+
| 40   | Joe  |
+------+------+
| 30   | Jill |
+      +------+
|      | Jack |
+      +------+
|      | Jane |
+------+------+
| 20   | John |
+------+------+
`))
}
