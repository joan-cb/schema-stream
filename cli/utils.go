package cli

import (
	"fmt"
	"time"

	"github.com/alperdrsnn/clime"
)

func getSpinner(message string) {
	spinner := clime.NewSpinner().
		WithStyle(clime.SpinnerClock).
		WithColor(clime.MagentaColor).
		WithMessage(message).
		Start()
	time.Sleep(1 * time.Second)
	spinner.Stop()
}

func printRainbow(message string) {
	fmt.Println(clime.Rainbow(message))
}

func printBold(message string) {
	fmt.Println(clime.BoldColor.Sprint(message))
}

func printInfo(message string) {
	fmt.Println(clime.Info.Sprint(message))
}

func printError(message string) {
	fmt.Println(clime.Error.Sprint(message))
}

func getProgressBar(message string) {
	// Single progress bar
	bar := clime.NewProgressBar(100).
		WithLabel(message).
		WithStyle(clime.ProgressStyleModern).
		WithColor(clime.GreenColor).
		ShowRate(true)

	for i := 0; i <= 20; i += 5 {
		bar.Set(int64(i))
		bar.Print()
		time.Sleep(50 * time.Millisecond)
	}
	bar.Finish()
}

func getBox(title string, message ...string) {
	box := clime.NewBox().
		WithTitle(title).
		WithBorderColor(clime.BlueColor).
		WithStyle(clime.BoxStyleRounded)

	for _, v := range message {
		box.AddLine(v)
	}

	box.AutoSize(true)
	box.Println()
}
