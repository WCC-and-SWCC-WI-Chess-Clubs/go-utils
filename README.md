
To build and run from the project root directory:

	go build -o chess-utils .
	./chess-utils postUpdate in.txt
	./chess-utils ratedEventByID 202602280803

Two intentional deviations from the Python:

ratedEventFromFile exits with the deprecation message instead of crashing (the underlying getSectionsFromHtml was already removed from reader.py)
All file-writing functions (clubEvents, winnersPage) use io.Writer parameters instead of redirecting stdout, which is idiomatic Go


Instructions for running the go utilities.   Finally all of the options explained :) There are two main use cases:
* Support the generation of HTML for tournament blog posts.
* Support the generation of HTML for website page content.

First step is to build; must have 'go' in your path: 

	go build -o chess-utils

The command line looks like this (from in the root dir):

`./chess-utils  <option> <parameter>`

All of the content is written to the console, so you can either copy the content from there, or redirect the output to a file. As always, you may have some cleanup to do if you grab content from the console...


`<option>` has the following choices:

* Weekly post generation options -
  * `postUpdate`
    * This is the main use case for these scripts - weekly update posts for in-flight events.
    * command line is `./chess-utils postUpdate <path to input text file>`
    * This expects an input file written by WinTD containing two sections (Open and Reserve, in that order) of X-table report data followed by two sections (Open and Reserve, in that order) of pairings data
    * To generate this file, within WinTD take the following steps:
      * Pair the upcoming round for both sections
      * Make sure the only open window in the app is the tournament window (check the "Windows" menu option)
      * Select both sections in the tournament window and click the "G" Games button in the ribbon (Two games windows will open)
      * Select the menu option "Reports -> Output to Window"
      * Select the menu option "Reports -> Print X-Table"
      * Select the menu option "Reports -> Print Game Windows"
    * All 4 sections will be in the report window (named something like "NONAME01.TXT")
      * Copy the contents to the file you use in the command line (I use a temp file called "in.txt" in the root directory, so my command line looks like: 
        * `./chess-utils postUpdate in.txt`
    * Run the script as noted above for `postUpdate` 
    * The table-formatted HTML for this week's post is written to the console or file if you've redirected it.
    * Paste that HTML into your blog post and you're done ! 
    *
    * NOTE 1 - This is smart enough to skip generating the pairings tables if you only send the 2 X-Table reports.  This is useful for last round updates when you don't have upcoming pairings.
    * NOTE 2 - This assumes 4 round tournaments.  This is hard-coded in the `processWinTDFile()` function in `main.py` (haven't had the need to make this cleaner.
  * `ratedEventByID`
    * Generates HTML for a completed, rated event.  Takes the event ID as its input parameter, and goes out to the USCF site to get the results.
      * `./chess-utils ratedEventByID 202507101092` (for example)
    *
    * NOTE 1 - This generates a complete page, including title and subtitle.  You will likely need to trim this for a blog post.
  * `ratedEventFromFile`
    * This option should not be used anymore, since the format of USCF website's crosstables changed (dramatically) in November 2025.
    * identical to `ratedEventByID` except pulls the USCF HTML source from a given input file, rather than getting it from the web.
      * The URL to get the HTML for an event from looks like this: https://www.uschess.org/msa/XtblMain.php?202507101092 where the value after the '?' character is your event ID
      * `./chess-utils ratedEventFromFile a-local-file-you-saved-for-event-xxxxx.html`
  * `ratedEventFromRatingReport`
    * similar to `ratedEventByID` and `ratedEventFromFile` except the input file is the text from a USCF generated rating report (not what's posted on the web) 
      * `./chess-utils file local-file-with-rating-report-text.txt`
  * `updatePostFromWinTdXTable`
    * This is one component of the `postUpdate` option that expects one section's winTD crosstable report file as input 
      * `./chess-utils winTD one-section-xtable-from-wintd.txt`
  * `updatePostFromWinTdPairings`
    * This is the other component of the `postUpdate` option that expects one section's winTD pairings report file as input 
      * `./chess-utils pairings one-section-of-pairings-from-wintd.txt`
* Static website page content options -
  * `clubEvents <clubAbbrev>`
    * This generates a page with links to all events rated for this affiliate - grouped by year.
    * The Affiliate ID used is determined by the club abbreviation passed in.
    * There is a facility to override displayed names, place the event ID and desired names in `data/<clubAbbrev>_event_names.json` as needed
    * The output is stored in `data/web/past_tournaments.html` - copy to the website repo
    * Needs to be periodically run to pick up recent events.
    * `./chess-utils clubEvents`
  * `winnersPage`
    * similar to `clubEvents` but builds a page grouped by year for all winners of the club championship and waukesha memorial
    * Uses the `data/winners.json` file to find the event IDs for these and get the winners from the web, if no ID is available, uses the winner listed in that file
      * This is because not all events are named consistently - this ensures we get all of them
    * The output is stored in `data/web/champions.html` - copy to the website repo
    * Needs to be periodically run to pick up recent events.
    * `./chess-utils winnersPage`
  * `generateEventsJS <clubAbbrev>`
    * Same idea as `clubEvents`, but for either WCC or SWCC, and emits a JS module instead of an HTML page
    * Takes the club abbreviation as its input parameter (e.g. `wcc` for WCC)
    * Same `data/<clubAbbrev>_event_names.json` override facility as `clubEvents`
    * The output is stored in `data/web/<clubAbbrev>_events.js` as `export const TOURNAMENTS = { '2026': [ { name, url }, ... ], ... }`, grouped by year, newest year and newest event first
    * `./chess-utils  generateEventsJS wcc`

# New USCF Ratings API
The new USCF ratings API is available. You can find documentation here: 

https://ratings-api.uschess.org/swagger/index.html

### Frequently Used Endpoints

* https://ratings-api.uschess.org/api/v1/rated-events/{eventId}
* https://ratings-api.uschess.org/api/v1/rated-events/202602280803

* https://ratings-api.uschess.org/api/v1/rated-events/{eventId}/sections/{sectionNumber}
* https://ratings-api.uschess.org/api/v1/rated-events/202602280803/sections/1

* https://ratings-api.uschess.org/api/v1/rated-events/{eventId}/sections/{sectionNumber}/standings

### Sample JSON

#### Get Rated Event
{
	"id": "202602280803",
	"name": "44th Annual Waukesha Memorial",
	"startDate": "2026-02-28",
	"endDate": "2026-02-28",
	"affiliate": {
		"id": "A5008948",
		"name": "WAUKESHA CHESS CLUB",
		"expirationDate": "2026-06-30",
		"status": "Active"
	},
	"sectionCount": 2,
	"playerCount": 54,
	"gameCount": 91,
	"stateCode": "WI",
	"city": "WAUKESHA",
	"zipCode": "53188",
	"countryCode": "USA",
	"isWomens": false,
	"status": "Rated",
	"sections": [
		{
			"id": "01KJKS52TXAHYKYDPVEM7SC803",
			"number": 1,
			"name": "Over 1200"
		},
		{
			"id": "01KJKS52TXNC4MB64ADB8G4C82",
			"number": 2,
			"name": "Reserve - under 1200"
		}
	],
	"firstRatedDate": "2026-02-28",
	"lastRatedDate": "2026-03-08",
	"createdDate": "2026-03-01"
}


#### Get Section for Rated Event
{
	"id": "01KJKS52TXAHYKYDPVEM7SC803",
	"name": "Over 1200",
	"number": 1,
	"playerCount": 36,
	"gameCount": 59,
	"roundCount": 4,
	"startDate": "2026-02-28",
	"endDate": "2026-02-28",
	"format": "Swiss",
	"timeControl": "G/60;d5",
	"ratingSystem": "D",
	"participantCoding": "NonScholastic",
	"isGrandPrix": true,
	"isExtendedGp": false,
	"isJuniorGp": false,
	"gpPoints": 0,
	"isTeamEvent": false,
	"isOnline": false,
	"shouldSubmitToFide": false,
	"isBlitz": false,
	"firstRatedDate": "2026-02-28",
	"lastRatedDate": "2026-03-07",
	"createdDate": "2026-03-01"
}

* https://ratings-api.uschess.org/rated-events/202403277302/sections/1/standings
