# API Reference

## Route

### /api/v1/player

#### Overview

Return user data and different from last query.

#### Param

- `player` player id or player username

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
    "play_count": "13",
    "rank": "1",
    "pp": "727",
    "acc": "0.01%",
    "total_play": "20"
  }
}
```

At example behind, you can see diff fields, diff is the
data different between each query. Data can store in Redis
or MySQL.

