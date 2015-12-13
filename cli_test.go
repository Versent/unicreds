package unicreds

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {

	tt := []struct {
		tableFormat int
		output      string
		headers     []string
		rows        [][]string
	}{
		{
			tableFormat: TableFormatTerm,
			output: `+------------+-----------+
|    NAME    |  VERSION  |
+------------+-----------+
| testlogin1 | testpass1 |
| testlogin2 | testpass2 |
+------------+-----------+
`,
			headers: []string{"Name", "Version"},
			rows:    [][]string{{"testlogin1", "testpass1"}, {"testlogin2", "testpass2"}},
		},
		{
			tableFormat: TableFormatCSV,
			output:      "Name,Version\ntestlogin1,testpass1\ntestlogin2,testpass2\n",
			headers:     []string{"Name", "Version"},
			rows:        [][]string{{"testlogin1", "testpass1"}, {"testlogin2", "testpass2"}},
		},
	}

	for _, tv := range tt {
		var b bytes.Buffer

		table := NewTable(&b)
		table.SetHeaders(tv.headers)
		table.SetFormat(tv.tableFormat)
		table.BulkWrite(tv.rows)
		table.Render()

		assert.Equal(t, tv.output, b.String())
	}

}
