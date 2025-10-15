# bluemonday-css

[![Go](https://github.com/dsh2dsh/bluemonday-css/actions/workflows/go.yml/badge.svg)](https://github.com/dsh2dsh/bluemonday-css/actions/workflows/go.yml)
[![GoDoc](https://godoc.org/github.com/dsh2dsh/bluemonday-css?status.png)](https://godoc.org/github.com/dsh2dsh/bluemonday-css)

bluemonday-css is a [bluemonday] helper for sanitizing inline styles. It
extracted from [bluemonday] into this repository for making it optional and
reduce dependencies of [bluemonday].

[bluemonday]: https://github.com/dsh2dsh/bluemonday

# Usage

Although it's possible to handle inline CSS using `AllowAttrs` with a `Matching`
rule, writing a single monolithic regular expression to safely process all
inline CSS which you wish to allow is not a trivial task. Instead of attempting
to do so, you can allow the `style` attribute on whichever element(s) you desire
and use style policies to control and sanitize inline styles.

It is strongly recommended that you use `Matching` (with a suitable regular
expression) `MatchingEnum`, or `MatchingHandler` to ensure each style matches
your needs, but default handlers are supplied for most widely used styles.

``` go
import (
	"github.com/dsh2dsh/bluemonday"
	css "github.com/dsh2dsh/bluemonday-css"
)
```

``` go
stylesPolicy := css.NewPolicy()

// Allow the 'text-decoration' property to be set to 'underline',
// 'line-through' or 'none' on 'span' elements only.
stylesPolicy.AllowStyles("text-decoration").
  MatchingEnum("underline", "line-through", "none").OnElements("span")

// Allow the 'color' property with valid RGB(A) hex values only on every HTML
// element that has been allowed.
stylesPolicy.AllowStyles("color").
  Matching(
    regexp.MustCompile("(?i)^#([0-9a-f]{3,4}|[0-9a-f]{6}|[0-9a-f]{8})$")).
  Globally()

// Default handler
stylesPolicy.AllowStyles("background-origin").Globally()

p := bluemonday.UGCPolicy().WithStyleHandler(stylesPolicy.Sanitize)

// Allow only 'span' and 'p' elements
p.AllowElements("span", "p", "strong")

// Only allow 'style' attributes on 'span' and 'p' elements
p.AllowAttrs("style").OnElements("span", "p")

// The span has an invalid 'color' which will be stripped along with other
// disallowed properties
html := p.Sanitize(`<p style="color:#f00;">
  <span style="text-decoration: underline; background-image: url(javascript:alert('XSS')); color: #f00ba; background-origin: invalidValue">
    Red underlined <strong style="text-decoration:none;">text</strong>
  </span>
</p>`)

fmt.Println(html)
```

Which outputs:

``` html
<p style="color: #f00">
  <span style="text-decoration: underline">
    Red underlined <strong>text</strong>
  </span>
</p>
```

If you need more specific checking, you can create a handler that takes in a
string and returns a bool to validate the values for a given property. The
string parameter has been converted to lowercase and unicode code points have
been converted.

``` go
myHandler := func(value string) bool {
  // Validate your input here
  return true
}

// Allow the 'color' property with values validated by the handler (on any
// element allowed a 'style' attribute)
p.AllowStyles("color").MatchingHandler(myHandler).Globally()
```
