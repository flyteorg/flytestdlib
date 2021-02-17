package viper

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/pkg/errors"

	stdLibErrs "github.com/flyteorg/flytestdlib/errors"

	"github.com/spf13/cobra"

	"github.com/flyteorg/flytestdlib/config"
	"github.com/flyteorg/flytestdlib/config/files"
	"github.com/flyteorg/flytestdlib/logger"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"

	"github.com/spf13/pflag"
	viperLib "github.com/spf13/viper"
)

const (
	keyDelim = "."
)

var (
	dereferencableKinds = map[reflect.Kind]struct{}{
		reflect.Array: {}, reflect.Chan: {}, reflect.Map: {}, reflect.Ptr: {}, reflect.Slice: {},
	}
)

type viperAccessor struct {
	// Determines whether parsing config should fail if it contains un-registered sections.
	strictMode bool
	viper      *CollectionProxy
	rootConfig config.Section
	// Ensures we initialize the file Watcher once.
	watcherInitializer *sync.Once
}

func (viperAccessor) ID() string {
	return "Viper"
}

func (viperAccessor) InitializeFlags(cmdFlags *flag.FlagSet) {
	// TODO: Implement?
}

func (v viperAccessor) InitializePflags(cmdFlags *pflag.FlagSet) {
	err := v.addSectionsPFlags(cmdFlags)
	if err != nil {
		panic(errors.Wrap(err, "error adding config PFlags to flag set"))
	}

	// Allow viper to read the value of the flags
	err = v.viper.BindPFlags(cmdFlags)
	if err != nil {
		panic(errors.Wrap(err, "error binding PFlags"))
	}
}

func (v viperAccessor) addSectionsPFlags(flags *pflag.FlagSet) (err error) {
	return v.addSubsectionsPFlags(flags, "", v.rootConfig)
}

func (v viperAccessor) addSubsectionsPFlags(flags *pflag.FlagSet, rootKey string, root config.Section) error {
	for key, section := range root.GetSections() {
		prefix := rootKey + key + keyDelim
		if asPFlagProvider, ok := section.GetConfig().(config.PFlagProvider); ok {
			flags.AddFlagSet(asPFlagProvider.GetPFlagSet(prefix))
		}

		if err := v.addSubsectionsPFlags(flags, prefix, section); err != nil {
			return err
		}
	}

	return nil
}

// Binds keys from all sections to viper env vars. This instructs viper to lookup those from env vars when we ask for
// viperLib.AllSettings()
func (v viperAccessor) bindViperConfigsFromEnv(root config.Section) (err error) {
	allConfigs, err := config.AllConfigsAsMap(root)
	if err != nil {
		return err
	}

	return v.bindViperConfigsEnvDepth(allConfigs, "")
}

func (v viperAccessor) bindViperConfigsEnvDepth(m map[string]interface{}, prefix string) error {
	errs := stdLibErrs.ErrorCollection{}
	for key, val := range m {
		subKey := prefix + key
		if asMap, ok := val.(map[string]interface{}); ok {
			errs.Append(v.bindViperConfigsEnvDepth(asMap, subKey+keyDelim))
		} else {
			errs.Append(v.viper.BindEnv(subKey, strings.ToUpper(strings.Replace(subKey, "-", "_", -1))))
		}
	}

	return errs.ErrorOrDefault()
}

