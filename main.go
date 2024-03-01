package main

import (
	"fmt"

	"github.com/abhishek-bhangalia-busy/detect-cycle/db"
	"github.com/abhishek-bhangalia-busy/detect-cycle/initializers"
	"github.com/abhishek-bhangalia-busy/detect-cycle/models"
	"github.com/go-pg/pg/v10"
)

var DB *pg.DB

func init() {
	initializers.LoadEnvVariables()
}

func main() {
	//connect to DB
	DB = db.ConnectToDB()
	defer DB.Close()

	//creating schema
	err := db.CreateSchema(DB)
	if err != nil {
		fmt.Println("error : ",err.Error())
	}
	var empID, managerID uint64
	fmt.Println("Enter the empID :")
	fmt.Scanln(&empID)
	fmt.Println("Enter the managerID :")
	fmt.Scanln(&managerID)

	var ans bool
	ans, err = canAddEdge(empID, managerID)
	if err != nil {
		fmt.Println("error : ", err.Error())
	}
	if ans {
		fmt.Println("Cycle does not exist. Employee can be managed by manager.")
	} else {
		fmt.Println("Cycle will exist. So Employee can't be managed by manager.")
	}
}


// this func detects cycle in whole graph 
// so if cycle can exist in the existed table also then we will use this method

// func detectCycle() (bool, error) {
// 	// Fetch all rows from the table
// 	var rows []models.EmployeeManager
// 	if err := DB.Model(&rows).Select(); err != nil {
// 		fmt.Println("Error querying database: %v\n", err)
// 		return false, err
// 	}

// 	// Execute query to retrieve max values
// 	var result struct {
// 		MaxID uint64
// 	}
// 	_, err := DB.QueryOne(&result, `SELECT GREATEST(MAX(employee_id) , MAX(manager_id) ) AS max_id FROM employee_managers`)
// 	if err != nil {
// 		fmt.Println("can't get max eid and max mid")
// 		return false, err
// 	}

// 	//Creating graph (adjacency list) from data
// 	graph := make(map[uint64][]uint64)
// 	for _, row := range rows {
// 		graph[row.EmployeeID] = append(graph[row.EmployeeID], row.ManagerID)
// 	}

// 	var empID, managerID uint64
// 	fmt.Println("Enter the empID :")
// 	fmt.Scanln(&empID)
// 	fmt.Println("Enter the managerID :")
// 	fmt.Scanln(&managerID)
// 	graph[empID] = append(graph[empID], managerID)
// 	if empID > result.MaxID {
// 		result.MaxID = empID;
// 	}
// 	if(managerID > result.MaxID){
// 		result.MaxID = managerID;
// 	}

// 	//checking if graph has cycle
// 	var vis = make([]bool, result.MaxID+1)
// 	var dfs_visit = make([]bool, result.MaxID+1)

// 	for k := range graph {
// 		if !vis[k] && hasCycle(k, vis, dfs_visit, graph) {
// 			return true, nil
// 		}
// 	}
// 	return false, nil
// }


//this func will not check whole graph for cycle. 
//so, if we are given that cycle will not exist in the table then we can use this method and save time as it will not traverse whole graph for cycle
//it will only check the part of graph which is connected to newly added edge

func canAddEdge(empID uint64, managerID uint64) (bool, error) {
	// Fetch all rows from the table
	var rows []models.EmployeeManager
	if err := DB.Model(&rows).Select(); err != nil {
		fmt.Println("Error querying database: %v\n", err)
		return true, err
	}

	graph := make(map[uint64][]uint64)
	for _, row := range rows {
		graph[row.EmployeeID] = append(graph[row.EmployeeID], row.ManagerID)
	}
	
	graph[empID] = append(graph[empID], managerID)

	var vis = make(map[uint64]bool)
	var dfs_visit = make(map[uint64]bool)

	if hasCycle(empID, vis, dfs_visit, graph) {
		return false, nil
	}
	return true, nil
}


func hasCycle(v uint64, vis map[uint64]bool, dfs_vis map[uint64]bool, g map[uint64][]uint64) bool {
	vis[v] = true
	dfs_vis[v] = true
	for _, nb := range g[v] {
		if !vis[nb] {
			if hasCycle(nb, vis, dfs_vis, g) {
				return true
			}
		} else if dfs_vis[nb] {
			return true
		}
	}
	dfs_vis[v] = false
	return false
}
