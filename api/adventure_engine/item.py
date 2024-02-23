from sqlalchemy import ForeignKey, String
from config import db
from sqlalchemy.orm import Mapped, mapped_column

class Item(db.Model):
    id = db.Column(db.String(10), primary_key = True)
    name = db.Column(db.String(40), unique = False, nullable = False)
    description = db.Column(db.String(200), unique = False, nullable = False)
    state = db.Column(db.String(40), unique = False, nullable = True)
    player_id = mapped_column(String(10), ForeignKey("player.id"))
    position_id = mapped_column(String(10), ForeignKey("position.id"))

    def get_serializable(self):
        return {
            "id": self.id,
            "name": self.name,
            "description": self.description,
            "state": self.state
        }