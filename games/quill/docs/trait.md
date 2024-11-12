# Trait

## Types

| **Name**   | **Description**                                                            | **Args**                    |
|-------------|----------------------------------------------------------------------------|----------------------------|
| `Aimless`   | Codex is randomized every turn.                                            |                            |
| `Assassin`  | Attacking a unit from behind deals `Amount` extra damage.                  | `Amount`                   |
| `BattleCry` | On place register `Hooks` and apply `Events`.                              | `Hooks`, `Events`          |
| `Berserk`   | Attacking and killing a unit resets cooldown to zero.                      |                            |
| `Buff`      | Apply buff to `Stat` by `Amount`.                                          | `Stat`, `Amount`           |
| `DeathCry`  | On death register `Hooks` and apply `Events`.                              | `Hooks`, `Events`          |
| `Debuff`    | Apply debuff to `Stat` by `Amount`.                                        | `Stat`, `Amount`           |
| `Dodge`     | Incoming attacks have a 1 in 3 chance of missing.                          |                            |
| `Enemies`   | Apply `Trait` to enemy `ChooseUnits`.                                      | `Trait`, `ChooseUnits`     |
| `Enrage`    | On taking damage register `Hooks` and apply `Events`.                      | `Hooks`, `Events`          |
| `Eternal`   | Item is passed to `ChooseUnit` if all `Conditions` are met.                | `Conditions`, `ChooseUnit` |
| `Execute`   | On attacking a unit if it's injured then kill it.                          |                            |
| `Friends`   | Apply `Trait` to friendly `ChooseUnits`.                                   | `Trait`, `ChooseUnits`     |
| `Gift`      | On attacking a unit gift the unit `Trait`.                                 | `Trait`                    |
| `Haste`     | On place or summon cooldown is set to zero.                                |                            |
| `Lobber`    | Ranged unit deals damage to target and all adjacent units.                 |                            |
| `Pillage`   | After damaging a base register `Hooks` and apply `Events`.                 | `Hooks`, `Events`          |
| `Poison`    | At owner turn end take `Amount` magic damage.                              | `Amount`                   |
| `Purity`    | Cannot be targeted by spells.                                              |                            |
| `Ranged`    | Can attack up to `Amount` spaces away. Do not take damage when attacking.  | `Amount`                   |
| `Recode`    | Apply `Code` to codex using `SetFunction`.                                 | `Code`, `SetFunction`      |
| `Shield`    | Mitigate `Amount` physical damage.                                         | `Amount`                   |
| `Spiky`     | Deal `Amount` extra damage when attacked.                                  | `Amount`                   |
| `Surge`     | Add mana amount to attack.                                                 |                            |
| `Thief`     | If attacked unit holds items then steal one randomly instead of attacking. |                            |
| `Tired`     | Does not cooldown at the end of owner's turn.                              |                            |
| `Ward`      | Mitigate `Amount` magic damage.                                            | `Amount`                   |

## Args

| **Name**           | **Requirements**                                                     |
|--------------------|----------------------------------------------------------------------|
| `Amount`           | An integer.                                                          |
| `Choose{X}`        | [Choose](./choose.md).                                               |
| `Conditions`       | list of [Condition](./condition.md).                                 |
| `Code`             | An eight character string containing only 0 or 1 i.e. `11001111`.    |
| `Events`           | List of [Event](./event.md).                                         |
| `Hooks`            | List of [Hook](./hook.md).                                           |
| `SetFunction`      | `Union`, `Intersect`, `Replace`.                                     |
| `Stat`             | `Cost`, `Attack`, `Health`, `BaseMovement`, `BaseCooldown`, `Range`. |
| `Trait`            | [Trait](./trait.md).                                                 
