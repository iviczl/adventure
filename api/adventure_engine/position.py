from typing import List

from sqlalchemy import ForeignKey, String
from adventure_engine.action import Action
from adventure_engine.item import Item
from sqlalchemy.orm import Mapped, relationship, mapped_column
from config import db

class Position(db.Model):
    id = db.Column(db.String(10), primary_key = True)
    description = db.Column(db.String(100), unique = False, nullable = False)
    available_actions: Mapped[List["Action"]] = relationship(overlaps="leaving_actions,entering_actions")
    entering_actions: Mapped[List[Action]] = relationship(overlaps="available_actions,leaving_actions")
    leaving_actions: Mapped[List[Action]] = relationship(overlaps="available_actions,entering_actions")
    items: Mapped[List[Item]] = relationship()
    adventure_id: Mapped[str] = mapped_column(ForeignKey("adventure.id"))

    def get_serializable(self):
        return {
            "id": self.id, 
            "description": self.description,
            "available_actions": list(map(lambda a: a.get_serializable(), self.available_actions)),
            "entering_actions": list(map(lambda a: a.get_serializable(), self.entering_actions)),
            "leaving_actions": list(map(lambda a: a.get_serializable(), self.leaving_actions)),
            "items":list(map(lambda i: i.get_serializable(), self.items))
        }