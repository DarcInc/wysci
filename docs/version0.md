# Motivation
It's 11:17 PM and I've come to the conclusion that the code base needs to be refactored.
I should go to bed and work on this with a good night's sleep and a clear head, but the project's already behind.
I'm adding a simple feature that should be two-point feature.
There's no good way to do it.
There's an ugly way to hack it in, but the hack creates a maintenance headache.
There's a correct way to integrate that feature, but it will result in a lot of code and digging into the details of the ORM we chose.
Had we understood the problem better three months ago, we would have made different choices.

I'm not the only one in this predicament.
The rest of the team is grumbling.
The next time we do this, we won't be in this miss.
We will better understand the problem and will make better choices.
The choices we made three months ago, in our first sprints, seemed reasonable at that time.
This was the first time this team ever worked on this type of application.
If it had been a regular Web application, we would have known exactly what to do.

Most software developers find themselves in that predicament at some point in their careers.
There's a good chance you have been that developer.
Maybe you have been in meetings, explaining to your customer, why simple features take weeks to implement.
You and your customers know it's not a big request, but you have to work in a legacy code base.
(By legacy I don't mean from the 1970's, I just mean it already exists.)
It's one thing for the issue to come up on code that's 10 or 15 years old, but it sometimes comes up on six-month old code.
Every decision you make early on in a project has the potential to bind future actions.

The big-design-up-front school (BDUF) tried to mitigate this issue with detailed design and analysis before coding begins.
But there are some problems with BDUF.
First is most customers won't pay for detailed designs that are obsolete days after they were created.
But more importantly, the assumptions made during design are often made using experience.
If a team has little experience with a technology, it's hard to make a good design.
Just as bad decisions early in an agile project create problems in later sprints, bad choices in the initial design can lead to a bad final design.

How do you take the risk out of those decisions?
Without a time machine, is there a way you can visit the future to see what kinds of issues you might face?
Can you deliver this without treating this project as a "learning exercise?"
The answer, of course, is to create a "Proof of Concept."
A POC (Proof of Concept) is a narrowed down version of what you intend to build.
By quickly iterating on the problem once, you expose the issues and decision points you may encounter in the future.

## The Proof of Concept
I sometimes view projects as being on a spectrum between the well known and the largely unknown.
At one end of the spectrum are projects where the scope and implemenation plan are well understood.
This project is similar to previous projects.
A form based Web application is a good example.  
I've written this type of application several times in several languages.
While there may be new requirements, I understand what needs to be done.
Making decisions in the solution space is based on experience or cookbook style from Stack Overflow or GitHub.

The second extreme is something truly novel.
This is unlike other projects to date, and experience is a poor guide.  
An example might be a GNOME application using Scala.
I've never written a GNOME application, I don't write desktop GUI applications, and I've never used Scala.
Chances are, there isn't a series of cookbook articles on writing GNOME desktop applications in Scala.
Visualizing what needs to be done during implementation isn't easy.
I have to make a lot of decisions based on assumptions and 'first principles,' rather than memory or documentation.

Most projects are somewhere between these extremes.
Even in familiar applications there is some novelty (otherwise you should just use or create a product).
The more novelty, the more unknown create real or perceived risk.
One way to become familiar with new tools or frameworks is to write a simple example program.
Even if it has nothing to do with the project, it gives a modicum of experience to ground future decisions.

There are times when the unknowns are not just one new tool for framework.
Maybe there are several new tools, or maybe this is a different kind of application.
What creates risk is now these new tools interact in a novel context.
What you really need is a simplified version of the problem.
The need for a POC is a function of you, your team, and the combined experience.
It doesn't really matter that someone out there in the world has done it, unless they also happen to be on your team.

