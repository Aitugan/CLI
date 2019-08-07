package main

import (
	"fmt"
)

func findWaitingTime(processes []int, n int, bt []int, wt []int, quantum int) {
	rem_bt := make([]int, n)
	for i := 0; i < n; i++ {
		rem_bt[i] = bt[i]
	}
	t := 0

	for {
		done := true
		for i := 0; i < n; i++ {
			if rem_bt[i] > 0 {
				done = false
				if rem_bt[i] > quantum {
					t += quantum
					rem_bt[i] -= quantum
				} else {
					t += rem_bt[i]
					wt[i] = t - bt[i]
					rem_bt[i] = 0
				}
			}
		}
		if done == true {
			break
		}
	}
}

func findTurnAroundTime(processes []int, n int, bt []int, wt []int, tat []int) {
	for i := 0; i < n; i++ {
		tat[i] = bt[i] + wt[i]
	}
}

func findAvgTime(processes []int, n int, bt []int, quantum int) {
	wt := make([]int, n)
	tat := make([]int, n)
	totalWt := 0
	totalTat := 0

	findWaitingTime(processes, n, bt, wt, quantum)
	findTurnAroundTime(processes, n, bt, wt, tat)
	fmt.Println("Processes  Burst time  Waiting time  Turn around time")

	for i := 0; i < n; i++ {
		totalWt = totalWt + wt[i]
		totalTat = totalTat + tat[i]
		fmt.Printf("%v\t\t%v\t %v\t\t%v\n", (i + 1), bt[i], wt[i], tat[i])
	}
	fmt.Printf("Avg waiting time = %v\n", totalWt/n)
	fmt.Printf("Avg turn arounf time = %v\n", totalTat/n)

}

// func main() {
// 	processes := []int{1, 2, 3}
// 	n := len(processes)

// 	burstTime := []int{10, 5, 8}
// 	quantum := 2
// 	findAvgTime(processes, n, burstTime, quantum)
// }
