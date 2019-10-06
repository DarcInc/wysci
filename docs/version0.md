# Motivation
It's 11:17 PM.
I've realized the code base I'm working on needs to be refactored.
I should go home, go to bed, and work on this after a good night's sleep and a clear head, but the project is already behind.
I'm adding a simple feature that should have been a day, or day and a half, of work.
There's no good way to do it.
There's an ugly hack, but the hack creates a maintenance headache.
The correct way to integrate the feature requires digging into the details of the ORM.
Had we understood the problem better three months ago, we would have made different choices.

I'm not the only one in this predicament.
The rest of the team is struggling with the same issues.
We won't be in this mess the next time do this.
We will better understand the problem and the tools.
We will make better choices.

We wouldn't be on a death-march if this had been a typical web application.
This was the first time this team ever worked on this type of application.
It was the first time we used this ORM.
We also needed a plugin framework, which we had to write and required some new language features.
The plugins sometimes confused the ORM, which was a side effect we hadn't anticipated.

Most software developers find themselves in a similar predicament at some point in their careers.
Maybe you have been in meetings, explaining to your customer, why simple features take weeks to implement.
You and your customers know it's not a big request, but you have to work in a legacy code base.
(Legacy doesn't mean from the 1970's, it just means already existing.)
Legacy issues to are expected on older code.  
Legacy issues should not appear on recent code.
In part, these issues arise because every decision made in a project sometimes binds future actions.

The big-design-up-front school (BDUF) tried to mitigate these issues with detailed designs before coding.
It was thought a detailed design would reveal potential problems before those problems were manifest in code.
But there are drawbacks with BDUF.
First, most customers won't pay for detailed designs they will not use and are almost immediately obsolete.
More importantly, the assumptions made during design are often made using experience.
If a team has little experience with a technology, problem domain, or system type, good designs are difficult.
Bad decisions early in an agile project create legacy issues in later sprints, just as bad choices in design can lead to an unworkable design.

How do you take the risk out of early decisions?
Without a time machine, can you can visit the future to understand the issues you might face?
The answer, of course, is to create a "Proof of Concept."
A POC (Proof of Concept) is a narrowed version of a core part of what you intend to build.
By quickly iterating on the problem once, you can expose issues and decision points you may encounter in the future.

