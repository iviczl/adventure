from flask import Response, jsonify, request, make_response, session
from adventure_engine.adventure import Adventure
from adventure_engine.adventure_helper import get_info
from config import app, db

if __name__ == "__main__":
    with app.app_context():
        db.create_all()
        app.run(debug = True)

def _adventure_files():
    return ["./adventures/dungeon.json", "./adventures/lost_world.json"]

def _adventures():
    adventures: list = []
    for adventure_file in _adventure_files():
        adventures.append({"path": adventure_file, "info": get_info(adventure_file)})
    return adventures

def _get_adventure(id: str):
    return next((a for a in _adventures() if a["info"]["id"] == id), None)

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
    adventure_id = request.json.get("gameId")

    if not player_name or not adventure_id:
        return (jsonify({"message": "player and gameId arguments are mandatory."}), 400)

    adventure = Adventure() # (player_name)
    Adventure.load(_get_adventure(adventure_id)["path"], adventure)
    adventure.player.name = player_name
    adventure.actual_position = next((p for p in adventure.positions if p.id == "0"), None)
    try:
        db.session.add(adventure)
        db.session.commit()
        session[adventure.id] = player_name # adventure.get_serializable()
    except Exception as e:
        return jsonify({ "message": str(e)}), 400

    response =  make_response(adventure.actual_position.get_serializable())
    return response, 201

@app.route("/do", methods=["POST","OPTIONS"])
def do():
    if request.method =="OPTIONS":
       return ""
    # POST
    adventure_id = request.json.get("gameId")
    action_id = request.json.get("actionId")

    if not adventure_id or not action_id:
        return (jsonify({"message": "actionId and gameId arguments are mandatory."}), 400)

    adventure = Adventure.build(session[adventure_id])
    print("OLD POSITION",adventure.actual_position.id)
    print("ACTION TO EXECUTE",action_id)
    session.pop(adventure_id)
    adventure.do(action_id)
    print("POSITION",adventure.actual_position.id)
    session[adventure.id] = adventure.get_serializable()
    # session.modified = True
    response =  make_response(adventure.actual_position.get_serializable())
    return response

@app.route("/quit", methods =["DELETE"])
def quit():
    adventure_id = request.json.get("gameId")

    if not adventure_id:
        return (jsonify({"message": "gameId argument is mandatory."}), 400)

    session.pop(adventure_id)

