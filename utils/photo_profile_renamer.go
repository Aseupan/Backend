package utils

func RenameLink(link string) string {
	var result string
	for _, char := range link {
		if char == ' ' {
			result += "%20"
		} else {
			result += string(char)
		}
	}
	return result
}
