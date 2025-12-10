file_name = "results.csv"
csv = ""

def main():
    with open(file_name, mode="+r") as f:
        csv = f.readlines()

    result = {}
    time_limit = 4
    ucb_c_key = 5
    rollout_capture_index = 10
    key_field_indexes = [time_limit, ucb_c_key, rollout_capture_index]
    i = 0
    for i, line in enumerate(csv):
        if i == 0:
            continue
        fields = line.replace("\n", "").split(",")
        time_limit = fields[4]
        ucb_c = fields[5]
        capture_reward = fields[8]
        key_str = ""
        for key_field in key_field_indexes:
            if len(fields)-1 >= key_field:
                key_str += fields[key_field]
        if result.get(key_str):
            if fields[0] == "POMCP":
                if int(fields[2]) > 0:
                    result[key_str] = (result[key_str][0] + 1, result[key_str][1])
                elif int(fields[2]) < 0:
                    result[key_str] = (result[key_str][0], result[key_str][1] + 1)
            else:
                if int(fields[2]) < 0:
                    result[key_str] = (result[key_str][0] + 1, result[key_str][1])
                elif int(fields[2]) > 0:
                    result[key_str] = (result[key_str][0], result[key_str][1] + 1)
        else:
            if fields[0] == "POMCP":
                if int(fields[2]) > 0:
                    result[key_str] = (1, 0)
                elif int(fields[2]) < 0:
                    result[key_str] = (0, 1)
            else:
                if int(fields[2]) < 0:
                    result[key_str] = (1, 0)
                elif int(fields[2]) > 0:
                    result[key_str] = (0, 1)
    print(result)

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