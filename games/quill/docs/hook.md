# Hook

Hooks are objects registered to the game engine that trigger `Before`/`After` some future is triggered [Event](./event.md). Hooks are composed of [conditons](./condition.md), [events](./event.md), and in many cases [chooses](./choose.md).

Below is an example hook:

```yaml
When: After # Before / After
Types: # list of event types that will trigger the hook
- PlayCard
Conditions: # list of additional conditions that must be met
- Type: ManaBelow
  Args:
    Amount: 1
    ChoosePlayer:
      Type: CurrentPlayer
Events: # list of events to apply to the game on trigger and when conditions are met
- Type: DrawCard
  Args:
    ChoosePlayer:
      Type: CurrentPlayer
ReuseConditions: # list of conditions that must be met to keep this hook from being removed after trigger
- Type: Fail
```

This hook is triggered when a card is played. The engine then checks whether or not the current player's mana is 0. If the condition passes then the current player draws a card, if it fails then no card is drawn. After triggering, no matter the result of conditions, the reuse conditions are checked. In this case the fail condition is present so the engine cleans up this hook.
