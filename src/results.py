import statistics
import re
import matplotlib.pyplot as plt
import numpy as np

import matplotlib
import matplotlib as mpl
from scipy.stats import binomtest
import pandas as pd

csv = []
files = [
         "results/dpc_uct.csv", "results/dpc_corrective.csv","results/dpc_ept.csv","results/dpc_evaluation_cut_off.csv",
         "results/dpc_evaluation_cut_off.csv",
         "results/dpc_greedy.csv", "results/dpc_ic.csv",
         "results/dpc_k_best.csv", "results/dpc_rollout_capture.csv",
         "results/DPC_IR.csv",
         "results/LAC_UCT.csv", "results/LAC_OM.csv", "results/LAC_CAPTURE_PREF.csv",
         "results/LAC_CORRECTIVE.csv", "results/LAC_EPT.csv",
         "results/LAC_EVALUATION_CUT_OFF.csv", "results/LAC_GREEDY.csv",
         "results/LAC_IC_2.csv", "results/LAC_IR.csv", "results/LAC_K_BEST.csv",
         "results/LAC_UCT.csv", "results/LAC_OM.csv", "results/LAC_GREEDY.csv",
         "results/LAC_CAPTURE_PREF.csv", "results/LAC_CORRECTIVE.csv", "results/LAC_EPT.csv",
         "results/LAC_EVALUATION_CUT_OFF.csv", "results/LAC_K_BEST.csv",  
         "results/DC_UCT.csv", "results/DC_OM.csv", "results/DC_GREEDY.csv",
         "results/DC_EPT_2.csv", "results/DC_CORRECTIVE.csv", "results/DC_ROLLOUT_CAPTURE.csv",
         "results/DC_EVAL_CUT_OFF_2.csv", "results/DC_K_BEST.csv",
         "results/LAC_IC_3.csv",
         "results/DC_EVAL_CUT_OFF_2.csv", "results/DC_CORRECTIVE_2.csv", "results/DC_GREEDY_2.csv",
         "results/DC_CAPTURE_PREF_2.csv", "results/DC_EPT_2.csv",
         "results/DC_K_BEST.csv", "results/DC_IR.csv", "results/DC_IC.csv",
         "results/DC_EVAL_CUT_OFF_2.csv", "results/DC_CORRECTIVE_2.csv", "results/DC_GREEDY_2.csv",
         "results/DC_CAPTURE_PREF_2.csv", "results/DC_EPT_2.csv",
         "results/DC_K_BEST.csv", "results/DC_IR.csv", "results/DC_IC.csv", "results/DPC_GREEDY_VS_BASELINE.csv",
         "results/DPC_GREEDY_VS_BASELINE_OM.csv", "results/DPC_GREEDY_VS_EVAL_CUT_OFF.csv",
         "results/DPC_GREEDY_VS_EVAL_CUT_OFF_OM.csv", "results/LAC_BASE_VS_GREEDY.csv",
         "results/LAC_EVAL_VS_GREEDY.csv", "results/DC_BASELINE_VS_GREEDY.csv",
         "results/DC_IR_VS_GREEDY.csv", "results/DC_IC_2.csv", "results/DC_IR_2.csv",
         "results/LAC_IR_2.csv", "results/DPC_IR_2.csv", "results/DPC_IC_2.csv",
        ]
    
