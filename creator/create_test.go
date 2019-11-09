package creator

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestCreateBase(t *testing.T) {
	os.Setenv("ZARUBA_TEMPLATE_DIR", "../playground/templates")
	os.Setenv("sender", "sender@gmail.com")
	os.Setenv("receiver", "receiver@gmail.com")
	target := path.Join("..", "playground", "projects", "test-create-base")
	err := Create("test", target, false)
	if err != nil {
		t.Errorf("%#v", err)
		return
	}

	// inspect readme.txt
	expectedReadmeContent := "# Test\nThis is a {{ .test }}"
	readmeContent, err := readGeneratedFile(target, "readme.txt")
	if err != nil {
		t.Error(err)
	} else if strings.Trim(readmeContent, "\n") != strings.Trim(expectedReadmeContent, "\n") {
		t.Errorf("Expected:\n%s\nActual:\n%s", expectedReadmeContent, readmeContent)
	}

	// inspect email/email.txt
	expectedEmailContent := "from: sender@gmail.com\nto: receiver@gmail.com\nHello,\nThis is an email from sender@gmail.com to receiver@gmail.com."
	emailContent, err := readGeneratedFile(target, "email/email.txt")
	if err != nil {
		t.Error(err)
	} else if strings.Trim(emailContent, "\n") != strings.Trim(expectedEmailContent, "\n") {
		t.Errorf("Expected:\n%s\nActual:\n%s", expectedEmailContent, emailContent)
	}

	// inspect hello.txt
	expectedHelloContent := "hello world"
	helloContent, err := readGeneratedFile(target, "hello.txt")
	if err != nil {
		t.Error(err)
	} else if strings.Trim(helloContent, "\n") != strings.Trim(expectedHelloContent, "\n") {
		t.Errorf("Expected:\n%s\nActual:\n%s", expectedHelloContent, helloContent)
	}

	// inspect special.txt
	_, err = readGeneratedFile(target, path.Join("special", "special.txt"))
	if err == nil {
		t.Error("Error expected")
	}
}

func TestCreateSpecial(t *testing.T) {
	os.Setenv("ZARUBA_TEMPLATE_DIR", "../playground/templates")
	os.Setenv("sender", "sender@gmail.com")
	os.Setenv("receiver", "receiver@gmail.com")
	target := path.Join("..", "playground", "projects", "test-create-special")
	err := Create("test:special", target, false)
	if err != nil {
		t.Errorf("%#v", err)
		return
	}

	// inspect readme.txt
	expectedReadmeContent := "# Test\nThis is a {{ .test }}"
	readmeContent, err := readGeneratedFile(target, "readme.txt")
	if err != nil {
		t.Error(err)
	} else if strings.Trim(readmeContent, "\n") != strings.Trim(expectedReadmeContent, "\n") {
		t.Errorf("Expected:\n%s\nActual:\n%s", expectedReadmeContent, readmeContent)
	}

	// inspect email/email.txt
	expectedEmailContent := "from: sender@gmail.com\nto: receiver@gmail.com\nHello,\nThis is an email from sender@gmail.com to receiver@gmail.com."
	emailContent, err := readGeneratedFile(target, "email/email.txt")
	if err != nil {
		t.Error(err)
	} else if strings.Trim(emailContent, "\n") != strings.Trim(expectedEmailContent, "\n") {
		t.Errorf("Expected:\n%s\nActual:\n%s", expectedEmailContent, emailContent)
	}

	// inspect hello.txt
	expectedHelloContent := "hello world"
	helloContent, err := readGeneratedFile(target, "hello.txt")
	if err != nil {
		t.Error(err)
	} else if strings.Trim(helloContent, "\n") != strings.Trim(expectedHelloContent, "\n") {
		t.Errorf("Expected:\n%s\nActual:\n%s", expectedHelloContent, helloContent)
	}

	// inspect special.txt
	expectedSpecialContent := "this is special"
	specialContent, err := readGeneratedFile(target, path.Join("special", "special.txt"))
	if err != nil {
		t.Error(err)
	} else if strings.Trim(specialContent, "\n") != strings.Trim(expectedSpecialContent, "\n") {
		t.Errorf("Expected:\n%s\nActual:\n%s", expectedSpecialContent, specialContent)
	}
}
func readGeneratedFile(target, filepath string) (string, error) {
	data, err := ioutil.ReadFile(path.Join(target, filepath))
	return string(data), err
}
