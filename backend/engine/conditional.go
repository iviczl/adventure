package engine

import "fmt"

func conditional(function map[string]interface{}, adventure *Adventure) bool {
	conditionsMet := true
	// Process conditions
	for _, conditionInterface := range function["conditions"].([]interface{}) {
		condition := conditionInterface.(map[string]interface{})
		if _, ok := condition["playerAttributeMustBe"].(bool); ok {
			attributeCode, exists := condition["attributeCode"].(string)
			if !exists {
				panic("Missing attribute code.")
			}
			attributeValue, error := adventure.Player.GetAttributeValue(attributeCode)
			if error != nil {
				panic(fmt.Sprintf("Attribute not found: %s", attributeCode))
			}
			equalsValue, exists := condition["equals"].(int)
			if exists {
				conditionsMet = conditionsMet && (attributeValue == equalsValue)
			} else if greaterThanValue, exists := condition["greaterThan"].(int); exists {
				conditionsMet = conditionsMet && (attributeValue > greaterThanValue)
			} else if lessThanValue, exists := condition["lessThan"].(int); exists {
				conditionsMet = conditionsMet && (attributeValue < lessThanValue)
			} else {
				panic("Missing comparison value for player attribute condition.")
			}
		} else if mustHave, ok := condition["playerMustHave"].(bool); ok {
			if _, exists := condition["itemCode"].(string); !exists {
				panic("Missing item code.")
			}
			itemCode := condition["itemCode"].(string)
			item := adventure.Player.GetItem(itemCode)
			hasItem := item != nil
			conditionsMet = conditionsMet && !(mustHave != hasItem)
			if state, exists := condition["itemState"]; exists {
				stateOk := item.State == state.(string)
				conditionsMet = conditionsMet && !(stateOk != hasItem)
			}
		} else if mustHave, ok := condition["positionMustHave"]; ok {
			if _, exists := condition["positionCode"]; !exists {
				panic("Missing position code.")
			}
			positionCode := condition["positionCode"].(string)
			position := adventure.GetPosition(positionCode)
			if position == nil {
				panic(fmt.Sprintf("Position not found: %s", positionCode))
			}
			if _, exists := condition["itemCode"]; !exists {
				panic("Missing item code.")
			}
			itemCode := condition["itemCode"].(string)
			hasItem := position.GetItem(itemCode) != nil
			conditionsMet = conditionsMet && !(mustHave.(bool) != hasItem)

		} else if positionCode, ok := condition["positionCode"]; ok {
			position := adventure.GetPosition(positionCode.(string))
			if position == nil {
				panic(fmt.Sprintf("Position not found: %s", positionCode))
			}
			if visited, ok := condition["visited"]; ok {
				conditionsMet = conditionsMet && (position.Visited == visited.(bool))
			} else if endPosition, ok := condition["endPosition"]; ok {
				conditionsMet = conditionsMet && (position.EndPosition == endPosition.(bool))
			}

		} else if itemCode, ok := condition["itemCode"]; ok {
			item := adventure.GetItem(itemCode.(string))
			if item == nil {
				panic(fmt.Sprintf("Not existing item %s.", itemCode))
			}
			if name, ok := condition["name"]; ok {
				conditionsMet = conditionsMet && (item.Name == name.(string))
			} else if state, ok := condition["state"]; ok {
				conditionsMet = conditionsMet && (item.State == state.(string))
			}

		} else if npcCode, ok := condition["npcCode"]; ok {
			npc := adventure.GetNpc(npcCode.(string))
			if npc == nil {
				panic(fmt.Sprintf("Not existing npc %s.", npcCode))
			}
			if name, ok := condition["name"]; ok {
				conditionsMet = conditionsMet && (npc.Name == name.(string))
			} else if state, ok := condition["state"].(string); ok {
				conditionsMet = conditionsMet && (npc.State == state)
			}
		} else {
			fmt.Println("NOTHING")
		}

		if !conditionsMet {
			break
		}
	}

	// Execute actions based on conditions
	if conditionsMet {
		for _, actionCodeInterface := range function["actionCodes"].([]interface{}) {
			actionCode := actionCodeInterface.(string)
			executable, err := adventure.findAction(actionCode)
			if err == nil && executable != nil {
				adventure.ExecuteAction(executable)
			}
		}
	} else {
		if elseActionCodes, ok := function["elseActionCodes"]; ok {
			for _, actionCodeInterface := range elseActionCodes.([]interface{}) {
				actionCode := actionCodeInterface.(string)
				executable, err := adventure.findAction(actionCode)
				if err == nil && executable != nil {
					adventure.ExecuteAction(executable)
				}
			}
		}
	}

	return conditionsMet
}
