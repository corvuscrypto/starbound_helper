#Starbound Server Helper
This program is just a way to retrieve a userlist and server status from the server.
The program just pings the server and reads the starbound server log file to determine
the state of the server and the users connected. This program was put together very quickly
so it's not the prettiest, but it works well.

## Configuring
For this there is a single configurable file "paths.cfg"

The variables within this are:

* starbound_log - the filepath to the starbound server log (should end in filename.log)
* output_directory - the directory path (should _not_ end in filename, just directory you want the log to be in)
* starbound_address - the host address of the starbound server
* starbound_port - the port on which the server is Running

## Running
Once the configuration is correct, run it from the shell as a background process:
```
sh starboundHelper&

or

./starboundHelper&
```

Either should work. You shouldn't need superuser privileges unless you are planning to write
the log to a protected folder. I recommend using stdin to log any fatal errors just in case.

## Killing the process

just use `kill` command to send an interrupt signal to it
I'm lazy to make a service CLI interface :P

## Retrieving JSON data

data is stored within a .json file named starbound_data.json. It will be in the specified output directory. Use your favorite web server to retrieve the file and use it or maybe even javascript to parse the JSON. Your call.
