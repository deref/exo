package cmdutil

import (
	"bytes"
	"fmt"
	"strings"
)

type ParsedArgs struct {
	Command string
	Args    []string
	Flags   map[string]string
}

func (pa *ParsedArgs) Dump() string {
	var buf bytes.Buffer
	buf.WriteString("Command: ")
	buf.WriteString(pa.Command)
	buf.WriteByte('\n')
	buf.WriteString("Args: \n")
	for _, arg := range pa.Args {
		buf.WriteString("\t")
		buf.WriteString(arg)
		buf.WriteByte('\n')
	}
	buf.WriteString("Flags: \n")
	for name, val := range pa.Flags {
		buf.WriteString("\t")
		buf.WriteString(name)
		buf.WriteString(" = ")
		buf.WriteString(val)
		buf.WriteByte('\n')
	}

	return buf.String()
}

type argParser struct {
	tokens []parseToken
	idx    int
	state  parseState
	out    *ParsedArgs
	err    error
}

type parseToken struct {
	tag parseTokenTag
	val string
	str string
	idx int
}

type parseTokenTag int

const (
	parseTokenWord parseTokenTag = iota
	parseTokenLongFlag
	parseTokenShortFlag
)

type parseState int

const (
	parseStateCommand parseState = iota
	parseStateArg
	parseStateFlag
)

func (p *argParser) parse() {
	p.out = &ParsedArgs{
		Args:  []string{},
		Flags: make(map[string]string),
	}
	p.err = nil
	if p.isEnd() {
		return
	}
	p.parseInvocation()
}

func (p *argParser) parseInvocation() {
	// This is safe because we check p.isEnd() in p.parse().
	p.out.Command = p.next().val
	p.parseAny()
}

func (p *argParser) parseAny() {
	tok := p.next()
	switch {
	case tok == nil:
		return

	case tok.tag == parseTokenWord:
		p.out.Args = append(p.out.Args, tok.val)
		p.parseAny()

	// There is currently no difference in how we handle short and long flags.
	case tok.tag == parseTokenLongFlag || tok.tag == parseTokenShortFlag:
		valTok := p.next()
		if valTok == nil {
			p.err = fmt.Errorf("expected value for flag: %s", tok.str)
			return
		}
		if valTok.tag != parseTokenWord {
			p.err = fmt.Errorf("unexpected value for flag: %s; got another flag instead: %s", tok.str, valTok.str)
			return
		}

		p.out.Flags[tok.val] = valTok.val
		p.parseAny()

	default:
		panic("unrecognized token")
	}
}

func (p *argParser) current() *parseToken {
	if p.isEnd() {
		return nil
	}
	return &p.tokens[p.idx]
}

func (p *argParser) next() *parseToken {
	tok := p.current()
	p.idx++

	return tok
}

func (p *argParser) isEnd() bool {
	return p.idx >= len(p.tokens)
}

func tokenizeArgs(args []string) []parseToken {
	tokens := make([]parseToken, 0, len(args))
	parsingFlags := true
	for idx, arg := range args {
		switch {
		case arg == "--":
			parsingFlags = false
		case parsingFlags && strings.HasPrefix(arg, "--"):
			flag := arg[2:]
			parts := strings.Split(flag, "=")
			switch len(parts) {
			case 1:
				tokens = append(tokens, parseToken{
					str: arg,
					idx: idx,
					tag: parseTokenLongFlag,
					val: flag,
				})
			case 2:
				tokens = append(tokens, parseToken{
					str: arg,
					idx: idx,
					tag: parseTokenLongFlag,
					val: parts[0],
				}, parseToken{
					str: arg,
					idx: idx,
					tag: parseTokenWord,
					val: parts[1],
				})
			}

		case parsingFlags && strings.HasPrefix(arg, "-"):
			tokens = append(tokens, parseToken{
				str: arg,
				idx: idx,
				tag: parseTokenShortFlag,
				val: arg[1:],
			})

		default:
			tokens = append(tokens, parseToken{
				str: arg,
				idx: idx,
				tag: parseTokenWord,
				val: arg,
			})
		}
	}

	return tokens
}

func ParseArgs(args []string) (*ParsedArgs, error) {
	tokens := tokenizeArgs(args)
	parser := &argParser{tokens: tokens}
	parser.parse()

	return parser.out, parser.err
}
