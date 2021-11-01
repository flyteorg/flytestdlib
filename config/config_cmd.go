package config

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"unsafe"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

const (
	PathFlag        = "file"
	StrictModeFlag  = "strict"
	CommandValidate = "validate"
	CommandDiscover = "discover"
	CommandDocs     = "docs"
)

type AccessorProvider func(options Options) Accessor

type printer interface {
	Printf(format string, i ...interface{})
	Println(i ...interface{})
}

func NewConfigCommand(accessorProvider AccessorProvider) *cobra.Command {
	opts := Options{}
	rootCmd := &cobra.Command{
		Use:       "config",
		Short:     "Runs various config commands, look at the help of this command to get a list of available commands..",
		ValidArgs: []string{CommandValidate, CommandDiscover, CommandDocs},
	}

	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validates the loaded config.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return validate(accessorProvider(opts), cmd)
		},
	}

	discoverCmd := &cobra.Command{
		Use:   "discover",
		Short: "Searches for a config in one of the default search paths.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return validate(accessorProvider(opts), cmd)
		},
	}

	docsCmd := &cobra.Command{
		Use:   "docs",
		Short: "Generate configuration documetation in rst format",
		RunE: func(cmd *cobra.Command, args []string) error {
			sections := GetRootSection().GetSections()
			orderedSectionKeys := getOrderedSectionKeys(sections)
			visitedSection := map[string]bool{}
			typeLink := map[reflect.Type]string{}
			for _, sectionKey := range orderedSectionKeys {
				if canPrint(sections[sectionKey].GetConfig()) {
					printConfigTable(sections[sectionKey], sectionKey, false, visitedSection, typeLink)
				}
			}
			return nil
		},
	}

	// Configure Root Command
	rootCmd.PersistentFlags().StringArrayVar(&opts.SearchPaths, PathFlag, []string{}, `Passes the config file to load.
If empty, it'll first search for the config file path then, if found, will load config from there.`)

	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(discoverCmd)
	rootCmd.AddCommand(docsCmd)

	// Configure Validate Command
	validateCmd.Flags().BoolVar(&opts.StrictMode, StrictModeFlag, false, `Validates that all keys in loaded config
map to already registered sections.`)

	return rootCmd
}

// Redirects Stdout to a string buffer until context is cancelled.
func redirectStdOut() (old, new *os.File) {
	old = os.Stdout // keep backup of the real stdout
	var err error
	_, new, err = os.Pipe()
	if err != nil {
		panic(err)
	}

	os.Stdout = new

	return
}

