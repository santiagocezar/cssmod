package cssmod

import (
	"fmt"
	"hash/adler32"
	"io"
	"regexp"
	"strings"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

func hash(b []byte) string {
	h := adler32.Checksum(b)
	return fmt.Sprintf("%x", h)
}

/*
Transform takes input from a Reader (r) with a filename.
Returns the transfromed CSS with namespaced classes, and
a map that indexes the new classes with the original
class names.
*/
func Transform(r io.Reader, filename string) (output []byte, classes map[string]string) {
	parser := parse.NewInput(r)
	lexer := css.NewLexer(parser)

	classes = make(map[string]string)

	re := regexp.MustCompile("(^[~!@$%^&*()+=,./';:\"?><[\\]\\\\{}|`#0-9])?[~!@$%^&*()+=,./';:\"?><[\\]\\\\{}|`#]+")

	safeName := re.ReplaceAllString(strings.TrimSuffix(filename, ".module.css"), "-")
	id := hash(parser.Bytes())

	parsingClass := false

	for {
		tt, text := lexer.Next()
		switch tt {
		case css.ErrorToken:
			if lexer.Err() != io.EOF {
				fmt.Printf("An error ocurred: %v", lexer.Err())
			}
			return
		case css.DelimToken:
			if len(text) == 1 && text[0] == 46 { // a period
				parsingClass = true
			}
			output = append(output, text...)
		case css.IdentToken:
			if parsingClass {
				parsingClass = false
				orig := string(text)
				gen := safeName + "_" + orig + "_" + id
				classes[orig] = gen
				output = append(output, gen...)
			} else {
				output = append(output, text...)
			}
		default:
			output = append(output, text...)
		}
	}

}