## The Application
Wysci exposes SQL queries as endpoints returning CSV data.
Returning data as CSV may be the most consistently requested feature I've built into search applications.
'How do I get this into Excel?'
Is the question asked sometime in a late-stage demo of projects on which I've worked.
You can be a visualization artist with [D3](https://d3js.org/), but some users prefer spreadsheets.

There are enough open questions I've decided to start wysci as a POC.
I'm familiar with parts of the problem, such as writing Go code and writing Web APIs.
I also have many years of SQL experience.
I'm not clear about some parts of the application:
1. JSON (my first go-to format for configuration) doesn't handle multi-line strings well.  I've heard TOML does a better job, but I've never used it.
2. I've used Postgres with Go, but I knew ahead of time the number of columns and their type in the query result.  How do I handle dynamic results?
3. Are there limits to the type of queries I can define in the configuration file?
4. How do I define and pass parameters?  Can I enforce typing or other constraints?
5. I'm very wary of concatenating text to form SQL queries.  Maybe I can template the queries?

## The Proof of Concept
I sometimes view projects on a spectrum between the well known and the largely unknown.
At one end are projects with well understood scope and implementation plans.
These are projects are similar to previous projects, like a forms-based Web application.  
I've written this type of application several times in several languages.
There may be new requirements, but I understand what needs to be done.
You can rely on experience when making decisions in a familiar solution space.
Abundant and high quality examples (maybe in GitHub or Stack Overflow) can complement personal experience.

The second extreme is a truly novel application.
An example for me might be a GNOME application using Scala.
I have never written a GNOME application, I don't regularly write desktop GUI applications, and I've never used Scala.
Chances are, there is no abundance of documentation on writing GNOME desktop applications in Scala.
Experience may be a poor guide.
Decisions I make in development are lower quality and more likely to create legacy issues in the future.

Most projects are between these extremes.
Even familiar applications have some novelty (otherwise you should just use or create a product).
The more novelty, the more real or perceived risk.


Sometimes the uncertainty is not limited to a single tool or framework.
There may be several new tools, or this may be a different kind of application.
Risk can be amplified as these new tools interact with a novel context.
Having already built a similar application with the same tools would be ideal.
Developers often work with simple example programs when learning new frameworks and tools.
Even simple examples can raise confidence and mitigate risk.
This is the function a POC performs.

### Limiting Scope
Scope management isn't unique to a POC, but it is critical.
Ideas will present themselves as you or your team start developing the POC.
For example, I thought I should implement content negotiation through the `Accept` header.
Some users may want JSON or XML results.
But that isn't a core part of the original wysci concept.

The key reason to limit scope in the POC is that the POC should be completed quickly.
Depending on the scope of the final application, a POC should last between a day an a week.
Adding scope to the proof of concept takes time away from developing the production code.
You also risk never completing the POC as new features are introduced.

### Testing
I will skip testing.
I am not lazy and I'm very aware of the benefits of testing.
Tests are easier to write when we can anticipate the structure and layering in a project.
As I work through the POC I expect to understand better how to structure my application.
Writing tests would slow me down without a clear mental map of how to separate layers and responsibilities.

If I make changes to method signatures, move around responsibilities, or change the layering, I'll have to stop and update my tests.
The goal of the POC is to quickly answer questions and validate assumptions.
A POC is not responsible for delivering production code.
I certainly believe in writing tests for production code.

### Incompleteness
Even within the limited scope I will not address every concern.
I will not support every data type that could be returned in a Postgres query.
There is some risk that ignoring a detail might result in a drastically different decision.
On the other hand, I feel the basic types are sufficient for the POC.  
Other types might never be used.
For example, unsigned byte arrays may not be as useful in a spreadsheet.

### Disposable
I've written POCs in the past that were pushed into production.
It was fine for version 1 of the application, but terrible when changes began to accumulate.
The need for speed early on favored simplified assumptions and ignoring complicating issues.
These issues and incorrect assumptions haunted me during versions 2, 3, and so on.
I don't have to literally throw away my POC, but I should be free to start from scratch.

### Functional
Everything in the POC should work.
If there is no time to deliver a working feature, it probably has little value.
Therefore, I will not include it.
Other team members may check out the POC and run it to get a better understanding of the project.
A broken application will result in wasted time.

### It's Private
If you're doing work for a client, and you show them the POC, they may assume you are done.
An architect might show a client a small-scale, table-top model of a building.
Unless you're Derek Zoolander, you probably won't confuse that model for a real building.
However, clients can view the POC as a deliverable product.
It's been my experience that disclaimers or warnings are not effective.
Rather than dealing with the questions the POC will generate, keep it within the team.

### Recipe - How to do it Wrong
1. Every idea is added into the POC.
2. Constantly re-writing tests as wholesale changes are made to the code.
3. Making sure each feature covers every edge case and corner case.
4. Planning to take the POC to production.
5. Lots of non-functional features in the POC for developers to debug.
6. Demoing it to end users who have many questions about the stability, performance, or appearance of the application.

### Recipe - How to do it Right
1. Fixate on implementing only what is needed to reduce project risk.  Keep a wiki of 'good ideas for later.'
2. Keep it small and simple enough to test by running it.
3. Focus on the core concepts, postpone details until after the POC.  Focus on the 80% (or less) that can be done quickly, not the 20% that will take most of the time.
4. Assume this is throw away code.
5. Everything basically works.  You've exercised the issue.
6. This is for you and your team.  If you need to 'take it to the client,' make some slides.

## Abandoning the POC
There may come a point where the POC is no longer needed.
You may start the POC and realize that the problem is simpler than originally believed.
Or maybe you decide the problem you're trying to solve isn't worth solving.
It's okay to abandon a POC because it's no longer valuable.
What is the benefit of completing something with will serve no purpose?

## Isn't This Just Agile?
The agile terminology might be a spike.
A spike allows a developer to investigate an issue to help determine its scope or possible solution.
Usually a spike is an afternoon or a couple of hours.  
A POC lasts a little longer.
Like spikes, a POC should be used whenever it's needed, even in later stages of the project.

## Conclusion
The world is full of complexity and unknowns.
As the man said, there are known unknowns and unknown unknowns.
The goal of a POC is to try to understand better the known unknowns and to expose some of the unknown unknowns.
The key is to do it quickly and efficiently.
Then take the learnings, but not necessarily the code, to design and build the production code.
With that in mind, these principles will help you keep the POC short, focused, and valuable.
