# wysci
Wysci is a library for executing simple SQL commands against a database.
The data can then be retrieved from the wysci server at the configured endpoint.
For the time being it will support Postgres.

## Motivation
This is a simple service that serves as a basis for exploring API design concepts.
It exposes a set of simple endpoints that retrieve data and return it using CSV.
Without being overly complicated, it provides an API that can illustrate good design practices

### Side Note
There's a certain percentage of projects I've worked on where the end users only wanted to get the data into Excel.
The web page and the application were nice, but what they really wanted to was to open the data in Excel.
Several times this need was expressed at the tail end of the project.
I wind up adding an "open in Excel" button onto the page.
So here is a simple service that takes a query and returns something that opens in Excel.

## Configuring
The wysci server is meant to be configuration driven.
The configuration file, expressed in [toml](https://github.com/toml-lang/toml), defines the queries and endpoints for the server.

### TOML File
The toml file is located (by default) in the current working directory and named `wysci.toml`.
To use a file located elsewhere (e.g. under something like `/etc/wysci/hr.toml`), you can pass the `-config` command line parameter.
You can also set the `CONFIG` environment variable.

### Database Configuration
Wysci connects to a Postgres datbase.  
The database connection information can either be placed in the toml file, passed through environment variables, or on the command line.
The recommendation for database username and password are to pass them on the command line or through environment variables.

| Setting          | Environment Variable | Command Line |
|------------------|----------------------|--------------|
|Database Host     |DBHOST                | -dbhost      |
|Database User     |DBUSER                | -dbuser      |
|Database Name     |DATABASE              | -database    |
|Database Password |DBPASS                | -dbpass      |
|Database Port     |DBPORT                | -dbport      |

### Queries
SQL queries as an SQL statement and possible parameters.
The query can include parameters. 
These are defined using the `$n` format where 1, 2, 3, etc. indicate the first, second, third parameters and so on, respectively.


### Endpoints
Endpoints define service endpoints for specific queries.
These are exposed as `/api/v1/\[name\]` endpoints.

#### Parameters
Defines the parameters to be passed to the endpoint.
These parameters are then passed to the configured query.

|Option  |Description                                     |
|--------|------------------------------------------------|
|type    |The type of the value (e.g. number, string,etc) |
|ordinal |The position in the query                       |

#### Headers
The endpoint can return additional headers.
For example, coercing the type to `text/csv` so the browser will open the spreadsheet program.
