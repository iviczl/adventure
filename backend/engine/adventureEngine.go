package engine

import (
	"fmt"
	"math/rand"
	"text-adventure/constants"
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

func (adventure *Adventure) findAction(actionCode string) (*models.Action, error) {
	action := adventure.GetAction(actionCode)
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
	item := adventure.GetItem(itemCode)
	if item == nil {
		return nil, fmt.Errorf("invalid itemCode %s", itemCode)
	}
	return item, nil
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

// Execute performs the given action on the adventure.
func (adventure *Adventure) ExecuteAction(action *models.Action) error {
	fmt.Printf("EXECUTING ACTION %s\n", action.Code)

	if !*action.Active || adventure.Phase == models.ENDED {
		return nil
	}

	switch action.Operation {
	case constants.CHANGE_POSITION:
		if action.PositionCode == "" {
			return fmt.Errorf("missing position in action %s", action.Code)
		}
		// Leaving actions
		adventure.executeLeavingActions()
		adventure.ActualPosition = adventure.GetPosition(action.PositionCode)
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

	case constants.CHANGE_POSITION_DESCRIPTION:
		var position *models.Position
		if action.PositionCode != "" {
			position = adventure.GetPosition(action.PositionCode)
		} else {
			position = adventure.ActualPosition
		}
		if position != nil {
			position.Description = action.PositionDescription
		}

	case constants.CHANGE_POSITION_TEMPORARY_DESCRIPTION:
		var position *models.Position
		if action.PositionCode != "" {
			position = adventure.GetPosition(action.PositionCode)
		} else {
			position = adventure.ActualPosition
		}
		if position != nil {
			position.TemporaryDescription = action.PositionDescription
		}

	case constants.APPEND_POSITION_TEMPORARY_DESCRIPTION:
		var position *models.Position
		if action.PositionCode != "" && action.PositionCode != adventure.ActualPosition.Code {
			position = adventure.GetPosition(action.PositionCode)
		} else {
			position = adventure.ActualPosition
		}
		if position != nil {
			position.TemporaryDescription = position.Description + " " + action.PositionDescription
		}

	case constants.PREPEND_POSITION_TEMPORARY_DESCRIPTION:
		var position *models.Position
		if action.PositionCode != "" && action.PositionCode != adventure.ActualPosition.Code {
			position = adventure.GetPosition(action.PositionCode)
		} else {
			position = adventure.ActualPosition
		}
		if position != nil {
			position.TemporaryDescription = action.PositionDescription + " " + position.Description
		}
		fmt.Println("PREPEND POSITION TEMPORARY DESCRIPTION ", action.Visible)

	case constants.CHANGE_POSITION_VISITED:
		var position *models.Position
		if action.PositionCode != "" {
			position = adventure.GetPosition(action.PositionCode)
		} else {
			position = adventure.ActualPosition
		}
		if position != nil {
			position.Visited = action.PositionVisited
		}

	case constants.CHANGE_ACTION_ACTIVE:
		adventure.changeActionActive(action.ActionCode, *action.ActionActive)

	case constants.CHANGE_ACTION_VISIBLE:
		adventure.changeActionVisible(action.ActionCode, *action.ActionVisible)

	case constants.CHANGE_ITEM_DESCRIPTION:
		return adventure.changeItemDescription(action.ItemCode, action.ItemDescription)

	case constants.CHANGE_ITEM_STATE:
		return adventure.changeItemState(action.ItemCode, action.Value)

	case constants.PICK_UP_ITEM:
		item := popItemFromList(&adventure.ActualPosition.Items, action.ItemCode)
		if item == nil {
			return fmt.Errorf("there is no such item at the actual position: %s:%s", adventure.ActualPosition.Code, action.ItemCode)
		}
		adventure.Player.Items = append(adventure.Player.Items, item)

	case constants.PUT_DOWN_ITEM:
		item := popItemFromList(&adventure.Player.Items, action.ItemCode)
		if item == nil {
			return fmt.Errorf("there is no such item in the player's inventory: %s", action.ItemCode)
		}
		adventure.ActualPosition.Items = append(adventure.ActualPosition.Items, item)

	case constants.MOVE_ITEM:
		item := popItemFromList(&adventure.ActualPosition.Items, action.ItemCode)
		if item == nil {
			return fmt.Errorf("there is no such item at the actual position: %s:%s", adventure.ActualPosition.Code, action.ItemCode)
		}
		newPosition := adventure.GetPosition(action.PositionCode)
		if newPosition == nil {
			return fmt.Errorf("there is no such position: %s", action.PositionCode)
		}
		newPosition.Items = append(newPosition.Items, item)

	case constants.ERASE_ITEM:
		item := popItemFromList(&adventure.ActualPosition.Items, action.ItemCode)
		if item == nil {
			return fmt.Errorf("there is no such item at the actual position: %s:%s", adventure.ActualPosition.Code, action.ItemCode)
		}

	case constants.CONDITIONAL:
		conditional(action.Function, adventure)

	case constants.SWITCH_CONDITIONAL:
		for _, function := range action.Functions {
			if conditional(function, adventure) {
				break
			}
		}

	case constants.LIST:
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

	case constants.RANDOM:
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
	action.Operation = constants.CHANGE_POSITION
	action.PositionCode = adventure.StartPositionCode
	activeVar := true
	action.Active = &activeVar
	adventure.Phase = models.STARTED
	return adventure.ExecuteAction(action)
}

func (adventure *Adventure) Do(actionCode string) error {
	var action *models.Action
	if adventure.ActualPosition != nil {
		action = adventure.ActualPosition.GetAction(actionCode)
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
