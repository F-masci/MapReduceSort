package utils

// MergeSortedArrays Merges multiple sorted arrays into a single sorted array using the merge sort algorithm
func MergeSortedArrays(input [][]int) []int {
	// Start by merging the first two arrays
	result := input[0]
	for _, array := range input[1:] {
		// Merge each subsequent array with the result using merge sort
		result = merge(result, array)
	}
	return result
}

// merge Merges two sorted arrays into one sorted array
func merge(left, right []int) []int {
	var result []int
	i, j := 0, 0

	// Compare elements from both arrays and add the smaller one to the result
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	// If there are remaining elements in the left array, append them
	for i < len(left) {
		result = append(result, left[i])
		i++
	}

	// If there are remaining elements in the right array, append them
	for j < len(right) {
		result = append(result, right[j])
		j++
	}

	return result
}
