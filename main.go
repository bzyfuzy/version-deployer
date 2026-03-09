package main

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DeployResult struct {
	Site     string
	Files    int
	Duration float64
	Status   string
}

func initDB() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "deploy_log.db")
	if err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS deploy_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		site TEXT NOT NULL,
		version TEXT NOT NULL,
		files INTEGER,
		duration REAL,
		status TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(query)

	return db, err
}

func logDeploy(db *sql.DB, site, version string, files int, duration float64, status string) {

	query := `
	INSERT INTO deploy_log(site,version,files,duration,status)
	VALUES(?,?,?,?,?)
	`

	_, err := db.Exec(query, site, version, files, duration, status)

	if err != nil {
		fmt.Println("Log error:", err)
	}
}

func copyFile(src, dst string) error {

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)

	return err
}

func copyDir(src, dst string) (int, error) {

	files := 0

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, os.ModePerm)
		}

		err = copyFile(path, target)
		if err == nil {
			files++
		}

		return err
	})

	return files, err
}

func findTargets(root string) ([]string, error) {
	var targets []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			match, _ := filepath.Match("*.nuudelms.mn", info.Name())
			if match {
				if info.Name() == "test.nuudelms.mn" {
					fmt.Println("Skipping:", path)
					return nil
				}

				targets = append(targets, path)
			}
		}

		return nil
	})

	return targets, err
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: deployer SOURCE_FOLDER PROJECTS_ROOT")
		return
	}

	source := os.Args[1]
	root := os.Args[2]

	version := filepath.Base(source)

	fmt.Println("Source:", source)
	fmt.Println("Root:", root)
	fmt.Println("Version:", version)
	fmt.Println()

	db, err := initDB()
	if err != nil {
		fmt.Println("DB error:", err)
		return
	}
	defer db.Close()

	targets, err := findTargets(root)
	if err != nil {
		fmt.Println("Scan error:", err)
		return
	}

	fmt.Println("Found", len(targets), "sites")
	fmt.Println()
	workers := 8
	jobs := make(chan string)
	results := make(chan DeployResult)
	var wg sync.WaitGroup

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for target := range jobs {
				start := time.Now()
				fmt.Println("Deploying ->", target)
				files, err := copyDir(source, target)
				duration := time.Since(start).Seconds()
				status := "SUCCESS"

				if err != nil {
					status = "FAILED"
					fmt.Println("Error:", err)
				}

				results <- DeployResult{
					Site:     target,
					Files:    files,
					Duration: duration,
					Status:   status,
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		for _, t := range targets {
			jobs <- t
		}
		close(jobs)
	}()

	for r := range results {
		logDeploy(db, r.Site, version, r.Files, r.Duration, r.Status)
		fmt.Printf(
			"Logged: %s | files=%d | time=%.2fs | %s\n",
			r.Site,
			r.Files,
			r.Duration,
			r.Status,
		)
	}
	fmt.Println()
	fmt.Println("Deployment finished 🚀")
}
