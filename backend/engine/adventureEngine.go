package engine

import (
	"fmt"
	"math/rand"
	"text-adventure/models"
)

type Adventure struct {
	models.Adventure
}

func (adventure *Adventure) executeLeavingActions() {
	if adventure.ActualPosition == nil {
		return
	}

	for i := range adventure.ActualPosition.LeavingActions {
		positionCode := adventure.ActualPosition.Code
		if err := adventure.ExecuteAction(adventure.ActualPosition.LeavingActions[i]); err != nil {
			fmt.Printf("Error executing leaving action: %v\n", err)
		}
		// Break if the position code changes
		if positionCode != adventure.ActualPosition.Code {
			break
		}
	}
}

func (adventure *Adventure) executeEnteringActions() {
	for i := range adventure.ActualPosition.EnteringActions {
		positionCode := adventure.ActualPosition.Code
		if err := adventure.ExecuteAction(adventure.ActualPosition.EnteringActions[i]); err != nil {
			fmt.Printf("Error executing entering action: %v\n", err)
		}
		// Break if the position code changes
		if positionCode != adventure.ActualPosition.Code {
			break
		}
	}
}

func GetPositionFromPositionList(positions []*models.Position, positionCode string) *models.Position {
	for i := range positions {
		if positions[i].Code == positionCode {
			return positions[i]
		}
	}
	return nil
}

// getAction retrieves an action by its code from the adventure.
func getAction(adventure *Adventure, actionCode string) *models.Action {
	// Try to get the action from the positions
	action := getActionFromPositionList(adventure.Positions, actionCode)
	if action == nil {
		// If not found, try to get it from the available actions
		for i := range adventure.AvailableActions {
			if adventure.AvailableActions[i].Code == actionCode {
				return adventure.AvailableActions[i]
			}
		}
	}
	return action
}

// getActionFromPositionList retrieves an action by its code from a list of positions.
func getActionFromPositionList(positions []*models.Position, actionCode string) *models.Action {
	for i := range positions {
		action := getActionFromPosition(positions[i], actionCode)
		if action != nil {
			return action
		}
	}
	return nil
}

// getActionFromPosition retrieves an action by its code from a single position.
func getActionFromPosition(position *models.Position, actionCode string) *models.Action {
	// Check available actions
	for i := range position.AvailableActions {
		if position.AvailableActions[i].Code == actionCode {
			return position.AvailableActions[i]
		}
	}
	// Check entering actions
	for i := range position.EnteringActions {
		if position.EnteringActions[i].Code == actionCode {
			return position.EnteringActions[i]
		}
	}
	// Check leaving actions
	for i := range position.LeavingActions {
		if position.LeavingActions[i].Code == actionCode {
			return position.LeavingActions[i]
		}
	}
	return nil
}

func (adventure *Adventure) findAction(actionCode string) (*models.Action, error) {
	action := getAction(adventure, actionCode)
	if action == nil {
		return nil, fmt.Errorf("invalid actionCode %v", actionCode)
	}
	return action, nil
}

// changeActionActive changes the active state of an action.
func (adventure *Adventure) changeActionActive(actionCode string, value bool) {
	action, err := adventure.findAction(actionCode)
	if err != nil {
		fmt.Printf("Error finding action: %v %v\n", actionCode, err)
		return
	}
	if action != nil {
		action.Active = &value
	}
}

// changeActionVisible changes the visibility of an action.
func (adventure *Adventure) changeActionVisible(actionCode string, value bool) {
	action, err := adventure.findAction(actionCode)
	if err != nil {
		return
	}
	if action != nil {
		valueVar := value
		action.Visible = &valueVar
	}
}

// _findItem finds an item in the adventure by its code.
func (adventure *Adventure) findItem(itemCode string) (*models.Item, error) {
	item := adventure.getItem(itemCode)
	if item == nil {
		return nil, fmt.Errorf("invalid itemCode %s", itemCode)
	}
	return item, nil
}

