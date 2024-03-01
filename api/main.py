from os import listdir
from flask import Response, jsonify, request, make_response, session
from adventure_engine.adventure import Adventure, Position
from adventure_engine.adventure_helper import get_info
from config import app, db

ADVENTURE_PATH = "./adventures/"

if __name__ == "__main__":
    with app.app_context():
        db.create_all()
        app.run(debug = True)

def _adventure_files():
    files = listdir(ADVENTURE_PATH)
    if len(files) == 0:
        return files
    return map(lambda adventure: ADVENTURE_PATH + adventure ,filter(lambda file: file.endswith(".json") ,files))
    # return ["./adventures/dungeon.json", "./adventures/lost_world.json"]

def _adventures():
    adventures: list = []
    for adventure_file in _adventure_files():
        adventures.append({"path": adventure_file, "info": get_info(adventure_file)})
    return adventures

def _get_adventure(id: str):
    return next((a for a in _adventures() if a["info"]["id"] == id), None)

ADVENTURE = "adventure"

def _authenticate(keys):
    adventure_id = session[ADVENTURE]
    if not adventure_id:
        return False, (jsonify({"message": "You should choose a game first."}), 400)
        
    values = {}
    for key in keys:
        value = request.json.get(key)
        if value:
            values[key] = value
        else:
            return False, (jsonify({"message": f"{key} argument is mandatory."}), 400)

    values[ADVENTURE] = adventure_id
    return True, None, values

@app.route("/")
def main():
    response =  make_response(list(map(lambda a: a["info"], _adventures())))
    return response

@app.route("/new", methods=["POST","OPTIONS"])
def new():
    if request.method =="OPTIONS":
        return ""
    # POST
    player_name = request.json.get("player")
    game_id = request.json.get("gameId")

    if not player_name or not game_id:
        return (jsonify({"message": "player and gameId arguments are mandatory."}), 400)

    adventure = Adventure() # (player_name)
    Adventure.load(_get_adventure(game_id)["path"], adventure)
    adventure.player.name = player_name
    try:
        db.session.add(adventure)
        print("ADVENTURE ID",adventure.id)
        db.session.commit()
        session[ADVENTURE] = adventure.id 
    except Exception as e:
        return jsonify({ "message": str(e)}), 400

    response =  make_response(adventure.actual_position.get_serializable())
    return response, 201

@app.route("/do", methods=["POST","OPTIONS"])
def do():
    if request.method =="OPTIONS":
       return ""
    # POST
    has_rights, response, values = _authenticate(["actionId"])
    if not has_rights:
        return response
    try:
        adventure = Adventure.query.get(values[ADVENTURE])
        print("OLD POSITION",adventure.actual_position.code)
        print("ACTION TO EXECUTE",values["actionId"])
        adventure.do(values["actionId"])
        db.session.commit()
        print("POSITION",adventure.actual_position.code)
    except Exception as e:
        return jsonify({ "message": str(e)}), 400

    response =  make_response(adventure.actual_position.get_serializable())
    return response

@app.route("/quit", methods =["DELETE"])
def quit():
    has_rights, response, values = _authenticate()
    if not has_rights:
        return response

    session.pop(values[ADVENTURE])

