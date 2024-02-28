# Physical Storage

The goal of this exercise is to design a custom row-oriented file format that can ultimately be plugged into our query
executor, and then implement a `FileScan` node for reading records from a file.

## Suggestions

Before writing any code, start by considering the design problem.  You might find it useful to answer some of the
following questions

* What does the data look like?
* How will it be accessed?  All at once?  In a "streaming" fashion?  Random access?
* Can we _generate_ the files in a streaming manner?  If not, what data do we need beforehand (e.g. number of rows, contents of rows)
* Can we change the files after they've been written, or are they immutable?
* Should the format be optimized for writes or reads?  Can we support both efficiently?
* Do the files need to be "self-describing", or can the schema be stored somewhere else?
* Should the file be divided into "chunks" or sections?
* Will the data be compressed?  How?  Does that have any implication on our ability to write and read the files?
* Do we need to guard against data corruption due to bit rot?  How much should we do so?

You don't need to answer all of these questions beforehand, but they're helpful things to keep in mind while designing your custom format.  For example, if you stored all the rows in a serialized JSON document with the following schema:
```json
{
  "rows": [
    {
      "col1": "val1",
      "col2": "val2"
    }
  ]
}
```
Would this be a good choice?  Why not?

Ideally the format you come up with should be "binary oriented", meaning that you do all the data manipulation manually
and don't rely on external encoding like JSON, Protobuf, CSV, etc.  For example, Prometheus stores all of its raw
"chunks" data (compressed time series values) in "chunk" files with a very simple binary format.  The file format for
indexing these chunks is slightly more complicated, but not by much.

After coming up with a simple binary schema for your file format, the first step is to implement a `Writer`.  The
`Writer` should ideally generate the files in a streaming fashion (although it's alright if you find this difficult at
first and buffer all the values in memory before flushing them all at once) by exposing an API to write a single row
at a time, as well as an additional method for "closing" the file once there are no more rows to be written.
After implementing the `Writer`, you should implement a `Reader` that can read the file format and materialize stored
rows.  At the very least the reader should expose an API that returns a single row at a time even if the implementation
reads the entire file into memory, although ideally the implementation is such that data can be read from the file in
pieces (or even one row at a time).

While it is possible to test the `Writer` and `Reader` independently, this process tends to be cumbersome and brittle.
We recommend that you test your `Writer` and `Reader` in unison by writing "round-trip" tests that write a series of
rows using the Writer and ensure that the exact same rows are returned by the `Reader` implementation.

Once you've completed implementing the `Writer` and `Reader` interfaces, the final step is to implement a `FileScan`
operator that implements the `Operator` interface from the previous query executor lab.  The `FileScan` operator should
be a thin wrapper around the `Reader` interface that allows us to plug the custom file format into the query executor
so that it can operate over files directly instead of data in memory.