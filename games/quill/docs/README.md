# Quill Docs

## Key Terms
- [Event](./event.md) - provides ways to affect change on the game state.
- [Condition](./condition.md) - provides requirements that must be passed before event(s) can happen.
- [Hook](./hook.md) - provides ways to trigger side affect event(s) before/after an event happens.
- [Choose](./choose.md) - provides a unified way to select units/tiles/players.
- [Trait](./trait.md) - added to cards and provide special characteristics to their holder.

## Card

Cards are created and stored in YAML files then parsed and converted into card objects in the game. This makes it easy to create and modify cards on the fly.

### Core

All cards contain the following values:

```yaml
ID: S0000
Name: Name
Description: "Description."
Cost: 3
Conditions: # list of conditions required to play this card.
Targets: # list of chooses that are validated against a list of targets passed by a user when playing a card.
Hooks: # list of hooks that are registered to the game engine on card play.
Events: # list of events that are performed on the game state on card play.
Traits: # list of special characteristics that the card holds.
```

### Items

Items are cards that may be held by units and have the following values:

```yaml
HeldTraits: # traits added to the unit that holds this item
```

### Spells

Spells perform affects on the game engine and state. They do no have any additional values.

### Units

Units are placed on the board. There are three types of units:
- Base: non-moving unit. Protect your bases while destroying your opponent's.
- Creature: a unit that may move, attack, and hold items.
- Structure: a unit that may not move, attack, or hold items. 

Units additionally have the following values:

```yaml
Type: Creature # Base / Creature / Structure
DamageType: Physical # Physical / Ranged (Physical SubType) / Magic / Poison (Magic SubType) / Pure
Attack: 2
Health: 2
Cooldown: 1 # how many turns until the unit may attack.
Movement: 1 # how many moves a unit may make a turn.
Codex: "11000000" # an eight character string representing directions a unit may move or attack `up|down|left|right|upper-left|lower-right|lower-left|upper-right`.
```
