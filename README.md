# OSU API SERVER

[![Build Status](https://drone.avmtn.net/api/badges/avimitin/osuapiserver/status.svg?ref=refs/heads/master)](https://drone.avmtn.net/avimitin/osuapiserver)

This project is planned to make a API backend server that store and calculate osu player's data.

## Configuration

Config should store in `~/.config/osuapi/config.json`, or specific config path by env `osu_conf_path`.

### Config Example

See [config.json](./example/config/config.json)

## Route

### /api/v1/player

#### Overvier

Return user data and different from last query.

#### Param

- `player`: player id or player username

#### Example

```bash
curl -X POST \
-H "Content-Type: application/x-www-form-urlencoded" \
-d "player=shigetora" \
"http://localhost:11451/api/v1/player"
```

```text
Response
--------
{
  "latest_data": {
    "user_id": "3660913",
    "username": "Shigetora",
    "join_date": "2013-11-27 21:13:46",
    "count300": "5591498",
    "count100": "774185",
    "count50": "64236",
    "playcount": "33617",
    "ranked_score": "11858742843",
    "total_score": "52492353585",
    "pp_rank": "19235",
    "level": "100.256",
    "pp_raw": "0",
    "accuracy": "93.22009086608887",
    "count_rank_ss": "1",
    "count_rank_ssh": "0",
    "count_rank_s": "563",
    "count_rank_sh": "81",
    "count_rank_a": "543",
    "country": "FI",
    "total_seconds_played": "1453964",
    "pp_country_rank": "26426",
    "events": []
  },
  "diff": {
    "play_count": "0",
    "rank": "0",
    "pp": "0.000",
    "acc": "0.00%",
    "total_play": "0"
  }
}
```

At example behind, you can see diff fields, diff is the data different between each query.

