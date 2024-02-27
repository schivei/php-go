package main

import (
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
)

func makeClean(extDir string) {
	// run "make clean"
	cmd := exec.Command("make", "clean")
	cmd.Dir = extDir
	cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
	cmd.Env = append(cmd.Env, os.Environ()...)
	_ = cmd.Run()
}

func phpizeClean(extDir string) {
	// run "phpize --clean"
	cmd := exec.Command("phpize", "--clean")
	cmd.Dir = extDir
	cmd.Env = append(cmd.Env, "CGO_ENABLED=1")
	cmd.Env = append(cmd.Env, os.Environ()...)
	_ = cmd.Run()
}

func cleanAll(extDir string) {
	makeClean(extDir)
	phpizeClean(extDir)

	_ = os.Remove(extDir + "/tests/fixtures/go/test.so")

	files, err := os.ReadDir(extDir)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if strings.HasSuffix(file.Name(), ".dep") || strings.HasSuffix(file.Name(), ".txt") {
			_ = os.Remove(extDir + "/" + file.Name())
		}
	}

	_ = os.RemoveAll(extDir + "/include")
}

func main() {
	pkgName := "github.com/schivei/php-go"

	log.Println("Building the package", pkgName, "as a PHP extension")

	pkgPath := exec.Command("go", "list", "-f", "{{.Dir}}", pkgName)
	pkgPath.Env = append(pkgPath.Env, "CGO_ENABLED=1")
	pkgPath.Env = append(pkgPath.Env, "GO111MODULE=on")
	pkgPath.Env = append(pkgPath.Env, os.Environ()...)

	var outb bytes.Buffer
	pkgPath.Stdout = &outb
	var stderr bytes.Buffer
	pkgPath.Stderr = &stderr

	err := pkgPath.Run()

	if err != nil {
		if stderr.Len() == 0 {
			panic(err)
		}

		msg := stderr.String()
		panic(errors.Join(err, errors.New(msg)))
	}

	pkgDir := outb.String()

	pkgDir = strings.ReplaceAll(pkgDir, "\n", "")
	pkgDir = strings.ReplaceAll(pkgDir, "\r", "")
	pkgDir = strings.ReplaceAll(pkgDir, "\t", "")

	if pkgDir[0] == '"' {
		pkgDir = pkgDir[1 : len(pkgDir)-1]
	}

	log.Println("The package", pkgName, "is located at", pkgDir)

	curPath := exec.Command("pwd")
	outb.Reset()
	curPath.Stdout = &outb
	err = curPath.Run()

	if err != nil {
		panic(err)
	}

	curDir := outb.String()

	log.Println("The current directory is", curDir)

	var destPath string

	if len(os.Args) <= 1 {
		os.Args = append(os.Args, "./")
	}

	destPath = os.Args[1]
	if _, err = os.Stat(destPath); os.IsNotExist(err) {
		log.Panicln("The directory", destPath, "does not exist")
	}

	var fileInfo os.FileInfo
	if fileInfo, err = os.Stat(destPath); os.IsNotExist(err) || !fileInfo.IsDir() {
		log.Panicln("The path", destPath, "is not a directory")
	}

	destPath = destPath + "/phpgo.so"

	log.Println("The destination path is", destPath)

	extPath := pkgDir + "/ext"

	defer func(extDir string) {
		defer cleanAll(extDir)

		if r := recover(); r != nil {
			log.Panicln(r)
		}
	}(extPath)

	log.Println("The extension path is", extPath)

	phpize := exec.Command("phpize")
	phpize.Dir = extPath
	phpize.Env = append(phpize.Env, "CGO_ENABLED=1")
	phpize.Env = append(phpize.Env, "GO111MODULE=on")
	phpize.Env = append(phpize.Env, os.Environ()...)
	phpize.Stdout = os.Stdout
	phpize.Stderr = os.Stderr
	err = phpize.Run()

	if err != nil {
		panic(err)
	}

	configure := exec.Command("./configure")
	configure.Dir = extPath
	configure.Env = append(configure.Env, "CGO_ENABLED=1")
	configure.Env = append(configure.Env, "GO111MODULE=on")
	configure.Env = append(configure.Env, os.Environ()...)
	configure.Stdout = os.Stdout
	configure.Stderr = os.Stderr
	err = configure.Run()

	if err != nil {
		var configLog []byte
		e := err
		configLog, err = os.ReadFile(extPath + "/config.log")
		if err != nil {
			panic(errors.Join(e, err))
		}

		msg := e.Error()
		msg += ":\n" + string(configLog)

		panic(msg)
	}

	cmake := exec.Command("make")
	cmake.Dir = extPath
	cmake.Env = append(cmake.Env, "CGO_ENABLED=1")
	cmake.Env = append(cmake.Env, "GO111MODULE=on")
	cmake.Env = append(cmake.Env, os.Environ()...)
	cmake.Stdout = os.Stdout
	cmake.Stderr = os.Stderr
	err = cmake.Run()

	if err != nil {
		panic(err)
	}

	makeTestFixture := exec.Command("make")
	makeTestFixture.Dir = extPath + "/tests/fixtures/go"
	makeTestFixture.Env = append(makeTestFixture.Env, "CGO_ENABLED=1")
	makeTestFixture.Env = append(makeTestFixture.Env, os.Environ()...)
	makeTestFixture.Stdout = os.Stdout
	makeTestFixture.Stderr = os.Stderr
	err = makeTestFixture.Run()

	if err != nil {
		panic(err)
	}

	makeTest := exec.Command("make", "test")
	makeTest.Dir = extPath
	makeTest.Env = append(makeTest.Env, "CGO_ENABLED=1")
	makeTest.Env = append(makeTest.Env, os.Environ()...)
	makeTest.Stdout = os.Stdout
	makeTest.Stderr = os.Stderr
	err = makeTest.Run()

	if err != nil {
		panic(err)
	}

	makeInstall := exec.Command("make", "install")
	makeInstall.Dir = extPath
	makeInstall.Env = append(makeInstall.Env, "CGO_ENABLED=1")
	makeInstall.Env = append(makeInstall.Env, os.Environ()...)
	makeInstall.Stdout = os.Stdout
	makeInstall.Stderr = os.Stderr
	_ = makeInstall.Run()

	php := exec.Command("php", "-d", "extension=./modules/phpgo.so", "-m")
	php.Dir = extPath
	php.Env = append(php.Env, "CGO_ENABLED=1")
	php.Env = append(php.Env, os.Environ()...)
	php.Stdout = os.Stdout
	php.Stderr = os.Stderr
	err = php.Run()

	if err != nil {
		panic(err)
	}

	extPath = extPath + "/modules/phpgo.so"

	// copy the extPath to the destPath using go not command
	ff, err := os.ReadFile(extPath)
	if err != nil {
		panic(err)
	}

	tt, err := os.Create(destPath)

	if err != nil {
		panic(err)
	}

	defer func() { _ = tt.Close() }()

	_, err = tt.Write(ff)

	if err != nil {
		panic(err)
	}

	log.Println("The extension has been built and installed at", destPath)
}
