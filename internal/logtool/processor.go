package logtool

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

func worker(ctx context.Context, jobs <-chan string, results chan<- Result) {
	for {
		select {
		case <-ctx.Done():
			return
		case filePath, ok := <-jobs:
			if !ok {
				return
			}

			file, err := os.Open(filePath)
			if err != nil {
				results <- Result{Err: err}
				continue
			}

			counts, err := CountLogLevels(ctx, file)
			closeErr := file.Close()

			if err != nil {
				results <- Result{Err: err}
				continue
			}

			if closeErr != nil {
				results <- Result{Err: closeErr}
				continue
			}

			results <- Result{Counts: counts}
		}
	}
}

func ProcessDirConcurrent(ctx context.Context, arg string, pool int) error {
	totals := make(map[string]int)
	jobs := make(chan string, pool*2)
	results := make(chan Result, pool*2)
	var wg sync.WaitGroup
	var collectorWg sync.WaitGroup

	collectorWg.Go(func() {
		for result := range results {
			if result.Err != nil {
				fmt.Printf("Could not read result, error: %v\n", result.Err)
				continue
			}
			MergeCounts(totals, result.Counts)
		}
	})

	for w := 1; w <= pool; w++ {
		wg.Go(func() {
			worker(ctx, jobs, results)
		})
	}

	err := filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".log" {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case jobs <- path:
		}

		return nil
	})

	close(jobs)
	wg.Wait()
	close(results)
	collectorWg.Wait()

	if err != nil {
		return err
	}

	PrintLogLevels(totals)

	return nil
}

func ProcessDirSync(ctx context.Context, arg string) error {
	totals := make(map[string]int)

	err := filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".log" {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		counts, err := CountLogLevels(ctx, file)
		closeErr := file.Close()

		if err != nil {
			return err
		}

		if closeErr != nil {
			return closeErr
		}

		MergeCounts(totals, counts)

		return nil
	})

	if err != nil {
		return err
	}

	PrintLogLevels(totals)

	return nil
}

func ProcessLogFile(ctx context.Context, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	counts, err := CountLogLevels(ctx, file)
	if err != nil {
		return err
	}

	PrintLogLevels(counts)

	return nil
}

func MergeCounts(dst, src map[string]int) {
	for key, value := range src {
		dst[key] += value
	}
}
