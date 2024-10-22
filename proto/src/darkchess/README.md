# Dark Chess

## Piece data in grid_position [0,1]

Turn down.

```json
{
  "board": [
    {
      "grid_position": {
        "x": 0, 
        "y": 1
      },
      "piece": "0",
      "is_revealed": false
    },
    {
      "grid_position": {
        "x": 0, 
        "y": 2
      },
      "piece": "0x275",
      "is_revealed": true
    }
  ]
}
```

Reveal a piece.

```json
{
  "board": [
    {
      "grid_position": {
        "x": 0,
        "y": 1
      },
      "piece": "0x110",
      "is_revealed": true
    },
    {
    "grid_position": {
      "x": 0,
      "y": 2
    },
    "piece": "0x275",
    "is_revealed": true
    }
  ]
}
```

After GENERAL_RED being captured by SOLDIER_BLACK.

And grid_position [0,1] is empty now.

```json
{
  "board": [
    {
      "grid_position": {
        "x": 0,
        "y": 1
      },
      "piece": "0x275",
      "is_revealed": true,
    }, {
      "grid_position": {
        "x": 0,
        "y": 2
      },
      "is_revealed": false
    }
  ]
}
```

## Reference

[Game State](https://game-soul-technology.atlassian.net/wiki/x/PgCIU)
