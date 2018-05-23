package gubrak

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsArray(t *testing.T) {
	assert.True(t, IsArray(
		[]string{"a", "b", "c", "d"},
	))
}

func TestIsArrayFail(t *testing.T) {
	assert.False(t, IsArray(
		make(map[string]interface{}),
	))
}

func TestIsBool(t *testing.T) {
	assert.True(t, IsBool(
		true,
	))
}

func TestIsBoolFail(t *testing.T) {
	assert.False(t, IsBool(
		"hello",
	))
}

func TestIsChannel(t *testing.T) {
	data := make(chan string)

	assert.True(t, IsChannel(
		data,
	))
}

func TestIsChannelFail(t *testing.T) {
	assert.False(t, IsChannel(
		"hello",
	))
}

func TestIsFloat32(t *testing.T) {
	assert.True(t, IsFloat(
		float32(24),
	))
}

func TestIsFloat64(t *testing.T) {
	assert.True(t, IsFloat(
		float64(24),
	))
}

func TestIsFloatFail(t *testing.T) {
	assert.False(t, IsFloat(
		"hello",
	))
}

func TestIsFunction(t *testing.T) {
	closure := func() string {
		return "hello"
	}

	assert.True(t, IsFunction(
		closure,
	))
}

func TestIsFunctionFail(t *testing.T) {
	assert.False(t, IsFunction(
		"hello",
	))
}

func TestIsInt(t *testing.T) {
	assert.True(t, IsInt(
		24,
	))
}

func TestIsIntFailString(t *testing.T) {
	assert.False(t, IsInt(
		"hello",
	))
}

func TestIsIntFailFloat32(t *testing.T) {
	assert.False(t, IsInt(
		float32(12),
	))
}

func TestIsMap(t *testing.T) {
	assert.True(t, IsMap(
		make(map[string]interface{}),
	))
}

func TestIsMapFail(t *testing.T) {
	assert.False(t, IsMap(
		make([]string, 0),
	))
}

func TestIsNil(t *testing.T) {
	assert.True(t, IsNil(
		nil,
	))
}

func TestIsNilEmptyFunction(t *testing.T) {
	var closure func(string) bool

	assert.True(t, IsNil(
		closure,
	))
}

func TestIsNilEmptyInterface(t *testing.T) {
	var data interface{}

	assert.True(t, IsNil(
		data,
	))
}

func TestIsNilEmptyPointer(t *testing.T) {
	var data *string

	assert.True(t, IsNil(
		data,
	))
}

func TestIsNilEmptySlice(t *testing.T) {
	var data []string

	assert.True(t, IsNil(
		data,
	))
}

func TestIsNilEmptyMap(t *testing.T) {
	var data map[string]interface{}

	assert.True(t, IsNil(
		data,
	))
}

func TestIsNilFailString(t *testing.T) {
	assert.False(t, IsNil(
		"hello",
	))
}

func TestIsNumericInt(t *testing.T) {
	assert.True(t, IsNumeric(
		12,
	))
}

func TestIsNumericFloat(t *testing.T) {
	assert.True(t, IsNumeric(
		float64(12),
	))
}

func TestIsNumericFail(t *testing.T) {
	assert.False(t, IsNumeric(
		"hello",
	))
}

func TestIsString(t *testing.T) {
	assert.True(t, IsString(
		"hello",
	))
}

func TestIsStringFail(t *testing.T) {
	assert.False(t, IsString(
		float64(24),
	))
}

func TestIsUint(t *testing.T) {
	assert.True(t, IsUint(
		uint(24),
	))
}

func TestIsUintFailInt(t *testing.T) {
	assert.False(t, IsUint(
		24,
	))
}

func TestIsUintFailString(t *testing.T) {
	assert.False(t, IsUint(
		"hello",
	))
}

func TestIsUintFailFloat32(t *testing.T) {
	assert.False(t, IsUint(
		float32(12),
	))
}

func TestIsPointer(t *testing.T) {
	data := "hello"

	assert.True(t, IsPointer(&data))
}

func TestIsPointerFail(t *testing.T) {
	assert.False(t, IsPointer(
		"hello",
	))
}

func TestIsObject(t *testing.T) {
	type SomeStruct struct {
		Name string
	}

	assert.True(t, IsObject(
		SomeStruct{},
	))
}

func TestIsObjectAnonymous(t *testing.T) {
	data := struct {
		Name string
	}{}

	assert.True(t, IsObject(
		data,
	))
}

func TestIsObjectFail(t *testing.T) {
	assert.False(t, IsObject(
		"hello",
	))
}
