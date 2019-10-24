# Starting the Application
The POC helped clarify several unknowns.
More importantly, I have [answers](https://github.com/DarcInc/wysci/blob/master/docs/version0.md#the-application) to some questions I had when starting the POC.
With the POC completed, I'm ready to start laying out the initial application.

The key problem was safely casting a slice of `interface{}` elements to the correct pointer types and correctly handling `null` values.
Essentially, I need to construct a buffer for the result.
The elements in the buffer need to be typed correctly for the SQL result type.
I started down one possible approach in the POC, but I don't intend to use it.
The code feels unnecessarily verbose.

Second to that problem was how to define endpoints and queries.
The solution seemed obvious at a high level.
When I tried validating parameters in the POC, it wasn't clear where that responsiblity should rest.
Adding parameter definitions in the endpoint simplifies validation.
However, it's the query that knows what types it should expect.
Also, my first attempts at defining type information were very verbose.
Clearly, this will take a little more analysis.

## Refining the Feature Set
When I started on the POC I had a rough idea of what I wanted to build.
I wasn't sure about the details, but I wanted to specify everything in the [TOML configuration file](https://github.com/DarcInc/wysci/blob/v0.0/docs/example.toml).
As I was iterating over the structure of the file, I thoought about how to integrate some missing components.
1. I need to define what type of security I will support and how to configure it.
2. How to embed parameters in URLs as opposed to using only query parameters.
3. The system should provide information on the configured endpoints and queries.
4. Handling filtering, sorting, and other optional parameters.
5. Logging and tracing to quickly identify problems.
6. A command line version to include in scripts or testing.

## Structuring the Project
The basic structure won't change. 
Go projects usually include main programs in the `cmd` or `cmds` subdirectory.
See [Dave Cheney's blog](https://dave.cheney.net/2014/12/01/five-suggestions-for-setting-up-a-go-project) for more information.
However, I'm debating the value of spliting the web server into its own project.
The UI assets (JavaScript, CSS, etc.) will make for a messy `cmds` subdirectory.

### Golang Level API
I've chosen to create a package with nested commands.
This implies will expose some functions and data structures.
I need to decide what should be public and what should be private.
I also want my API to be clear to other parties and safe to use.

### Dependencies
For the POC I chose a set of bare bones packages.
For the production version I may need to add or substitute new dependencies.
My personal philosopohy is to start simple.
Prefer a few simple, well understood tools than large and complex tools.
I tried not to include anything I didn't need.

|Project                             |Version |Reason                         |
|------------------------------------|--------|-------------------------------|
|github.com/BurntSushi/toml          |v0.3.1  |TOML parser                    |
|github.com/julienschmidt/httprouter |v1.2.0  |HTTP path handling and routing |
|github.com/lib/pq                   |v1.2.0  |Database access                |

I don't see any need, at the present, to use a higher level database package.
I might need to change the routing package to better accomodate a UI.
Since the UI might be a separate project, I'm not sure I'll need to switch.
I will probably need to add packages for security and monitoring.
The security requirements might impcat my choice of request router.

## First Iteration
The temptation is to 'just start.'
I have to narrow down my tasks for the first iteration of the production version.
I need to document what I expect to accomplish.
Even though I'm working alone, defining tasks requires thinking about the tasks.
Just like writing this requires me to think through my methodology.

### Tabled For Now
I intend to postpone security for the first iteration.
There are several ways to secure web applications.
In  corporate environments applications normally sit behind proxies that enforce security policies.
These proxies normally pass authentication and roles in the HTTP headers to the proxied application.
This is radiaclly different from Oauth2, JWT or other mechanisms.
Maybe I should include multiple security models in the application.
I think I'm not ready to tackle security at this time.

Because security is sometimes tied to routing, I should probably postpone changing the routing.
I don't want to suck in a larger package, hoping its features fix my issues.
I can leave the existing router in place for now.
I think I also need to make some decisions about a UI and where that UI might live.
This is also related to the question of responsiblity for parameter validation.

### Iteration Contents 
For this first iteration, it might be better to focus on the internals of wysci.
Specifically, the portion that dynamically executes a query.
I need to pass parameters to that query, return types, cast responses, format those response, expand the number of sql types, etc.
There is plenty of work.
My goals for this sprint are:
1. Improve the number of data types.
2. Better handle matching SQL types to Go types.
3. Better interrogate queries about their returned types.
4. Better communicate errors:
   1. Unknown type returned in response
   2. Missing parameters or wrong type
5. Passing context.

I have decided to integrate a formatted logging tool rather than tracing.
Tracing (as opposed to logging) is meant to capture the flow of operations and the telemetry as those operations execute.
Tracing usually requires a backend to interpret the traces (for example, traces have child traces).
While there are project like [OpenTracing](https://opentracing.io/), tracing requires a backend to interpret the trace.
Most of the users of wysci will probably prefer simple logs.
1. Log entry and exit points for operations.
2. Implement levels of verbosity

## Results
While working on items 1 and 2 (SQL data types and matching those types to Go types), it became apparent that the `NullString` type is sufficient.
Since the library is meant to output CSV, only the representation of the data is important.
Passing pointers to `NullString` to `Scan` captures both numeric and text values as strings.
The `NullString` type also allows distinguishing no value from empty strings.

Which raises an originally overlooked problem, how are `null` values displayed.
Options might be the text 'NULL' or nothing between commas.
This is something that should be configurable.

