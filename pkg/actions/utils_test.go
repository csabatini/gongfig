package actions

import (
	"os/exec"
	"os"
	"strings"
	"testing"
)

// Run particular test in separate thread and test it exits with non zero value
// such test implementation is needed as tested function does not return with error
// but simply stops the execution of the whole program (os.Exit)
func runExit(testName string) error {
	cmd := exec.Command(os.Args[0], strings.Join([]string{"-test.run=", testName}, ""))
	cmd.Env = append(os.Environ(), "CHECK_EXIT=1")
	err := cmd.Run()

	return err
}

func getRoutesURL() string {
	//Create path /services/<service name>/routes
	routesPathElements := []string{ServicesPath, TestEmailService.Name, RoutesPath}
	return strings.Join(routesPathElements, "/")
}

func TestGetFullPath(t *testing.T) {
	rootAdminPath := "http://localhost:8001"
	fullApiPath := getFullPath(rootAdminPath, []string{ServicesPath})
	expectedApiPath :=  "http://localhost:8001/services"

	if fullApiPath != expectedApiPath {
		t.Fatal("expected", expectedApiPath, ", got", fullApiPath)
	}

	subAdminPath := "http://localhost:8001/kong-admin"
	fullApiPath = getFullPath(subAdminPath, []string{ServicesPath})
	expectedApiPath =  "http://localhost:8001/kong-admin/services"

	if fullApiPath != expectedApiPath {
		t.Fatal("expected", expectedApiPath, ", got", fullApiPath)
	}

}
