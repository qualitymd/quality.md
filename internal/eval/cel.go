package eval

import (
	"fmt"

	"github.com/google/cel-go/cel"
)

// evalCEL compiles and evaluates a CEL predicate. A boolean true result passes.
//
// The environment is intentionally empty for now; as the spec model grows we
// will declare variables (file metrics, prior step results, etc.) here.
func evalCEL(factor, name, expr string) Result {
	r := Result{Factor: factor, Requirement: name}

	env, err := cel.NewEnv()
	if err != nil {
		return celErr(r, "env", err)
	}

	ast, iss := env.Compile(expr)
	if iss != nil && iss.Err() != nil {
		return celErr(r, "compile", iss.Err())
	}

	prg, err := env.Program(ast)
	if err != nil {
		return celErr(r, "program", err)
	}

	val, _, err := prg.Eval(map[string]any{})
	if err != nil {
		return celErr(r, "eval", err)
	}

	if pass, ok := val.Value().(bool); ok && pass {
		r.Status = StatusPass
		r.Detail = "expression evaluated true"
		return r
	}

	r.Status = StatusFail
	r.Detail = fmt.Sprintf("expression evaluated to %v", val.Value())
	return r
}

func celErr(r Result, stage string, err error) Result {
	r.Status = StatusFail
	r.Detail = fmt.Sprintf("cel %s: %v", stage, err)
	return r
}
