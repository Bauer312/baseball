# Baseball Statistics

Baseball is a great game for lovers of statistics.  Major League Baseball has a great website filled with raw statistics, and the goal of this project is to make it fairly easy to use them.  This project provides you with tools that allow you to download the raw data and then parse it into a format that is useful for direct statistical analysis with common data science tools, or to import it into a database.

## Overview
This is a set of command line tools written in the Go language and licensed under the Apache 2.0 open source license.  Please enjoy and remember that Major League Baseball keeps this amazing wealth of data on a wide-open web server, asking only that you abide by their terms of service.

## Tools
* getDate - This tool downloads the overview page for a date or dates.  You can specify the date or dates that you want using flags:
    * -date - use words, such as 'today', 'yesterday', 'lastweek'
    * -beg - beginning date, in YYYYMMDD format (20170401)
    * -end - ending date, in YYYYMMDD format (20170407)
* getGames - This tool downloads 4 game detail pages for each game on each date or dates.  The game files are 'game.xml', 'game_events.xml', 'inning_all.xml', and 'inning_hit.xml'.  You can specify the date or dates that you want using flags:
    * -date - use words, such as 'today', 'yesterday', 'lastweek'
    * -beg - beginning date, in YYYYMMDD format (20170401)
    * -end - ending date, in YYYYMMDD format (20170407)