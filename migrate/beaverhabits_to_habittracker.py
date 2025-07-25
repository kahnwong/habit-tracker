import json

with open("/Users/kahnwong/Downloads/beaverhabits.json", "r") as f:
    data = json.load(f)

for habit in data["habits"]:
    print(f"habit-tracker create {habit['name']}")

    for activity in habit["records"]:
        print(f"habit-tracker do {habit['name']} {activity['day']}")
