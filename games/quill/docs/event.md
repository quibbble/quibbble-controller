# Event

## Types

| **Name**               | **Description**                                                                                        | **Args**                                                  | **Hookable** |
|------------------------|-------------------------------------------------------------------------------------------------------|-----------------------------------------------------------|--------------|
| `AddItemToUnit`        | Adds `ChooseItem` from `ChoosePlayer`'s hand to `ChooseUnit`.                                         | `ChoosePlayer`, `ChooseItem`, `ChooseUnit`                | T            |
| `AddTraitToCard`       | Adds `Trait` to `ChooseCard`.                                                                         | `Trait`, `ChooseCard`                                     | F            |
| `AttackUnit`           | `ChooseUnit` attacks `ChooseDefender`.                                                                | `ChooseUnit`, `ChooseDefender`                            | T            |
| `BurnCard`             | Trash the card at the top of `ChoosePlayer`'s deck.                                                   | `ChoosePlayer`                                            | T            |
| `Cooldown`             | decrement `ChooseUnits` cooldowns by 1.                                                               | `ChooseUnits`                                             | F            |
| `DamageUnit`           | Apply `Amount` damage of type `DamageType` to `ChooseUnit`.                                           | `Amount`, `DamageType`, `ChooseUnit`                      | T            |
| `DamageUnits`          | Apply `Amount` damage of type `DamageType` to all `ChooseUnits`.                                      | `Amount`, `DamageType`, `ChooseUnits`                     | F            |
| `DiscardCard`          | Discard `ChooseCard` in `ChoosePlayer`'s hand.                                                        | `ChoosePlayer`, `ChooseCard`                              | T            |
| `DrainBaseMana`        | Drain `ChoosePlayer`'s base mana by `Amount`                                                          | `ChoosePlayer`, `Amount`                                  | T            |
| `DrainMana`            | Drain `ChoosePlayer`'s mana by `Amount`.                                                              | `ChoosePlayer`, `Amount`                                  | T            |
| `DrawCard`             | Place the top card in `ChoosePlayer`'s deck into their hand.                                          | `ChoosePlayer`                                            | T            |
| `EndGame`              | Sets the the winner as `ChooseWinner` and ends the game.                                              | `ChooseWinner`                                            | F            |
| `EndTurn`              | Ends the current turn.                                                                                |                                                           | T            |
| `GainBaseMana`         | Adds `Amount` to `ChoosePlayer`'s base mana.                                                          | `ChoosePlayer`, `Amount`                                  | T            |
| `GainMana`             | Adds `Amount` to `ChoosePlayer`'s mana.                                                               | `ChoosePlayer`, `Amount`                                  | T            |
| `HealUnit`             | Heals `ChooseUnit` by `Amount`                                                                        | `ChooseUnit`, `Amount`                                    | T            |
| `HealUnits`            | Heals `ChooseUnits` by `Amount`                                                                       | `ChooseUnits`, `Amount`                                   | F            |
| `KillUnit`             | Removes `ChooseUnit` from the board and resets the card.                                              | `ChooseUnit`                                              | T            |
| `ModifyUnit`           | Modifies `ChooseUnit`'s `Stat` by `Amount`.                                                           | `ChooseUnit`, `Stat`, `Amount`                            | T            |
| `ModifyUnits`          | Modifies all `ChooseUnits` `Stat`s by `Amount`.                                                       | `ChooseUnits`, `Stat`, `Amount`                           | F            |
| `MoveUnit`             | Moves `ChooseUnit` to `ChooseTile`. Decrements `ChooseUnit` if `UnitMovement` is true.                | `ChooseUnit`, `ChooseTile`, `UnitMovement`                | T            |
| `PlaceUnit`            | Move `ChooseUnit` from `ChoosePlayer`'s hand and places it on `ChooseTile` with option `InPlayRange`. | `ChoosePlayer`, `ChooseUnit`, `ChooseTile`, `InPlayRange` | T            |
| `PlayCard`             | Plays `ChooseCard` in `ChoosePlayer`'s hand.                                                          | `ChoosePlayer`, `ChooseCard`                              | T            |
| `RecallUnit`           | Removes and resets `ChooseUnit` on board and adds it back to owner's hand/                            | `ChooseUnit`                                              | T            |
| `RecycleDeck`          | Shuffles `ChoosePlayer`'s discard and sets it to their deck.                                          | `ChoosePlayer`                                            | T            |
| `RefreshMovement`      | Sets `ChooseUnits` movement back to base values.                                                      | `ChooseUnits`                                             | F            |
| `RemoveItemFromUnit`   | Removes `ChooseItem` from `ChooseUnit`.                                                               | `ChooseItem`, `ChooseUnit`                                | T            |
| `RemoveTraitFromCard`  | Removes `ChooseTrait` from `ChooseCard`.                                                              | `ChooseTrait`, `ChooseCard`                               | F            |
| `RemoveTraitsFromCard` | Removes all `ChooseTraits` from `ChooseCard`.                                                         | `ChooseTraits`, `ChooseCard`                              | F            |
| `SackCard`             | Discards `ChooseCard` from `ChoosePlayer`'s hand and applies `SackOption`.                            | `ChoosePlayer`, `ChooseCard`, `SackOption`                | T            |
| `SummonUnit`           | Summons card `ChooseID` for `ChoosePlayer` on `ChooseTile` with option `InPlayRange`.                 | `ChoosePlayer`, `ChooseTile`, `ChooseID`, `InPlayRange`   | T            |
| `SwapStats`            | Swaps `ChooseCardA` and `ChooseCardB` `Stat`.                                                         | `ChooseCardA`, `ChooseCardB`, `Stat`                      | T            |
| `SwapUnits`            | Swaps `ChooseUnitA` and `ChooseUnitB` locations.                                                      | `ChooseUnitA`, `ChooseUnitB`                              | T            |

## Args

| **Name**           | **Requirements**                                             |
|--------------------|--------------------------------------------------------------|
| `Amount`           | An integer.                                                  |
| `Choose{X}`        | [Choose](./choose.md).                                       |
| `DamageType`       | `Physical`, `Magic`, `Pure`.                                 |
| `ID`               | A five character string representing a card i.e. `U0001`.    |
| `InPlayRange`      | A boolean.                                                   |
| `Stat`             | `Cost`, `Attack`, `Health`, `Movement`, `Cooldown`.          |
| `Trait`            | [Trait](./trait.md).                                         |
