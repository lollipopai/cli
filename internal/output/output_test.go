package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorizeValue_Null(t *testing.T) {
	// In non-TTY mode (test), colors are disabled
	// Just test the logic paths
	result := colorizeValue("null")
	assert.Contains(t, result, "null")
}

func TestColorizeValue_Bool(t *testing.T) {
	result := colorizeValue("true")
	assert.Contains(t, result, "true")

	result = colorizeValue("false")
	assert.Contains(t, result, "false")
}

func TestColorizeValue_String(t *testing.T) {
	result := colorizeValue(`"hello"`)
	assert.Contains(t, result, `"hello"`)
}

func TestColorizeValue_Number(t *testing.T) {
	result := colorizeValue("42")
	assert.Contains(t, result, "42")

	result = colorizeValue("3.14")
	assert.Contains(t, result, "3.14")
}

func TestColorizeValue_Brackets(t *testing.T) {
	assert.Equal(t, "{", colorizeValue("{"))
	assert.Equal(t, "}", colorizeValue("}"))
	assert.Equal(t, "[", colorizeValue("["))
	assert.Equal(t, "]", colorizeValue("]"))
	assert.Equal(t, "{}", colorizeValue("{}"))
	assert.Equal(t, "[]", colorizeValue("[]"))
}

func TestColorizeValue_Trailing(t *testing.T) {
	result := colorizeValue("42,")
	assert.Contains(t, result, "42")
	assert.True(t, result[len(result)-1] == ',')
}

func TestColorizeJSON_KeyValue(t *testing.T) {
	raw := `{
  "name": "test",
  "count": 5,
  "active": true,
  "data": null
}`
	result := colorizeJSON(raw)
	assert.Contains(t, result, "name")
	assert.Contains(t, result, "test")
	assert.Contains(t, result, "count")
	assert.Contains(t, result, "5")
}
