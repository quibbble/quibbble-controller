# Quibbble K8s Controller

The project allows Quibbble games to be run in a K8s cluster. Games are spun up as individual pods with their entire lifecycle being by the Quibbble controller. This allows Quibbble to take advantage of the power of K8s, primarily the ability to scale and to seperate and self heal in the event of any single game failure.

## TODO

- Add storage

## Supported Games
- [carcassonne](/games/carcassonne/)
- [connect4](/games/connect4/)
- [indigo](/games/indigo/)
- [stratego](/games/stratego/)
- [tictactoe](/games/tictactoe/)
- [tsuro](/games/tsuro/)

## Architecture

There are three main processes in this system.
1. `Controller Service` - Processes game creation requests.
2. `Controller Cronjob` - Periodically searches and cleans stale games.
3. `Game Service` - Runs a game instance.

## Flows

### Game Creation
- Send `POST https://api.quibbble.com/create` with some `qqn` such as `[key "tictactoe"][id "example"][teams "red, blue"]`.
- `Controller Service` processes the request and, if valid, creates K8s ConfigMap, Pod, Service, and Ingress resources.
- Game can now be accessed at `https://api.quibbble.com/tictactoe/example`.

### Game Connection
- Join a game by connection to `wss://api.quibbble.com/tictactoe/example/connect` with websockets.
- Connection should be open to a `Game Service` instance and relevant game messages should be recieved.

### Game Cleanup
- `Controller Cronjob` will kick off every `X` timeperiod. 
- Job requests data from all live games by calling `GET https://api.quibbble.com/<key>/<id>/active` for each game.
- If there are no connected players and no recent updates then all K8s related resources are deleted.

## REST API

<details>
 <summary><code>POST</code> <code><b>/create</b></code> <code>(create a game)</code></summary>

##### Parameters

> | name      |  type     | data type                          | description                                                       |
> |-----------|-----------|------------------------------------|-------------------------------------------------------------------|
> | None      |  required | object ([QGN](/pkg/qgn/README.md)) | QGN descibing the game to create                                  |


##### Responses

> | http code     | content-type                      | response                                                            |
> |---------------|-----------------------------------|---------------------------------------------------------------------|
> | `201`         | `text/plain;charset=UTF-8`        | `Created`                                                           |
> | `400`         | `text/plain;charset=UTF-8`        | `Bad Request`                                                       |
> | `409`         | `text/plain;charset=UTF-8`        | `Conflict`                                                          |
> | `500`         | `text/plain;charset=UTF-8`        | `Internal Server Error`                                             |

##### Example cURL

> ```javascript
>  curl -X POST -H "Content-Type: application/qgn" --data @post.qgn https://api.quibbble.com/create
> ```
</details>

<details>
 <summary><code>WEBSOCKET</code> <code><b>/{key}/{id}/connect</b></code> <code>(connect to a game)</code></summary>

##### Parameters

> | name      |  type     | data type                          | description                                                       |
> |-----------|-----------|------------------------------------|-------------------------------------------------------------------|
> | key       |  required | string                             | The name of the game i.e. `tictactoe` or `connect4`               |
> | id        |  required | string                             | The unique id of the game instance to join                        |


##### Responses

> None

##### Example wscat

> ```javascript
>  wscat -c wss://api.quibbble.com/{key}/{id}/connect
> ```
</details>

<details>
 <summary><code>GET</code> <code><b>/{key}/{id}/snapshot?format={format}</b></code> <code>(get game snapshot)</code></summary>

##### Parameters

> | name      |  type     | data type                          | description                                                       |
> |-----------|-----------|------------------------------------|-------------------------------------------------------------------|
> | key       |  required | string                             | The name of the game i.e. `tictactoe` or `connect4`               |
> | id        |  required | string                             | The unique id of the game instance to join                        |
> | format    |  required | one of `json` or `qgn`             | The type of data to return                                        |


##### Responses

> | http code     | content-type                            | response                                                            |
> |---------------|-----------------------------------------|---------------------------------------------------------------------|
> | `200`         | `application/json` or `application/qgn` | JSON or [QGN](/pkg/qgn/README.md)                                   |
> | `400`         | `text/plain;charset=UTF-8`              | `Bad Request`                                                       |
> | `404`         | `text/plain;charset=UTF-8`              | `Not Found`                                                         |
> | `500`         | `text/plain;charset=UTF-8`              | `Internal Server Error`                                             |

##### Example cURL

> ```javascript
>  curl -X GET https://api.quibbble.com/{key}/{id}/snapshot?format=json
> ```
</details>

<details>
 <summary><code>GET</code> <code><b>/{key}/{id}/active</b></code> <code>(get game activity)</code></summary>

##### Parameters

> | name      |  type     | data type                          | description                                                       |
> |-----------|-----------|------------------------------------|-------------------------------------------------------------------|
> | key       |  required | string                             | The name of the game i.e. `tictactoe` or `connect4`               |
> | id        |  required | string                             | The unique id of the game instance to join                        |


##### Responses

> | http code     | content-type                            | response                                                            |
> |---------------|-----------------------------------------|---------------------------------------------------------------------|
> | `200`         | `application/json`                      | JSON data describing player count and last update time              |
> | `404`         | `text/plain;charset=UTF-8`              | `Not Found`                                                         |

##### Example cURL

> ```javascript
>  curl -X GET https://api.quibbble.com/{key}/{id}/active
> ```
</details>


## Websocket Messaging

### Sendable Messages

<details>
 <summary><code><b>join</b></code> <code>(join a team)</code></summary>

##### Message

```json
{
    "type": "join",
    "details": "$TEAM"
}
```
</details>

<details>
 <summary><code><b>action</b></code> <code>(perform a game action)</code></summary>

##### Message

```json
{
    "type": "$ACTION",
    "details": {...}
}
```
</details>

<details>
 <summary><code><b>chat</b></code> <code>(send a chat message)</code></summary>

##### Message

```json
{
    "type": "chat",
    "details": "$MESSAGE"
}
```
</details>


### Recievable Messages

<details>
 <summary><code><b>snapshot</b></code> <code>(retrieve a snapshot of the game)</code></summary>

##### Details

Message sent to all players on every game state change.

##### Message

```json
{
    "type": "snapshot",
    "details": {...}
}
```
</details>

<details>
 <summary><code><b>connection</b></code> <code>(retrieve player connection information)</code></summary>

##### Details

Message sent to all players on every player connection, drop, or team change.

##### Message

```json
{
    "type": "snapshot",
    "details": {
        "uid": "$UID1",
        "players": {
            "$UID1": "$TEAM1",
            "$UID2": "$TEAM2",
            "$UID3": null
        }
    }
}
```
</details>

<details>
 <summary><code><b>chat</b></code> <code>(retrieve chat message)</code></summary>

##### Details

Message sent to all players on every sent chat message.

##### Message

```json
{
    "type": "chat",
    "details": {
        "uid": "$UID",
        "team": "$TEAM",
        "message": "$MESSAGE",
    }
}
```
</details>

<details>
 <summary><code><b>error</b></code> <code>(retrieve error message)</code></summary>

##### Details

Message sent to origin player on failed action message.

##### Message

```json
{
    "type": "error",
    "details": "$MESSAGE"
}
```
</details>
