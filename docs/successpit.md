# The Pit of Succcess
One of the by-products of this project is I'm thinking in "long-form" about how I write code.
After twenty years of writing code there are many decisions I make without thinking.
I short-cut to the solution without showing the work.
Eventually I forget about why I do what I do.
Many of these decisions are exmaples of 'thinking fast.'
This is what experience buys you.

Some fraction of this might be worth documenting.
Of course there's the 'for the benefit of future generations of programmers.'
But realistically, no one may ever read this.
I'm really writing this to expose and question my thinking.
With over twenty years of writing code, I don't agonize over many decisions..

Experience (and the ability to think fast through the problem) comes when you have already visited a problem.
Of course it's never the _same_ problem.  
It's a related problem, but in a different language or with a different environment.
The downside is your experience may not be relevant.
My C++ programming intuition can lead me down some bad paths in JavaScript.
Taking apart my thinking will hopefully expose problem areas when thinking about Go.

## Once Upon a Time...
Once upon a time there was a developer faced with an unfamiliar API.
The developer started with the documentation.
The examples were really complete, even too complete.
In addition to the features she needed, the example covered many other fancy topics.
Burried in the really, really complete example was the correct way to use the API.
However, every time the developer tried to write a test program with that API there was nothing but failure.

"Falling into the pit of success" is my new guide-star for API development.
A naive user should be able to skim an example have a fair chance of success.
Tools I make should 'just work'.
It doesn't matter if the tools is a command line application, a REST API, or library.

An example of what not to do might be the Java libraries I used in the early 2000's.
Just to put 'Hello World' on a web page required me to write multiple configuration files.
The project was too complex to 'just run,' so I had to set up my build tooling.
It felt like I was yak shaving more than actually solving a problem.

## Intuition and Layering
My first thoughts on the post-POC code was to create two layers for extracting and formatting results.
The lower layer would extract the typed data from the database.
The igher layer would be responsible for formatting the result.
I could reuse the lower layer to implement JSON responses (for example).

The lower layer has a `Process` function.
I first thought I needed to return 'correct' Go types from the query response.
It turns out for my purposes that string results are sufficient.
The `Process` function takes database rows and returns a slice of string slices.
The writer can then iteratate over the slice of string slices and write those to the output.

```
    func Process(rows *sql.Rows) ([][]string, error)
```

## What About NULLs?
As I started to write the tests, I quickly realized I had a `NULL` problem.
Databases can return `NULL` values in results.
For example, not every address has a suite or apartment number and a database design might make the `SUITE_OR_APT` column nullable.
If I `SELECT * FROM ADDRESSES`, I would sometimes get a `NULL` for the `SUITE_OR_APT` column.

A `NULL` from the database is closer in spirit to a `nil` in Go.
If you have a pointer to a string, and it's `nil`, it's definitely not the same as an empty string.
The `nil` indicates there is no string while an empty string `''` is definitely a string.
Representing a `NULL` value is a little use-case dependent.
Delimited output might skip that field, outputting two consecutive delimiters.
In other cases, it might be good to output a string like `(NULL)` or `[NUll]`.

Returning string slices implies a configuration to set the desired `NULL` representation.
I could also handle this by returning slices of string pointers.
The `NULL` values from the database would simply be `nil` pointers in the returned slice.
But returning pointers would require my user to check before derferencing a value.

If query processor injected the `NULL` string value, users could write naive implementations to format the returned values.
They could simply range over the slice without having to worry about `nil` pointers.
They could also use functions like `Join`.

Let's create a struct to anchor the implemention with a `NullString` field and a `Process` method.
When it comes across a database `null`, it inserts the value of `NullString`.

```
    type QueryProcessor struct {
        NullString: string
    }

    func (qp QueryProcessor) Process (rows *sql.Rows) ([][]string, error) {
        // Implement here will check if the column in a row is NULL and 
        // insert the qp.NullString value as appropriate.
        return nil, nil
    }
```

We added a configuration item but there's no way to require the user to set it.
Fortunately, the 'zero value' for a string is the empty string.
The behavior is still reasonable If they user forgets to set the default value.
If we were printing out a CSV file, using something like `Join`, there would be two adjacent commas.
If they need to 'advertise' the `NULL` values, they can use a string like `(NULL)`.

```
# The second record has some missing fields that are NULL in the database
1230,Jane User,345 Musical Way,Suite 300,Space Junction,MO,USA
1234,Joe User,123 Melody Ln,,Jump City,,USA
```

## That's Fine, But What About Our 27,343,982 rows?
But there's a huge problem with returning slices (even if they are `string` or `*string`).
The in memory slice implies I've loaded the entire result set into memory.
Imagine querying a database that has millions of rows.
The server winds up bringing tens or hundreds of megabytes into memory.
It also increases the latency between the request and the first data returned.

