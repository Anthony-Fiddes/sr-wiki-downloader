# sr-wiki-downloader

This is a **s**ub**r**eddit **w**iki **d**ownloader. I started working on it
when I realized that the API existed (referenced in this
[post](https://www.reddit.com/r/DataHoarder/comments/ga2p8y/comment/foxdvju/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button)
), but before I realized that Reddit will rate limit you if you make more than
10 requests per minute when you're not logged in. That means it's pretty slow to
download a subreddit's wiki, but if you just want to have a copy on your
computer, it works.

# Usage

You can just clone the repo and use go run.

Example:

```
❯ go run . germany ./test
2024/03/01 19:42:59 Successfully downloaded r/germany/wiki/index
2024/03/01 19:43:05 Successfully downloaded r/germany/wiki/american-dream
2024/03/01 19:43:11 Successfully downloaded r/germany/wiki/assistantbot_statistics
2024/03/01 19:43:17 Attempt 0 to request r/germany/wiki/autobahn_safety was rate limited
2024/03/01 19:43:24 Successfully downloaded r/germany/wiki/autobahn_safety
2024/03/01 19:43:30 Successfully downloaded r/germany/wiki/benefits
2024/03/01 19:43:36 Successfully downloaded r/germany/wiki/black
2024/03/01 19:43:42 Successfully downloaded r/germany/wiki/blue-card
2024/03/01 19:43:49 Attempt 0 to request r/germany/wiki/brexit was rate limited
2024/03/01 19:43:55 Attempt 1 to request r/germany/wiki/brexit was rate limited
2024/03/01 19:44:01 Successfully downloaded r/germany/wiki/brexit
2024/03/01 19:44:07 Attempt 0 to request r/germany/wiki/children was rate limited
2024/03/01 19:44:13 Successfully downloaded r/germany/wiki/children
2024/03/01 19:44:20 Successfully downloaded r/germany/wiki/citizenship
2024/03/01 19:44:26 Attempt 0 to request r/germany/wiki/citizenship-detour was rate limited
2024/03/01 19:44:32 Successfully downloaded r/germany/wiki/citizenship-detour
2024/03/01 19:44:39 Successfully downloaded r/germany/wiki/config/description
...
```