## The Application
The basic idea behind wysci is to expose SQL queries as endpoints that return CSV data.
In the past I've written applications that were search and data centric.
One request that always comes at the end of the project is to be able to download the data into Excel.
In some cases there's no discussion on this requirement until a later stage demo and a user asks 'how do I get this into Excel?'
No matter how expertly you create visualizations in [D3](https://d3js.org/), there are users that just want to open in Excel.

I'm starting wysci with a proof of concept.
I'm familiar with parts of the problem.
I like Go and I've written a program or two in Go.
I've wrtten plenty of Web API endpoints.
I also have many years of SQL experience.
But there are some things which are open questions, among which are:
1. JSON (my first go-to format for configuration) doesn't handle multi-line strings well.  I've heard TOML does a better job, but I've never used it.
2. I understand how to use Postgres and Go, but when I know ahead of time the number of columns in the result and their type.  How do I handle dynamic results?
3. What kinds of endpoints and queries can I define in the configuration file?  Will there be any serious limitations?
4. How do I define and pass parameters between queries and the front end endpoints?  Can I enforce typing or other constraints?
5. I'm very wary of concatinating text to form SQL queries.  Maybe I can template the queries?

There are a fair number of unknowns.
I suspect it may involve some unsafe casting.
This also represents the guts of the application, safely exposing SQL queries defined in configuration files.
I have enough open questions that I would rather start with a POC.
What are the bounds on the POC?
What should I include and what should I leave for the actual application.

### Limiting Scope
Above all I will need to limit scope.
This is harder than it looks, even for this fairly simple API.
For example, I could let the user dictate the format of the response through the `Accept` header.
In fact, a good API should allow for allowing the user to specify the content types they are willing to accept.
But that isn't a core part of the original idea, which is to provde a CSV data endpoint.

The key reason to limit scope is the proof of concept should be a short lived activity.
Depending on the scope of the final application, somewhere between a day an a week.
Adding scope to the proof of concept takes time away from the actual application.
You also risk never completing the POC and answering your questions.

### Testing
I will skip testing.
This is not laziness or a lack of intention.
Tests are much easier to write when the structure and layering are more clear.
As I work through the POC I expect to understand better how to break up my application.
With no clear idea how to separate layers or responsibilities, writing tests would just slow me down.

If I make changes to method signatures, move around responsibilities, or change the layering, I'll have to stop and update my tests.
The goal of the POC is to quickly answer questions, not deliver production code.
I certainly believe in writing tests for production code.

### Incompleteness
Even within the limited scope I don't cover *everything*.
I don't have to support every type that could be returned in a Postgres query.
There maybe some risk that ignoring a type might result in a reachitecting of how queries are handled, but that's not likely.
But what if those types never become important?
Would someone looking to pull the results of a query into Excel really care about arrays of raw bytes?
Maybe all we ever need are strings, numbers, and dates.
I may have to addresses these types in the production app, along with security, monitoring, caching, etc. etc.
However, these don't help answer any of my core questions.

### Disposable
I've written POCs in the past that were pushed into production.
It was fine for version 1 of the application, but terrible when changes began to arrive.
The need for speed allowed me to simplify assumptions and ignore certain issues.
These issues and incorrect assumptions haunted me during versions 2, 3, and so on.
I would rather spend very little time on the POC and then retired it.
I expect to re-write large chunks of it, or even just check it into an 'examples' folder.

### Functional
Everything in the POC should work.
If there's no time to make it work, it probably has little value.
Therefore, don't even include it.
Especially when working inside a team, the ability to check out the POC and run it helps others understand how the software will be constructed.

### It's Private
If you're doing work for a client, and you show them the POC, they will assume you are done.
An architect might show a client a small, scale, table-top model of their building.
Unless you're Derek Zoolander, you probably won't confuse the model for the real thing.
However, clients often confuse the POC application with a deliverable application.
It's been my experience that clients forget or ignore any disclaimers or warnings before the demo.
Rather than explain why the 'data is wrong,' or 'who gave you our data,' or 'when will you fix all the bugs,' keep this as an internal effort.

### Recipe - How to do it Wrong
1. Every idea you have along the way is added into the POC.
2. Spend half the time re-writing tests because all the interfaces are in flux.
3. Try to make sure every option and every feature is included.
4. Because of all the work you're putting into it, this must be the production code.
5. Lots of non-functional features in the POC for user to click on.
6. Demoing it to people that won't understand the difference between the model on the table-top and the actual building.

### Recipe - How to do it Right
1. What do you need clarified?  Fixate on just that.  Keep a wiki of 'good ideas for later.'
2. Keep it small so you can test just by running it.
3. Only focus on the things you really need to clarify the problem.  Nothing else is important right now.
4. Assume this is throw away code.
5. Everything basically works.  You've exercised the issue.
6. This is for you and your team, like private data inside a class.  If you need to 'take it to the client,' make some slides.

## Abandoning the POC
There may come a point where the POC is no longer needed.
You start the POC and realize that you can do this with existing tools, or the solution is actually much simpler than first thought.
Or maybe you decide problem you're trying to solve isn't worth solving.
It's okay to abandon a POC because it's no longer valuable.

## Isn't This Just Agile?
The agile terminology might be a spike.
A spike is where you go off and investigate something.
Usually a spike is an afternoon or a couple of hours.  
A POC lasts a little longer.
It could be up to a week, but a day or two is often sufficient.

While you make create a POC during other parts of the project, you will generally create it at the very begining.
Spikes occur throughout the project as issues and need arise.
For example, I had never used RabbitMQ until one project so I took a spike to setup and test out the RabbitMQ Java API.
That helped me better estimate the points for the follow on work.

## Conclusion
The world is full of complexity and unknowns.
As the man said, there are known unknowns and unknown unknowns.
The goal of a POC is to try to understand better the known unknowns and to expose some of the unkown unknowns.
The key is to do it quickly and efficiently.
Then take the learnings, but not necessarily the code, to the actual problem.
With that in mind, these principles will help you keep it short, focused, and valuable.
