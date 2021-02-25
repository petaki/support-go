package forms

// Bag type.
type Bag map[string][]string

// Add function.
func (b Bag) Add(field, message string) {
	b[field] = append(b[field], message)
}
