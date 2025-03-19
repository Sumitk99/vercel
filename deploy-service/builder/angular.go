package builder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func BuildAngularProject(projectPath string) (*string, error) {
	log.Println("Building project at path: ", projectPath)
	cmdInstall := exec.Command("npm", "install")
	cmdInstall.Dir = projectPath
	cmdInstall.Stdout = os.Stdout
	cmdInstall.Stderr = os.Stderr

	err := cmdInstall.Run()
	if err != nil {
		fmt.Println("Error installing dependencies:", err)
		return nil, err
	}

	// Step 2: Build Angular
	cmdBuild := exec.Command("npm", "run", "build", "--", "--output-path=dist")
	cmdBuild.Dir = projectPath
	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr

	err = cmdBuild.Run()
	if err != nil {
		fmt.Println("Error building Angular project:", err)
		return nil, err
	}

	fmt.Println("✅ Angular build completed successfully! Output stored in repo/2uUEJ/dist/")
	buildPath := projectPath + "/dist"
	return &buildPath, nil
}
