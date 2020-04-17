# Minesweeper 
This API provides services to create and play the Minesweeper game.

## Getting started

Run in a local environment
```
clone this repository
go run main.go
```

Run with docker
```
docker build -t minesweeper .
docker run -p 8080:8080 -d minesweeper
```

## API

### Create
Creates a new game with a specific configuration

Method: POST 

    /games
do
Body
```json
{
    "rows": 10,
    "columns": 10,
    "bombs": 30
}
```

Response

```json
{
    "id": "fba225fb-07bc-4ad5-9ab9-6756d135ef9c",
    "board": {
        "squares": [
            [
                {
                    "type": 0,
                    "revealed": false,
                    "marked": false
                },
                ...
            ],
            ...
        ],
        "status": "new"
    },
    "started_at": 0,
    "elapsed_time": 0
}
```

| id                        | string                                         | the game unique id                                                                                                                                       |   |   |
|---------------------------|------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------|---|---|
| board                     | object                                         | the board with the squares                                                                                                                               |   |   |
| board.squares             | matrix of objects                              | a matrix of squares                                                                                                                                      |   |   |
| board.squares[x].type     | int enum {1 | 2}                               | 1. represents an empty square 2. represents a square with a bomb                                                                                         |   |   |
| board.squares[x].revealed | bool                                           | indicates whether the square has been revealed                                                                                                           |   |   |
| board.squares[x].marked   | bool                                           | indicates whether the square has been marked with a question symbol                                                                                      |   |   |
| board.status              | string enum {"new", "won", "lost", "on_going"} | - new: the game has not been started yet - won: the game has been won  - lost: the game has been lost  - on_going: the game has started but not finished |   |   |
| game.started_at           | int                                            | the timestamp when the game has started                                                                                                                  |   |   |
| game.elapsed_time         | int                                            | the seconds that has been elapsed since the game began                                                                                                                 |   |   |

### Get
Get a game by id

Method: GET 

    /games/:id

### Mark square
Add a question mark to a square

Method: UPDATE 

    /games/:id/mark-square

Body
```json
{
    "row": 2,
    "column": 3
}
```

### Play square
This is the endpoint to start playing the game. Reveals the square and set the game status.  

Method: UPDATE 

    /games/:id/play-square

Body
```json
{
    "row": 2,
    "column": 3
}
```

## Notes
- I adopted an hexagonal architecture approach to separate the different layers. 
- Due to de lack of time, the persistance layer has been implemented as a local key value store. It can be easily changed to a DynamoDB by implementing the game Storage interface.
- All the endpoints returns an obfuscated json for the game entity to hide the bombs locations and internal state.