from typing import Optional
from sqlalchemy import ForeignKey, String, Integer
from config import db
from sqlalchemy.orm import Mapped, mapped_column

class Item(db.Model):
    id = db.Column(db.Integer(), primary_key = True, autoincrement = True)
    code = db.Column(db.String(10), unique = False, nullable = False)
    name = db.Column(db.String(40), unique = False, nullable = False)
    description = db.Column(db.String(200), unique = False, nullable = False)
    state = db.Column(db.String(40), unique = False, nullable = True)
    player_id: Mapped[Optional[int]] = mapped_column(Integer(), ForeignKey("player.id"))
    position_id: Mapped[Optional[int]] = mapped_column(Integer(), ForeignKey("position.id"))

    def get_serializable(self):
        return {
            "id": self.id,
            "code": self.code,
            "name": self.name,
            "description": self.description,
            "state": self.state
        }