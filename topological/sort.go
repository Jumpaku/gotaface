package topological

// Sort performs topological sort on a directed graph represented by the adjacency list.
// It takes a 2D integer array 'graph' where graph[u] contains the indices of the destinations
// that can be reached from vertex u. This function returns the sorted vertices
// and a boolean flag indicating whether a valid topological order exists.
//
// If a valid topological order exists, this function returns (orders for each vertex, true).
// If a cycle is detected in the graph, this function returns (nil, false).
//
// Example:
// order, ok := topological.Sort([][]int{{5}, {3, 6}, {5, 7}, {0, 7}, {1, 2, 6}, {}, {7}, {0}})
// println(order, ok) // []int{4, 1, 1, 2, 0, 5, 2, 3} true
//
// order, ok := topological.Sort([][]int{{1}, {0}})
// println(order, ok) // nil false
func Sort(graph [][]int) ([]int, bool) {
	// count in-degrees
	inDegrees := make([]int, len(graph))
	for _, vs := range graph {
		for _, v := range vs {
			inDegrees[v]++
		}
	}

	// init queue
	queue := []int{}
	queueHead := 0
	for u := range graph {
		if inDegrees[u] == 0 {
			queue = append(queue, u)
		}
	}

	order := make([]int, len(graph))
	// perform topological sort
	for queueHead < len(queue) {
		u := queue[queueHead]
		queueHead++
		for _, v := range graph[u] {
			inDegrees[v]--
			if inDegrees[v] == 0 {
				order[v] = order[u] + 1
				queue = append(queue, v)
			}
		}
	}

	if queueHead != len(graph) { // cycle detected
		return nil, false
	}

	// no cycles
	return order, true
}
