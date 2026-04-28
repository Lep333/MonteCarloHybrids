# Monte Carlo Tree Search-Hybrids
Dieses Projekt testet Monte Carlo Tree Search Hybride
in Spielen mit partiellen Informationen (siehe Abschlussarbeit.pdf). Dafür wurden 
die Spiele Dark Pawn Chess (5x5), Dark Los Alamos Chess 
(6x6) und Dark Chess (8x8) implementiert. Diese Spiele
befinden sich in src/chess_variation und implementieren
die Schnittstelle ChessVariation. 
ChessVariation bietet folgende Funktionen:
- GetPossibleMoves(): mögliche Züge der aktuellen Spielerin
- Hash(), ViewHash(): geben einen Zobrist Hash zurück
- ExecuteMove(move): führt den mitgegebenen move auf dem Spielbrett 
- GameOver(): überprüft ob das Spielende erreicht wurde

In src/player befindet sich POMCP und die Hybride.
POMCP befindet sich in der Datei: pomcp.go, die Hybride befinden sich in selection.go.

In main.go können folgende Einstellungen getroffen
werden:
- die Spielauswahl
- POMCP Einstellungen:
    - Termination_parameter: die verfügbare Zeit pro
    Zugwahl
    - Hybride: settings.Early_playout_termination, settings.Rollout_selection: nutzen um POMCP mit
    Hybriden zu kombinieren
- Simulationsanzahl: iterations

Mit `go build .` das Projekt kompilieren.

Dann mit `./monte_carlo_hybrids`die Simulationsausführung starten. Die Simulationen 
und ihre Einstellungen werden dann als csv Datei in 
/results gespeichert, dabei hängt der Name der 
Datei von der Variable `name` aus main.go ab.

Die Python-Datei `results.py` erzeugt bei Ausführung
ein Siegesraten Diagramm und liefert in der Konsole
Informationen zu Sieg/Niederlagen und Rollouts/s, 
außerdem werden die Daten als Latex Tabelle in appendix.txt gespeichert.

Bevor `results.py` mit `python3 results.py` ausgeführt werden
kann, müssen die benötigten Abhängigkeiten mit 
`pip3 install -r requirements.txt` installiert werden.
Dabei muss die aktive Directory /src sein! 

## English
This is a project that researches the usecase of
MCTS in games with partial informations. The different
games are implemented in src/chess_variation. The player
models are in src/players.

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