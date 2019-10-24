package config

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	PathFlag        = "file"
	StrictModeFlag  = "strict"
	CommandValidate = "validate"
	CommandDiscover = "discover"
	CommandGenerate = "generate"
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
		ValidArgs: []string{CommandValidate, CommandDiscover, CommandGenerate},
	}

	rootCmd.PersistentFlags().BoolVar(&opts.StrictMode, StrictModeFlag, false, `Validates that all keys in loaded config
map to already registered sections.`)

	validateCmd := &cobra.Command{
		Use:   CommandValidate,
		Short: "Validates the loaded config.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return validate(accessorProvider(opts), cmd)
		},
	}

	discoverCmd := &cobra.Command{
		Use:   CommandDiscover,
		Short: "Searches for a config in one of the default search paths.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return validate(accessorProvider(opts), cmd)
		},
	}

	var outputPath string

	generateCmd := &cobra.Command{
		Use:   CommandGenerate,
		Short: "Generates a single config file by compiling all default values and optionally provided config files.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generate(accessorProvider(opts), cmd, outputPath)
		},
	}

	generateCmd.Flags().StringVarP(&outputPath, "output-path", "o", "", "Where to write"+
		" the generated config")

	// Configure Root Command
	rootCmd.PersistentFlags().StringArrayVar(&opts.SearchPaths, PathFlag, []string{}, `Passes the config file to load.
If empty, it'll first search for the config file path then, if found, will load config from there.`)

	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(discoverCmd)
	rootCmd.AddCommand(generateCmd)

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

func generate(accessor Accessor, p printer, outputPath string) error {
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
	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		p.Println(red("Failed to validate config file."))
		return err
	}

	printInfo(p, accessor)
	green := color.New(color.FgGreen).SprintFunc()
	p.Println(green("Validated config file successfully."))

	p.Printf("Generating config at: %v\r\n", outputPath)
	rootSection := GetRootSection()
	m, err := createMap(rootSection)
	if err != nil {
		return err
	}

	bytes, err := yaml.Marshal(m)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(filepath.Join(outputPath, "config_all.yaml"), bytes, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func createMap(section Section) (res map[string]interface{}, err error) {
	res = map[string]interface{}{}
	if rootConfig := section.GetConfig(); rootConfig != nil {
		res, err = toMap(rootConfig)
		if err != nil {
			return nil, err
		}
	}

	for subName, sub := range section.GetSections() {
		res[subName], err = createMap(sub)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func toMap(config Config) (map[string]interface{}, error) {
	bytes, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	res := map[string]interface{}{}
	return res, json.Unmarshal(bytes, &res)
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
