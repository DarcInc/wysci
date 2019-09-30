# Motivation
It's just past eleven PM and I've come to the unshakable conclusion that this code needs to be refactored.
I'm adding a simple feature and there's no good way to do it.
I can brute force it into a function/handler/class, but it will be ugly and a maintenance nightmare.
If I try to be elegant, it will be a lot of code, additional data types, and digging through the guts of some of the libraries.
This isn't just a simple extraction of a class, the decisions we made three months ago were short-sighted.  
We need to gut and rebuild.
Refactoring would take precious days that can no longer be accomodated in schedule or cost.

And I'm not the only one in this predicament.
The rest of the team is grumbling, as well.
If only we had just thought this through a little better, we wouldn't be in this mess.
The project will regardlessly be delayed if we keep going down this road.
Costs will go up, if not for development then for maintanence.
We are sadled with too many poor decisions that were made early on, before we really understood the project.
The code I'm writing tonight, instead of going to bed and getting a good night's sleep, will just add to the codebase pinning us down.

You don't have to be a software developer for long before you may be faced with a similar problem.
Maybe you have been in meetings, trying to explain to your users, why a feature will take weeks and not days to implement.
You and they know it's not a big request, but you have to work in a legacy code base.
(And by legacy I don't mean from the 1970's, I just mean it already exists.)
Every decision you make early on in a project has the potential to bind future actions.

Unless you have a time machine to be able to go back and alter your old choices, how do you remove risk from early decisions?
The answer, of course, is to create a "Proof of Concept."
A POC (Proof of Concept) is a narrowed down version of what you intend to build.
A POC is often a good starting point to clarify your thinking and eliminate unknowns.

## The Proof of Concept
I sometimes view projects as being on a spectrum between the well known and the largely unknown.
At one end of the spectrum are projects that are well understood in scope and likely implementation.
This project is very similar to something I've done on previous projects.
An example might be a form-based Web application.  
This is something I've done several times in several languages.
While there may be new requirements or tools, I largely understand what needs to be done.
I can clearly visualize how to implement correct solution.
Making decisions in the solution space is largely remembering what I did last time or reviewing one of many examples on the Stack Overflow or GitHub.

The second extreme is something truly novel.
This is unlike other projects I've created.  
An example might be a GNOME application using Scala.
I've never written a GNOME application, I don't write desktop GUI applications, and I've never used Scala.
In some cases it won't be something that's common, where the solution can be found on Stack Overflow or GitHub.
Visualizing the correct solution is hard, and the intermediate steps are fuzzy.
A lot of decisions I make in the solution space are based on 'first principles' rather than memory or an example.

Most projects are somewhere between these two extremes.
Even in largely similar applications there is usually some novelty (otherwise you should just use or create a product).
Depending on the degree of novelty and the risk it introduces, we sometimes write simple example programs.
Sometimes these are just simple programs to make sure we understand a clear, working example before integrating it into a larger code base.
If there is enough novelty, you could create a 'proof of concept.'
The need for a POC is a function of you, your team, and the combined experience.

This application is an example of a proof of concept.
It is more novel than not, and some of it involves interface and pointer magic.
Building interfaces based on HTTP is not new for me, but creating dynamic queries in Go is new to me.
I've written Go code that uses understood types of queries, but not where the type is unknown until run-time.
In addition, instead of having logic for specific endpoints, I want the system to be configuration driven.
Since this is the majority of the application, I'll consider it a proof of concept.

### Limiting Scope
The key element to a good proof of concept is to limit scope.
This is harder than it looks.
Even for this fairly simple API.
For example, one of my early ideas was to all the user to specify the type of response.
I could return comma separated values (as originally intended), but why not JSON, XML, or other formats?
In fact, a good API should allow for that.
Rather than opening that can of worms, I tucked the idea in my back pocket for now.

There are multiple reasons for keeping the scope down but the key reason is the proof of concept should be a short lived activity.
A day or two of coding is ideal.
The more features and enhancements made during this phase will mean the proof of concept will drag on and on.
Maybe it will never finish.
This is sometimes not a bad thing.
It's okay to abandon a POC if you have all the answers you need and there's no point in completing it other than work.
But not finishing because you're piling on features is failure.

### Testing
In the version 0 code base there is a no testing.
This is not laziness or a lack of intention.
I had no idea how I work through some of the issues and was unsure about how to break up the application.
With no clear idea how to separate layers or responsibilities, writing tests would just slow me down.

