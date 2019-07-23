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
	//dialogProvider := dialogProvider{}

	channel := plugin.NewMethodChannel(messenger, channelName, plugin.StandardMethodCodec{})
	//channel.HandleFunc("openDirectory", p.filePicker(dialogProvider, true))
	/*
		channel.HandleFunc("ANY", p.filePicker(dialogProvider, false, "*"))
		channel.HandleFunc("IMAGE", p.filePicker(dialogProvider, false, "*"))
		channel.HandleFunc("AUDIO", p.filePicker(dialogProvider, false, "*"))
		channel.HandleFunc("VIDEO", p.filePicker(dialogProvider, false, "*"))
	*/
	channel.CatchAllHandleFunc(p.handleFilePicker)
	return nil
}

func (p *FilePickerPlugin) handleFilePicker(methodCall interface{}) (reply interface{}, err error) {
	var filter string

	method := methodCall.(plugin.MethodCall)
	multipleSelection := method.Arguments.(bool)

	switch method.Method {
	case "ANY":
		filter = "*"
	case "IMAGE":
		filter = "*"
	case "AUDIO":
		filter = "*"
	case "VIDEO":
		filter = "*"
	default:
		if strings.HasPrefix(method.Method, "__CUSTOM_") {
			resolveType := strings.Split(method.Method, "__CUSTOM_")
			filter = "Files (*." + resolveType[1] + ")\x00*." + resolveType[1] + "\x00All Files (*.*)\x00*.*\x00\x00"
			fmt.Println("handleFilePicker fileExtension:" + filter)
		} else {
			filter = "*"
		}
	}

	dialogProvider := dialogProvider{}
	fileDescriptor, err := p.filePicker(dialogProvider, false, filter, multipleSelection)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open dialog picker")
	}

	// return the randomized Method Name
	return fileDescriptor, nil
}

func (p *FilePickerPlugin) filePicker(dialog dialog, isDirectory bool, filter string, multipleSelection bool) (reply interface{}, err error) {
	fmt.Println("file Picker")

	switch multipleSelection {
	case false:
		fmt.Println("filePicker fileExtension:" + filter)
		fileDescriptor, _, err := dialog.File("select file", filter, isDirectory)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open dialog picker")
		}
		return fileDescriptor, nil

	case true:
		fileDescriptors, _, err := dialog.FileMulti("select files", filter)
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
		fileDescriptor, _, err := dialog.File("select file", filter, isDirectory)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open dialog picker")
		}
		return fileDescriptor, nil
	}

}
