from typing import List

from sqlalchemy import ForeignKey, String
from adventure_engine.item import Item
from config import db
from sqlalchemy.orm import Mapped, relationship
from sqlalchemy.orm import Mapped, mapped_column

class Player(db.Model):
    id = db.Column(db.String(10), primary_key = True)
    name = db.Column(db.String(40), unique = False, nullable = False)
    items: Mapped[List["Item"]] = relationship()
    adventure_id: Mapped[str] = mapped_column(String(10), ForeignKey("adventure.id"))
    adventure: Mapped["Adventure"] = relationship(back_populates="player", single_parent=True)


    def __init__(self, name: str):
        self.name = name
        # self.items = []
    
    def get_serializable(self):
        return {
            "name": self.name,
            "items": list(map(lambda i: i.get_serializable(), self.items))
        }