def main():
    for file in files:
        with open(file, mode="+r") as f:
            csv.extend(f.readlines())

    result_dict = {}
    rollout_capture_index = 10
    #key_field_indexes = [time_limit, ucb_c_key, rollout_capture_index]
    i = 0
    for i, line in enumerate(csv):
        if i == 0 or line == "":
            continue
        print(i)
        fields = line.replace("\n", "").split(",")
        player1, player2, result, settings1, settings2, moves, no_rollouts, no_beliefs = fields
        no_rollouts = no_rollouts[1:]
        roll = no_rollouts.replace("[", "").replace("]", "").split(" ")
        player1_rolls = player2_rolls = []
        for i, ro in enumerate(roll):
            if i%2==0:
                player1_rolls.append(int(ro))
            else:
                player2_rolls.append(int(ro))
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
                result_dict[settings1][1],
                result_dict[settings1][2] + player1_rolls
            )
        else:
            result_dict[settings1] = (white,[0,0,0],player1_rolls)
        if result_dict.get(settings2):
            result_dict[settings2] = (
                result_dict[settings2][0],
                [
                    result_dict[settings2][1][0]+black[0],
                    result_dict[settings2][1][1]+black[1],
                    result_dict[settings2][1][2]+black[2],
                ],
                result_dict[settings2][2] + player2_rolls
                )
        else:
            result_dict[settings2] = ([0,0,0],black,player2_rolls)
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
            "DPC_UCT",
            [r'\"Ucb_c\":([0-9]+(?:\.[0-9]+)?)'],
            "POMCP Gewinnrate mit verschiedenen UCB Konstante c Werten",
            "c"
        ),
        (
            "DPC_CORRECTIVE",
            [r'\"Rollout_selection\":\{\"Bound\":([0-9]+(?:\.[0-9]+)?)'],
            "Korrektur",
            "Bewertungsschwellwert für sofortige Zugwahl"
        ),
        (
            "DPC_GREEDY",
            [r'\"Rollout_selection\":\{\"Epsilon\":([0-9]+(?:\.[0-9]+)?)'],
            "Greedy",
            "Wahrscheinlichkeit für zufällige Zugwahl"
        ),
        (
            "DPC_EPT",
            [r'\"Early_playout_termination\":\{\"Max_depth\":([0-9]+(?:\.[0-9]+)?)'],
            "Frühzeitiger Abbruch",
            "Rollout-Tiefe"
        ),
        (
            "DPC_EVALUATION_CUT_OFF",
            [r'\"Early_playout_termination\":\{\"Threshold\":([0-9]+(?:\.[0-9]+)?)'],
            "Bewertungsabbrüche",
            "Bewertungsschwellwert für sofortigen Abbruch"
        ),
        (
            "DPC_ROLLOUT_PREF",
            [r'\"Rollout_capture\":([0-9]+(?:\.[0-9]+)?)'],
            "Schlagpräferenz",
            "Schlagwahrscheinlichkeit"
        ),
        (
            "DPC_K_BEST",
            [r'\"Rollout_selection\":\{"K":([0-9]+(?:\.[0-9]+)?)'],
            "K-Beste",
            "Anzahl berücksichtigter Züge"
        ),
        (
            "DPC_GREEDY_VS_BASELINE",
            [r'\"Rollout_selection\":\{"K":([0-9]+(?:\.[0-9]+)?)'],
            "K-Beste",
            "k"
        ),
        (
            "DPC_GREEDY_VS_BASELINE_OM",
            [r'\"Rollout_selection\":\{"K":([0-9]+(?:\.[0-9]+)?)'],
            "K-Beste",
            "k"
        ),
        (
            "DPC_GREEDY_VS_EVAL_CUT_OFF",
            [r'\"Ucb_c\":([0-9]+(?:\.[0-9]+)?)'],
            "POMCP Gewinnrate mit verschiedenen UCB Konstante c Werten",
            "c"
        ),
        (
            "DPC_GREEDY_VS_EVAL_CUT_OFF_OM",
            [r'\"Ucb_c\":([0-9]+(?:\.[0-9]+)?)'],
            "POMCP Gewinnrate mit verschiedenen UCB Konstante c Werten",
            "c"
        ),
        (
            "LAC_BASE_VS_GREEDY",
            [r'\"Ucb_c\":([0-9]+(?:\.[0-9]+)?)'],
            "POMCP Gewinnrate mit verschiedenen UCB Konstante c Werten",
            "c"
        ),
        (
            "LAC_EVAL_VS_GREEDY",
            [r'\"Ucb_c\":([0-9]+(?:\.[0-9]+)?)'],
            "POMCP Gewinnrate mit verschiedenen UCB Konstante c Werten",
            "c"
        ),
        (
            "DC_IR_VS_GREEDY",
            [r'\"Ucb_c\":([0-9]+(?:\.[0-9]+)?)'],
            "POMCP Gewinnrate mit verschiedenen UCB Konstante c Werten",
            "c"
        ),
        (
            "DC_BASELINE_VS_GREEDY",
            [r'\"Ucb_c\":([0-9]+(?:\.[0-9]+)?)'],
            "POMCP Gewinnrate mit verschiedenen UCB Konstante c Werten",
            "c"
        ),
        (
            "DPC_IC_2",
            [
                r'\"Early_playout_termination\":\{\"Max_depth\":([0-9]+(?:\.[0-9]+)?)',
                r'\"Search_depth\":([0-9]+(?:\.[0-9]+)?)',
            ],
            "Vielversprechende Abbrüche",
            ["Minimax-Suchtiefe", "Rollout-Tiefe"]
        ),
        (
            "DPC_IR_2",
            [
                r'\"Rollout_selection\":\{\"Search_depth\":([0-9]+(?:\.[0-9]+)?)',
                r'\"Search_depth\":[0-9]+ \"Epsilon\":([0-9]+(?:\.[0-9]+)?)',
             ],
            "Vielversprechende Rollouts",
            ["Wahrscheinlichkeit für Minimax-Suche", "Minimax-Suchtiefe"]
        ),
        (
            "LAC_UCT",
            [r'\"Ucb_c\":([0-9]+(?:\.[0-9]+)?)'],
            "POMCP Gewinnrate mit verschiedenen UCB Konstante c Werten",
            "c"
        ),
        (
            "LAC_OM",
            [r'\"OM_Threshold\":([0-9]+(?:\.[0-9]+)?)'],
            "POMCP Gewinnrate mit verschiedenen Gegnermodellierungen",
            "c"
        ),
        (
            "DC_UCT",
            [r'\"Ucb_c\":([0-9]+(?:\.[0-9]+)?)'],
            "POMCP Gewinnrate mit verschiedenen UCB Konstante c Werten",
            "c"
        ),
        (
            "DC_OM",
            [r'\"OM_Threshold\":([0-9]+(?:\.[0-9]+)?)'],
            "POMCP Gewinnrate mit verschiedenen UCB Konstante c Werten",
            "c"
        ),
        (
            "DC_GREEDY_2",
            [r'\"Rollout_selection\":\{\"Epsilon\":([0-9]+(?:\.[0-9]+)?)'],
            "Greedy",
            "Wahrscheinlichkeit für zufällige Zugwahl"
        ),
        (
            "DC_CORRECTIVE_2",
            [r'\"Rollout_selection\":\{\"Bound\":([0-9]+(?:\.[0-9]+)?)'],
            "Korrektur",
            "Bewertungsschwellwert für sofortige Auswahl"
        ),
        (
            "DC_CAPTURE_PREF_2",
            [r'\"Rollout_capture\":([0-9]+(?:\.[0-9]+)?)'],
            "Schlagpräferenz",
            "Schlagwahrscheinlichkeit"
        ),
        (
            "DC_EPT_2",
            [r'\"Early_playout_termination\":\{\"Max_depth\":([0-9]+(?:\.[0-9]+)?)'],
            "Frühzeitige Abbrüche",
            "Rollout-Tiefe"
        ),
        # (
        #     "DC_K_BEST",
        #     [r'\"Rollout_selection\":\{"K":([0-9]+(?:\.[0-9]+)?)'],
        #     "K-Beste",
        #     "k"
        # ),
        (
            "DC_EVAL_CUT_OFF_2",
            [r'\"Early_playout_termination\":\{\"Threshold\":([0-9]+(?:\.[0-9]+)?)'],
            "Bewertungsabbrüche",
            "Bewertungsschwellwert für sofortigen Abbruch"
        ),
        (
            "DC_IR",
            [
                r'"Search_depth":\s*\d+.*?"Epsilon":\s*([\d.]+)',
            ],
            "Vielversprechende Rollouts",
            "Wahrscheinlichkeit für Minimax-Suche"
        ),
        (
            "DC_IC",
            [
                r'"Max_depth":([0-9]+(?:\.[0-9]+)?)',
            ],
            "Vielversprechende Abbrüche",
            "Rollout-Tiefe"
        ),
        (
            "DC_IC_2",
            [
                r'"Max_depth":([0-9]+(?:\.[0-9]+)?)',
            ],
            "Vielversprechende Abbrüche",
            "Rollout-Tiefe"
        ),
        (
            "DC_K_BEST",
            [r'\"Rollout_selection\":\{"K":([0-9]+(?:\.[0-9]+)?)'],
            "K-Beste",
            "Anzahl berücksichtigter Züge"
        ),
        (
            "DC_IR",
            [
                r'"Search_depth":\s*\d+.*?"Epsilon":\s*([\d.]+)',
            ],
            "Vielversprechender-Rollout",
            "Abbruchschwellwert"
        ),
        (
            "DC_IR_2",
            [
                r'"Search_depth":\s*\d+.*?"Epsilon":\s*([\d.]+)',
            ],
            "Vielversprechender Rollout",
            "Wahrscheinlichkeit für Minimax-Suche"
        ),
        (
            "DC_IC",
            [
                r'"Max_depth":([0-9]+(?:\.[0-9]+)?)',
            ],
            "Vielversprechende-Abbrüche",
            "Abbruchtiefe"
        ),
        (
            "LAC_GREEDY",
            [r'\"Rollout_selection\":\{\"Epsilon\":([0-9]+(?:\.[0-9]+)?)'],
            "Greedy",
            "Wahrscheinlichkeit für zufällige Zugwahl"
        ),
        (
            "LAC_EPT",
            [r'\"Early_playout_termination\":\{\"Max_depth\":([0-9]+(?:\.[0-9]+)?)'],
            "Frühzeitige Abbrüche",
            "Rollout-Tiefe"
        ),
        (
            "LAC_CAPTURE_PREF",
            [r'\"Rollout_capture\":([0-9]+(?:\.[0-9]+)?)'],
            "Schlagpräferenz",
            "Schlagwahrscheinlichkeit"
        ),
        (
            "LAC_CORRECTIVE",
            [r'\"Rollout_selection\":\{\"Bound\":([0-9]+(?:\.[0-9]+)?)'],
            "Korrektur",
            "Bewertungsschwellwert für sofortige Zugwahl"
        ),
        (
            "LAC_EVALUATION_CUT_OFF",
            [r'\"Early_playout_termination\":\{\"Threshold\":([0-9]+(?:\.[0-9]+)?)'],
            "Bewertungsabbrüche",
            "Bewertungsschwellwert für sofortigen Abbruch"
        ),
        (
            "LAC_K_BEST",
            [r'\"Rollout_selection\":\{"K":([0-9]+(?:\.[0-9]+)?)'],
            "K-Beste",
            "Anzahl berücksichtigter Züge"
        ),
        (
            "LAC_IC",
            [
                r'\"Early_playout_termination\":\{\"Max_depth\":([0-9]+(?:\.[0-9]+)?)',
                r'\"Search_depth\":([0-9]+(?:\.[0-9]+)?)',
            ],
            "Vielversprechende Abbrüche",
            ["Minimax-Suchtiefe", "Rollout-Tiefe"]
        ),
        
        (
            "LAC_IC_3",
            [
                r'\"Early_playout_termination\":\{\"Max_depth\":([0-9]+(?:\.[0-9]+)?)',
                r'\"Search_depth\":([0-9]+(?:\.[0-9]+)?)',
            ],
            "Vielversprechende Abbrüche",
            ["Minimax-Suchtiefe", "Rollout-Tiefe"]
        ),
        (
            "LAC_IR",
            [
                r'\"Rollout_selection\":\{\"Search_depth\":([0-9]+(?:\.[0-9]+)?)',
                r'\"Search_depth\":[0-9]+ \"Epsilon\":([0-9]+(?:\.[0-9]+)?)',
             ],
            "Vielversprechende Rollouts",
            ["Wahrscheinlichkeit für Minimax-Suche", "Minimax-Suchtiefe"]
        ),
        (
            "LAC_IR_2",
            [
                r'\"Rollout_selection\":\{\"Search_depth\":([0-9]+(?:\.[0-9]+)?)',
                r'\"Search_depth\":[0-9]+ \"Epsilon\":([0-9]+(?:\.[0-9]+)?)',
             ],
            "Vielversprechende Rollouts",
            ["Wahrscheinlichkeit für Minimax-Suche", "Minimax-Suchtiefe"]
        ),
    ]
    # table header
    # param1, w-w, r, w-b, median_rollouts
    for name, reg, title, x_axis_label in diagrams:
        ww = []
        r = []
        bw = []
        rollouts = []
        param = []
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
                ww.append(value[0][0])
                bw.append(value[1][0])
                r.append(value[0][1]+value[1][1])
                param.append(ucb_c)
                medi = statistics.median(value[2])
                rollouts.append(medi)
                print(name, ucb_c, win_percentage, medi)
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
        plt.rcParams.update({
            'font.size': 10,        
            'axes.titlesize': 10,
            'axes.labelsize': 9,
            'xtick.labelsize': 8,
            'ytick.labelsize': 8,
            'legend.fontsize': 8,
            'figure.figsize': (3.51, 2.5)
        })
        plt.savefig(f"{name}.pdf", format="pdf", bbox_inches="tight")
        plt.close()
        data = pd.DataFrame({
                x_axis_label: param,
                "W-W": ww,
                "R": r,
                "B-W": bw,
                r"Rollout\textbackslash s": rollouts
        })
        latex_table = data.to_latex(
            index=False,
            caption=title, # Dein Titel
            float_format="%.2f",
            label="tab:mcts_ergebnisse", # Für Querverweise im Text (\ref{tab:mcts_ergebnisse})
            position="ht",               # Platzierung in LaTeX (here, top)
            column_format="lrrrr",         # Ausrichtung der Spalten   
        )
        with open("appendix.txt", "a") as f:
            f.write(latex_table)
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
            print(name, param1, param2, win_percentage, statistics.median(value[2]))
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
    # dpc ir = 0.04 dpc ic = 0.15
    # lac ir = 0.012
    width = 0.012  # total horizontal spread
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
    plt.legend(title=x_axis_label[1],loc='upper center', bbox_to_anchor=(0.5, -0.3), 
        ncol=2, fontsize=8)
    plt.savefig(f"{name}.pdf", format="pdf", bbox_inches="tight")
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