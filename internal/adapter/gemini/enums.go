package gemini

import "fmt"

func (r Role) String() string {
	return string(r)
}

func (r Role) IsValid() bool {
	switch r {
	case RoleUser, RoleModel:
		return true
	default:
		return false
	}
}

func ParseRole(s string) (Role, error) {
	role := Role(s)
	if !role.IsValid() {
		return "", fmt.Errorf("invalid role: %s", s)
	}
	return role, nil
}

func (f FinishReason) String() string {
	return string(f)
}

func (f FinishReason) IsValid() bool {
	switch f {
	case FinishReasonStop, FinishReasonMaxTokens, FinishReasonSafety,
		FinishReasonRecitation, FinishReasonOther, FinishReasonBlocklist,
		FinishReasonProhibitedContent, FinishReasonSpii:
		return true
	default:
		return false
	}
}

func (f FinishReason) IsError() bool {
	switch f {
	case FinishReasonSafety, FinishReasonRecitation, FinishReasonBlocklist,
		FinishReasonProhibitedContent, FinishReasonSpii:
		return true
	default:
		return false
	}
}
