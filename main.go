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
	fmt.Println(detectCycle2())
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
		fmt.Println("can't get max eid and max mid");
		return false, err
	}

	//Creating graph (adjacency list) from data
	graph := make([][]uint64, result.MaxID+1)
	for i:=0; i<=int(result.MaxID); i++{
		graph[i] = make([]uint64, result.MaxID+1)
	}
	for _, row := range rows {
		graph[row.EmployeeID] = append(graph[row.EmployeeID], row.ManagerID)
	}

	//checking if graph has cycle
	var vis = make([]bool, result.MaxID+1)
	var dfs_visit = make([]bool, result.MaxID+1)
	
	var i uint64
	for i=1; i < result.MaxID; i++ {
		if !vis[i] && hasCycle(i, vis, dfs_visit, graph) {
			return true, nil
		}
	}
	return false, nil
}


func detectCycle2() (bool, error) {
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
		fmt.Println("can't get max eid and max mid");
		return false, err
	}

	//Creating graph (adjacency list) from data
	graph := make(map[uint64][]uint64)
	for _, row := range rows {
		graph[row.EmployeeID] = append(graph[row.EmployeeID], row.ManagerID)
	}
	
	//checking if graph has cycle
	var vis = make([]bool, result.MaxID+1)
	var dfs_visit = make([]bool, result.MaxID+1)
	
	for k := range graph{
		if !vis[k] && hasCycle2(k, vis, dfs_visit, graph) {
			return true, nil
		}
	}
	return false, nil
}

func hasCycle(v uint64, vis []bool, dfs_vis []bool, g [][]uint64) bool {
	for _,nb := range g[v] {
		if nb == 0 {
			continue;
		}
		if dfs_vis[nb] {
			return true
		}
		vis[nb] = true;
		dfs_vis[nb] = true;
		if hasCycle(nb, vis, dfs_vis, g) {
			dfs_vis[nb] = false
			return true
		}
		dfs_vis[nb] = false;
	}
	return false
}



func hasCycle2(v uint64, vis []bool, dfs_vis []bool, g map[uint64][]uint64) bool {
	for _,nb := range g[v] {
		if nb == 0 {
			continue;
		}
		if dfs_vis[nb] {
			return true
		}
		vis[nb] = true;
		dfs_vis[nb] = true;
		if hasCycle2(nb, vis, dfs_vis, g) {
			dfs_vis[nb] = false
			return true
		}
		dfs_vis[nb] = false;
	}
	return false
}
