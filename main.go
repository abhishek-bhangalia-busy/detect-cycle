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
		panic(err)
	}
	fmt.Println(detectCycle())
}

func detectCycle() (bool, error) {
	// Fetch all rows from the table
	var rows []models.EmployeeManager
	if err := DB.Model(&rows).Select(); err != nil {
		fmt.Println("Error querying database: %v\n", err)
		return false, err
	}

	// Execute query to retrieve max values
	var result struct {
		MaxID uint64
	}
	_, err := DB.QueryOne(&result, `SELECT GREATEST(MAX(employee_id) , MAX(manager_id) ) AS max_id FROM employee_managers`)
	if err != nil {
		fmt.Println("can't get max eid and max mid")
		return false, err
	}

	//Creating graph (adjacency list) from data
	graph := make(map[uint64][]uint64)
	for _, row := range rows {
		graph[row.EmployeeID] = append(graph[row.EmployeeID], row.ManagerID)
	}

	var empID, managerID uint64
	fmt.Println("Enter the empID :")
	fmt.Scanln(&empID)
	fmt.Println("Enter the managerID :")
	fmt.Scanln(&managerID)
	graph[empID] = append(graph[empID], managerID)
	if empID > result.MaxID {
		result.MaxID = empID;
	}
	if(managerID > result.MaxID){
		result.MaxID = managerID;
	}

	//checking if graph has cycle
	var vis = make([]bool, result.MaxID+1)
	var dfs_visit = make([]bool, result.MaxID+1)

	for k := range graph {
		if !vis[k] && hasCycle(k, vis, dfs_visit, graph) {
			return true, nil
		}
	}
	return false, nil
}

func hasCycle(v uint64, vis []bool, dfs_vis []bool, g map[uint64][]uint64) bool {
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
