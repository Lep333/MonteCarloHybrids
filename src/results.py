import re
import matplotlib.pyplot as plt
import numpy as np

import matplotlib
import matplotlib as mpl
from scipy.stats import binomtest

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
        player1, player2, result, settings1, settings2, moves, no_rollouts, no_beliefs = fields
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
    # print(result_dict)
    # print_heatmap(result_dict)

    # print hybrids
    print_hybrids(result_dict)

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

def print_hybrids(result_dict: dict[list]):
    diagrams = [
        (
            "DPC_OM",
            [r'\"OM_Threshold\":([0-9]+(?:\.[0-9]+)?)'],
            "POMCP Gewinnrate mit verschiedenen UCB Konstante c Werten",
            "c"
        ),
        (
            "DPC_GREEDY",
            [r'\"Rollout_selection\":\{\"Epsilon\":([0-9]+(?:\.[0-9]+)?)'],
            "Greedy-Hybrid",
            "Epsilon"
        ),
        (
            "DPC_EPT",
            [r'\"Early_playout_termination\":\{\"Max_depth\":([0-9]+(?:\.[0-9]+)?)'],
            "Early-Playout-Termination",
            "Abbruchstiefe"
        ),
        (
            "DPC_CAPTURE_PREF",
            [r'\"Rollout_capture\":([0-9]+(?:\.[0-9]+)?)'],
            "Schlagpräfarenz",
            "Schlagwahrscheinlichkeit"
        ),
        (
            "DPC_EVALUATION_CUT_OFF",
            [r'\"Early_playout_termination\":\{\"Threshold\":([0-9]+(?:\.[0-9]+)?)'],
            "Frühzeitige Abbrüche",
            "Abbruchsschwellwert"
        ),
        (
            "DPC_K_BEST",
            [r'\"Rollout_selection\":\{"K":([0-9]+(?:\.[0-9]+)?)'],
            "K-Beste",
            "k"
        ),
        (
            "DPC_IC",
            [
                r'\"Early_playout_termination\":\{\"Max_depth\":([0-9]+(?:\.[0-9]+)?)',
                r'\"Search_depth\":([0-9]+(?:\.[0-9]+)?)',
            ],
            "Vielversprechende Abbrüche",
            ["Suchtiefe", "Abbruchstiefe"]
        ),
        (
            "LAC_OM",
            [r'\"OM_Threshold\":([0-9]+(?:\.[0-9]+)?)'],
            "POMCP Gewinnrate mit verschiedenen UCB Konstante c Werten",
            "c"
        ),
        (
            "LAC_GREEDY",
            [r'\"Rollout_selection\":\{\"Epsilon\":([0-9]+(?:\.[0-9]+)?)'],
            "Greedy-Hybrid",
            "Epsilon"
        ),
        (
            "LAC_EPT",
            [r'\"Early_playout_termination\":\{\"Max_depth\":([0-9]+(?:\.[0-9]+)?)'],
            "Early-Playout-Termination",
            "Abbruchstiefe"
        ),
        (
            "LAC_CAPTURE_PREF",
            [r'\"Rollout_capture\":([0-9]+(?:\.[0-9]+)?)'],
            "Schlagpräfarenz",
            "Schlagwahrscheinlichkeit"
        ),
        (
            "LAC_CORRECTIVE",
            [r'\"Rollout_selection\":\{\"Epsilon\":([0-9]+(?:\.[0-9]+)?)'],
            "Korrektur",
            "Schwellwert Zug zu wählen"
        ),
        (
            "LAC_EVALUATION_CUT_OFF",
            [r'\"Early_playout_termination\":\{\"Threshold\":([0-9]+(?:\.[0-9]+)?)'],
            "Frühzeitige Abbrüche",
            "Abbruchsschwellwert"
        ),
        (
            "LAC_K_BEST",
            [r'\"Rollout_selection\":\{"K":([0-9]+(?:\.[0-9]+)?)'],
            "K-Beste",
            "k"
        ),
        (
            "LAC_IC",
            [
                r'\"Early_playout_termination\":\{\"Max_depth\":([0-9]+(?:\.[0-9]+)?)',
                r'\"Search_depth\":([0-9]+(?:\.[0-9]+)?)',
            ],
            "Vielversprechende Abbrüche",
            ["Suchtiefe", "Abbruchstiefe"]
        ),
        (
            "LAC_IR",
            [
                r'\"Rollout_selection\":\{\"Search_depth\":([0-9]+(?:\.[0-9]+)?)',
                r'\"Search_depth\":[0-9]+ \"Epsilon\":([0-9]+(?:\.[0-9]+)?)',
             ],
            "Vielversprechende Rollouts",
            ["Epsilon", "Suchtiefe"]
        ),
    ]

    for name, reg, title, x_axis_label in diagrams:
        x = []
        y = []
        e_low = []
        e_high = []
        if len(reg) > 1:
            two_parameters(result_dict, name, reg, title, x_axis_label)
            continue
        else:
            reg = reg.pop()
        for key, value in result_dict.items():
            pomcp_name = re.search(r'POMCP_name"\s*:\s*"([^"]*)', key)
            if pomcp_name and name == pomcp_name.group(1):
                ucb_c = ""
                if found := re.search(reg, key):
                    ucb_c = float(found.group(1))
                player_wins = value[0][0] + value[1][0]
                player_draws = (value[0][1] + value[1][1]) * 0.5
                player_games = sum(value[0]) + sum(value[1])
                win_percentage = 100 * (player_wins + player_draws) / player_games

                res = binomtest(int(player_wins+player_draws), player_games)
                ci = res.proportion_ci(1 - 0.05, method="wilson")
                ci_low_pct = 100 * ci.low
                ci_high_pct = 100 * ci.high
                lower_err = win_percentage - ci_low_pct
                upper_err = ci_high_pct - win_percentage
                print(name, ucb_c, win_percentage, value)
                x.append(ucb_c)
                y.append(win_percentage)
                e_low.append(lower_err)
                e_high.append(upper_err)
    
        x_label, y, e_low, e_high = zip(*sorted(zip(x, y, e_low, e_high)))
        x = 0.5 + np.arange(len(y))
        fig, ax = plt.subplots()
        indices = np.arange(len(y))
        ax.errorbar(indices, y, [e_low, e_high], fmt='o', linewidth=2, capsize=6)
        ax.set_xticks(indices)
        ax.set_xticklabels(x_label)
        ax.set_title(title)
        ax.set_ylabel("Siegesrate in %")
        ax.set_xlabel(x_axis_label)
        ax.grid(True, which='major', axis='y', linestyle='--', alpha=0.6)
        ax.set_ylim(0,100)
        ax.set_yticks(np.arange(0, 101, 10))
        fig.tight_layout()
        plt.savefig(f"{name}.png", dpi=300)
        plt.close()

