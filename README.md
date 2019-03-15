# Baseball Statistics

Baseball is a great game for lovers of statistics.  Major League Baseball has a great website filled with raw statistics, and the goal of this project is to make it fairly easy to use them.  This project provides you with tools that allow you to download the raw data and then parse it into a format that is useful for direct statistical analysis with common data science tools, or to import it into a database.

## Overview
This is a set of command line tools written in the Go language and licensed under the Apache 2.0 open source license.  Please enjoy and remember that Major League Baseball keeps this amazing wealth of data on a wide-open web server, asking only that you abide by their terms of service.

## Tools
* getStats - This tool downloads and processes the stats for you.  You can specify the date or dates that you want using flags:
    * -beg - the beginning date, in YYYYMMDD format (20170401)
    * -end - the ending date, in YYYYMMDD format (20170401)

### game.xml

The game.xml file contains basic information about each game, such as the time and location of the game as well as the teams involved.  Included in this file is a MLB-provided "primary key" for that game.  For our primary key, we'll take that and the date.  There are 3 output files when game.xml is parsed:
* gameInfo.dat
* gameTeams.dat
* gameStadium.dat

##### TODO

* Load the master_scoreboard.xml file for each day and use this as the source for all game data.  I probably don't even need to grab the game.xml file at all.  This appears to have the date and PK for each game, which can be passed to the other game files that don't include these values, such as game_events.xml
Example: http://gd2.mlb.com/components/game/mlb/year_2017/month_07/day_23/master_scoreboard.xml

* Combine the dateConvert and the dateInput stages into a single stage