// getItem retrieves an item by its code from the adventure.
// This function needs to be implemented based on your data structure.
func (adventure *Adventure) getItem(itemCode string) *models.Item {
	// Check the player's items
	for i := range adventure.Player.Items {
		if adventure.Player.Items[i].Code == itemCode {
			return adventure.Player.Items[i]
		}
	}
	// Check the items in the current position
	if adventure.ActualPosition != nil {
		for i := range adventure.ActualPosition.Items {
			if adventure.ActualPosition.Items[i].Code == itemCode {
				return adventure.ActualPosition.Items[i]
			}
		}
	}
	return nil
}

// changeItemDescription changes the description of an item.
func (adventure *Adventure) changeItemDescription(itemCode string, value string) error {
	item, err := adventure.findItem(itemCode)
	if err != nil {
		return err
	}
	item.Description = value
	return nil
}

// changeItemState changes the state of an item.
func (adventure *Adventure) changeItemState(itemCode string, value string) error {
	item, err := adventure.findItem(itemCode)
	if err != nil {
		return err
	}
	item.State = value
	return nil
}

// popItemFromList removes an item with the given code from the list and returns it.
func popItemFromList(items *[]*models.Item, itemCode string) *models.Item {
	for i, item := range *items {
		if item.Code == itemCode {
			// Remove the item from the slice
			removedItem := (*items)[i]
			*items = append((*items)[:i], (*items)[i+1:]...)
			return removedItem
		}
	}
	return nil
}

// getItemFromPlayer retrieves an item by its code from the player's inventory.
func getItemFromPlayer(player *models.Player, itemCode string) *models.Item {
	for i := range player.Items {
		if player.Items[i].Code == itemCode {
			return player.Items[i]
		}
	}
	return nil
}

// getItemFromPosition retrieves an item by its code from a single position.
func getItemFromPosition(position *models.Position, itemCode string) *models.Item {
	for i := range position.Items {
		if position.Items[i].Code == itemCode {
			return position.Items[i]
		}
	}
	return nil
}

