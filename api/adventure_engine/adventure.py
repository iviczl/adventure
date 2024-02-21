# import json
from adventure_engine.action import Action
from adventure_engine.item import Item
from adventure_engine.player import Player
from adventure_engine.position import Position
from adventure_engine.adventure_helper import *
from enum import IntEnum

class AdventurePhase(IntEnum):
    NOT_LOADED = -1
    LOADED = 0
    STARTED = 1
    ENDED = 2

class Adventure():
    id: str =  ""
    title: str = ""
    description: str = ""
    player: Player
    positions: list[Position] = []
    actual_position: Position = None
    phase: AdventurePhase = AdventurePhase.NOT_LOADED
    available_actions: list[Action] = []

    def __init__(self, player_name: str):
        # self.id = ""
        # self.title = ""
        # self.description = ""
        self.player = Player(player_name) 
        # self.positions: list[Position] = []
        # self.actual_position: Position = None
        # self.phase = AdventurePhase.NOT_LOADED
        # self.available_actions: list[Action] = []

    def get_serializable(self):
        return {
            "id": self.id,
            "title": self.title,
            "description": self.description,
            "positions": list(map(lambda p: p.get_serializable() ,self.positions)),
            "actual_position": self.actual_position.get_serializable(),
            "player": self.player.get_serializable(),
            "phase": str(self.phase),
            "available_actions": list(map(lambda a: a.get_serializable(), self.available_actions))
        }

    def load(self, file_path: str):
        data = load_adventure_json(file_path)

        self.id = data["id"]
        self.title = data["title"]
        self.description = data["description"]
        self.phase = AdventurePhase.LOADED
        if "items" in data["player"]:
            # print("!player items!",data["player"]["items"])
            # return
            self.player.items = list(map(lambda i: Item(**i), data["player"]["items"]))
        self.available_actions = list(map(lambda a: Action(**a), data["available_actions"]))

        for position_data in data["positions"]:
            position = Position(position_data["id"], position_data["description"])
            self.positions.append(position)
            if "items" in position_data:
                position.items = list(map(lambda i: Item(**i), position_data["items"]))
            if "available_actions" in position_data:
                position.available_actions = list(map(lambda a: Action(**a), position_data["available_actions"]))
            if "entering_actions" in position_data:
                position.entering_actions = list(map(lambda a: Action(**a), position_data["entering_actions"]))
            if "leaving_actions" in position_data:
                position.leaving_actions = list(map(lambda a: Action(**a), position_data["leaving_actions"]))

        self.actual_position = get_position_from_position_list(self.positions, "0")

    def do(self, action_id: str):
        action: Action = None
        if self.actual_position != None:
            action = get_action_from_position_list(self.actual_position.available_actions, action_id)
        if action == None and len(self.available_actions) > 0:
            action = get_action_from_position_list(self.available_actions, action_id)
        if action == None:
            raise RuntimeError(f"Invalid action id.")
        action.execute(self)