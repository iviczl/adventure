from adventure_engine.action import Action
from adventure_engine.item import Item

class Position():
    def __init__(self, id: str, description: str, available_actions: list[Action] = [], entering_actions: list[Action] = [], leaving_actions: list[Action] = [], items: list[Item] = []):
        self.id = id
        self.description = description
        self.available_actions = available_actions
        self.entering_actions = entering_actions
        self.leaving_actions = leaving_actions
        self.items = items

    def get_serializable(self):
        return {
            "id": self.id, 
            "description": self.description,
            "available_actions": list(map(lambda a: a.get_serializable(), self.available_actions)),
            "entering_actions": list(map(lambda a: a.get_serializable(), self.entering_actions)),
            "leaving_actions": list(map(lambda a: a.get_serializable(), self.leaving_actions)),
            "items":list(map(lambda i: i.get_serializable(), self.items))
        }