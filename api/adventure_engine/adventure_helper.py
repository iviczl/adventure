# from adventure import Adventure
# from position import Position

import json

def file_content(file_path: str):
    with open(file_path) as adventure_file:
        return adventure_file.read()

def load_json(file_path: str):
    # try:
        with open(file_path) as adventure_file:
            file_content = adventure_file.read()
            return json.loads(file_content)
    # except Exception as error:
    #     print(f"An error occured: {error}")

def get_info(file_path: str):
    data: dict = load_json(file_path)
    return {"id": data["id"], "title": data["title"], "description" : data["description"]}

def get_position_from_position_list(positions: list, position_id: str):
    return next((p for p in positions if p.id == position_id), None)

def get_action_from_position_list(positions: list, action_id: str):
    for position in positions:
        return get_action_from_position(position, action_id)
    
def get_action_from_position(position, action_id: str):
    action = next((a for a in position.available_actions if a.id == action_id), None)
    if action == None:
        action = next((a for a in position.entering_actions if a.id == action_id), None)
    if action == None:
        action = next((a for a in position.leaving_actions if a.id == action_id), None)

    return action
    
def get_item(adventure, item_id: str):
    item = next((i for i in adventure.actual_position.items if i.id == item_id), None)
    if item == None:
        item = next((i for i in adventure.player.items if i.id == item_id), None)

    return item