func printConfigTable(section Section, sectionName string, isSubsection bool, visitedSection map[string]bool, typeLink map[reflect.Type]string) {
	val := reflect.Indirect(reflect.ValueOf(section.GetConfig()))

	if val.Kind() == reflect.Slice {
		val = reflect.Indirect(reflect.ValueOf(val.Index(0).Interface()))
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Type", "Default Value", "Description"})
	table.SetAlignment(3)
	table.SetRowLine(true)

	subsections := make(map[string]interface{})
	for i := 0; i < val.Type().NumField(); i++ {
		t := val.Type().Field(i)
		tagType := t.Type
		if tagType.Kind() == reflect.Ptr {
			tagType = t.Type.Elem()
		}

		fieldName := t.Name
		fieldType := tagType.Kind().String()
		fieldDefaultValue := fmt.Sprintf("%v", val.Field(i))
		fieldDefaultValue = strings.Replace(fieldDefaultValue, "_", "\\_", -1)
		fieldDescription := ""

		if tagType.Kind() == reflect.Map || tagType.Kind() == reflect.Slice || tagType.Kind() == reflect.Struct {
			fieldType = tagType.String()
			// Set default value to field type, and user can check its default value in subsection table
			fieldDefaultValue = fieldType
		}

		if jsonTag := t.Tag.Get("json"); len(jsonTag) > 0 && !strings.HasPrefix(jsonTag, "-") {
			var commaIdx int
			if commaIdx = strings.Index(jsonTag, ","); commaIdx < 0 {
				commaIdx = len(jsonTag)
			}
			if jsonTag[:commaIdx] != "" {
				fieldName = jsonTag[:commaIdx]
			}
		}

		if pFlag := t.Tag.Get("pflag"); len(pFlag) > 0 && !strings.HasPrefix(pFlag, "-") {
			var commaIdx int
			if commaIdx = strings.Index(pFlag, ","); commaIdx < 0 {
				commaIdx = -1
			}
			if pFlag[commaIdx+1:] != "" {
				fieldDescription = strings.TrimPrefix(pFlag[commaIdx+1:], " ")
			}
		}

		if tagType.Kind() == reflect.Struct {
			f := val.Field(i)
			// In order to get value from unexported field
			if f.Kind() == reflect.Ptr {
				f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
			} else {
				f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr()))
			}
			// Remove the package name of the type
			ss := strings.Split(fieldType, ".")
			if len(ss) > 0 {
				fieldType = ss[len(ss)-1]
			}

			if canPrint(f.Interface()) {
				if visitedSection[fieldType] {
					if typeLink[tagType] == "" {
						// Some types have the same name, but they are different type.
						// Add field name at the end to tell the difference between them.
						fieldType = fmt.Sprintf("%s (%s)", fieldType, fieldName)
						subsections[fieldType] = f.Interface()
					}
				} else {
					visitedSection[fieldType] = true
					subsections[fieldType] = f.Interface()
				}
				fieldType = fmt.Sprintf("`%s`_", fieldType)
				typeLink[tagType] = fieldType
			}
		}
		data := []string{fieldName, fieldType, fieldDefaultValue, fieldDescription}
		table.Append(data)
	}

	if section != nil {
		sections := section.GetSections()
		orderedSectionKeys := getOrderedSectionKeys(sections)
		for _, sectionKey := range orderedSectionKeys {
			fieldName := sectionKey
			t := reflect.TypeOf(sections[sectionKey].GetConfig())
			fieldType := t.String()
			if t.Kind() == reflect.Ptr {
				fieldType = t.Elem().String()
			}
			fieldDefaultValue := fieldType
			fieldDescription := ""

			if visitedSection[fieldType] {
				if typeLink[t] == "" {
					// Some types have the same name, but they are different type.
					// Add field name at the end to tell the difference between them.
					fieldType = fmt.Sprintf("%s (%s)", fieldType, fieldName)
					subsections[fieldType] = sections[sectionKey].GetConfig()
				}
			} else {
				visitedSection[fieldType] = true
				subsections[fieldType] = sections[sectionKey].GetConfig()
			}
			fieldType = fmt.Sprintf("`%s`_", fieldType)
			typeLink[t] = fieldType
			data := []string{fieldName, fieldType, fieldDefaultValue, fieldDescription}
			table.Append(data)
		}
	}

	if isSubsection {
		fmt.Println(sectionName)
		fmt.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
	} else {
		fmt.Println("Section:", sectionName)
		fmt.Println("-----------------------------------------")
	}
	table.Render()
	fmt.Println()

	for k, v := range subsections {
		printConfigTable(NewSection(v, nil), k, true, visitedSection, typeLink)
	}
}

func getOrderedSectionKeys(sectionMap SectionMap) []string {
	var orderedSectionKeys []string
	for s := range sectionMap {
		orderedSectionKeys = append(orderedSectionKeys, s)
	}
	sort.Strings(orderedSectionKeys)
	return orderedSectionKeys
}

// Print out config docs if and only if the section is struct or slice
func canPrint(b interface{}) bool {
	val := reflect.Indirect(reflect.ValueOf(b))
	if val.Kind() == reflect.Struct || val.Kind() == reflect.Slice {
		return true
	}
	return false
}

func validate(accessor Accessor, p printer) error {
	// Redirect stdout
	old, n := redirectStdOut()
	defer func() {
		err := n.Close()
		if err != nil {
			panic(err)
		}
	}()
	defer func() { os.Stdout = old }()

	err := accessor.UpdateConfig(context.Background())

	printInfo(p, accessor)
	if err == nil {
		green := color.New(color.FgGreen).SprintFunc()
		p.Println(green("Validated config file successfully."))
	} else {
		red := color.New(color.FgRed).SprintFunc()
		p.Println(red("Failed to validate config file."))
	}

	return err
}

func printInfo(p printer, v Accessor) {
	cfgFile := v.ConfigFilesUsed()
	if len(cfgFile) != 0 {
		green := color.New(color.FgGreen).SprintFunc()

		p.Printf("Config file(s) found at: %v\n", green(strings.Join(cfgFile, "\n")))
	} else {
		red := color.New(color.FgRed).SprintFunc()
		p.Println(red("Couldn't find a config file."))
	}
}
