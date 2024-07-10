package main

import (
	"_01"
	"engine"
)

// TODO: split each type into it's own submodule
/////////////////////////////////////////////////

func main() {
	// We first want to delete the log file if it exists, but we won't want to do this in production
	engine.NewFile("log.txt")
	err := engine.InitLogger(engine.DefaultConfig)

	if err != nil {
		engine.Log(engine.ErrorLevel, "Error initializing logger: %v", err)
		return
	}

	defer engine.Close()

	forest := _01.NewForest("forest")
	health := forest.GetTree("health")
	health.PopulateBranches("health::workout::[core({crunches,planks}),lower({squats,lunges}),upper({pushups,pullups})]")
	health.PopulateBranches("health::workout::cardio::{burpees,jump rope,jumping jacks,mountain climbers,running}")
	health.PopulateBranches("health::workout::[upper({arms::{pushups,pullups},chest::{bench press,flies},back::{rows,deadlifts},shoulders::{presses,shrugs})]")
	health.PopulateBranches("health::sleep::{nap,night}")
	health.PopulateBranches("health::hygiene::{brush teeth,shower}")
	health.PopulateBranches("health::diet::{breakfast,lunch,dinner,snacks}")
	health.PopulateBranches("health::mental::{meditation,reading}")

	forest.PopulateBranches("education::{reading,writing,math,science}")
	forest.PopulateBranches("finance::{budgeting,investing,saving}")
	forest.PopulateBranches("relationships::{family,romantic,friends}")

	//forest.Log()
}
