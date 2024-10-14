package tools

func AddTool(numbers []int) int {
	result := 0
	for _, num := range numbers {
		result = result + num
	}
	return result
}

func SubTool(numbers []int) int {
	result := 0
	for i, num := range numbers {
		if i == 0 {
			result = num
		} else {
			result = result - num
		}
	}
	return result
}
