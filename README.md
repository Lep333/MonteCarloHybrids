# Monte Carlo Tree Search-Hybrids
This is a project that researches the usecase of
MCTS in games with partial informations. The different
games are implemented in /chess_variation. The player
models are in /players.

## Run it
To run different chess simulations do the following.
1) Change directory into src:
```
cd src/
```
2) Then compile the go program:
```
go build .
```
3) Then run the program:
```
./monte_carlo_hybrids
```
The statistics of the games will be saved in
the results.csv. The python program results.py
will create diagrams based on your results.

## Settings
You can modify the main.go to play the chess variation
with the players you want.