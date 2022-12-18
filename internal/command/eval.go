package command

import (
	"fmt"

	"github.com/hashibuto/artillery"
	"github.com/hashibuto/shellfire/internal/utils"
)

var EvalCommand = &artillery.Command{
	Name:        "eval",
	Description: "evaluate basic arithmatic expression",
	Arguments: []*artillery.Argument{
		{
			Name:        "expression",
			Description: "arithmatic expression",
		},
	},
	Options: []*artillery.Option{
		{
			Name:        "dec",
			ShortName:   'd',
			Description: "output result in decimal format",
			Type:        artillery.Bool,
			Value:       true,
		},
	},
	OnExecute: evalExpressionCmd,
}

func evalExpressionCmd(n artillery.Namespace, p *artillery.Processor) error {
	var args struct {
		Expression string
		Dec        bool
	}
	err := artillery.Reflect(n, &args)
	if err != nil {
		return err
	}

	value, err := evalExpression(args.Expression)
	if err != nil {
		return err
	}

	if args.Dec {
		fmt.Printf("%d\n", value)
	} else {
		fmt.Printf("\\x%08x\n", value)
	}

	return nil
}

func evalExpression(exp string) (int, error) {
	tokens := []string{}
	curTok := ""
	for i := 0; i < len(exp); i++ {
		char := exp[i]
		if char == '+' || char == '-' {
			if len(curTok) > 0 {
				tokens = append(tokens, curTok)
				curTok = ""
			}

			tokens = append(tokens, string(char))
		} else if char == ' ' || char == '\t' {
			if len(curTok) > 0 {
				tokens = append(tokens, curTok)
				curTok = ""
			}
		} else {
			curTok += string(char)
		}
	}
	if len(curTok) > 0 {
		tokens = append(tokens, curTok)
		curTok = ""
	}

	if len(tokens) == 0 {
		return 0, fmt.Errorf("Empty expression")
	}

	prevTok := ""
	for _, tok := range tokens {
		if tok == "-" || tok == "+" {
			if prevTok == "-" || prevTok == "+" {
				return 0, fmt.Errorf("Invalid expression, multiple adjacent operators")
			}
			if prevTok == "" {
				return 0, fmt.Errorf("Expression cannot begin with an operator")
			}
		}
		prevTok = tok
	}
	if prevTok == "-" || prevTok == "+" {
		return 0, fmt.Errorf("Expression cannot end with an operator")
	}

	// Tokens sequence is fine, now evaluate the expression and break out if an error
	op := ""
	total := 0
	for _, tok := range tokens {
		if tok == "-" || tok == "+" {
			op = tok
		} else {
			val, err := utils.ParseNumber(tok)
			if err != nil {
				return 0, err
			}
			switch op {
			case "-":
				total -= val
			case "+":
				total += val
			default:
				total = val
			}
		}
	}

	return total, nil
}
