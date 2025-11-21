file_name = "results.csv"
csv = ""

def main():
    with open(file_name, mode="+r") as f:
        csv = f.readlines()

    result = {"c": {}, "time": {}, "capture": {}}
    i = 0
    for line in csv:
        fields = line.replace("\n", "").split(",")
        time_limit = fields[4]
        ucb_c = fields[5]
        capture_reward = fields[8]
        sum_result("time", result, fields, 4)
        sum_result("c", result, fields, 5)
        sum_result("capture", result, fields, 9)
        if i == 199:
            print(result)
            result = {"c": {}, "time": {}, "capture": {}}
            i = 0
        i+=1

def sum_result(category: str, results: dict, line: list, field_no: int):
    if results[category].get(line[field_no]):
        if line[0] == "POMCP":
            if int(line[2]) > 0:
                results[category][line[field_no]] = (results[category][line[field_no]][0] + 1, results[category][line[field_no]][1])
            elif int(line[2]) < 0:
                results[category][line[field_no]] = (results[category][line[field_no]][0], results[category][line[field_no]][1] + 1)
        else:
            if int(line[2]) < 0:
                results[category][line[field_no]] = (results[category][line[field_no]][0] + 1, results[category][line[field_no]][1])
            elif int(line[2]) > 0:
                results[category][line[field_no]] = (results[category][line[field_no]][0], results[category][line[field_no]][1] + 1)
    else:
        if line[0] == "POMCP":
            if int(line[2]) > 0:
                results[category][line[field_no]] = (1, 0)
            elif int(line[2]) < 0:
                results[category][line[field_no]] = (0, 1)
        else:
            if int(line[2]) < 0:
                results[category][line[field_no]] = (1, 0)
            elif int(line[2]) > 0:
                results[category][line[field_no]] = (0, 1)

if __name__ == "__main__":
    main()