func conditional(function map[string]interface{}, adventure *Adventure) bool {
	conditionsMet := true

	// Process conditions
	for _, conditionInterface := range function["conditions"].([]interface{}) {
		condition := conditionInterface.(map[string]interface{})

		if mustHave, ok := condition["playerMustHave"]; ok {
			if _, exists := condition["itemCode"]; !exists {
				panic("Missing item code.")
			}
			itemCode := condition["itemCode"].(string)
			hasItem := getItemFromPlayer(adventure.Player, itemCode) != nil
			conditionsMet = conditionsMet && !(mustHave.(bool) != hasItem)

		} else if mustHave, ok := condition["positionMustHave"]; ok {
			if _, exists := condition["positionCode"]; !exists {
				panic("Missing position code.")
			}
			positionCode := condition["positionCode"].(string)
			position := GetPositionFromPositionList(adventure.Positions, positionCode)
			if position == nil {
				panic(fmt.Sprintf("Position not found: %s", positionCode))
			}
			if _, exists := condition["itemCode"]; !exists {
				panic("Missing item code.")
			}
			itemCode := condition["itemCode"].(string)
			hasItem := getItemFromPosition(position, itemCode) != nil
			conditionsMet = conditionsMet && !(mustHave.(bool) != hasItem)

		} else if positionCode, ok := condition["positionCode"]; ok {
			position := GetPositionFromPositionList(adventure.Positions, positionCode.(string))
			if position == nil {
				panic(fmt.Sprintf("Position not found: %s", positionCode))
			}
			if visited, ok := condition["visited"]; ok {
				conditionsMet = conditionsMet && (position.Visited == visited.(bool))
			} else if endPosition, ok := condition["endPosition"]; ok {
				conditionsMet = conditionsMet && (position.EndPosition == endPosition.(bool))
			}

		} else if itemCode, ok := condition["itemCode"]; ok {
			item := adventure.getItem(itemCode.(string))
			if item == nil {
				panic(fmt.Sprintf("Not existing item %s.", itemCode))
			}
			if name, ok := condition["name"]; ok {
				conditionsMet = conditionsMet && (item.Name == name.(string))
			} else if state, ok := condition["state"]; ok {
				conditionsMet = conditionsMet && (item.State == state.(string))
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

// Execute performs the given action on the adventure.
func (adventure *Adventure) ExecuteAction(action *models.Action) error {
	fmt.Printf("EXECUTING ACTION %s\n", action.Code)

	if !*action.Active || adventure.Phase == models.ENDED {
		return nil
	}

	switch action.Operation {
	case models.CHANGE_POSITION:
		if action.PositionCode == "" {
			return fmt.Errorf("missing position in action %s", action.Code)
		}
		// Leaving actions
		adventure.executeLeavingActions()
		adventure.ActualPosition = GetPositionFromPositionList(adventure.Positions, action.PositionCode)
		if adventure.ActualPosition == nil {
			return fmt.Errorf("position not found: %s", action.PositionCode)
		}
		fmt.Printf("POSITION CHANGED TO %s\n", adventure.ActualPosition.Code)
		// Entering actions
		adventure.executeEnteringActions()
		adventure.ActualPosition.Visited = true
		if adventure.ActualPosition.EndPosition {
			adventure.Phase = models.ENDED
		}

	case models.CHANGE_POSITION_DESCRIPTION:
		var position *models.Position
		if action.PositionCode != "" {
			position = GetPositionFromPositionList(adventure.Positions, action.PositionCode)
		} else {
			position = adventure.ActualPosition
		}
		if position != nil {
			position.Description = action.PositionDescription
		}

	case models.CHANGE_POSITION_TEMPORARY_DESCRIPTION:
		var position *models.Position
		if action.PositionCode != "" {
			position = GetPositionFromPositionList(adventure.Positions, action.PositionCode)
		} else {
			position = adventure.ActualPosition
		}
		if position != nil {
			position.TemporaryDescription = action.PositionDescription
		}

	case models.APPEND_POSITION_TEMPORARY_DESCRIPTION:
		var position *models.Position
		if action.PositionCode != "" && action.PositionCode != adventure.ActualPosition.Code {
			position = GetPositionFromPositionList(adventure.Positions, action.PositionCode)
		} else {
			position = adventure.ActualPosition
		}
		if position != nil {
			position.TemporaryDescription = position.Description + " " + action.PositionDescription
		}

	case models.PREPEND_POSITION_TEMPORARY_DESCRIPTION:
		var position *models.Position
		if action.PositionCode != "" && action.PositionCode != adventure.ActualPosition.Code {
			position = GetPositionFromPositionList(adventure.Positions, action.PositionCode)
		} else {
			position = adventure.ActualPosition
		}
		if position != nil {
			position.TemporaryDescription = action.PositionDescription + " " + position.Description
		}
		fmt.Println("PREPEND POSITION TEMPORARY DESCRIPTION ", action.Visible)

	case models.CHANGE_POSITION_VISITED:
		var position *models.Position
		if action.PositionCode != "" {
			position = GetPositionFromPositionList(adventure.Positions, action.PositionCode)
		} else {
			position = adventure.ActualPosition
		}
		if position != nil {
			position.Visited = action.PositionVisited
		}

	case models.CHANGE_ACTION_ACTIVE:
		adventure.changeActionActive(action.ActionCode, *action.ActionActive)

	case models.CHANGE_ACTION_VISIBLE:
		adventure.changeActionVisible(action.ActionCode, *action.ActionVisible)

	case models.CHANGE_ITEM_DESCRIPTION:
		return adventure.changeItemDescription(action.ItemCode, action.ItemDescription)

	case models.CHANGE_ITEM_STATE:
		return adventure.changeItemState(action.ItemCode, action.Value)

	case models.PICK_UP_ITEM:
		item := popItemFromList(&adventure.ActualPosition.Items, action.ItemCode)
		if item == nil {
			return fmt.Errorf("there is no such item at the actual position: %s:%s", adventure.ActualPosition.Code, action.ItemCode)
		}
		adventure.Player.Items = append(adventure.Player.Items, item)

	case models.PUT_DOWN_ITEM:
		item := popItemFromList(&adventure.Player.Items, action.ItemCode)
		if item == nil {
			return fmt.Errorf("there is no such item in the player's inventory: %s", action.ItemCode)
		}
		adventure.ActualPosition.Items = append(adventure.ActualPosition.Items, item)

	case models.MOVE_ITEM:
		item := popItemFromList(&adventure.ActualPosition.Items, action.ItemCode)
		if item == nil {
			return fmt.Errorf("there is no such item at the actual position: %s:%s", adventure.ActualPosition.Code, action.ItemCode)
		}
		newPosition := GetPositionFromPositionList(adventure.Positions, action.PositionCode)
		if newPosition == nil {
			return fmt.Errorf("there is no such position: %s", action.PositionCode)
		}
		newPosition.Items = append(newPosition.Items, item)

	case models.ERASE_ITEM:
		item := popItemFromList(&adventure.ActualPosition.Items, action.ItemCode)
		if item == nil {
			return fmt.Errorf("there is no such item at the actual position: %s:%s", adventure.ActualPosition.Code, action.ItemCode)
		}

	case models.CONDITIONAL:
		// var function map[string]interface{} // func(*Adventure) bool
		// if err := json.Unmarshal([]byte(action.Function), &function); err != nil {
		// 	return err
		// }
		conditional(action.Function, adventure)

	case models.SWITCH_CONDITIONAL:
		// var functions []map[string]interface{} //[]func(*Adventure) bool
		// if err := json.Unmarshal([]byte(action.Functions), &functions); err != nil {
		// 	return err
		// }
		for _, function := range action.Functions {
			if conditional(function, adventure) {
				break
			}
		}

	case models.LIST:
		// var actionCodes []string
		// if err := json.Unmarshal([]byte(action.ActionCodes), &actionCodes); err != nil {
		// 	return err
		// }
		for _, actionCode := range action.ActionCodes {
			executable, err := adventure.findAction(actionCode)
			if err == nil && executable != nil {
				adventure.ExecuteAction(executable)
			}
			if executable.Code == "15" {
				fmt.Println("ACTION AFTER CHANGE", action.Code, *action.Active)
				foundAction, err := adventure.findAction("10")
				if err == nil && foundAction != nil {
					fmt.Println("FOUND ACTION", foundAction.Code, *foundAction.Active)
					fmt.Println("ACTIONS ARE EQUAL", action == foundAction)
				}
			}
		}

	case models.RANDOM:
		// var actionCodes []string
		// if err := json.Unmarshal([]byte(action.ActionCodes), &actionCodes); err != nil {
		// 	return err
		// }
		if len(action.ActionCodes) > 0 {
			randomAction, err := adventure.findAction(action.ActionCodes[rand.Intn(len(action.ActionCodes))])
			if err == nil && randomAction != nil {
				adventure.ExecuteAction(randomAction)
			}
		}

	default:
		// No operation
	}

	return nil
}

func (adventure *Adventure) Start() error {
	action := &models.Action{}
	action.Operation = models.CHANGE_POSITION
	action.PositionCode = adventure.StartPositionCode
	activeVar := true
	action.Active = &activeVar
	adventure.Phase = models.STARTED
	return adventure.ExecuteAction(action)
}

func (adventure *Adventure) Do(actionCode string) error {
	var action *models.Action
	if adventure.ActualPosition != nil {
		action = getActionFromPosition(adventure.ActualPosition, actionCode)
	}
	if action == nil && len(adventure.AvailableActions) > 0 {
		executableAction, err := adventure.findAction(actionCode)
		if err != nil {
			return fmt.Errorf("action not found: %s", actionCode)
		}
		action = executableAction
	}
	err := adventure.ExecuteAction(action)
	if err != nil {
		return fmt.Errorf("error executing action: %s", err.Error())
	}
	return nil
}
