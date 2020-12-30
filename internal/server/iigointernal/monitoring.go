package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type monitor struct {
	judgeID           shared.ClientID
	speakerID         shared.ClientID
	presidentID       shared.ClientID
	internalIIGOCache []shared.Accountability
}

func (m *monitor) addToCache(roleToMonitorID shared.ClientID, variables []rules.VariableFieldName, values [][]float64) {
	pairs := []rules.VariableValuePair{}
	for index, variable := range variables {
		pairs = append(pairs, rules.MakeVariableValuePair(variable, values[index]))
	}
	m.internalIIGOCache = append(m.internalIIGOCache, shared.Accountability{
		ClientID: roleToMonitorID,
		Pairs:    pairs,
	})
}

func (m *monitor) monitorRole(roleAccountable baseclient.Client) (bool, bool) {
	roleToMonitor, roleName := m.findRoleToMonitor(roleAccountable)
	decideToMonitor := roleAccountable.MonitorIIGORole(roleName)
	evaluationResult := false
	if decideToMonitor {
		evaluationResult = m.evaluateCache(roleToMonitor)
	}
	return decideToMonitor, evaluationResult
}

func (m *monitor) evaluateCache(roleToMonitorID shared.ClientID) bool {
	performedRoleCorrectly := true
	for _, entry := range m.internalIIGOCache {
		if entry.ClientID == roleToMonitorID {
			variablePairs := entry.Pairs
			var rulesAffected []string
			for _, variable := range variablePairs {
				valuesToBeAdded, foundRules := rules.PickUpRulesByVariable(variable.VariableName, rules.RulesInPlay)
				if foundRules {
					rulesAffected = append(rulesAffected, valuesToBeAdded...)
				}
				rules.UpdateVariable(variable.VariableName, variable)
			}
			for _, rule := range rulesAffected {
				evaluation, err := rules.BasicBooleanRuleEvaluator(rule)
				if err != nil {
					continue
				}
				performedRoleCorrectly = evaluation && performedRoleCorrectly
			}
		}
	}
	return performedRoleCorrectly
}

func (m *monitor) findRoleToMonitor(roleAccountable baseclient.Client) (shared.ClientID, baseclient.Role, error) {
	switch roleAccountable.GetID() {
	case m.speakerID:
		return m.presidentID, baseclient.President, nil
	case m.presidentID:
		return m.judgeID, baseclient.Judge, nil
	case m.judgeID:
		return m.speakerID, baseclient.Speaker, nil
	}
default:
	return nil, nil, 
}