func (v viperAccessor) updateConfig(ctx context.Context, r config.Section) error {
	// Binds all keys to env vars.
	err := v.bindViperConfigsFromEnv(r)
	if err != nil {
		return err
	}

	v.viper.AutomaticEnv() // read in environment variables that match

	shouldWatchChanges := true
	// If a config file is found, read it in.
	if err = v.viper.ReadInConfig(); err == nil {
		logger.Printf(ctx, "Using config file: %+v", v.viper.ConfigFilesUsed())
	} else if asErrorCollection, ok := err.(stdLibErrs.ErrorCollection); ok {
		shouldWatchChanges = false
		for i, e := range asErrorCollection {
			if _, isNotFound := errors.Cause(e).(viperLib.ConfigFileNotFoundError); isNotFound {
				logger.Printf(ctx, "[%v] Couldn't find a config file [%v]. Relying on env vars and pflags.",
					i, v.viper.underlying[i].ConfigFileUsed())
			} else {
				return err
			}
		}
	} else if reflect.TypeOf(err) == reflect.TypeOf(viperLib.ConfigFileNotFoundError{}) {
		shouldWatchChanges = false
		logger.Printf(ctx, "Couldn't find a config file. Relying on env vars and pflags.")
	} else {
		return err
	}

	if shouldWatchChanges {
		v.watcherInitializer.Do(func() {
			// Watch config files to pick up on file changes without requiring a full application restart.
			// This call must occur after *all* config paths have been added.
			v.viper.OnConfigChange(func(e fsnotify.Event) {
				fmt.Printf("Got a notification change for file [%v]\n", e.Name)
				v.configChangeHandler()
			})
			v.viper.WatchConfig()
		})
	}

	return v.RefreshFromConfig(ctx, r, true)
}

func (v viperAccessor) UpdateConfig(ctx context.Context) error {
	return v.updateConfig(ctx, v.rootConfig)
}

// Checks if t is a kind that can be dereferenced to get its underlying type.
func canGetElement(t reflect.Kind) bool {
	_, exists := dereferencableKinds[t]
	return exists
}

// This decoder hook tests types for json unmarshaling capability. If implemented, it uses json unmarshal to build the
// object. Otherwise, it'll just pass on the original data.
func jsonUnmarshallerHook(_, to reflect.Type, data interface{}) (interface{}, error) {
	unmarshalerType := reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	if to.Implements(unmarshalerType) || reflect.PtrTo(to).Implements(unmarshalerType) ||
		(canGetElement(to.Kind()) && to.Elem().Implements(unmarshalerType)) {

		ctx := context.Background()
		raw, err := json.Marshal(data)
		if err != nil {
			logger.Printf(ctx, "Failed to marshal Data: %v. Error: %v. Skipping jsonUnmarshalHook", data, err)
			return data, nil
		}

		res := reflect.New(to).Interface()
		err = json.Unmarshal(raw, &res)
		if err != nil {
			logger.Printf(ctx, "Failed to umarshal Data: %v. Error: %v. Skipping jsonUnmarshalHook", data, err)
			return data, nil
		}

		return res, nil
	}

	return data, nil
}

// Parses RootType config from parsed Viper settings. This should be called after viper has parsed config file/pflags...etc.
func (v viperAccessor) parseViperConfig(root config.Section) error {
	// We use AllSettings instead of AllKeys to get the root level keys folded.
	return v.parseViperConfigRecursive(root, v.viper.AllSettings())
}

func (v viperAccessor) parseViperConfigRecursive(root config.Section, settings interface{}) error {
	errs := stdLibErrs.ErrorCollection{}
	var mine interface{}
	myKeysCount := 0
	if asMap, casted := settings.(map[string]interface{}); casted {
		myMap := map[string]interface{}{}
		for childKey, childValue := range asMap {
			if childSection, found := root.GetSections()[childKey]; found {
				errs.Append(v.parseViperConfigRecursive(childSection, childValue))
			} else {
				myMap[childKey] = childValue
			}
		}

		mine = myMap
		myKeysCount = len(myMap)
	} else if asSlice, casted := settings.([]interface{}); casted {
		mine = settings
		myKeysCount = len(asSlice)
	} else {
		mine = settings
		if settings != nil {
			myKeysCount = 1
		}
	}

	if root.GetConfig() != nil {
		c, err := config.DeepCopyConfig(root.GetConfig())
		errs.Append(err)
		if err != nil {
			return errs.ErrorOrDefault()
		}

		errs.Append(decode(mine, defaultDecoderConfig(c, v.decoderConfigs()...)))
		errs.Append(root.SetConfig(c))

		return errs.ErrorOrDefault()
	} else if myKeysCount > 0 {
		// There are keys set that are meant to be decoded but no config to receive them. Fail if strict mode is on.
		if v.strictMode {
			errs.Append(errors.Wrap(
				config.ErrStrictModeValidation,
				fmt.Sprintf("strict mode is on but received keys [%+v] to decode with no config assigned to"+
					" receive them", mine)))
		}
	}

	return errs.ErrorOrDefault()
}

