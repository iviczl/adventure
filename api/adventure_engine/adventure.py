# import json
from typing import List, Optional
from adventure_engine.action import Action
from adventure_engine.item import Item
from adventure_engine.player import Player
from adventure_engine.position import Position
from adventure_engine.adventure_helper import *
from enum import IntEnum
from sqlalchemy.orm import Mapped, relationship
from config import db

class AdventurePhase(IntEnum):
    NOT_LOADED = "-1"
    LOADED = "0"
    STARTED = "1"
    ENDED = "2"

class Adventure(db.Model):
    id = db.Column(db.String(10), primary_key = True)
    title: str = db.Column(db.String(60), unique = False, nullable = False)
    description = db.Column(db.String(100), unique = False, nullable = False)
    player: Mapped["Player"] = relationship(back_populates="adventure")
    positions:  Mapped[List["Position"]] = relationship()
    actual_position: Mapped[Optional["Position"]] = relationship(overlaps="positions")
    phase: Mapped[str] = db.Column(db.String(), server_default=str(AdventurePhase.NOT_LOADED))
    available_actions: Mapped[Optional[List["Action"]]] = relationship()

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
    
    @staticmethod
    def build(data, adventure = None):
        if adventure == None:
            adventure = Adventure()
        adventure.id = data["id"]
        adventure.player =  Player(data["player"]["name"])
        adventure.title = data["title"]
        adventure.description = data["description"]
        adventure.phase = AdventurePhase.LOADED
        if "items" in data["player"]:
            # print("!player items!",data["player"]["items"])
            # return
            adventure.player.items = list(map(lambda i: Item(**i), data["player"]["items"]))
        adventure.available_actions = list(map(lambda a: Action(**a), data["available_actions"]))

        for position_data in data["positions"]:
            position = Position()
            position.id = position_data["id"]
            position.description = position_data["description"]
            adventure.positions.append(position)
            if "items" in position_data:
                position.items = list(map(lambda i: Item(**i), position_data["items"]))
            if "available_actions" in position_data:
                position.available_actions = list(map(lambda a: Action(**a), position_data["available_actions"]))
            if "entering_actions" in position_data:
                position.entering_actions = list(map(lambda a: Action(**a), position_data["entering_actions"]))
            if "leaving_actions" in position_data:
                position.leaving_actions = list(map(lambda a: Action(**a), position_data["leaving_actions"]))

        actual_position_id = "0"
        if "actual_position" in data and data["actual_position"] != None:
            actual_position_id = data["actual_position"]["id"]
        adventure.actual_position = get_position_from_position_list(adventure.positions, actual_position_id)
        return adventure

    @staticmethod
    def load(file_path: str, adventure):
        rawData = load_json(file_path)
        return Adventure.build(rawData, adventure)

    def do(self, action_id: str):
        action: Action = None
        if self.actual_position != None:
            action = get_action_from_position(self.actual_position, action_id)
        if action == None and len(self.available_actions) > 0:
            action = next((a for a in self.available_actions if a.id == action_id), None)
        if action == None:
            raise RuntimeError(f"Invalid action id in position: {self.actual_position.id + ":" + action_id}")
        print("EXECUTING ACTION",action.id)
        Action.execute(action, self)