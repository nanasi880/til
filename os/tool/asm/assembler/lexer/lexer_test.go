package lexer

import (
	"testing"

	"go.nanasi880.dev/xtesting"
)

func TestSplitToken(t *testing.T) {

	testCases := []struct {
		s     string
		wants []Token
	}{
		{
			s:     `Token1 Token2, Token3, "Token,4"`,
			wants: []Token{"Token1", "Token2", "Token3", `"Token,4"`},
		},
		{
			s:     `"Token\"5" "Token6"`,
			wants: []Token{`"Token\"5"`, `"Token6"`},
		},
		{
			s:     `"Token\\\"7\\`,
			wants: []Token{`"Token\\\"7\\`},
		},
		{
			s:     `tok "Invalid Token\\\" \ \ "`,
			wants: nil,
		},
	}

	for i, tt := range testCases {

		tokens, err := SplitToken([]rune(tt.s))
		if tt.wants != nil && err != nil {
			t.Fatal(err)
		}
		if tt.wants == nil && err != nil {
			continue
		}
		if tt.wants == nil && err == nil {
			t.Fatal(tokens)
		}

		if len(tt.wants) != len(tokens) {
			t.Fatal(i, " ", tokens)
		}

		for ii := range tokens {
			if tt.wants[ii] != tokens[ii] {
				t.Fatal(i, " ", ii, " ", tokens)
			}
		}
	}
}

func TestAnalyze(t *testing.T) {

	src := xtesting.MustOpen(t, "testdata/asm.txt")
	defer xtesting.MustClose(t, src)

	file, err := Analyze(src)
	if err != nil {
		t.Fatal(err)
	}

	if len(file) != 2 {
		t.Fatal(file)
	}
	if file[0][0] != "THIS_IS_LABEL:" {
		t.Fatal(file)
	}
	if file[1][0] != "THIS_IS_MNEMONIC" {
		t.Fatal(file)
	}
	if file[1][1] != "PARAM1" {
		t.Fatal(file)
	}
	if file[1][2] != `"PARAM2"` {
		t.Fatal(file)
	}
	if file[1][3] != `"PARAM3;"` {
		t.Fatal(file)
	}
}
