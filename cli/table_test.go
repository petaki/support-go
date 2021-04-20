package cli

import (
	"testing"
)

func TestTableCreate(t *testing.T) {
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
