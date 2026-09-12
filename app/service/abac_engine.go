package service

import (
	"fmt"
	"regexp"
	"strings"

	"gocms_v3/app/db"
	"gocms_v3/app/model/ngac"
)

type ABACEngine struct{}

func NewABACEngine() *ABACEngine { return &ABACEngine{} }

func (e *ABACEngine) Evaluate(ruleExpr string, subjectAttrs, objectAttrs map[string]string) bool {
	if ruleExpr == "" {
		return true
	}
	evaluatedExpr := e.expandAttributes(ruleExpr, subjectAttrs, objectAttrs)
	return e.evaluate(evaluatedExpr)
}

func (e *ABACEngine) evaluate(expr string) bool {
	expr = strings.TrimSpace(expr)
	if expr == "" || expr == "true" || expr == "TRUE" || expr == "1" {
		return true
	}
	if expr == "false" || expr == "FALSE" || expr == "0" {
		return false
	}
	if strings.HasPrefix(strings.ToUpper(expr), "NOT ") {
		return !e.evaluate(strings.TrimSpace(expr[4:]))
	}
	if parts := splitByOperator(expr, "OR"); len(parts) > 1 {
		for _, part := range parts {
			if e.evaluate(strings.TrimSpace(part)) {
				return true
			}
		}
		return false
	}
	if parts := splitByOperator(expr, "AND"); len(parts) > 1 {
		for _, part := range parts {
			if !e.evaluate(strings.TrimSpace(part)) {
				return false
			}
		}
		return true
	}
	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		return e.evaluate(expr[1 : len(expr)-1])
	}
	return e.evaluateComparison(expr)
}

func splitByOperator(expr, operator string) []string {
	var parts []string
	depth := 0
	current := &strings.Builder{}
	inString := false
	for _, r := range expr {
		switch r {
		case '(':
			if !inString {
				depth++
			}
			current.WriteRune(r)
		case ')':
			if !inString {
				depth--
			}
			current.WriteRune(r)
		case '"':
			inString = !inString
			current.WriteRune(r)
		default:
			if !inString && depth == 0 && strings.HasPrefix(strings.ToUpper(expr[current.Len():]), strings.ToUpper(operator)) {
				parts = append(parts, current.String())
				current.Reset()
				continue
			}
			current.WriteRune(r)
		}
	}
	parts = append(parts, current.String())
	return parts
}

func (e *ABACEngine) evaluateComparison(expr string) bool {
	ops := []string{"==", "!=", ">=", "<=", ">", "<", "=~", "!~", "^=", "contains", "startswith", "endswith"}
	for _, op := range ops {
		parts := strings.Split(expr, op)
		if len(parts) == 2 {
			return e.compare(strings.TrimSpace(parts[0]), strings.TrimSpace(strings.Trim(parts[1], `"'`)), op)
		}
	}
	return false
}

func (e *ABACEngine) compare(attrName, expectedValue, op string) bool {
	allAttrs := make(map[string]string)
	for k, v := range e.GetAllAttributes() {
		allAttrs[k] = v
	}
	actualValue := allAttrs[attrName]
	if actualValue == "" {
		return false
	}
	switch strings.ToLower(op) {
	case "==":
		return actualValue == expectedValue
	case "!=":
		return actualValue != expectedValue
	case ">=":
		return actualValue >= expectedValue
	case "<=":
		return actualValue <= expectedValue
	case ">":
		return actualValue > expectedValue
	case "<":
		return actualValue < expectedValue
	case "=~":
		m, err := regexp.MatchString(expectedValue, actualValue)
		return err == nil && m
	case "!~":
		m, err := regexp.MatchString(expectedValue, actualValue)
		return err == nil && !m
	case "^=":
		return strings.HasPrefix(actualValue, expectedValue)
	case "contains":
		return strings.Contains(actualValue, expectedValue)
	case "startswith":
		return strings.HasPrefix(actualValue, expectedValue)
	case "endswith":
		return strings.HasSuffix(actualValue, expectedValue)
	default:
		return false
	}
}

func (e *ABACEngine) expandAttributes(expr string, subjectAttrs, objectAttrs map[string]string) string {
	allAttrs := make(map[string]string)
	for k, v := range subjectAttrs {
		allAttrs["subject."+k] = v
		allAttrs[k] = v
	}
	for k, v := range objectAttrs {
		allAttrs["object."+k] = v
		if _, exists := allAttrs[k]; !exists {
			allAttrs[k] = v
		}
	}
	for key, value := range allAttrs {
		expr = strings.ReplaceAll(expr, "{"+key+"}", value)
	}
	return expr
}

func (e *ABACEngine) GetAllAttributes() map[string]string { return make(map[string]string) }

func (e *ABACEngine) GetSubjectAttributes(userID uint) map[string]string {
	var attrs []ngac.UserAttribute
	db.GetDB().Where("user_id = ?", userID).Find(&attrs)
	result := make(map[string]string)
	for _, attr := range attrs {
		key := fmt.Sprintf("subject.%s", attr.Name)
		result[key] = attr.Value
		result[attr.Name] = attr.Value
	}
	return result
}

func (e *ABACEngine) GetObjectAttributes(resourceType string, resourceID uint) map[string]string {
	var attrs []ngac.ObjectAttribute
	db.GetDB().Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).Find(&attrs)
	result := make(map[string]string)
	for _, attr := range attrs {
		key := fmt.Sprintf("object.%s", attr.Name)
		result[key] = attr.Value
		result[attr.Name] = attr.Value
	}
	return result
}

func (e *ABACEngine) GetContextAttributes() map[string]string {
	var attrs []ngac.ContextAttribute
	db.GetDB().Where("is_global = true").Find(&attrs)
	result := make(map[string]string)
	for _, attr := range attrs {
		key := fmt.Sprintf("context.%s", attr.Name)
		result[key] = attr.Value
		result[attr.Name] = attr.Value
	}
	return result
}

var globalABACEngine = NewABACEngine()

func GetABACEngine() *ABACEngine { return globalABACEngine }
