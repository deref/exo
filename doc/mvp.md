Initial MVP: Foreman replacement.

Prereq: You've already got a Procfile.

```
> cat Procfile
< api: node-dev ./server.js
< gui: vite
```

To run in the terminal, killed on `^c`:

```
> expo
> exo run Procfile
*** NCURSES APP RENDERS HERE ***
```

To start as a Deamon:

```
> exo start Procfile
< Console UI is available at http://localhost:5000/
```

```
> exo cui
*** NCURSES APP RENDERS HERE ***
```

When you open the web UI, you see:

![img_0196](https://user-images.githubusercontent.com/119164/125343720-406d9180-e30b-11eb-83ad-380036f7cbaa.jpg)

Which has the following features:

- List of Processes
- Shows running/stopped status. Click to toggle (stop/start).
- Shows process logs visibility. Click to toggle (hide/show).
- Shows logs, you can scroll.
- If you're at the bottom of the logs, live updates.
- If you scroll up, and new messages come in, some indicator shows you that there are new messages.


API:

```
> http /start-process name=web

> http /stop-process name=worker

> http /describe-processes
< { "processes": [
<   {"name": "web", "running": true},
<   {"name": "worker", "running": false}
< ] }

> http /get-events 'processes:=["web", "worker"]' before:=null
< { "events": [
<     { "sid": "12345",
<       "timestamp": "2021-07-12T20:18:39.965Z", 
<       "log": "web",
<       "message": "hello"
<     }
<   ],
<   "cursor": "abc123" }
```
