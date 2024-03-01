from enum import StrEnum
from json import loads
from typing import Optional

from sqlalchemy import ForeignKey
from config import db
from adventure_engine.adventure_helper import get_item, get_position_from_position_list, get_action
from sqlalchemy.orm import Mapped, mapped_column, relationship


class Operation(StrEnum):
   CHANGE_POSITION = "cp"
   CHANGE_ITEM_STATE = "cis"
  #  CHANGE_ADVENTURE_STATE = "cas"
   CONDITIONAL = "con"

class Action(db.Model):
    id = db.Column(db.Integer(), primary_key = True, autoincrement = True)
    code = db.Column(db.String(10), unique = False, nullable = False)
    description = db.Column(db.String(100), unique = False, nullable = False)
    operation = db.Column(db.String(10), unique = False, nullable = False)
    position_code = db.Column(db.String(10), unique = False, nullable = True)
    item_code = db.Column(db.String(10), unique = False, nullable = True)
    value = db.Column(db.String(40), unique = False, nullable = True)
    function = db.Column(db.String(500), unique = False, nullable = True)
    active = db.Column(db.Boolean, default = True, nullable = False )
    available_actions_position_id: Mapped[Optional[int]] = mapped_column(ForeignKey("position.id"))
    entering_actions_position_id: Mapped[Optional[int]] = mapped_column(ForeignKey("position.id"))
    leaving_actions_position_id: Mapped[Optional[int]] = mapped_column(ForeignKey("position.id"))
    # available_actions_position: Mapped[Optional["Position"]] = relationship("Position", back_populates= "available_actions")
    # entering_actions_position: Mapped[Optional["Position"]] = relationship("Position", back_populates= "entering_actions")
    # leaving_actions_position: Mapped[Optional["Position"]] = relationship("Position", back_populates= "leaving_actions")
    # available_actions_adventure_id: Mapped[Optional[int]] = mapped_column(ForeignKey("adventure.id"))
    available_actions_adventure_id: Mapped[Optional[int]] = mapped_column(ForeignKey("adventure.id"))

    def get_serializable(self):
        return { 
            # "id": self.id, 
            "code": self.code,
            "description": self.description,
            # "operation": str(self.operation),
            # "position_code": self.position_code,
            # "value": self.value,
            # "item_code": self.item_code,
            # "function": self.function,
            # "active": self.active
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
    def _change_item_description(adventure, item_code: str, value: str):
        item = get_item(adventure,item_code)
        if item == None:
            raise ValueError(f"Invalid item_code {item_code}.") 
        item.description = value
        Action._execute_entering_actions(adventure)

    @staticmethod
    def _change_item_state(adventure, item_code: str, value: str):
        item = get_item(adventure,item_code )
        if item == None:
            raise ValueError(f"Invalid item_code {item_code}.") 
        item.state = value
        Action._execute_entering_actions(adventure)

    @staticmethod
    def execute(action, adventure):
        if action.operation == Operation.CHANGE_POSITION:
            if action.position_code == None:
                raise ValueError(f"Not existing position in action {action.code}.")
            # leaving actions
            Action._execute_leaving_actions(adventure)
            adventure.actual_position = get_position_from_position_list(adventure.positions, action.position_code)
            adventure.actual_position.visited = True
            print("POSITION CHANGED TO", adventure.actual_position.code)
            # entering actions
            Action._execute_entering_actions(adventure)
        elif action.operation == Operation.CHANGE_ITEM_STATE:
            Action._change_item_state(adventure, action.item_code, action.value)
        elif action.operation == Operation.CONDITIONAL:
            condition_met = True
            function = loads(action.function)
            item = None
            for condition_item in function["condition"]:
                if "item_code" in condition_item:
                    item = get_item(adventure, condition_item["item_code"] )
                    if not item:
                        raise ValueError(f"Not existing item {condition_item["item_code"]}.")
                    if "name" in condition_item:
                        condition_met = condition_met and item.name == condition_item["name"]
                    elif "state" in condition_item:
                        condition_met = condition_met and item.state == condition_item["state"]
                
                if not condition_met:
                    break
            else:
                for change_item in function["change"]:
                    if "position_code" in change_item:
                        position = get_position_from_position_list(adventure.positions, change_item["position_code"])
                        if not position:
                            raise ValueError(f"Invalid position code: {change_item["position_code"]}")
                        if "description" in change_item:
                            position.description = change_item["description"]
                        if "visited" in change_item:
                            position.visited = change_item["visited"]
                        if len(change_item.keys()) == 1: # there is no changeable property given, so it must be a jump to the position
                            adventure.actual_position = position
                            print("POSITION CHANGED TO", adventure.actual_position.code)
                    elif "action_code" in change_item:
                        changeable_action = get_action(adventure, change_item["action_code"])
                        if not changeable_action:
                            raise ValueError(f"Invalid action_code: {change_item["action_code"]}")
                        if "active" in change_item:
                            changeable_action.active = change_item["active"]
                        else:
                            raise ValueError(f"Invalid action property: {change_item["action_code"]}")
                    elif "item_code" in change_item:
                        if "description" in change_item:
                            Action._change_item_description(adventure, change_item["item_code"], change_item["description"] )
                        if "state" in change_item:
                            Action._change_item_state(adventure, change_item["item_code"], change_item["state"])
                    else:
                        raise ValueError(f"Invalid change identifier in the action: {str(**change_item)}")
                        


                    
          # "function": {
          #   "condition": [{ "item_code": 0, "state": "not lit" }],
          #   "change": [
          #     {
          #       "position_code": "-2"
          #     }
          #   ]
          # },

