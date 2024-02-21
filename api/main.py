from flask import Flask, Request, Response, request, make_response, session
from adventure_engine.adventure import Adventure
from adventure_engine.adventure_helper import get_info

app = Flask(__name__)
app.secret_key = "Q8KKJxUmRg"

def _adventure_files():
    return ["./adventures/dungeon.json", "./adventures/lost_world.json"]

def _adventures():
    adventures: list = []
    for adventure_file in _adventure_files():
        adventures.append({"path": adventure_file, "info": get_info(adventure_file)})
    return adventures

def _get_adventure(id: str):
    return next((a for a in _adventures() if a["info"]["id"] == id), None)

def _add_cors_header(response: Response):
    response.headers["Access-Control-Allow-Origin"] ="http://localhost:5173"

def _add_access_control_header(response: Response):
    response.headers["Access-Control-Allow-Headers"] ="content-type"

def prefetch_response():
    response = make_response()
    _add_cors_header(response)
    _add_access_control_header(response)
    return response


@app.route("/")
def main():
    # adv = _adventures()
    # print(adv)
    response =  make_response(list(map(lambda a: a["info"], _adventures())))
    _add_cors_header(response)
    return response

@app.route("/new", methods=["POST","OPTIONS"])
def new():
    if request.method =="OPTIONS":
        return prefetch_response()
    # POST
    data = request.get_json()
    adventure = Adventure(data["player"])
    adventure.load(_get_adventure(data["gameId"])["path"])
    session[adventure.id] = adventure.get_serializable()
    adventure.actual_position = next((p for p in adventure.positions if p.id == "0"), None)

    response =  make_response(adventure.actual_position.get_serializable())
    _add_cors_header(response)
    return response

@app.route("/do", methods=["POST","OPTIONS"])
def do():
    if request.method =="OPTIONS":
       return prefetch_response()
    # POST
    data = request.get_json()
    adventure_id = data["gameId"]
    adventure: Adventure = session[adventure_id]
    adventure.do(data["actionId"])
    # adventure.actual_position = next((p for p in adventure.positions if p.id == "0"), None)

    response =  make_response(adventure.actual_position.get_serializable())
    _add_cors_header(response)
    return response

@app.route("/quit")
def quit():
    session.pop(Request.get_json().name)

