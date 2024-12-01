package helper

import (
	"fmt"
	"github.com/adriein/hastypal/internal/hastypal/shared/exception"
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
			return exception.
				New(fmt.Sprintf("Env var %s not set", key)).
				Trace("os.LookupEnv", "env-var-checker.go")
		}
	}

	return nil
}
