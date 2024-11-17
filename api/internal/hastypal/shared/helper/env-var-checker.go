package helper

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/shared/types"
	"os"
)

type EnvVarChecker struct {
	envVars []string
}

func NewEnvVarChecker(vars ...string) *EnvVarChecker {
	return &EnvVarChecker{
		envVars: vars,
	}
}

func (evc *EnvVarChecker) Check() error {
	for _, key := range evc.envVars {
		_, exists := os.LookupEnv(key)

		if !exists {
			return types.ApiError{
				Msg:      fmt.Sprintf("Env var %s not set", key),
				Function: "Check",
				File:     "env-var-checker.go",
			}
		}
	}

	return nil
}
