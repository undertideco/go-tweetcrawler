# go-tweetcrawler

This is a tweet crawler written in Go. It's simple and is probably nothing fancy. We're pretty new to Go and are open to any suggestions to improve our code so do feel free to submit an issue/PR if you see parts of the code that can be better.

## What Does It Do
The tweet crawler does only one thing - you enter a username, and it will crawl your target's tweets (to the latest 3200 due to a limit with Twitter's API) and then listens every 5 minutes to attempt a new crawl on the person's timeline.

The tweets are then populated under a csv, with the tweet's id, content and lastly date the tweet was posted.

## How To Use
1. Clone The Repository
2. Copy config.json.example to config.json
3. Fill up consumer/access secret and keys in the config.json taken from [Twitter](http://apps.twitter.com/)
4. Run
5. $$$

## License
The MIT License (MIT)

Copyright (c) [2015] [Undertide LLP]

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
