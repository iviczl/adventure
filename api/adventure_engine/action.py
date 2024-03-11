from enum import StrEnum
from json import loads
from typing import Optional
from random import choice

from sqlalchemy import ForeignKey
from adventure_engine.adventure_phase import AdventurePhase
from config import db
from adventure_engine.adventure_helper import get_item, get_item_from_player, get_item_from_position, pop_item_from_list, get_position_from_position_list, get_action
from sqlalchemy.orm import Mapped, mapped_column


class Operation(StrEnum):
   CHANGE_POSITION = "cp"
   CHANGE_POSITION_DESCRIPTION = "cpd"
   CHANGE_POSITION_TEMPORARY_DESCRIPTION = "cptd"
   APPEND_POSITION_TEMPORARY_DESCRIPTION = "aptd"
   CHANGE_POSITION_VISITED = "cpv"
   CHANGE_ACTION_ACTIVE = "caa"
   CHANGE_ACTION_VISIBLE = "cav"
   CHANGE_ITEM_DESCRIPTION = "cad" 
   CHANGE_ITEM_STATE = "cis"
  #  CREATE_ITEM = "cri"
   PICK_UP_ITEM = "pui"
   PUT_DOWN_ITEM = "pdi"
   MOVE_ITEM = "moi"
   ERASE_ITEM = "eri"
   CONDITIONAL = "con" # conditional action execution
   SWITCH_CONDITIONAL = "scon" # first matching conditional action execution
   LIST = "ls"
   RANDOM = "ran"
  #  CHANGE_ADVENTURE_STATE = "cas"

