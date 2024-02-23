# from adventure_engine.adventure import Adventure
from enum import StrEnum

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
    id = db.Column(db.String(10), primary_key = True)
    description = db.Column(db.String(100), unique = False, nullable = False)
    operation = db.Column(db.String(10), unique = False, nullable = False)
    position_id = db.Column(db.String(10), unique = False, nullable = True)
    value = db.Column(db.String(40), unique = False, nullable = True)
    item_id = db.Column(db.String(10), unique = False, nullable = True)
    function = db.Column(db.String(500), unique = False, nullable = True)
    active = db.Column(db.Boolean, default = True, nullable = True )
    available_actions_position_id: Mapped[str] = mapped_column(ForeignKey("position.id"))
    available_actions_adventure_id: Mapped[str] = mapped_column(ForeignKey("adventure.id"))

    def get_serializable(self):
        return { 
            "id": self.id, 
            "description": self.description,
            "operation": str(self.operation),
            "position_id": self.position_id,
            "value": self.value,
            "item_id": self.item_id,
            "function": self.function,
            "active": self.active
        }

    @staticmethod
    def execute(action, adventure):
        if action.operation == Operation.CHANGE_POSITION:
            # leaving actions
            for leaving_action in adventure.actual_position.leaving_actions:
                if leaving_action.position_id == None:
                    raise f"Not existing position in action {leaving_action.id}."
                Action.execute(leaving_action, adventure)
            adventure.actual_position = get_position_from_position_list(adventure.positions, action.position_id)
            # entering actions
            for entering_action in adventure.actual_position.entering_actions:
                if entering_action.position_id == None:
                    raise f"Not existing position in action {entering_action.id}."
                Action.execute(entering_action, adventure)
        elif action.operation == Operation.CHANGE_ITEM_STATE:
            item = get_item(adventure, action.item_id )
            if item == None:
               raise f"Not existing item {action.item_id}." 
            item.state = action.value
        elif action.operation == Operation.CONDITIONAL:
            condition_met = True
            for condition_item in action.function["condition"]:
                if "item_id" in condition_item:
                    item = get_item(adventure, condition_item["item_id"] )
                    if "name" in condition_item:
                        condition_met = condition_met and item.name == condition_item["name"]
                    elif "state" in condition_item:
                        condition_met = condition_met and item.state == condition_item["state"]
                
                if not condition_met:
                    break
            else:
                for change_item in action.function["change"]:
                    pass


                    
          # "function": {
          #   "condition": [{ "item_id": 0, "state": "not lit" }],
          #   "change": [
          #     {
          #       "position_id": "-2"
          #     }
          #   ]
          # },
