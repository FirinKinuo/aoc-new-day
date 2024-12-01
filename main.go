package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

var (
	day        = flag.Int("day", 0, "Day of problem in year (required)")
	year       = flag.Int("year", 0, "Year of Advent of Code (required)")
	readmeOnly = flag.Bool("readme-only", false, "Generate only Readme.md file")

	session = os.Getenv("AOC_SESSION")
)

func validateFlags() {
	if *day == 0 {
		log.Fatal("day flag is required")
	}

	if *year == 0 {
		log.Fatal("year flag is required")
	}
}

func validateEnvs() {
	if session == "" {
		log.Warn("AOC_SESSION env is empty. Part 2 of the day and problem input will not be received!")
	}
}

func validateFull() {
	validateFlags()
	validateEnvs()
}

func main() {
	flag.Parse()

	validateFull()

	log.Infof("Fetching day %02d info from Advent of Code %d", *day, *year)
	aocClient := NewAOCApi(session)
	aocDay, err := aocClient.GetDayInfo(*day, *year)
	if err != nil {
		log.Fatalf("Error getting aocDay info: %v", err)
	}

	dayName := fmt.Sprintf(
		"day%02d_%s",
		*day,
		strings.ReplaceAll(strings.ToLower(aocDay.Title), " ", "_"),
	)

	log.Infof("Creating directory %d/%s", *year, dayName)
	targetDirectory, err := CreateTargetDirectory(*year, dayName)
	if err != nil {
		log.Fatalf("Error creating target directory: %v", err)
	}

	log.Infof("Creating %s/Readme.md with day description", targetDirectory)
	err = NewReadme(aocDay.Description).CreateFile(targetDirectory)
	if err != nil {
		log.Fatalf("Error creating readme: %v", err)
	}

	if *readmeOnly {
		log.Infof("Done")
		return
	}

	log.Infof("Creating %s/%s.go template for solving", targetDirectory, dayName)
	err = NewGo(dayName, *year).CreateFile(targetDirectory)
	if err != nil {
		log.Fatalf("Error creating new go template: %v", err)
	}

	log.Infof("Creating %s/%s input for testing. Maybe wrong, check!", targetDirectory, TestInputType.Filename())
	err = NewInput(TestInputType, aocDay.TestInput).CreateFile(targetDirectory)
	if err != nil {
		log.Fatalf("Error creating new test input file: %v", err)
	}

	log.Infof("Creating %s/%s input for solving", targetDirectory, ProblemInputType.Filename())
	err = NewInput(ProblemInputType, aocDay.ProblemInput).CreateFile(targetDirectory)
	if err != nil {
		log.Fatalf("Error creating new problem input file: %v", err)
	}

	log.Infof("Done. Let's solve this problem!")
}
