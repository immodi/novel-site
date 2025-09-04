package helpers

// Helper function to get ranking style
func GetRankStyle(rank int) string {
	switch rank {
	case 1:
		return "bg-yellow-500"
	case 2:
		return "bg-gray-400"
	case 3:
		return "bg-amber-700"
	default:
		return "bg-[#353535] dark:bg-gray-300 text-gray-300 dark:text-gray-700"
	}
}
