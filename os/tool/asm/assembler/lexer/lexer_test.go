package lexer

import "testing"

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
			s:     `"Token\"5"`,
			wants: []Token{`"Token\"5"`},
		},
	}

	for i, tt := range testCases {

		tokens, err := SplitToken([]byte(tt.s))
		if tt.wants != nil && err != nil {
			t.Fatal(err)
		}
		if tt.wants == nil {
			continue
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
