# Condition

## Types

| **Name**          | **Description**                                                             | **Args**                       |
|-------------------|-----------------------------------------------------------------------------|--------------------------------|
| `Contains`        | Pass when choice from `Choose` in choices from `ChooseChain`.               | `ChooseChain`, `Choose`        |
| `Fail`            | Condition always fails.                                                     |                                |
| `ManaAbove`       | Pass when `ChoosePlayer` mana is above `Amount`.                            | `ChoosePlayer`, `Amount`       |
| `ManaBelow`       | Pass when `ChoosePlayer` mana is below `Amount`.                            | `ChoosePlayer`, `Amount`       |
| `Match`           | Pass when `ChooseA` matches `ChooseB`.                                      | `ChooseA`, `ChooseB`           |
| `MatchDamageType` | Pass when `DamageType` from  `EventContext` matchest provided `DamageType`. | `EventContext`, `DamageType`   |
| `StatAbove`       | Pass when `ChooseCard`'s `Stat` is above `Amount`.                          | `ChooseCard`, `Stat`, `Amount` |
| `StatBelow`       | Pass when `ChooseCard`'s `Stat` is below `Amount`.                          | `ChooseCard`, `Stat`, `Amount` |
| `UnitMissing`     | Pass when `ChooseUnit` is not on the board.                                 | `ChooseUnit`                   |

## Args

| **Name**           | **Requirements**                                             |
|--------------------|--------------------------------------------------------------|
| `Amount`           | An integer.                                                  |
| `ChooseChain`          | List of [Choose](./choose.md).                           |
| `Choose{X}`        | [Choose](./choose.md).                                       |
| `Stat`             | `Cost`, `Attack`, `Health`, `Movement`, `Cooldown`, `Range`. |
