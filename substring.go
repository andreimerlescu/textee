package textee

// Len is part of sort.Interface.
func (sq SortedStringQuantities) Len() int {
	return len(sq)
}

// Swap is part of sort.Interface.
func (sq SortedStringQuantities) Swap(i, j int) {
	sq[i], sq[j] = sq[j], sq[i]
}

// Less is part of sort.Interface. We use it to sort the slice by Quantity in descending order.
func (sq SortedStringQuantities) Less(i, j int) bool {
	return sq[i].Quantity > sq[j].Quantity
}
