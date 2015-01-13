Roadmap
=======

2015 Nov 12
-----------
Just completed CRD operations on Workouts/Exercises/Sets; updating can wait.
The next step should be the fastest possible route to storing data. The three
identified data sources are through the UI and client-side scripts[1]. Options
will be judged by their ability to start including customers, and the quality
of their feedback loop.

#### Client library
Pros:
* Can be used in acceptance tests
* Faster to write than a working UI
* If the .wkt format is finished, then I can upload my existing workouts.
* Jon would probably try it out.
Cons:
* Exclusive to me (and maybe Jon)

Steps
* Parse workout in .wkt syntax
    * Parse .wkt syntax
* Send to server
    * Write python client library
    * Write swagger.yaml
    * Write server's API handlers

#### UI
Pros:
* Allows me to start demoing to potential customers.
Cons:
* Slow start; need to learn front-end

Steps
* Create basic page
* Write Swagger.yaml entry
* Write js call to server
* Write server API handler

[1]: With the API still immature, writing an Android client will provide slow
feedback. Also, every experiment with Android has been slow-going (much like
Go). Also, Android is exclusionary. Mobile will not be considered at this time.
