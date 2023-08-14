package main

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"
)

var (
	timerActive     bool
	workDuration    = 25 * time.Minute
	breakDuration   = 5 * time.Minute
	timeRemaining   = workDuration
	timerTicker     *time.Ticker
	currentMenuItem *systray.MenuItem
)

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("Pomodoro Timer")
	systray.SetTooltip("Simple Pomodoro Timer")

	startWork := systray.AddMenuItem("Start Work", "Start a work session")
	stopWork := systray.AddMenuItem("Stop Work", "Stop the current work session")
	startBreak := systray.AddMenuItem("Start Break", "Start a break session")
	stopBreak := systray.AddMenuItem("Stop Break", "Stop the current break session")
	exit := systray.AddMenuItem("Exit", "Exit the application")

	for {
		select {
		case <-startWork.ClickedCh:
			startTimer(workDuration)
			currentMenuItem = startWork
		case <-stopWork.ClickedCh:
			stopTimer()
		case <-startBreak.ClickedCh:
			startTimer(breakDuration)
			currentMenuItem = startBreak
		case <-stopBreak.ClickedCh:
			stopTimer()
		case <-exit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func onExit() {
	// Clean up here if needed
}

func startTimer(duration time.Duration) {
	if timerActive {
		return
	}

	timeRemaining = duration
	timerTicker = time.NewTicker(1 * time.Second)
	timerActive = true

	go func() {
		for range timerTicker.C {
			if timerActive {
				updateTimer()
			}
		}
	}()

	updateTimer()
}

func stopTimer() {
	if !timerActive {
		return
	}

	timerTicker.Stop()
	timerActive = false
	if currentMenuItem != nil {
		currentMenuItem.Uncheck()
	}
}

func updateTimer() {
	minutes := int(timeRemaining.Minutes())
	seconds := int(timeRemaining.Seconds()) % 60

	systray.SetTitle(fmt.Sprintf("%02d:%02d", minutes, seconds))

	if timeRemaining <= 0 {
		stopTimer()
	}
	timeRemaining -= time.Second
}
