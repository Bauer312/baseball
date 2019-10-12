# Baseball Statistics

Baseball is a great game for lovers of statistics.  Major League Baseball has a great website filled with raw statistics, and the goal of this project is to make it fairly easy to use them.  This project provides you with tools that allow you to download the raw data and then parse it into a format that is useful for direct statistical analysis with common data science tools, or to import it into a database.

## Overview
This is a command line tool written in the Go language and licensed under the Apache 2.0 open source license.  Please enjoy and remember that Major League Baseball keeps this amazing wealth of data on a wide-open web server, asking only that you abide by their terms of service.

## Baseball
This tool downloads or processes data for you.  MLB has two data sites, Savant (the newest) and Gameday.  Specify which you want to pull data from along with information about desired dates and where you'd like the data to be stored.