def two_parameters(result_dict: dict, name: str, reg: list[str], title: str, x_axis_label: str):
    x = {}
    y = {}
    e_low = {}
    e_high = {}
    for key, value in result_dict.items():
        pomcp_name = re.search(r'POMCP_name"\s*:\s*"([^"]*)', key)
        if pomcp_name and name == pomcp_name.group(1):
            param1 = float(re.search(reg[0], key).group(1))
            param2 = float(re.search(reg[1], key).group(1))
            player_wins = value[0][0] + value[1][0]
            player_games = sum(value[0]) + sum(value[1])
            win_percentage = 100 * player_wins / player_games
            res = binomtest(player_wins, player_games)
            ci = res.proportion_ci(1 - 0.05, method="wilson")
            print(name, param1, param2, win_percentage, value)
            if not x.get(param1):
                x[param1] = []
                y[param1] = []
                e_low[param1] = []
                e_high[param1] = []
            x[param1].append(param2)
            y[param1].append(win_percentage)
            ci_low_pct = 100 * ci.low
            ci_high_pct = 100 * ci.high
            lower_err = win_percentage - ci_low_pct
            upper_err = ci_high_pct - win_percentage
            e_low[param1].append(lower_err)
            e_high[param1].append(upper_err)
    fig, ax = plt.subplots()
    groups = sorted(x.keys())
    n_groups = len(groups)
    width = 0.15  # total horizontal spread
    offsets = np.linspace(-width, width, n_groups)
    for offset, param1 in zip(offsets, groups):
        #x_label, y, e_low, e_high = zip(*sorted(zip(x, y, e_low, e_high)))
        #x = 0.5 + np.arange(len(y))
        x_shifted = np.array(x[param1]) + offset
        ax.errorbar(x_shifted, y[param1], [e_low[param1], e_high[param1]], fmt='o', linewidth=2, capsize=6, label=f"{param1}")

    ax.set_title(title)
    ax.set_ylabel("Siegesrate in %")
    ax.set_xlabel(x_axis_label[0])
    ax.legend(title=x_axis_label[1])
    ax.grid(True, which='major', axis='y', linestyle='--', alpha=0.6)
    ax.set_ylim(0,100)
    ax.set_yticks(np.arange(0, 101, 10))
    unique_x = sorted({v for vals in x.values() for v in vals})
    ax.set_xticks(unique_x)
    fig.tight_layout()
    plt.savefig(f"{name}.png", dpi=300)
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