// Adds any specific configs controlled by this viper accessor instance.
func (v viperAccessor) decoderConfigs() []viperLib.DecoderConfigOption {
	return []viperLib.DecoderConfigOption{
		func(config *mapstructure.DecoderConfig) {
			config.ErrorUnused = v.strictMode
		},
	}
}

// defaultDecoderConfig returns default mapsstructure.DecoderConfig with support
// of time.Duration values & string slices
func defaultDecoderConfig(output interface{}, opts ...viperLib.DecoderConfigOption) *mapstructure.DecoderConfig {
	c := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		TagName:          "json",
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			jsonUnmarshallerHook,
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// A wrapper around mapstructure.Decode that mimics the WeakDecode functionality
func decode(input interface{}, config *mapstructure.DecoderConfig) error {
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

func (v viperAccessor) configChangeHandler() {
	ctx := context.Background()
	err := v.RefreshFromConfig(ctx, v.rootConfig, false)
	if err != nil {
		// TODO: Retry? panic?
		logger.Printf(ctx, "Failed to update config. Error: %v", err)
	} else {
		logger.Printf(ctx, "Refreshed config in response to file(s) change.")
	}
}

func (v viperAccessor) RefreshFromConfig(ctx context.Context, r config.Section, forceSendUpdates bool) error {
	err := v.parseViperConfig(r)
	if err != nil {
		return err
	}

	v.sendUpdatedEvents(ctx, r, forceSendUpdates, "")

	return nil
}

func (v viperAccessor) sendUpdatedEvents(ctx context.Context, root config.Section, forceSend bool, sectionKey config.SectionKey) {
	for key, section := range root.GetSections() {
		if !section.GetConfigChangedAndClear() && !forceSend {
			logger.Infof(ctx, "Config section [%v] hasn't changed.", sectionKey+key)
		} else if section.GetConfigUpdatedHandler() == nil {
			logger.Infof(ctx, "Config section [%v] updated. No update handler registered.", sectionKey+key)
		} else {
			logger.Infof(ctx, "Config section [%v] updated. Firing updated event.", sectionKey+key)
			section.GetConfigUpdatedHandler()(ctx, section.GetConfig())
		}

		v.sendUpdatedEvents(ctx, section, forceSend, sectionKey+key+keyDelim)
	}
}

func (v viperAccessor) ConfigFilesUsed() []string {
	return v.viper.ConfigFilesUsed()
}

// Creates a config accessor that implements Accessor interface and uses viper to load configs.
func NewAccessor(opts config.Options) config.Accessor {
	return newAccessor(opts)
}

func newAccessor(opts config.Options) viperAccessor {
	vipers := make([]Viper, 0, 1)
	configFiles := files.FindConfigFiles(opts.SearchPaths)
	for _, configFile := range configFiles {
		v := viperLib.New()
		v.SetConfigFile(configFile)

		vipers = append(vipers, v)
	}

	// Create a default viper even if we couldn't find any matching files
	if len(configFiles) == 0 {
		v := viperLib.New()
		vipers = append(vipers, v)
	}

	r := opts.RootSection
	if r == nil {
		r = config.GetRootSection()
	}

	return viperAccessor{
		strictMode:         opts.StrictMode,
		rootConfig:         r,
		viper:              &CollectionProxy{underlying: vipers},
		watcherInitializer: &sync.Once{},
	}
}

// Gets the root level command that can be added to any cobra-powered cli to get config* commands.
func GetConfigCommand() *cobra.Command {
	return config.NewConfigCommand(NewAccessor)
}
