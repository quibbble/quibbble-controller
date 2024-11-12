package errors

import (
	"fmt"
	"runtime"
	"strings"

	"slices"
)

const enabled = true

// Stack is an special error that creates a stacktrace
type Stack struct {
	Cause   error
	Callers []string
}

func Errorf(format string, args ...any) error {
	err := fmt.Errorf(format, args...)
	if !enabled {
		return err
	}
	pc, file, no, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()
	trim := strings.Split(name, "/")
	return &Stack{
		Cause:   err,
		Callers: []string{fmt.Sprintf("%s\n  %s:%d\n", trim[len(trim)-1], file, no)},
	}
}

func Wrap(err error) error {
	if !enabled {
		return err
	}
	stack, ok := err.(*Stack)
	if !ok {
		pc, file, no, _ := runtime.Caller(1)
		name := runtime.FuncForPC(pc).Name()
		trim := strings.Split(name, "/")
		return &Stack{
			Cause:   err,
			Callers: []string{fmt.Sprintf("%s\n  %s:%d\n", trim[len(trim)-1], file, no)},
		}
	}
	pc, file, no, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()
	trim := strings.Split(name, "/")
	stack.Callers = append(stack.Callers, fmt.Sprintf("%s\n  %s:%d\n", trim[len(trim)-1], file, no))
	return stack
}

func (s *Stack) Error() string {
	stack := s.Cause.Error() + "\n"
	slices.Reverse(s.Callers)
	for _, caller := range s.Callers {
		stack += caller
	}
	return stack[:len(stack)-1]
}
