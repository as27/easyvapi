package easyvapi

// pagedResponse is the generic envelope returned by paginated easyVerein endpoints.
type pagedResponse[T any] struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []T     `json:"results"`
}

// fetchFunc is the signature of the function that Iterator uses to load each page.
// It receives the full URL (including query parameters) and returns the slice of
// results for that page, the URL of the next page (or nil if this is the last page),
// and any error that occurred during the fetch.
type fetchFunc[T any] func(url string) ([]T, *string, error)

// Iterator provides lazy, page-by-page iteration over a paginated API resource.
// Pages are fetched on-demand only when the current buffer is exhausted, making
// this memory-efficient for large result sets. Use [Iterator.Next] to advance and
// [Iterator.Value] to read the current item. Check [Iterator.Err] after iteration
// to detect any errors that occurred during fetching.
//
// Example:
//
//	iter := client.Members.List(ctx, nil)
//	for iter.Next() {
//		member := iter.Value()
//		fmt.Printf("Member: %s\n", member.MembershipNumber)
//	}
//	if err := iter.Err(); err != nil {
//		log.Fatal(err)
//	}
type Iterator[T any] struct {
	fetch   fetchFunc[T]
	nextURL *string

	buf   []T
	pos   int
	err   error
	done  bool
}

// newIterator creates an Iterator that starts loading from startURL.
func newIterator[T any](startURL string, fn fetchFunc[T]) *Iterator[T] {
	url := startURL
	return &Iterator[T]{
		fetch:   fn,
		nextURL: &url,
	}
}

// Next advances the iterator to the next item. It returns true when a value
// is available and false when iteration is complete or an error has occurred.
// Pages are fetched lazily: a new HTTP request is only made once the current
// buffer is exhausted. This makes iteration memory-efficient even for very
// large result sets.
//
// Example:
//
//	for iter.Next() {
//		item := iter.Value()
//		// process item
//	}
//	if err := iter.Err(); err != nil {
//		log.Fatal(err)
//	}
func (it *Iterator[T]) Next() bool {
	if it.err != nil || it.done {
		return false
	}
	// Advance within the current buffer.
	if it.pos < len(it.buf) {
		it.pos++
		return true
	}
	// Buffer exhausted – do we have a next page?
	if it.nextURL == nil {
		it.done = true
		return false
	}
	results, next, err := it.fetch(*it.nextURL)
	if err != nil {
		it.err = err
		return false
	}
	it.buf = results
	it.nextURL = next
	it.pos = 0
	if len(it.buf) == 0 {
		it.done = true
		return false
	}
	it.pos = 1
	return true
}

// Value returns the current item. It must only be called after a successful
// call to [Iterator.Next]. Calling Value when Next has not been called or
// when Next returned false results in undefined behavior.
//
//	if iter.Next() {
//		item := iter.Value()  // safe to call
//		fmt.Println(item)
//	}
func (it *Iterator[T]) Value() T {
	return it.buf[it.pos-1]
}

// Err returns the first error encountered during iteration, if any.
// Call Err after iteration completes (when [Iterator.Next] returns false)
// to check if iteration ended normally or due to an error.
//
//	for iter.Next() {
//		// process items
//	}
//	if err := iter.Err(); err != nil {
//		log.Printf("iteration failed: %v", err)
//	}
func (it *Iterator[T]) Err() error {
	return it.err
}
