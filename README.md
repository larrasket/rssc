rssc provides a real-time, self-hostable regex-oriented customization service for atom,
rss and json feeds.


# Usage

rssc uses get request parameters to read user&rsquo;s preferences. following are the
supported parameters:

-   `src`: the source URL of the feeds. Example: <https://hnrss.org/newest>
-   `t`: the type of desired feeds (that rssc will return), can be `rss`, `atom` or `json`.
-   `descriptionf`: regular expression filter to filter feeds based on the description property.
-   `titlef`: regular expression filter to filter feeds based on the title property.
-   `contentf`: regular expression filter to filter feeds based on the content property.
-   `linkf`: regular expression filter to filter feeds based on the link property.
-   `net`: boolean (1 or 0) (default to 0 when omitted) whether to use the .NET
    engine regex instead of Go&rsquo;s.

All regex should be valid [RE2 regex syntax](https://github.com/google/re2/wiki/Syntax) (in case of `net`, see [MS&rsquo; Regular Expression Language](https://learn.microsoft.com/en-us/dotnet/standard/base-types/regular-expression-language-quick-reference)).


<a id="ex"></a>

# Examples

Replace <https://rssc.fly.dev/> in the following examples with your instance of
rssc (can be a `localhost:8080` if self-hosted).

-   Get all new feeds that contain the word &ldquo;emacs&rdquo; but not the word &ldquo;vim&rdquo;, and
    not posted by user lr0 (hnrss.org does not provide a creator/author property,
    so the author filter here is only for a demonstration purpose.):
    
    [https://rssc.fly.dev/rss?src=https://hnrss.org/newest?q=emacs&titlef=^(?!.\*\bvim\b).\*\bEmacs\b.\*&net=1&authorf=^(?!.\*\blr0\b).\*$](https://rssc.fly.dev/rss?src=https://hnrss.org/newest?q=emacs&titlef=^(?!.*\bvim\b).*\bEmacs\b.*&net=1&authorf=^(?!.*\blr0\b).*$)
    
    Do note that `net` is used here because RE2 does not support lookarounds.

-   Get geopolitical updates from your favorite hacker and skip all free software
    movement blessed propaganda:
    
    [https://rssc.fly.dev/rss?src=https://stallman.org/rss/rss.xml&contentf=^(?=.\*(?:palestine|syria|egypt|iraq|israel|algeria|morocco))(?!.\*(?:linux|gnu|software|programming|program)).\*$&net=1](https://rssc.fly.dev/rss?src=https://stallman.org/rss/rss.xml&contentf=^(?=.*(?:palestine|syria|egypt|iraq|israel|algeria|morocco))(?!.*(?:linux|gnu|software|programming|program)).*$&net=1)
-   Get BBC middle east updates, only ones that are related to Palestine:
    
    [https://rssc.fly.dev/rss?src=http://feeds.bbci.co.uk/news/world/middle\_east/rss.xml&titlef=(?i)palestine|palestinian|gaza](https://rssc.fly.dev/rss?src=http://feeds.bbci.co.uk/news/world/middle_east/rss.xml&titlef=(?i)palestine|palestinian|gaza)


# Self-host

To host rssc locally, install it using Go:

    go install github.com/larrasket/rssc@latest

Set environment variable for `PORT`, if user doesn&rsquo;t set a `PORT` value, rssc will use `:8080` by default.

    PORT=4040 rssc

Now rssc should be running at `localhost:4040`. You can use the same [examples](#ex)
provided, with replacing fly domain with your localhost.


# Notes

-   RE2 is much safer than the .NET Regex, that&rsquo;s to say, all the good features
    that .NET engine enables you to use come with a risk cost and ends up
    enabling [catastrophic backtracking](https://github.com/dlclark/regexp2#catastrophic-backtracking-and-timeouts), therefore a timeout of 5 seconds is currently
    enabled on the fly.dev instance when using `net=1`.
-   rssc is pre-alpha. Please report bugs [here](mailto:~lr0/public-inbox@lists.sr.ht) ([url](https://lists.sr.ht/~lr0/public-inbox))

