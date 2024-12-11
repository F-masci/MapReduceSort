package utils

// DistributeNumbersByValue Distributes the input numbers into chunks based on value range for each node
func DistributeNumbersByValue(input []int, nodes int) [][]int {

	// Calculate the maximum value in the input array
	maxValue := 0
	for _, num := range input {
		if num > maxValue {
			maxValue = num
		}
	}

	// Calculate the range for each node
	rangeSize := (maxValue + 1) / nodes

	// Initialize the result
	result := make([][]int, nodes)

	// Distribute numbers into the nodes
	for _, num := range input {
		// Calculate the node to which the number belongs
		node := num / rangeSize
		if node >= nodes {
			node = nodes - 1 // Ensure the maximum value is placed in the last node
		}
		result[node] = append(result[node], num)
	}

	return result
}