class Action(db.Model):
    id = db.Column(db.Integer(), primary_key = True, autoincrement = True)
    code = db.Column(db.String(10), unique = False, nullable = False)
    description = db.Column(db.String(100), unique = False, nullable = False)
    operation = db.Column(db.String(10), unique = False, nullable = False)
    position_code = db.Column(db.String(10), unique = False, nullable = True)
    position_description = db.Column(db.String(500), unique = False, nullable = True)
    position_visited = db.Column(db.Boolean, default = True, nullable = True )
    item_code = db.Column(db.String(10), unique = False, nullable = True)
    action_code = db.Column(db.String(10), unique = False, nullable = True)
    action_codes = db.Column(db.String(100), unique = False, nullable = True)
    action_visible = db.Column(db.Boolean, default = True, nullable = True )
    action_active = db.Column(db.Boolean, default = True, nullable = True )
    value = db.Column(db.String(40), unique = False, nullable = True)
    function = db.Column(db.String(500), unique = False, nullable = True)
    functions = db.Column(db.String(1000), unique = False, nullable = True)
    active = db.Column(db.Boolean, default = True, nullable = False )
    visible = db.Column(db.Boolean, default = True, nullable = False )
    available_actions_position_id: Mapped[Optional[int]] = mapped_column(ForeignKey("position.id"))
    entering_actions_position_id: Mapped[Optional[int]] = mapped_column(ForeignKey("position.id"))
    leaving_actions_position_id: Mapped[Optional[int]] = mapped_column(ForeignKey("position.id"))
    available_actions_adventure_id: Mapped[Optional[int]] = mapped_column(ForeignKey("adventure.id"))

    def get_serializable(self):
        return { 
            "code": self.code,
            "description": self.description,
        }

    @staticmethod
    def _execute_entering_actions(adventure):
        for entering_action in adventure.actual_position.entering_actions:
            position_code = adventure.actual_position.code
            Action.execute(entering_action, adventure)
            if position_code != adventure.actual_position.code:
                break

    @staticmethod
    def _execute_leaving_actions(adventure):
        for leaving_action in adventure.actual_position.leaving_actions:
            position_code = adventure.actual_position.code
            Action.execute(leaving_action, adventure)
            if position_code != adventure.actual_position.code:
                break

    @staticmethod
    def _find_action(adventure, action_code: str):
        action: Action = get_action(adventure, action_code)
        if not action:
            raise ValueError(f"Invalid action_code {action_code}.") 
        return action

    @staticmethod
    def _change_action_active(adventure, action_code: str, value: bool):
        action: Action = Action._find_action(adventure,action_code)
        action.active = value

    @staticmethod
    def _change_action_visible(adventure, action_code: str, value: bool):
        action: Action = Action._find_action(adventure,action_code)
        action.visible = value

    @staticmethod
    def _find_item(adventure,item_code):
        item = get_item(adventure,item_code)
        if not item:
            raise ValueError(f"Invalid item_code {item_code}.")
        return item

    @staticmethod
    def _change_item_description(adventure, item_code: str, value: str):
        item = Action._find_item(adventure,item_code)
        item.description = value

    @staticmethod
    def _change_item_state(adventure, item_code: str, value: str):
        item = Action._find_item(adventure,item_code)
        item.state = value

    @staticmethod
    def _conditional(function, adventure):
        conditions_met = True
        for condition in function["conditions"]:
            if "player_must_have" in condition:
                must_have = condition["player_must_have"]
                if not "item_code" in condition:
                    raise ValueError(f"Missing item code.")
                has_item = bool(get_item_from_player(adventure.player, condition["item_code"]))
                conditions_met = conditions_met and not(must_have ^ has_item)
            elif "position_must_have" in condition:
                print("POSITION_MUST_HAVE")
                must_have = condition["position_must_have"] == "true"
                if not "position_code" in condition:
                    raise ValueError(f"Missing position code.")
                position = get_position_from_position_list(adventure.positions, condition["position_code"])
                if not "item_code" in condition:
                    raise ValueError(f"Missing item code.")
                has_item = bool(get_item_from_position(position, condition["item_code"]))
                conditions_met = conditions_met and not(must_have ^ has_item)
            elif "item_code" in condition:
                print("ITEM_CODE")
                item = get_item(adventure, condition["item_code"] )
                if not item:
                    raise ValueError(f"Not existing item {condition["item_code"]}.")
                if "name" in condition:
                    conditions_met = conditions_met and item.name == condition["name"]
                elif "state" in condition:
                    conditions_met = conditions_met and item.state == condition["state"]
            else:
                print("NOTHING")

            if not conditions_met:
                break

        if conditions_met:
            for action_code in function["action_codes"]:
                executable: Action = Action._find_action(adventure, action_code)
                Action.execute(executable, adventure)
        else:
            if not "else_action_codes" in function:
                return
            
            for action_code in function["else_action_codes"]:
                executable: Action = Action._find_action(adventure, action_code)
                Action.execute(executable, adventure)

        return conditions_met

    @staticmethod
    def execute(action, adventure):
        if not action.active:
            return
        
        match action.operation:
            case Operation.CHANGE_POSITION:
                if action.position_code == None:
                    raise ValueError(f"Not existing position in action {action.code}.")
                # leaving actions
                Action._execute_leaving_actions(adventure)
                adventure.actual_position = get_position_from_position_list(adventure.positions, action.position_code)
                adventure.actual_position.visited = True
                print("POSITION CHANGED TO", adventure.actual_position.code)
                # entering actions
                Action._execute_entering_actions(adventure)
                if adventure.actual_position.end_position:
                    adventure.phase = AdventurePhase.ENDED
            case Operation.CHANGE_POSITION_DESCRIPTION:
                if action.position_code:
                    position = get_position_from_position_list(adventure.positions, action.position_code)
                else:
                    position = adventure.actual_position
                position.description = action.position_description
            case Operation.CHANGE_POSITION_TEMPORARY_DESCRIPTION:
                if action.position_code:
                    position = get_position_from_position_list(adventure.positions, action.position_code)
                else:
                    position = adventure.actual_position
                position.temporary_description = action.position_description
            case Operation.APPEND_POSITION_TEMPORARY_DESCRIPTION:
                if action.position_code:
                    position = get_position_from_position_list(adventure.positions, action.position_code)
                else:
                    position = adventure.actual_position
                position.temporary_description = position.description + " " + action.position_description
            case Operation.CHANGE_POSITION_VISITED:
                if action.position_code:
                    position = get_position_from_position_list(adventure.positions, action.position_code)
                else:
                    position = adventure.actual_position
                position.visited = action.position_visited
            case Operation.CHANGE_ACTION_ACTIVE:
                Action._change_action_active(adventure, action.action_code if action.action_code else action.code, action.action_active)
            case Operation.CHANGE_ACTION_VISIBLE:
                Action._change_action_visible(adventure, action.action_code if action.action_code else action.code, action.action_visible)
            case Operation.CHANGE_ITEM_DESCRIPTION:
                Action._change_item_description(adventure, action.item_code, action.item_description)
            case Operation.CHANGE_ITEM_STATE:
                Action._change_item_state(adventure, action.item_code, action.value)
            case Operation.PICK_UP_ITEM:
                item = pop_item_from_list(adventure.actual_position.items, action.item_code)
                if not item:
                    raise ValueError(f"There is no such item at the actual position: {adventure.actual_position.code}:{action.item_code}")
                adventure.player.items.append(item)
            case Operation.PUT_DOWN_ITEM:
                item = pop_item_from_list(adventure.player.items, action.item_code)
                if not item:
                    raise ValueError(f"There is no such item at the actual position: {adventure.actual_position.code}:{action.item_code}")
                adventure.actual_position.items.append(item)
            case Operation.MOVE_ITEM:
                item = pop_item_from_list(adventure.actual_position.items, action.item_code)
                if not item:
                    raise ValueError(f"There is no such item at the actual position: {adventure.actual_position.code}:{action.item_code}")
                new_position = get_position_from_position_list(adventure.positions, action.position_code)
                if not new_position:
                    raise ValueError(f"There is no such position: {action.position_code}")
                new_position.items.append(item)
            case Operation.ERASE_ITEM:
                item = pop_item_from_list(adventure.actual_position.items, action.item_code)
                if not item:
                    raise ValueError(f"There is no such item at the actual position: {adventure.actual_position.code}:{action.item_code}")
            case Operation.CONDITIONAL:
                function = loads(action.function)
                Action._conditional(function, adventure)
            case Operation.SWITCH_CONDITIONAL:
                functions = loads(action.functions)
                for function in functions:
                    executed = Action._conditional(function, adventure)
                    if executed:
                        break
            case Operation.LIST:
                action_codes = loads(action.action_codes)
                for action_code in action_codes:
                    executable = Action._find_action(adventure, action_code)
                    Action.execute(executable, adventure)
            case Operation.RANDOM:
                action_codes = loads(action.action_codes)
                action = Action._find_action(adventure, choice(action_codes))
                Action.execute(action,adventure)
            case _:
                pass