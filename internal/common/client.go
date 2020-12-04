package common

import (
	"errors"
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"gonum.org/v1/gonum/mat"
)

// ClientID is an enum for client IDs
type ClientID int

const (
	// Team1 ID
	Team1 ClientID = iota
	// Team2 ID
	Team2
	// Team3 ID
	Team3
	// Team4 ID
	Team4
	// Team5 ID
	Team5
	// Team6 ID
	Team6
)

// TeamIDs contain sequential IDs of all teams
var TeamIDs = [...]ClientID{Team1, Team2, Team3, Team4, Team5, Team6}

// Client is a base interface to be implemented by each client struct.
type Client interface {
	Echo(s string) string
	GetID() ClientID
	Logf(format string, a ...interface{})
}

// RegisteredClients contain all registered clients, exposed for the server.
var RegisteredClients = map[ClientID]Client{}

// RegisterClient registers clients into RegisteredClients
func RegisterClient(id ClientID, c Client) {
	// prevent double registrations
	if _, ok := RegisteredClients[id]; ok {
		// OK to panic here, as this is a _crucial_ step.
		panic(fmt.Sprintf("Duplicate client ID %v in RegisterClient!", id))
	}
	RegisteredClients[id] = c
}

// BasicRuleEvaluator implements a basic version of the Matrix rule evaluator, provides single boolean output (and error if present)
func BasicRuleEvaluator(ruleName string) (bool, error) {
	var variableVect []float64
	if rm, ok := rules.AvailableRules[ruleName]; ok {
		variables := rm.RequiredVariables

		for _, v := range variables {
			if val, varOk := rules.VariableMap[v]; varOk {
				variableVect = append(variableVect, val.Values...)
			} else {
				return false, errors.New(fmt.Sprintf("Variable: '%v' not found in global variable cache", v))
			}
		}

		variableVect = append(variableVect, 1)

		//Checking dimensions line up
		nRows, nCols := rm.ApplicableMatrix.Dims()

		if nCols != len(variableVect) {
			return false, errors.New(
				fmt.Sprintf("dimension mismatch in evaluating rule: '%v' rule matrix has '%v' columns, while we sourced '%v' variables",
					ruleName,
					nCols,
					len(variableVect),
				))
		}

		variableFormalVect := mat.NewVecDense(len(variableVect), variableVect)

		actual := make([]float64, nRows)

		c := mat.NewVecDense(nRows, actual)

		c.MulVec(&rm.ApplicableMatrix, variableFormalVect)
		aux := rm.AuxiliaryVector

		var resultVect []bool

		for i := 0; i < nRows; i++ {
			res := false
			switch interpret := aux.AtVec(i); interpret {
			case 0:
				res = c.AtVec(i) == 0
			case 1:
				res = c.AtVec(i) > 0
			case 2:
				res = c.AtVec(i) >= 0
			case 3:
				res = c.AtVec(i) != 0
			default:
				return false, errors.New(fmt.Sprintf("At auxillary vector entry: '%v' aux value outside of 0-3: '%v' was found", i, interpret))
			}
			resultVect = append(resultVect, res)
		}

		var finalBool = true

		for _, v := range resultVect {
			finalBool = finalBool && v
		}

		return finalBool, nil
	} else {
		return false, errors.New(fmt.Sprintf("rule name: '%v' provided doesn't exist in global rule list", ruleName))
	}
}
