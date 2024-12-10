package utils

func MergeAndRemoveDuplicates[T comparable](slices ...[]T) []T {
	// Crea una mappa per tracciare gli elementi unici
	uniqueMap := make(map[T]bool)
	var result []T

	// Itera su ogni slice passata
	for _, slice := range slices {
		// Aggiungi gli elementi della slice
		for _, num := range slice {
			if !uniqueMap[num] { // Aggiungi solo se non esiste già
				uniqueMap[num] = true
				result = append(result, num)
			}
		}
	}

	return result
}
