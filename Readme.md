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
Then, if you know one of your activity id (i.e. 205727472), you can try to download it in Fit format.
The Powershell process: 

    git clone https://github.com/aaaasmile/mysportsToFit.git
    cd mysportsToFit
    go mod tidy
    go build
    mkdir dest
    cp config_example.toml config.tom
    .\mysportsToFit.exe -idonly 205727472

Hard coded ids could be downloaded with (function UseHardCodedIds):

    .\mysportsToFit.exe -hardcoded
Result is

    go run .\main.go -hardcoded -chunksize 3  
    
    > using hard coded array activity ids 4
    > DownloadFit using target dir  ./dest
    > OnRequest CB for login...
    > Login ok
    > Request for 4 activities in download chunk size 3
    > [ix => 2] blocking in chunk 1 with 3 download processes
    > [205727472] start processing
    > [531746638] start processing
    > [108044668] start processing
    > [531746638] response received status 200 size 2700
    > [531746638] file written: ./dest/act_531746638.fit
    > [108044668] response received status 200 size 71909
    > [108044668] file written: ./dest/act_108044668.fit
    > [205727472] response received status 200 size 119081
    > [205727472] file written: ./dest/act_205727472.fit
    > [ix => 2] chunk 1 OK (size 3)
    > [ix => 3] blocking in chunk 2 with 1 download processes
    > [108044654] start processing
    > [108044654] response received status 200 size 71877
    > [108044654] file written: ./dest/act_108044654.fit
    > [ix => 3] chunk 2 OK (size 1)
    > Processed 4 Fit activities, time duration 723.1546ms
    > That's all folks!
## Status
In progress.