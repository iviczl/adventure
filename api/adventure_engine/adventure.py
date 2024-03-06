from typing import List, Optional

from sqlalchemy import ForeignKey, Integer
from adventure_engine.action import Action
from adventure_engine.item import Item
from adventure_engine.player import Player
from adventure_engine.position import Position
from adventure_engine.adventure_helper import *
from enum import IntEnum
from sqlalchemy.orm import Mapped, relationship, mapped_column
from config import db

class AdventurePhase(IntEnum):
    NOT_LOADED = "-1"
    LOADED = "0"
    STARTED = "1"
    ENDED = "2"

class Adventure(db.Model):
    id = db.Column(db.Integer(), primary_key = True, autoincrement = True)
    code = db.Column(db.String(10), unique = False, nullable = False)
    title: str = db.Column(db.String(60), unique = False, nullable = False)
    description = db.Column(db.String(100), unique = False, nullable = False)
    player: Mapped["Player"] = relationship("Player", back_populates="adventure")
    # position_id: Mapped[int] = mapped_column(ForeignKey("position.id"))
    positions:  Mapped[List["Position"]] = relationship(foreign_keys="[Position.adventure_id]")# "Position",foreign_keys=[position_id] back_populates="positions_adventure")# overlaps="adventure", 
    # actual_position_id: Mapped[int] = mapped_column(ForeignKey("position.id"))
    actual_position: Mapped["Position"] = relationship( foreign_keys="[Position.actual_position_adventure_id]", lazy="immediate") # , post_update= True, back_populates="actual_position_adventure", ,overlaps="adventure, positions", uselist= False
    phase: Mapped[str] = db.Column(db.String(), server_default=str(AdventurePhase.NOT_LOADED))
    available_actions: Mapped[Optional[List["Action"]]] = relationship(foreign_keys="[Action.available_actions_adventure_id]")

    # def get_serializable(self):
    #     return {
    #         "id": self.id,
    #         "title": self.title,
    #         "description": self.description,
    #         # "positions": list(map(lambda p: p.get_serializable() ,self.positions)),
    #         # # "actual_position": self.actual_position.get_serializable(),
    #         "actual_position_id": self.actual_position_id,
    #         "player": self.player.get_serializable(),
    #         "phase": str(self.phase),
    #         "available_actions": list(map(lambda a: a.get_serializable(), self.available_actions))
    #     }
    
    @staticmethod
    def build(data, adventure = None):
        if adventure == None:
            adventure = Adventure()
        adventure.code = data["code"]
        adventure.player =  Player(data["player"]["name"])
        adventure.title = data["title"]
        adventure.description = data["description"]
        adventure.phase = AdventurePhase.LOADED
        if "items" in data["player"]:
            # adventure.player.items = list(map(lambda i: Item(**i), prepare_items(data["player"]["items"])))
            adventure.player.items = list(map(lambda i: Item(**i), data["player"]["items"]))
        adventure.available_actions = list(map(lambda a: Action(**a), prepare_actions(data["available_actions"])))

        for position_data in data["positions"]:
            position = Position()
            position.code = position_data["code"]
            position.description = position_data["description"]
            adventure.positions.append(position)
            if "items" in position_data:
                position.items = list(map(lambda i: Item(**i), position_data["items"]))
            if "available_actions" in position_data:
                position.available_actions = list(map(lambda a: Action(**a), prepare_actions(position_data["available_actions"])))
            if "entering_actions" in position_data:
                position.entering_actions = list(map(lambda a: Action(**a), prepare_actions(position_data["entering_actions"])))
            if "leaving_actions" in position_data:
                position.leaving_actions = list(map(lambda a: Action(**a), prepare_actions(position_data["leaving_actions"])))


        adventure.actual_position = get_position_from_position_list(adventure.positions, "0")
        return adventure

    @staticmethod
    def load(file_path: str, adventure):
        rawData = load_json(file_path)
        return Adventure.build(rawData, adventure)

    def do(self, action_code: str):
        action: Action = None
        if self.actual_position != None:
            action = get_action_from_position(self.actual_position, action_code)
        if action == None and len(self.available_actions) > 0:
            action = next((a for a in self.available_actions if a.code == action_code), None)
        if action == None:
            raise RuntimeError(f"Invalid action code in position: {self.actual_position.code}:{action_code}")
        print("EXECUTING ACTION",action.code)
        Action.execute(action, self)