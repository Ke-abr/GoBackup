package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type BackupConfig struct {
	SourcePath string
	DestinationPath string
	BackupType string //full or incremental
}

type BackupReport struct {
	Timestamp string
	FilesBackedUp int
	TotalSize int64 //in bytes
	Duration string
	Success bool
	ErrorMessage string
}

type BackupUtility struct {
	Report []BackupReport
}

func getFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func getDirectoryFiles(dirPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error accessing file/directory at:", path)
			return nil
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, data, 0644)
	if err != nil {
		return err
	}

	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	srcRights := info.Mode()
	err = os.Chmod(dst, srcRights)
	if err != nil {
		return err
	}

	return nil
}

func createBackupDirectory(backupPath string) error {
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		err := os.MkdirAll(backupPath, 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {}
	return nil
}

func generateBackupFileName(backupType string) string {
	if backupType == "full" {
		return "full_backup_" + time.Now().Format("2006-01-02_15-04-05") + ".tar"
	} else if backupType == "incremental" {
		return "incremental_backup_" + time.Now().Format("2006-01-02_15-04-05") + ".tar"
	} else {
		return "Wrong backup type"
	}
}

func (bu *BackupUtility) performFullBackup(config BackupConfig) BackupReport {
	var report BackupReport
	report.Success = true
	currentTime := time.Now()

	// Create initial backup directory
	err := createBackupDirectory(config.DestinationPath)
	if err != nil {
    	fmt.Printf("Error creating backup dir: %v\n", err)
		report.ErrorMessage = err.Error()
		report.Success = false
	} else {
    	fmt.Println("Backup dir creating successfully")
	}

	// Get all the files from the source
	files, err := getDirectoryFiles(config.SourcePath)
	if err != nil {
		fmt.Printf("Error getting files: %v\n", err)
		report.ErrorMessage = err.Error()
		report.Success = false
	}


	for _, file := range files {
		// Get the final path to preserve directory structure
		relativePath, err := filepath.Rel(config.SourcePath, file)
		if err != nil {
			report.ErrorMessage = err.Error()
			report.Success = false
			break
		}
		destPath := filepath.Join(config.DestinationPath, relativePath)
		destDir := filepath.Dir(destPath)
		err = createBackupDirectory(destDir)
		if err != nil {
			report.ErrorMessage = err.Error()
			report.Success = false
			break			
		}
		err = copyFile(file, destPath)
		if err != nil {
    		fmt.Printf("Copy error: %v\n", err)
			report.ErrorMessage = err.Error()
			report.Success = false
			break
		}

		report.FilesBackedUp++
		fileSize, err := getFileSize(file)
		if err != nil {
			fmt.Printf("Error getting file size: %v\n", err)
			report.ErrorMessage = err.Error()
			report.Success = false
			break
		}
		report.TotalSize += fileSize
	}

	duration := time.Since(currentTime)
	report.Duration = fmt.Sprintf("%.2f seconds", duration.Seconds())
	report.Timestamp = time.Now().Format("2006-01-02 15:04:05")
	return report
}

func main() {
	// testPath := "/home/freedfox/git/GoBackup/"
	// testFile := "README.md"
	// var err error

	// size, err := getFileSize(testPath+testFile)
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	return
	// }
	// fmt.Printf("File size for %s%s: %d bytes\n", testPath, testFile, size)

	// directoryFiles, err := getDirectoryFiles(testPath)
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// 	return
	// }
	// fmt.Printf("List of files for path %s: %v", testPath, directoryFiles)

	// err = copyFile(testPath+testFile, testPath+"README_copy.md")
	// if err != nil {
    // 	fmt.Printf("Copy error: %v\n", err)
	// } else {
    // 	fmt.Println("File copied successfully!")
	// }

	// err = createBackupDirectory("/home/freedfox/git/GoBackup/testBackup/directory/")
	// if err != nil {
    // 	fmt.Printf("Error creating backup dir: %v\n", err)
	// } else {
    // 	fmt.Println("File copied successfully!")
	// }

	// fmt.Println(generateBackupFileName("full"))
	// fmt.Println(generateBackupFileName("incremental"))
	// fmt.Println(generateBackupFileName("invalid"))

	backupUtil := BackupUtility{}
	config := BackupConfig{
    SourcePath:      "/home/freedfox/git/GoBackup/",  // Your current directory
    DestinationPath: "/tmp/test_backup/", // New backup folder
    BackupType:      "full",
    }

	report := backupUtil.performFullBackup(config)
	fmt.Printf("Backup Success: %v\n", report.Success)
    fmt.Printf("Files Backed Up: %d\n", report.FilesBackedUp)
    fmt.Printf("Total Size: %d bytes\n", report.TotalSize)
    fmt.Printf("Duration: %s\n", report.Duration)
    fmt.Printf("Timestamp: %s\n", report.Timestamp)
    if report.ErrorMessage != "" {
        fmt.Printf("Error: %s\n", report.ErrorMessage)
    }
}