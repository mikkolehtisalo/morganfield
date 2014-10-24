package morganfield

import "testing"
import "fmt"

type SimpleStruct struct {
	Username string `json:"username"`
}

type SimpleIntegerStruct struct {
	Number1 int `json:"number1"`
	Number2 int `json:"number2"`
	Number3 int `json:"number3"`
}

type SimpleFloatStruct struct {
	Number1 float64 `json:"number1"`
}

// Go's JSON parser will not parse anything with missing brackets
func Test_missing_end_bracket(t *testing.T) {
	instr := "{\"username\": \"test\""
	_, err := UnMarshal(instr, &SimpleStruct{})
	if err == nil {
		t.Error("Test_missing_end_bracket parsed input with missing end bracket!")
	}
}

// Go's JSON parser does not parse unquoted keys
func Test_unquoted_key(t *testing.T) {
	instr := "{ username: justme }"
	_, err := UnMarshal(instr, &SimpleStruct{})
	if err == nil {
		t.Error("Test_unquoted_key parsed unquoted key")
	}
}

// Go's JSON parser does not parse unquoted values
func Test_unquoted_value(t *testing.T) {
	instr := "{ \"username\": justme }"
	_, err := UnMarshal(instr, &SimpleStruct{})
	if err == nil {
		t.Error("Test_unquoted_value parsed unquoted value")
	}
}

// Go's JSON parser does not parse empty input
func Test_whitespace(t *testing.T) {
	instr := "       "
	_, err := UnMarshal(instr, &SimpleStruct{})
	if err == nil {
		t.Error("Test_whitespace parsed empty input")
	}
}

// Go's JSON parser does not parse single quotes
func Test_single_quote(t *testing.T) {
	instr := "{'username':'admin'}"
	_, err := UnMarshal(instr, &SimpleStruct{})
	if err == nil {
		t.Error("Test_single_quote parsed single quotes")
	}
}

// Go's JSON parses converts both hex and octal escapes
func Test_escapes(t *testing.T) {
	// First hex, then octal escape
	instr := "{\"username\":\"test\x61test\141\"}"
	k, err := UnMarshal(instr, &SimpleStruct{})
	if err != nil {
		t.Error(fmt.Sprintf("Test_escapes UnMarshal: %v", err))
	}

	i, err := Marshal(k)
	if err != nil {
		t.Error(fmt.Sprintf("Test_escapes Marshal: %v", err))
	}

	if i != "{\"username\":\"testatesta\"}" {
		t.Error("Test_escapes was not able to convert escapes correctly")
	}
}

// Go's JSON parser will not allow comments
func Test_comments(t *testing.T) {
	instr := "{\n//Comment 1\n\"username\":\"test\"//Comment2\n}"
	_, err := UnMarshal(instr, &SimpleStruct{})
	if err == nil {
		t.Error("Test_comments parsed comments")
	}
}

// Go's JSON parser does not parse Hex or integer literals
func Test_integer_literals(t *testing.T) {
	instr1 := "{ \"number1\": 0xAB}"
	_, err := UnMarshal(instr1, &SimpleIntegerStruct{})
	if err == nil {
		t.Error("Test_integer_literals parsed hex integer literal")
	}

	instr2 := "{ \"number1\": 012 }"
	_, err = UnMarshal(instr2, &SimpleIntegerStruct{})
	if err == nil {
		t.Error("Test_integer_literals parsed octal integer literal")
	}
}

// -0.5 : ok
// +0.5 : fail
//  0.5 : ok
//   .5 : fail
//  +.5 : fail
//  -.5 : fail
func Test_decimal_numbers(t *testing.T) {
	// -0.5 is fine
	instr := "{ \"number1\": -0.5}"
	_, err := UnMarshal(instr, &SimpleFloatStruct{})
	if err != nil {
		t.Error(err)
	}

	instr = "{ \"number1\": +0.5}"
	_, err = UnMarshal(instr, &SimpleFloatStruct{})
	if err == nil {
		t.Error("Test_decimal_numbers parsed +0.5")
	}

	instr = "{ \"number1\": 0.5}"
	_, err = UnMarshal(instr, &SimpleFloatStruct{})
	if err != nil {
		t.Error(err)
	}

	instr = "{ \"number1\": .5}"
	_, err = UnMarshal(instr, &SimpleFloatStruct{})
	if err == nil {
		t.Error("Test_decimal_numbers parsed .5")
	}

	instr = "{ \"number1\": +.5}"
	_, err = UnMarshal(instr, &SimpleFloatStruct{})
	if err == nil {
		t.Error("Test_decimal_numbers parsed +.5")
	}

	instr = "{ \"number1\": -.5}"
	_, err = UnMarshal(instr, &SimpleFloatStruct{})
	if err == nil {
		t.Error("Test_decimal_numbers parsed -.5")
	}
}

// Go's JSON parser does not allow ellisions
func Test_ellisions(t *testing.T) {
	instr := "{ \"number1\": 10,,\"number2\": 20}"
	_, err := UnMarshal(instr, &SimpleIntegerStruct{})
	if err == nil {
		t.Error("Test_ellisions parsed ellision")
	}
}

// Go's JSON parser does not allow trailing commas
func Test_trailing_commas(t *testing.T) {
	instr := "{ \"number1\": 10,\"number2\": 20,}"
	_, err := UnMarshal(instr, &SimpleIntegerStruct{})
	if err == nil {
		t.Error("Test_trailing_commas parsed trailing comma")
	}
}

func Test_grouping_parenthesss(t *testing.T) {
	instr := "{ (\"number1\": 10,\"number2\": 20) }"
	_, err := UnMarshal(instr, &SimpleIntegerStruct{})
	if err == nil {
		t.Error("Test_trailing_commas parsed grouping parentheses")
	}
}

// Extra elements get silently removed
func Test_extra_element(t *testing.T) {
	instr := "{ \"number1\": 10,\"duck\":5,\"number2\": 20 }"
	_, err := UnMarshal(instr, &SimpleIntegerStruct{})
	if err != nil {
		t.Error(err)
	}
}

// Missing elements return default value
func Test_missing_element(t *testing.T) {
	instr := "{ \"number1\": 10,\"duck\":5,\"number2\": 20 }"
	k, err := UnMarshal(instr, &SimpleIntegerStruct{})
	if err != nil {
		t.Error(err)
	}

	x := k.(*SimpleIntegerStruct)
	if x.Number3 != 0 {
		t.Error("Test_missing_element: Number3 was not 0")
	}
}
