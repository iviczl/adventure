package engine

import (
	"fmt"
	"math/rand"
	"text-adventure/models"
)

func FightWithNpc(adventure *Adventure, action *models.Action) error {
	won, err := autoFight(adventure, action.NpcCode)
	if err != nil {
		return fmt.Errorf("error during fight: %s", err.Error())
	}
	if won {
		for _, actionCode := range action.ActionCodes {
			executable, err := adventure.findAction(actionCode)
			if err != nil {
				return fmt.Errorf("error finding action %s: %s", actionCode, err.Error())
			}
			if executable != nil {
				err = adventure.ExecuteAction(executable)
				if err != nil {
					return fmt.Errorf("error executing action %s: %s", actionCode, err.Error())
				}
			}
		}
	} else {
		for _, actionCode := range action.ElseActionCodes {
			executable, err := adventure.findAction(actionCode)
			if err != nil {
				return fmt.Errorf("error finding else action %s: %s", actionCode, err.Error())
			}
			if executable != nil {
				err = adventure.ExecuteAction(executable)
				if err != nil {
					return fmt.Errorf("error executing else action %s: %s", actionCode, err.Error())
				}
			}
		}
	}
	return nil
}

func autoFight(adventure *Adventure, npcCode string) (bool, error) {
	npc := adventure.GetNpc(npcCode)
	if npc == nil {
		return false, fmt.Errorf("npc not found: %s", npcCode)
	}
	npcVitalAttribute := npc.GetVitalAttribute()
	if npcVitalAttribute == nil {
		return false, fmt.Errorf("npc %s has no vital attribute", npcCode)
	}
	player := adventure.Player
	playerVitalAttribute := player.GetVitalAttribute()
	if playerVitalAttribute == nil {
		return false, fmt.Errorf("player has no vital attribute")
	}

	for npcVitalAttribute.Value > 0 && playerVitalAttribute.Value > 0 {
		var playerChosenAttackAbility *models.Ability
		var playerCostAttribute *models.Attribute
		for _, ability := range player.Abilities {
			playerCostAttribute = player.GetAttribute(ability.Attribute)
			if !ability.Active || ability.TargetDamage <= 0 || ability.Cost >= playerCostAttribute.Value {
				continue
			}
			if playerChosenAttackAbility == nil || ability.Chance > playerChosenAttackAbility.Chance {
				playerChosenAttackAbility = ability
			}
		}
		if playerChosenAttackAbility == nil {
			return false, fmt.Errorf("player has no active abilities with sufficient cost")
		}
		var npcChosenAttackAbility *models.Ability
		var npcCostAttribute *models.Attribute
		for _, ability := range npc.Abilities {
			npcCostAttribute = npc.GetAttribute(ability.Attribute)
			if !ability.Active || ability.TargetDamage <= 0 || ability.Cost >= npcCostAttribute.Value {
				continue
			}
			if npcChosenAttackAbility == nil || ability.Chance > npcChosenAttackAbility.Chance {
				npcChosenAttackAbility = ability
			}
		}
		if npcChosenAttackAbility == nil {
			return false, fmt.Errorf("npc has no active abilities with sufficient cost")
		}
		playerHitConnects := playerChosenAttackAbility.Chance == 1 || rand.Float32() <= playerChosenAttackAbility.Chance
		npcHitConnects := npcChosenAttackAbility.Chance == 1 || rand.Float32() <= npcChosenAttackAbility.Chance
		if !playerHitConnects {
			if npcHitConnects {
				playerVitalAttribute.Value = min(playerVitalAttribute.Value-npcChosenAttackAbility.TargetDamage, 0)
				npcCostAttribute.Value -= npcChosenAttackAbility.Cost
			}
		} else if !npcHitConnects {
			if playerHitConnects {
				npcVitalAttribute.Value = min(npcVitalAttribute.Value-playerChosenAttackAbility.TargetDamage, 0)
				playerCostAttribute.Value -= playerChosenAttackAbility.Cost
			}
		} else { // Both hits connect
			playerVitalValue := min(playerVitalAttribute.Value-npcChosenAttackAbility.TargetDamage, 0)
			npcVitalValue := min(npcVitalAttribute.Value-playerChosenAttackAbility.TargetDamage, 0)
			if playerVitalValue == 0 && npcVitalValue == 0 { // Both would be dead
				playerHits := rand.Intn(2) == 0 // Randomly decide who hits, this way only one will die
				if playerHits {
					playerVitalValue = playerVitalAttribute.Value
					playerCostAttribute.Value -= playerChosenAttackAbility.Cost
				} else { // NPC hits
					npcVitalValue = npcVitalAttribute.Value
					npcCostAttribute.Value -= npcChosenAttackAbility.Cost
				}
			}
			playerVitalAttribute.Value = playerVitalValue
			npcVitalAttribute.Value = npcVitalValue
		}
	}
	if npcVitalAttribute.Value == 0 {
		npc.State = npc.LoseState
	}
	return true, nil
}
