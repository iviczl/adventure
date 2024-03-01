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
    return {"id": data["code"], "title": data["title"], "description" : data["description"]}

def get_position_from_position_list(positions: list, position_code: str):
    return next((p for p in positions if p.code == position_code), None)

def get_action(adventure, action_code: str):
    action = get_action_from_position_list(adventure.positions, action_code)
    if not action:
        action = next((a for a in adventure.available_actions if a.code == action_code), None)
    return action

def get_action_from_position_list(positions: list, action_code: str):
    for position in positions:
        action = get_action_from_position(position, action_code)
        if action:
            return action
    
def get_action_from_position(position, action_code: str):
    action = next((a for a in position.available_actions if a.code == action_code), None)
    if action == None:
        action = next((a for a in position.entering_actions if a.code == action_code), None)
    if action == None:
        action = next((a for a in position.leaving_actions if a.code == action_code), None)

    return action
    
def get_item(adventure, action_code: str):
    item = next((i for i in adventure.actual_position.items if i.code == action_code), None)
    if item == None:
        item = next((i for i in adventure.player.items if i.code == action_code), None)

    return item

def stringify_action_functions(actions):
    for action in actions:
        if action["function"] != None:
            action["function"] = json.dumps(action["function"])
    return actions

def change_dict_keys(dictionary: dict, changeable: dict):
    for key in changeable.keys():
        change_key(dictionary, key, changeable[key])

def change_key(dictionary: dict, old_key: str, new_key:str ):
    if old_key in dictionary:
        value = dictionary[old_key]
        dictionary.pop(old_key)
        dictionary[new_key] = value

# def prepare_actions(actions: list):
#     stringify_action_functions(actions)
#     for action in actions:
#         change_dict_keys(action, { "id": "code", "position_id": "position_code", "item_id": "item_code"})

#     return actions

# def prepare_items(items: list):
#     for item in items:
#         change_dict_keys(item, { "id": "code"})

#     return items