### First Solution - Good Old Paging
One answer might be to return results in pages.
The user calles the `Proccess` method until it returns a zero length slice.

```
    qp := QueryProcessor{}
    res, err := qp.Process(bigQuery); 
    for err == nil && len(res) > 0 {
        // Do something with res

        res, err = qp.Process(bigQuery);
    }
```

Sometimes the user will need to tune the size of the page.  
Configuring page size implies a second option for the maximum number of rows to process per call to `qp.Process`.
If the user doesn't set it, the value would be 0 (the zero-value for integers).
I could simply assume a 'safe' value like 100, if the user doesn't explicitly set the property.

Alternatively, we could add a second parameter to `Process`, which is the batch size.
I don't like that alternative because it adds additional cognative load to the library user.
They now have another parameter they need to decide upon and may be the same for 90% of their code.
Let's explore adding the parameter to the `QueryProcessor` struct first.

```
    type QueryProcessor struct {
        NullString string
        MaxRows int
    }

    func (qp QueryProcessor) Process(rows *sql.Rows) ([][]string, error) {
        maxRows := 100
        if qp.MaxRows > 0 {
            maxRows = qp.MaxRows
        }

        // Process at most maxRows rows.

        return nil, nil
    }
```

There are two possible 'gotcha' issues with this approach.
The first is that nothing in the signature indicates we're returning a fixed batch size.
(Although the option of adding a parameter for page size would fix that.)
The second is we could have a small number of unusually large rows.
Maybe a batch size of 100 or even 1,000 is fine for 99% of rows, but 1% have huge text or JSON values.
In that case the better batch size would be 10, or even 1.

A user running wysci would notice that it breaks every once in a while with a memory error.
Assuming I can communicate the issue to them, or they debug it, they wind up setting the batch size to 1.
Then wysci is just really slow.
That's not the experience I want my users to have.

### Finally, The Dao of Go
The answer, of course, is not to return an in memory slice of string slices.
It would be better to take a `Writer` and write out the response as I go along.
For example, the `Writer` could point to a file.
I could return the number of bytes written.

```
    func (qp QueryProcessor) Process(rows *sql.Rows, w io.Writer) (int, error) {
        // New Implementation
    }
```

This looks like the kind of signature you'd find the standard library.
It feels right from a Go perspective.
In the internals of `Process`, the rows are scanned into memory one at a time.
Even if I were paging, I would still be scanning one row at a time.

Could there still be an issue?  
Absolutely!
There could be a row with massive text columns or maybe an embeded movie.
The user would see data arrive into their `io.Writer` until it suddenly stops.
The code should continue to work, even if there's a long pause to write that data to the writer.

Now that `Process` is responsible for writing the result, the user also has to pass in their delimiter.
Since this is CSV, it's safe to assume `,` is a default.

```
    type QueryProcessor struct {
        NullString, Delimiter string
    }
```

If the user just creates a `QueryProcesor` without setting the `NullString`, I get reasonable behavior.
If the user forgets to set `Delimiter`, I can default to `,`.
I think I've created an interface that a user is likely to get right the first time.

### There's Still an Issue
But there is another problem.
The `QueryProcessor` now has two responsiblities.
It is responsible for decoding the database result and formatting the output.
Good practice is to give `QueryProcessor` a single responsiblity.
Now, if we want to support structured (e.g. JSON) instead of delimited strings, we need to hack the `Process` function. 

What we need to do is delegate the formatting to another function.
We can start by providing a delimited output implementation.
We also want users to easily add their own implementations.
It should be optional so that it can default ot simple CSV.
If we pass a reference to it, then a NULL allows us to assume the CSV implementation.

```
    type Formatter interface {
        Format(parts []string) string
    }

    type QueryProcessor {
        NullString, Delimiter string
        RowFormatter Formatter
    }

    func (qp QueryProcessor) Process(rows *sql.Rows, w io.Writer) (int err) {
        // Implementation
    }
```

It also no longer makes sense to put the `NullString` and `Delimiter` members in the `QueryProcessor`.
They are implemenation details of the `Formatter`.
We can create a default `CSVFormatter` that implements the `Format` interface 

```
    type Formatter interface {
        Format(parts []string) string
    }

    type CSVFormatter struct {
        Delimiter, NullString string
    }

    func (c CSVFormatter) Format(parts []string) string {
        // If Delimiter is empty string assume ','
        // If NullString is empty string - that's okay.
    }

    type QueryProcessor {
        RowFormatter Formatter
    }

    func (qp QueryProcessor) Process(rows *sql.Rows, w io.Writer) (int err) {
        // Implementation defaults to CSV formatter if qp.RowFormatter is nil
    }
```

