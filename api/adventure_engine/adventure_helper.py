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

def get_item_from_player(player, item_code):
    return next((i for i in player.items if i.code == item_code), None)

def get_item_from_positions(positions, item_code):
    for position in positions:
        item = get_item_from_position(position, item_code)
        if item:
            return item

def pop_item_from_list(items, item_code):
    for i in range(len(items)):
        if items[i].code == item_code:
            return items.pop(i)

def get_item_from_position(position, item_code):
    return next((i for i in position.items if i.code == item_code), None)

def get_item(adventure, item_code: str):
    item = get_item_from_positions(adventure.positions, item_code)
    if not item:
        item = get_item_from_player(adventure.player, item_code)

    return item

def prepare_actions(actions):
    return stringify_action_codes(stringify_action_functions(set_action_defaults(actions)))

def set_action_defaults(actions):
    for action in actions:
        if not "visible" in action:
            action["visible"] = True
        if not "active" in action:
            action["active"] = True
    return actions

def stringify_action_functions(actions):
    for action in actions:
        if "functions" in action and action["functions"]:
            action["functions"] = json.dumps(action["functions"])
        if "function" in action and action["function"]:
            action["function"] = json.dumps(action["function"])
    return actions

def stringify_action_codes(actions):
    for action in actions:
        if "action_codes" in action and action["action_codes"]:
            action["action_codes"] = json.dumps(action["action_codes"])
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
