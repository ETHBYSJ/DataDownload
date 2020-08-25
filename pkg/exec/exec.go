package exec

import (
	"os/exec"
)


func ExecPython(script string, imageDir string, newDir string, labelDir string, code string) error {
	// _, err := exec.Command("python", "D:\\storage\\medical-image\\example.py", "D:\\storage\\medical-image").Output()
	_, err := exec.Command("python3", script, imageDir, newDir, labelDir, code).Output()
	return err
}

/*
func main() {
	execPython("D:\\storage\\medical-image\\handle.py", "D:\\storage\\medical-image\\chest_Xray_test", "new2")
}
*/