package service

import (
	"time"

	"gocms_v3/app/db"
	"gocms_v3/app/model/ngac"
)

type PDPConfig struct {
	LogEnabled      bool
	CacheEnabled    bool
	DefaultDecision string
}

var pdpConfig = PDPConfig{LogEnabled: true, CacheEnabled: false, DefaultDecision: "deny"}

func InitializePDP(cfg map[string]interface{}) {
	if v, ok := cfg["log_enabled"]; ok {
		pdpConfig.LogEnabled = v.(bool)
	}
	if v, ok := cfg["cache_enabled"]; ok {
		pdpConfig.CacheEnabled = v.(bool)
	}
	if v, ok := cfg["default_decision"]; ok {
		pdpConfig.DefaultDecision = v.(string)
	}
}

func EvaluateAuthz(req ngac.AuthzRequest) *ngac.AuthzResponse {
	start := time.Now()
	resp := &ngac.AuthzResponse{Allowed: false, Decision: pdpConfig.DefaultDecision, Reason: "default deny"}

	var userRoles []ngac.UserRole
	db.GetDB().Where("user_id = ?", req.UserID).Find(&userRoles)
	if len(userRoles) == 0 {
		return resp
	}

	roleIDs := make([]uint, len(userRoles))
	for i, ur := range userRoles {
		roleIDs[i] = ur.RoleID
	}

	var perms []ngac.Permission
	db.GetDB().Where("id IN ?", roleIDs).Find(&perms)

	permIDs := make(map[uint]bool)
	for _, rp := range perms {
		permIDs[rp.ID] = true
	}

	var matchedRules []ngac.RuleMatchInfo
	abacEngine := GetABACEngine()
	subjectAttrs := abacEngine.GetSubjectAttributes(req.UserID)
	objectAttrs := abacEngine.GetObjectAttributes(req.ResourceType, req.ResourceID)

	for _, perm := range perms {
		if perm.Action != req.Action {
			continue
		}
		attrs := make(map[string]string)
		for k, v := range subjectAttrs {
			attrs[k] = v
		}
		for k, v := range objectAttrs {
			attrs[k] = v
		}

		matched := abacEngine.Evaluate(perm.RuleExpr, subjectAttrs, objectAttrs)
		info := ngac.RuleMatchInfo{PermissionID: perm.ID, PermissionName: perm.Name, Code: perm.Code, Action: perm.Action, Allowed: perm.Allowed, Priority: perm.Priority}
		if matched {
			if perm.Allowed {
				resp.Allowed = true
				resp.Decision = "allow"
				resp.Reason = "permission allowed by rule"
			} else {
				resp.Allowed = false
				resp.Decision = "deny"
				resp.Reason = "permission denied by explicit rule"
			}
			matchedRules = append(matchedRules, info)
		}
	}

	resp.RulesMatched = matchedRules
	resp.EvaluationTime = time.Since(start).Milliseconds()

	if pdpConfig.LogEnabled {
		log := ngac.AccessLog{UserID: req.UserID, Username: req.Username, ResourceType: req.ResourceType, ResourceID: req.ResourceID, Action: req.Action, Decision: resp.Decision, EvaluationTime: resp.EvaluationTime, Reason: resp.Reason}
		db.GetDB().Create(&log)
	}

	return resp
}

func GetPDPConfig() PDPConfig { return pdpConfig }
