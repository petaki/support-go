package cli

import (
	"testing"
	"unicode/utf8"
)

func TestTable_Create(t *testing.T) {
	tests := []struct {
		name     string
		expected []string
		got      []string
	}{
		{
			"Single Row Table",
			[]string{
				"+------------------+",
				"| Code | Message   |",
				"+------------------+",
				"| 404  | Not Found |",
				"+------------------+",
			},
			(&Table{
				Headers: []string{
					"Code",
					"Message",
				},
				Rows: [][]string{
					{
						"404",
						"Not Found",
					},
				},
			}).Create(),
		},
		{
			"Multi-Row Table",
			[]string{
				"+----------------------------------------+",
				"| Code | Message                         |",
				"+----------------------------------------+",
				"| 400  | Bad Request                     |",
				"| 401  | Unauthorized                    |",
				"| 404  | Not Found                       |",
				"| 405  | Method Not Allowed              |",
				"| 431  | Request Header Fields Too Large |",
				"+----------------------------------------+",
			},
			(&Table{
				Headers: []string{
					"Code",
					"Message",
				},
				Rows: [][]string{
					{
						"400",
						"Bad Request",
					},
					{
						"401",
						"Unauthorized",
					},
					{
						"404",
						"Not Found",
					},
					{
						"405",
						"Method Not Allowed",
					},
					{
						"431",
						"Request Header Fields Too Large",
					},
				},
			}).Create(),
		},
		{
			"Empty Table",
			[]string{
				"+----------------+",
				"| Code | Message |",
				"+----------------+",
				"| No items       |",
				"+----------------+",
			},
			(&Table{
				Headers: []string{
					"Code",
					"Message",
				},
			}).Create(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.expected) != len(tt.got) {
				t.Errorf("expected: %v, got: %v", tt.expected, tt.got)
			} else {
				for index, tr := range tt.expected {
					if tr != tt.got[index] {
						t.Errorf("expected: %v, got: %v", tr, tt.got[index])
					}
				}
			}
		})
	}
}

func TestTable_CreateWithWideRowAndUnicode(t *testing.T) {
	table := &Table{
		Headers: []string{"Name"},
		Rows:    [][]string{{"Árvíztűrő", "extra"}},
	}

	rows := table.Create()

	if utf8.RuneCountInString(rows[0]) != utf8.RuneCountInString(rows[3]) {
		t.Errorf("expected aligned rows, got: %v and %v", rows[0], rows[3])
	}
}

func TestTable_CreateNarrow(t *testing.T) {
	for _, table := range []*Table{
		{Headers: []string{"ID"}},
		{},
		{Headers: []string{"A", "B"}, Rows: [][]string{{"x"}}},
	} {
		rows := table.Create()

		for _, row := range rows {
			if utf8.RuneCountInString(row) != utf8.RuneCountInString(rows[0]) {
				t.Errorf("expected aligned rows, got: %q", rows)
			}
		}
	}
}
