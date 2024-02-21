from adventure_engine.item import Item


class Player():
    name: str
    items: list[Item] = []

    def __init__(self, name: str):
        self.name = name
        # self.items = []
    
    def get_serializable(self):
        return {
            "name": self.name,
            "items": list(map(lambda i: i.get_serializable(), self.items))
        }