import re

file_name = "results.csv"
csv = ""

def main():
    with open(file_name, mode="+r") as f:
        csv = f.readlines()

    result_dict = {}
    rollout_capture_index = 10
    #key_field_indexes = [time_limit, ucb_c_key, rollout_capture_index]
    i = 0
    for i, line in enumerate(csv):
        if i == 0:
            continue
        fields = line.replace("\n", "").split(",")
        player1, player2, result, settings1, settings2, moves, no_rollouts = fields
        white = [0,0,0]
        black = [0,0,0]
        if int(result) == 1:
            white = [1,0,0]
            black = [0,0,1]
        elif int(result) == -1:
            white = [0,0,1]
            black = [1,0,0]
        else:
            white = [0,1,0]
            black = [0,1,0]

        if result_dict.get(settings1):
            result_dict[settings1] = ([
                result_dict[settings1][0][0]+white[0], 
                result_dict[settings1][0][1]+white[1],
                result_dict[settings1][0][2]+white[2]
                ],
                result_dict[settings1][1])
        else:
            result_dict[settings1] = (white,[0,0,0])
        
        if result_dict.get(settings2):
            result_dict[settings2] = (
                result_dict[settings2][0],
                [
                    result_dict[settings2][1][0]+black[0],
                    result_dict[settings2][1][1]+black[1],
                    result_dict[settings2][1][2]+black[2],
                ]
                )
        else:
            result_dict[settings2] = ([0,0,0],black)
    print(result_dict)
    print_heatmap(result_dict)

def print_heatmap(result_dict: dict[list]):
    ucb_label = set()
    capture_reward_label = set()
    for key, val in result_dict.items():
        ucb_c = int(re.search(r"Ucb_c:(\d+)", key).group(1))
        capture_reward = float(re.search(r"Rollout_capture:([\d.]+)", key).group(1))
        if capture_reward == 0:
            continue
        ucb_label.add(ucb_c)
        capture_reward_label.add(capture_reward)
    sorted_ucb_label = list(ucb_label)
    sorted_ucb_label.sort()
    sorted_capture_label = list(capture_reward_label)
    sorted_capture_label.sort(reverse=True)

    # modified from here: https://matplotlib.org/stable/gallery/images_contours_and_fields/image_annotated_heatmap.html
    import matplotlib.pyplot as plt
    import numpy as np

    import matplotlib
    import matplotlib as mpl

    heatmap = np.zeros((len(sorted_capture_label), len(sorted_ucb_label)))
    for key, val in result_dict.items():
        termination_param = float(re.search(r"Termination_parameter:([\d.]+)", key).group(1))
        ucb_c = int(re.search(r"Ucb_c:(\d+)", key).group(1))
        ucb_index = sorted_ucb_label.index(ucb_c)

        capture_reward = float(re.search(r"Rollout_capture:([\d.]+)", key).group(1))
        if capture_reward == 0 or termination_param != 1000:
            continue
        capture_index = sorted_capture_label.index(capture_reward)

        winrate = (val[0][0] + val[1][0]) / (sum(val[0]) + sum(val[1]))
        heatmap[capture_index][ucb_index] = round(winrate, 2)


    fig, ax = plt.subplots()
    im = ax.imshow(heatmap)

    # Show all ticks and label them with the respective list entries
    ax.set_xticks(range(len(sorted_ucb_label)), labels=sorted_ucb_label,
                rotation=45, ha="right", rotation_mode="anchor")
    ax.set_yticks(range(len(sorted_capture_label)), labels=sorted_capture_label)

    # Loop over data dimensions and create text annotations.
    for i in range(len(sorted_capture_label)):
        for j in range(len(sorted_ucb_label)):
            text = ax.text(j, i, heatmap[i, j],
                        ha="center", va="center", color="w")

    ax.set_title("Siegesrate abhängig von UCB_C und Schlagpräferenz")
    fig.tight_layout()
    plt.savefig("winrate_heatmap_dark_pawn_game_1sec.png", dpi=300)
    plt.close()

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