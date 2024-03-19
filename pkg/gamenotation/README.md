# Quibbble Game Notation - QGN

Many different standards exist for simplified human and machine readable notation of games. Chess is a common example with notations such as Portable Game Notation (PGN) or Forsyth–Edwards Notation (FEN). While useful for chess, these notations have difficulties being applied to other games. Currently no single standard exists that can be applied to any game hence the introduction of Quibbble Game Notation (QGN) as a potential solution to this problem.

## Format

QGN is structured into two distinct sections, tags and actions.

### Tags

Tags are key value pairs used to describe the initial game state as well as to store additional meta data about the game.

#### Tag Format
```
[key "value"]
```

#### Tag Example
```
[key "carcassonne"]
```

#### Tag Requirements
There are two tags necessary for any game, Game and Teams.
- `key`: the name of the game represented. Ex: `[key "Carcassonne"]`
- `teams`: the list of teams playing the game. Ex: `[teams "a, b"]`

### Actions

Actions are an ordered list of actions teams take to create the current game state. Actions require the team, the action done, and any additional details needed to perform the action.

#### Action Format
```
{team index}{action character}&{action detail 1}.{action detail 2}...
```

#### Action Examples

```
0a&1.2 // team at index 0 of list in Teams Tag does action a with details 1 and 2
1b     // team at index 1 of list in Teams Tag does action b
```

### QGN Example
```
[key "carcassonne"]
[teams "a, b"]
[seed "123"]
[completed "false"]
[date "10-31-2021"]

0c 0a&1.2 0b&1.2.k.b 1c 1c 1c 1a&0.1 1b&0.1.m {you can add
comments like so} 0a&2.2 0b&2.2.t.l
```

## Usage

This package provides both a method of creating as well as parsing QGN text.

### Create QGN

```go
notation := &Snapshot{
    Tags: map[string]string{"key": "carcassonne", "teams": "a, b", "seed": "123"}
    Actions: []Action{
        {
            Index: 0,
            Action: 'a',
            Details: []string{"1", "2"}
        },
        {
            Index: 1,
            Action: 'c',
        },
    }
}
raw := qgn.String()
```

### Parse QGN
```go
raw := "[key \"carcassonne\"][teams \"a, b\"][seed \"123\"]0c 0a&1.2"
notation, err := Parse(raw)
```
