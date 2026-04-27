package ugrpc

import (
	"fmt"
	"strings"
)

type FullMethod struct {
	Server  string
	Service string
	Method  string
}

func ParseFullMethod(fullMethod string) (*FullMethod, error) {
	if len(fullMethod) == 0 {
		return nil, fmt.Errorf("invalid full method: %s", fullMethod)
	}
	fullMethod = fullMethod[1:]
	strs := strings.Split(fullMethod, "/")
	if len(strs) != 2 {
		return nil, fmt.Errorf("invalid full method: %s", fullMethod)
	}
	serverAndService := strs[0]
	method := strs[1]

	strs = strings.Split(serverAndService, ".")
	if len(strs) != 2 {
		return nil, fmt.Errorf("invalid full method: %s", fullMethod)
	}

	return &FullMethod{
		Server:  strs[0],
		Service: strs[1],
		Method:  method,
	}, nil
}
