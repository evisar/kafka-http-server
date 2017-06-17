#!flask/bin/python
from flask import Flask

app = Flask(__name__)

@route("/stream")
def stream():
    def eventStream():
        while True:
            # Poll data from the database
            # and see if there's a new message
            if len(messages) > len(previous_messages):
                yield "data: {}\n\n".format(messages[len(messages)-1)])"
    
    return Response(eventStream(), mimetype="text/event-stream")


if __name__ == '__main__':
    app.run(debug=True)