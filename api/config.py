from flask import Flask
from flask_cors import CORS
from flask_sqlalchemy import SQLAlchemy

app = Flask(__name__)

app.config["CORS_ORIGINS"] = ["http://localhost:5173"]
app.config["CORS_ALLOW_HEADERS"] = ["content-type"]
app.config["CORS_SUPPORTS_CREDENTIALS"] = True
CORS(app)

app.secret_key = "Q8KKJxUmRg"
app.config["SESSION_COOKIE_SAMESITE"] = "None"
app.config["SESSION_COOKIE_SECURE"] = "True"
# TODO
# app.config["SESSION_COOKIE_PARTITIONED"] = "True"
print("ROOT",app.root_path)

# app.config["SQLALCHEMY_DATABASE_URI"] = "sqlite://" # memory only db
app.config["SQLALCHEMY_DATABASE_URI"] = "sqlite:///" + app.root_path + "\\db\\adventures.db"
app.config["SQLALCHEMY_TRACK_MODIFICATIONS"] = False
app.config["SQLALCHEMY_ECHO"] = True

db = SQLAlchemy(app)