
class Item():
    def __init__(self, id: str, name: str, description: str, state: str):
        self.id = id
        self.name = name
        self.description = description
        self.state = state
    
    def get_serializable(self):
        return {
            "id": self.id,
            "name": self.name,
            "description": self.description,
            "state": self.state
        }