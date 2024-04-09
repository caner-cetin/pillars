# Pillars

or, Kartaca Case Study for Caner Cetin. Pillars is shorter, and cooler.

- [Pillars](#pillars)
  - [Disclaimer](#disclaimer)
  - [Installation](#installation)
  - [Pitfalls](#pitfalls)
  - [Demos, Explanations, and Screenshots](#demos-explanations-and-screenshots)
  - [Why?](#why)
    - [Why MongoDB?](#why-mongodb)
    - [Why website refreshes itself too often after starting?](#why-website-refreshes-itself-too-often-after-starting)


## Disclaimer

You will probably see that web page is completely normal on a computer screen, but mobile page is a mess. I am sorry for that, not the best frontend developer. 

## Installation

```bash
# you might need to make the script executable
# chmod +x scripts/start.sh (optional)
bash scripts/start.sh
```

that is it. now go to `localhost:333`. 

## Pitfalls 

- You need latest or modern browsers, I am using [js-bson](https://github.com/mongodb/js-bson) library of MongoDB which requires top-level-await feature. [Check browser compability table right here](https://caniuse.com/?search=top%20level%20await).

- If you are having issues mounting the volumes, especially on OSX, you might need to create `container-data` and `keys` folders manually and add them to the `file sharing` in Docker settings.

- If you want to see the data live, you can do 

```bash
❯ docker exec mongodb mongosh
Current Mongosh Log ID:	66139dd81bff5e5777403e98
Connecting to:		mongodb://127.0.0.1:27017/?directConnection=true&serverSelectionTimeoutMS=2000&appName=mongosh+2.2.2
``` 

to get the connection string and use it in MongoDB Compass along with the user&password in the `.env` file.

## Demos, Explanations, and Screenshots

When you first navigate to `localhost:333`, you will see our one and only page, the index:

![alt text](./examples/image.png)

Frontend is really really simple. Screen is divided into two. 

- Google Maps
- Info Panel
  - Wipe Earthquakes
  - Submit Earthquakes
  - Connected / Reconnect
  - Latest Earthquakes Table.

Only 6 main components. Starting from `Submit Earthquakes`,

![alt text](./examples/image-2.png)

We have two script. First, `Feed a Single Earthquake` and the second, `Feed Earthquakes from USGS Dumps`. Yes, there is no `random earthquake input in random intervals`. I cannot simulate good enough data to make our screenshots, or your viewing experience better.  

Submitting our first dialog:



## Why?

### Why MongoDB?

- change streams

as the Apache Beam pipeline works and Flink executes, we will update the Mongo collection. Now, frontend works by connecting to a websocket endpoint to a route. How does that route knows data in the earthquake collection is updated and should send the new data to the frontend? Well, MongoDB change streams!

https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/read-operations/changestream/

### Why website refreshes itself too often after starting?

Quasar w/ Vite does that in cold start, unfortunately. See => https://github.com/quasarframework/quasar/issues/12933