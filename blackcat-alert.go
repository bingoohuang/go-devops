package main

import "fmt"

func blackcatAlertExLog(exLogResult *ExLogCommandResult) {
	fmt.Println("===exLog:", exLogResult)
}

func blackcatAlertAgent(agentResult *AgentCommandResult) {
	fmt.Println("===agent:", agentResult)
}