For example, I rewrote the method that takes a result set and creates a new buffer to store the result from the database.
When I rewrote the method I made broad changes to the signature of the method and split it into multiple methods.
Without a stable interface I would have had to rewrite my tests with each iteration.
With a clear and stable interface it's easy to write tests.

### Incompleteness
Even within the decided scope I don't cover *everything*.
There are several Postgres types I do not cover.
Could those types cause me to rethink some of my ideas in the future?
Absolutely!
But what if those types never become important?
Would someone looking to pull the results of a query into Excel really care about byte arrays?
Maybe all we ever need are strings, numbers, and dates.

This is a little different from scope creep.
Scope creep adds new things to do, while incompleteness simply says we can ignore some details for now.
For example, I did not address securing the API.  
That requires me to take care of lots of detail work that will provide few insights to the basic problem..
Incompleteness and scope creep both result in excess and unnecessary work, resulting in lots of late nights, and will jeopardize completing the proof of concept.

### You're Throwing It Away
If you want to dig yourself into a hole, write a proof of concept for production.
There are decisions I made in the proof of concept that won't stand in the production version.
Rather than anchor my code base to those early decisions, I will throw them away.
I just need to answer some basic questions.

Throwing it away doesn't necessarily mean starting a whole new project.
It means that you will re-write large chunks of the core functionality, based on what you learned in the proof of concept.
Nothing should be sacrosanct, although you may find it convenient to copy and paste portions of the old code.
It could also mean that you take the POC and check it into a folder where you keep examples and documentation.

### Nothing Non-Functional
Everything in the POC should work.
If you don't have the time to make it work, we can assume it would have provide little or no value to you.
Therefore, don't even include it.

### Avoid Clients and Non-Technical Managers
If you're doing work for a client, and you show them the POC, they will assume you are done.
An architect might show a client a small, scale, table-top model of their building.
Unless you're Derek Zoolander, you probably won't confuse the model for the real thing.
However, clients often confuse the POC application with a deliverable application.

I've explained to clients what they are about to see is a proof of concept, without real data, and lacking functionality.
I've then watched the presentation deteriorate as the client wants me to explain 'how I got their data but explain why it isn't correct.'
Managers arguing with each other if they can field an application that has these huge missing features.
The risk is that the engagement manager or sales person will attempt to smooth the situation but agreeing to whatever the client wants.
You may come out of the demonstration looking like you failed.
My recommendation is to show folks wire frames in power-point slides, if the client needs to 'see something.'

## Recipe - How to do it Wrong
1. Every idea you have along the way is added into the POC.
2. Churn on your testing because all the interfaces are in flux.
3. Try to make sure every option and every feature is included.
4. Because of all the work you're putting into it, this must be the production code.
5. Lots of non-functional features in the POC for user to click on.
6. Demoing it to people that won't understand the difference between the model on the table-top and the actual building.

## Recipe - How to do it Right
1. What do you need clarified?  Fixate on just that.  Keep a wiki of 'good ideas for later.'
2. If you really, really, really need to write a test for something go ahead, but you might not even keep this code, so don't bother.
3. Only focus on the things you really need to clarify the problem.  Skip everything but main-line, happy path cases.
4. Assume this is throw away code.  Maybe create it under the 'examples' folder.
5. Everything basically works (even if it doesn't cover every scenario).  You've completely exercised the issue.
6. This is for you and your team, it's like private data inside a class.  If you need to 'take it to the client,' make some slides.

## Is this not Agile?
The agile terminology might be a spike.
A spike is where you go off and investigate something.
In a two week sprint you might need to integrate RabbitMQ messaging, but you never worked with RabbitMQ.
Maybe because of a language choice or environmental issue, there's no cookbook on how to do this.
You grab a docker container (even though your production environment may not use docker) and try a few simple commands.
You write a simple test program and you raise your confidence and answer questions you might have.
However, your work is not checked into the main branch.

## Conclusion
Sometimes you just a much simpler solution of part of your problem because you're not sure how to solve the whole problem in all its complexity.
The key is to do it quickly and efficiently, so you don't spend a lot of time or energy on it.
You will apply the learnings, but not necessarily the code, to the actual problem at hand.
With that in mind, these principles will help you keep it short, focused, and valuable.
However, if you get to a point where you no longer need it, feel free to abandon your POC.

