# from adventure_engine.adventure import Adventure
from enum import StrEnum
from adventure_engine.adventure_helper import get_item

# from adventure_engine.position import Position

class Operation(StrEnum):
   CHANGE_POSITION = "cp"
   CHANGE_ITEM_STATE = "cis"
  #  CHANGE_ADVENTURE_STATE = "cas"
   CONDITIONAL = "con"

class Action():
    def __init__(self, id: str, description: str, operation: Operation, position_id: str = None, value: str = None, item_id: str = None, function: dict = None, active: bool = True):
        self.id = id
        self.description = description
        self.operation = operation
        self.position_id = position_id
        self.value = value
        self.item_id = item_id
        self.function = function
        self.active = active

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

    def execute(self, adventure):
        if self.operation == Operation.CHANGE_POSITION:
            # leaving actions
            for action in adventure.actual_position.leaving_actions:
                if action.position == None:
                    raise f"Not existing position in action {action.id}."
                action.execute(adventure)
            adventure.actual_position = self.position
            # entering actions
            for action in adventure.actual_position.entering_actions:
                if action.position == None:
                    raise f"Not existing position in action {action.id}."
                action.execute(adventure)
        elif self.operation == Operation.CHANGE_ITEM_STATE:
            item = get_item(adventure, self.item.id )
            if item == None:
               raise f"Not existing item {self.item.id}." 
            item.state = self.value
        elif self.operation == Operation.CONDITIONAL:
            condition_met = True
            for condition_item in self.function.condition:
                if "item_id" in condition_item:
                    item = get_item(adventure, condition_item["item_id"] )
                    if "name" in condition_item:
                        condition_met = condition_met and item.name == condition_item["name"]
                    elif "state" in condition_item:
                        condition_met = condition_met and item.state == condition_item["state"]
                
                if not condition_met:
                    break
            else:
                for change_item in self.function.change:
                    pass


                    
          # "function": {
          #   "condition": [{ "item_id": 0, "state": "not lit" }],
          #   "change": [
          #     {
          #       "position_id": "-2"
          #     }
          #   ]
          # },

