


#### Potential enhancements
- Shut down workers when there's no tasks and spin them back up when necessary
- Use test data rather than executing actual gets
- Make two filter functions and make n not configurable to preallocate slice size and avoid reallocation
- Potentially shard the store (although I saw worse results as sorting didn't seem to have a massive overhead when benchmarking)
- How much do we care about accurate results? Could we batch sorting / processing to reduce overhead on fetching sorted lists
- Better worker pool & HTTP tests, most time spent on the store