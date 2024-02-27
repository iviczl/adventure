# from adventure_engine.adventure import Adventure
from enum import StrEnum
from json import loads
from typing import Optional

from sqlalchemy import ForeignKey
from config import db
from adventure_engine.adventure_helper import get_item, get_position_from_position_list
from sqlalchemy.orm import Mapped, mapped_column, relationship

# from adventure_engine.position import Position

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
    value = db.Column(db.String(40), unique = False, nullable = True)
    item_code = db.Column(db.String(10), unique = False, nullable = True)
    function = db.Column(db.String(500), unique = False, nullable = True)
    active = db.Column(db.Boolean, default = True, nullable = False )
    available_actions_position_id: Mapped[Optional[int]] = mapped_column(ForeignKey("position.id"))
    available_actions_adventure_id: Mapped[Optional[int]] = mapped_column(ForeignKey("adventure.id"))

    def get_serializable(self):
        return { 
            "id": self.id, 
            "code": self.code,
            "description": self.description,
            "operation": str(self.operation),
            "position_code": self.position_code,
            "value": self.value,
            "item_code": self.item_code,
            "function": self.function,
            "active": self.active
        }

    @staticmethod
    def execute(action, adventure):
        if action.operation == Operation.CHANGE_POSITION:
            # leaving actions
            for leaving_action in adventure.actual_position.leaving_actions:
                if leaving_action.position_code == None:
                    raise f"Not existing position in action {leaving_action.id}."
                Action.execute(leaving_action, adventure)
            adventure.actual_position = get_position_from_position_list(adventure.positions, action.position_code)
            # entering actions
            for entering_action in adventure.actual_position.entering_actions:
                if entering_action.position_code == None:
                    raise f"Not existing position in action {entering_action.id}."
                Action.execute(entering_action, adventure)
        elif action.operation == Operation.CHANGE_ITEM_STATE:
            item = get_item(adventure, action.item_code )
            if item == None:
               raise f"Not existing item {action.item_code}." 
            item.state = action.value
        elif action.operation == Operation.CONDITIONAL:
            condition_met = True
            function = loads(action.function)
            for condition_item in function["condition"]:
                if "item_id" in condition_item:
                    item = get_item(adventure, condition_item["item_id"] )
                    if "name" in condition_item:
                        condition_met = condition_met and item.name == condition_item["name"]
                    elif "state" in condition_item:
                        condition_met = condition_met and item.state == condition_item["state"]
                
                if not condition_met:
                    break
            else:
                for change_item in function["change"]:
                    pass


                    
          # "function": {
          #   "condition": [{ "item_id": 0, "state": "not lit" }],
          #   "change": [
          #     {
          #       "position_id": "-2"
          #     }
          #   ]
          # },

