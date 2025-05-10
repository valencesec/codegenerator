package codegenerator

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

// var inlineTemplateRegex = regexp.MustCompile(`inline template: data file:`)
var inlineTemplateRegex = regexp.MustCompile(`\/\/ inline template: data file: "(.*)" *\r?\n((.|[\r\n])*?)\/\/ inline template end *\r?\n`)
var fileTemplateRegex = regexp.MustCompile(`\/\/ file template: data file: "(.*)" template file: "(.*)" *\r?\n`)
var globalDataRegex = regexp.MustCompile(`(\/\/|#) global data file: "(.*)" *\r?\n`)
var inlineDataRegex = regexp.MustCompile(`# inline data for template: template file: "(.*)" *\r?\n`)

func ScanFolder(rootFolder string, extension string) error {
	return filepath.WalkDir(rootFolder, func(filename string, d os.DirEntry, err error) error {
		base := filepath.Base(filename)
		if base == "node_modules" || base == ".git" {
			return filepath.SkipDir
		}
		if strings.HasSuffix(filename, extension) {
			target := filename[:len(filename)-len(extension)]
			log.Println("Generating ", target)
			return SingleFile(filename, target)
		}
		return nil
	})
}

func SingleFile(inFilename string, outFilename string) error {
	contentsBytes, err := os.ReadFile(inFilename)
	if err != nil {
		return err
	}
	globalData := make(map[string]interface{})
	for {
		match := globalDataRegex.FindSubmatchIndex(contentsBytes)
		if len(match) == 0 {
			break
		}
		upToCommand := contentsBytes[:match[0]]
		afterCommand := contentsBytes[match[1]:]
		contentsBytes = append(upToCommand, afterCommand...)

		dataFilenameGlobExpression := string(contentsBytes[match[4]:match[5]])
		dataFilenames, err := filepath.Glob(dataFilenameGlobExpression)
		if err != nil {
			log.Printf("While attempting to glob expression %s out of %s at position %d", dataFilenameGlobExpression, inFilename, match[0])
			return err
		}

		for _, dataFilename := range dataFilenames {
			log.Println("Reading global data file", dataFilename)
			dataContents, err := os.ReadFile(dataFilename)
			if err != nil {
				return err
			}
			var data map[string]interface{}
			err = yaml.Unmarshal(dataContents, &data)
			if err != nil {
				return err
			}
			for k, v := range data {
				globalData[k] = v
			}
		}
	}
	for {
		match := fileTemplateRegex.FindSubmatchIndex(contentsBytes)
		if len(match) == 0 {
			break
		}
		dataFilenameGlobExpression := string(contentsBytes[match[2]:match[3]])
		dataFilenames, err := filepath.Glob(dataFilenameGlobExpression)
		if err != nil {
			log.Printf("While attempting to glob expression %s out of %s at position %d", dataFilenameGlobExpression, inFilename, match[0])
			return err
		}
		templateFilename := string(contentsBytes[match[4]:match[5]])
		templateContents, err := os.ReadFile(templateFilename)
		if err != nil {
			log.Printf("While attempting to read template %s out of %s at position %d, template contents %s", templateFilename, inFilename, match[0], templateContents)
			return err
		}
		parsedTemplate, err := template.New(fmt.Sprintf("%s:%d", inFilename, match[0])).Funcs(AuxilirayFunctions()).Parse(string(templateContents))
		if err != nil {
			log.Printf("While attempting to parse %s at position %d, template contents %s", inFilename, match[0], templateContents)
			return err
		}
		outputBuffer := bytes.Buffer{}
		outputWriter := bufio.NewWriter(&outputBuffer)
		for _, dataFilename := range dataFilenames {
			log.Println("Generating using data file", dataFilename)
			dataContents, err := os.ReadFile(dataFilename)
			if err != nil {
				return err
			}
			var data map[string]interface{}
			err = yaml.Unmarshal(dataContents, &data)
			if err != nil {
				return err
			}
			finalData := map[string]interface{}{}
			for k, v := range globalData {
				finalData[k] = v
			}
			for k, v := range data {
				finalData[k] = v
			}
			err = parsedTemplate.Execute(outputWriter, finalData)
			if err != nil {
				return err
			}
			outputWriter.Flush()
		}
		upToCommand := contentsBytes[:match[0]]
		afterCommand := contentsBytes[match[1]:]
		contentsBytes = append(upToCommand, append(outputBuffer.Bytes(), afterCommand...)...)
	}
	for {
		match := inlineTemplateRegex.FindSubmatchIndex(contentsBytes)
		if len(match) == 0 {
			break
		}
		dataFilename := string(contentsBytes[match[2]:match[3]])
		templateContents := string(contentsBytes[match[4]:match[5]])
		parsedTemplate, err := template.New(fmt.Sprintf("%s:%d", inFilename, match[0])).Funcs(AuxilirayFunctions()).Parse(templateContents)
		if err != nil {
			log.Printf("While attempting to parse %s at position %d, template contents %s", inFilename, match[0], templateContents)
			return err
		}
		dataContents, err := os.ReadFile(dataFilename)
		if err != nil {
			return err
		}
		var data map[string]interface{}
		err = yaml.Unmarshal(dataContents, &data)
		if err != nil {
			return err
		}
		outputBuffer := bytes.Buffer{}
		outputWriter := bufio.NewWriter(&outputBuffer)
		finalData := map[string]interface{}{}
		for k, v := range globalData {
			finalData[k] = v
		}
		for k, v := range data {
			finalData[k] = v
		}
		err = parsedTemplate.Execute(outputWriter, finalData)
		if err != nil {
			return err
		}
		outputWriter.Flush()
		upToCommand := contentsBytes[:match[0]]
		afterCommand := contentsBytes[match[1]:]
		contentsBytes = append(upToCommand, append(outputBuffer.Bytes(), afterCommand...)...)
	}
	match := inlineDataRegex.FindSubmatchIndex(contentsBytes)
	if len(match) > 0 {
		templateFilename := string(contentsBytes[match[2]:match[3]])
		templateContents, err := os.ReadFile(templateFilename)
		if err != nil {
			log.Printf("While attempting to read template2 %s out of %s at position %d, template contents %s", templateFilename, inFilename, match[0], templateContents)
			return err
		}
		parsedTemplate, err := template.New(fmt.Sprintf("%s:%d", inFilename, match[0])).Funcs(AuxilirayFunctions()).Parse(string(templateContents))
		if err != nil {
			log.Printf("While attempting to parse %s at position %d, template contents %s", inFilename, match[0], templateContents)
			return err
		}
		var data map[string]interface{}
		err = yaml.Unmarshal(contentsBytes, &data)
		if err != nil {
			return err
		}

		outputBuffer := bytes.Buffer{}
		outputWriter := bufio.NewWriter(&outputBuffer)
		finalData := map[string]interface{}{}
		for k, v := range globalData {
			finalData[k] = v
		}
		for k, v := range data {
			finalData[k] = v
		}
		err = parsedTemplate.Execute(outputWriter, finalData)
		if err != nil {
			return err
		}
		outputWriter.Flush()
		contentsBytes = outputBuffer.Bytes()
	}
	err = os.WriteFile(outFilename, contentsBytes, 0644)
	if err != nil {
		return err
	}
	if strings.HasSuffix(outFilename, ".go") {
		_, err = exec.Command("gofmt", "-w", outFilename).CombinedOutput()
		if err != nil {
			log.Println("gofmt failed", err)
			return err
		}
	}
	return nil
}
