# Mysports to Fit
A small tool to download an activity in Fit format from the closing service mysports.tomtom.com.

## Motivation
I have more then 670 acitivities stored on TomTom Mysports and this service will close at 30.09.2023.
I am interested on the fit export, but it seems only available if you click each activity and then the share button. The full export of the App looks like that exports only gpx files. 
Converting the gpx file to a Fit one looses some information, for example the lap information.

I don't want to download hundred of activities manually, so I have created this small command line utility in golang. 

## Activity list
MysportsToFit needs an array of activity numbers so that the Scraper Colly
can download all the Fit files automatically.
Where is this activity list? I got it using the Network Har log file from the browser.
The Har file will be generated if you scroll the activity list using the Network tab inside the Developer Tool of the browser (F12 key to activate it in my Brave browser). 
You scroll from the top to the bottom inside the activity list and your Network list will be filled out with all activity requests. Then export the Network session into an Har file.
In this file the interesting part are lines like this:

    "url": "https://mysports.tomtom.com/service/webapi/v2/activity/253066833/thumbnail.png?dv=3.4&width=54&height=54&colour=cccccc&metric=position&thickness=2",
Here the interesting information is the id number 253066833, the activity id. 
This id is needed to download an activity in Fit format.

Why not scraping directly the Web App mysports.tomtom.com? Because the Web App is built using Angular
and Javascript. Scraping Angular applications is tricky and only few days are left for the download.
Parsing the Hal file in order to create an activity list is a faster choice for me.  

## Build and use
Clone the repository and create a config.toml file with your mysports credentials (email and password).
Then, if you know one of your activity ids, you can try to download it in Fit format.
The Powershell process: 

    git clone https://github.com/aaaasmile/mysportsToFit.git
    cd mysportsToFit
    go mod tidy
    go build
    mkdir dest
    cp config_example.toml config.tom
    .\mytom.exe -idonly 205727472

## Status
In progress.