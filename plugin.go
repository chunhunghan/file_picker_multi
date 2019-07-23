package file_picker_multi

import (
	"fmt"
	"strings"

	"github.com/go-flutter-desktop/go-flutter"
	"github.com/go-flutter-desktop/go-flutter/plugin"
	"github.com/pkg/errors"
)

const channelName = "file_picker"

type FilePickerPlugin struct{}

var _ flutter.Plugin = &FilePickerPlugin{} // compile-time type check

func (p *FilePickerPlugin) InitPlugin(messenger plugin.BinaryMessenger) error {
	fmt.Println("InitPlugin")
	dialogProvider := dialogProvider{}

	channel := plugin.NewMethodChannel(messenger, channelName, plugin.StandardMethodCodec{})
	//channel.HandleFunc("openDirectory", p.filePicker(dialogProvider, true))
	channel.HandleFunc("ANY", p.filePicker(dialogProvider, false, "*"))
	channel.HandleFunc("IMAGE", p.filePicker(dialogProvider, false, "*"))
	channel.HandleFunc("AUDIO", p.filePicker(dialogProvider, false, "*"))
	channel.HandleFunc("VIDEO", p.filePicker(dialogProvider, false, "*"))
	channel.CatchAllHandleFunc(p.fallBack)
	return nil
}

func (p *FilePickerPlugin) fallBack(methodCall interface{}) (reply interface{}, err error) {
	method := methodCall.(plugin.MethodCall)

	resolveType := strings.Split(method.Method, "__CUSTOM_")
	fileExtension := resolveType[1]
	fmt.Println("fallBack fileExtension:" + fileExtension)

	dlgProvider := dialogProvider{}
	p.filePicker(dlgProvider, false, fileExtension)
	// return the randomized Method Name
	return method.Method, nil
}

func (p *FilePickerPlugin) filePicker(dialog dialog, isDirectory bool, fileExtension string) func(arguments interface{}) (reply interface{}, err error) {
	return func(arguments interface{}) (reply interface{}, err error) {

		fmt.Println("file Picker")

		switch arguments.(bool) {
		case false:
			fmt.Println("filePicker fileExtension:" + fileExtension)
			fileDescriptor, _, err := dialog.File("select file", fileExtension, isDirectory)
			if err != nil {
				return nil, errors.Wrap(err, "failed to open dialog picker")
			}
			return fileDescriptor, nil

		case true:
			fileDescriptors, _, err := dialog.FileMulti("select files", fileExtension)
			if err != nil {
				return nil, errors.Wrap(err, "failed to open dialog picker")
			}

			//type []string is not supported by StandardMessageCodec
			sliceFileDescriptors := make([]interface{}, len(fileDescriptors))
			for i, file := range fileDescriptors {
				sliceFileDescriptors[i] = file
			}

			return sliceFileDescriptors, nil

		default:
			fileDescriptor, _, err := dialog.File("select file", fileExtension, isDirectory)
			if err != nil {
				return nil, errors.Wrap(err, "failed to open dialog picker")
			}
			return fileDescriptor, nil
		}

	}
}
