# mysports to Fit
A small tool to download an activity in Fit format from the closing service mysports.tomtom.com.

## Motivation
I have more then 600 acitivities stored on mysports. I am interested on the
fit export, but it seems only available if you click each activity and the the share button.
The full export looks like only for gpx files. Converting the gpx to fit looses some information (i.e. Lap).

## Usage
Create the file settings.json with your credential (email and password).
Then you need an array of activity numbers. In this case the Scraper Colly
can download the fit file automatically
Where is the activity list? I get it using the Network Har log file from the browser.
Then you can look for this string:

    "url": "https://mysports.tomtom.com/service/webapi/v2/activity/
and you get the activity number after that string

## Status
In progress.