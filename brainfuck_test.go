package brainfuck

import "testing"

func TestValidate(t *testing.T) {
	code := "["
	err := Validate(code)
	if err == nil {
		t.Errorf("\"%s\" is invalid code but marked as valid\n", code)
	}

	code = "]"
	err = Validate(code)
	if err == nil {
		t.Errorf("\"%s\" is invalid code but marked as valid\n", code)
	}

	code = "[]"
	err = Validate(code)
	if err != nil {
		t.Errorf("\"%s\" is valid code but marked as invalid: %s\n", code, err.Error())
	}

	code = "[[[][[[]]]]]"
	err = Validate(code)
	if err != nil {
		t.Errorf("\"%s\" is valid code but marked as invalid: %s\n", code, err.Error())
	}

	code = "[[[[[]]][][]]][]]]]]]][]]][[]"
	err = Validate(code)
	if err == nil {
		t.Errorf("\"%s\" is invalid code but marked as valid\n", code)
	}
}
