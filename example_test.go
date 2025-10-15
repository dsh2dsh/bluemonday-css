package css_test

import (
	"fmt"
	"regexp"

	css "github.com/dsh2dsh/bluemonday-css"
)

func Example() {
	p := css.NewPolicy()

	// Allow the 'text-decoration' property to be set to 'underline',
	// 'line-through' or 'none' on 'span' elements only.
	p.AllowStyles("text-decoration").
		MatchingEnum("underline", "line-through", "none").OnElements("span")

	// Allow the 'color' property with valid RGB(A) hex values only on every HTML
	// element that has been allowed.
	p.AllowStyles("color").
		Matching(
			regexp.MustCompile("(?i)^#([0-9a-f]{3,4}|[0-9a-f]{6}|[0-9a-f]{8})$")).
		Globally()

	// Default handler
	p.AllowStyles("background-origin").Globally()

	fmt.Println(p.Sanitize("p", "color:#f00;"))
	fmt.Println(p.Sanitize("span",
		"text-decoration: underline; background-image: url(javascript:alert('XSS')); color: #f00ba; background-origin: invalidValue"))
	fmt.Println(p.Sanitize("strong", "text-decoration:none;"))
	// Output:
	// color: #f00
	// text-decoration: underline
}
