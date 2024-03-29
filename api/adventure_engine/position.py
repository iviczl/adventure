from typing import List, Optional

from sqlalchemy import ForeignKey, event
from adventure_engine.action import Action
from adventure_engine.item import Item
from sqlalchemy.orm import Mapped, relationship, mapped_column
from config import db

class Position(db.Model):
    id = db.Column(db.Integer(), primary_key = True, autoincrement = True)
    code = db.Column(db.String(10), unique = False, nullable = False)
    description = db.Column(db.String(500), unique = False, nullable = False)
    visited: Mapped[bool] = mapped_column(default= False)
    end_position: Mapped[bool] = mapped_column(default= False)
    available_actions: Mapped[List["Action"]] = relationship(foreign_keys="[Action.available_actions_position_id]")
    entering_actions: Mapped[List["Action"]] = relationship(foreign_keys="[Action.entering_actions_position_id]")
    leaving_actions: Mapped[List["Action"]] = relationship(foreign_keys="[Action.leaving_actions_position_id]")
    items: Mapped[List["Item"]] = relationship(foreign_keys="[Item.position_id]")
    adventure_id: Mapped[int] = mapped_column(ForeignKey("adventure.id"))
    actual_position_adventure_id: Mapped[Optional[int]] = mapped_column(ForeignKey("adventure.id"))

    def __init__(self) -> None:
        self.temporary_description: str = ""

    def get_serializable(self):
        return {
            # "id": self.id, 
            "code": self.code,
            "description": self.temporary_description if self.temporary_description else self.description,
            "available_actions": list(map(lambda a: a.get_serializable(), filter(lambda action: action.active and action.visible, self.available_actions))), # should not return with inactive or not visible actions
            "end_position": self.end_position
            # "entering_actions": list(map(lambda a: a.get_serializable(), filter(lambda action: action.active and action.visible, self.entering_actions))),
            # "leaving_actions": list(map(lambda a: a.get_serializable(), filter(lambda action: action.active and action.visible, self.leaving_actions))),
            # "items":list(map(lambda i: i.get_serializable(), self.items))
        }

@event.listens_for(Position, "load")
def receive_load(target, context):
    target.temporary_description = ""