That feels better, even though there's something not right about the `Format(parts []string) string` signature.
In a sense it feels like we're going back to the start of this article.
We've no got the dreaded in-memory strings (even though it's just the strings from one result rwo scanned into memory).

Maybe the formatter should take a `io.Writer` and return bytes written and an error.
That made the processor more Go-ish.
We also have the issue of telling the `Formatter` if the individual strings are nil or not.
We could pass `*string` instead of `string`, but I think there's a better way.

The database API has a type `sql.NullString`, which we can pass to `Format` instead of `string`.
Because we're scanning nullable strings, we need to use the `sql.NullString` type to hold those strings.
Now the `Format` function would know if a string is `NULL` or not.

```
    type Formatter interface {
        Format(parts []sql.NullString, w io.Writer) (int64 error)
    }

    type CSVFormatter struct {
        Delimiter, NullString string
    }

    func (c CSVFormatter) Format(parts []sql.NullString, w io.Writer) (int64 error) {
        // If Delimiter is empty string assume ','
        // If NullString is empty string - that's okay.
    }
```

The `[]sql.NullString` parameter is a compromise.
It does not work with methods like `strings.Join(...)`.
But user's are likely to use the default formatter or maybe one of the included formatters.
If they're willing to write their own formatter, I'm less worried about their ability to write a simple implementation.
The `[]sql.NullString` is also more explicit than `[]*string`, regarding the nullability of the value.

### Almost A Happier Place
This seems like a good way to structure this part of wysci.
While we can't prevent a user from sticking full-length HD movies into their database, we handle even very large column sizes.
The user can simply create a `QueryProcessor` and get reasonable behavior.
The can customize that behavior by creating their own custom `Formatter` implementation.
If I, or a user, decides to create a JSON formatter, they can implement the `Formatter` interface.

Each part has one responsibility.
The `QueryProcoessor` takes the rows from the query and decodes them into types to pass to the `Formatter`.
The `Formatter` is responsible for formatting output.
There's one, last lingering problem.

Sometimes I'll need to know the names of the fields and their types.
Even the simple CSV formatter should print a header row so we know what those headers might be.
Also, type information may be useful (so a JSON represenation isn't all text fields).
Should this be an implementation detail of the formatter or part of the formatter API?

```
    type Formatter interface {
        Format(names []string, types []*sql.ColumnType, values []sql.NullString, w io.Writer) (int64, error)
    }
```

Since the `Process` method will call `Format` multiple times per row.
It should be responsible for passing on the type information and column headings.
(It will also be more efficient to extract that information once).
The other option is to add that information as fields on a formatter object.
Then it's up to the impelementation to decide problems like printing the headers just once.

```
    type Formatter interface {
        Format(values []sql.NullString, w io.Writer) (int64, error)
    }

    type CSVFomatter struct {
        Delimiter, NullString string
        Names []string
        Types []*sql.ColumnType
        didPrintHeaders bool
    }

    func (csv *CSVFormatter) Format(values []sql.NullString, w io.Writer) (int64, error) {
        // if !didPrintHeaders print headers and set didPrintHeaders to true.
        // print the row
    }

```

At first I didn't like this.
It felt like the `Formatter` had two responsiblities.
The first is writing the output in the desired format.
The second is tracking the output state.
They overlap but aren't the same thing.
You could have two different formatters that have the same state machine but totally different formats.

But I also don't want to over-engineer this.
Building a state transition machine seems like overkill.
Certainly the types and column headers are part of the state of the formatter.
As it stands, however, a user might try to reuse a `Formatter`.
As soon as I moved the `didPrintHeaders` into the implementation - I destroyed my ability to reuse it.

Let's make a constructor from the query result.
That way the users has to create a new `Formatter` each time.
I remove the risk the user tries to avoid an allocation and creates a bug for themselves.

```
    type CSVFormatter {
        Delimiter, NullString string
        names []string
        types []*sql.ColumnType
        didPrintHeader bool
    }

    func NewCSVFormatter(result *sql.Rows) (*CSVFormatter, error)
    func (csv *CSVFormatter) Format(values []sql.NullString, w io.Writer) (int64, error)
```

Note that I moved the `names` and `types` as being internal to the formatter.
That means they are now part of the internal state.
If the user tries to use a `QueryProcessor` without a formatter, I'll default to a new instance of `CSVFormatter`.

## Nearly Happy
I'm good with things for now.
I need to wrap up the first iteration and move on to other features.
There are some compromises in the API, but it feels like idiomatic Go code.

Most importantly, someone using my code will be more likely to 'get it right.'
I don't want someone struggling for an afternoon to reuse this part of wysci, only to find they'd made a simple mistake.
Or, for someone to start the server and think it's buggy because it crashes on their database of HD movies.
I want the user to fall into the pit of success.

