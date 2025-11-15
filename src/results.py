file_name = "results.csv"
csv = ""

with open(file_name, mode="+r") as f:
    csv = f.readlines()

result = {}
for line in csv:
    fields = line.split(",")
    if val := result.get(fields[5]):
        if fields[0] == "POMCP":
            if int(fields[2]) > 0:
                result[fields[5]] = (result[fields[5]][0] + 1, result[fields[5]][1])
            elif int(fields[2]) < 0:
                result[fields[5]] = (result[fields[5]][0], result[fields[5]][1] + 1)
        else:
            if int(fields[2]) < 0:
                result[fields[5]] = (result[fields[5]][0] + 1, result[fields[5]][1])
            elif int(fields[2]) > 0:
                result[fields[5]] = (result[fields[5]][0], result[fields[5]][1] + 1)
    else:
        if fields[0] == "POMCP":
            if int(fields[2]) > 0:
                result[fields[5]] = (1, 0)
            elif int(fields[2]) < 0:
                result[fields[5]] = (0, 1)
        else:
            if int(fields[2]) < 0:
                result[fields[5]] = (1, 0)
            elif int(fields[2]) > 0:
                result[fields[5]] = (0, 1)

for key, val in result.items():
    print(f"{key}: {val}")