package stats

// Stat is the base interface that the tray utility will call to display to
// display the user-selected system statistic. Stat only requires a String() method
// which displays the string representation of the implementing system statistic
type Stat interface {
	// Init creates the readed and sets it's update interval to interval * time.Milliseconds
	Start(c chan string) error
	Stop()
	String() string
	ChangeInterval(interval int)
}
