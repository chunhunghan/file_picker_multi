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
	var fileExtension string

	method := methodCall.(plugin.MethodCall)
	multipleSelection := method.Arguments.(bool)

	switch method.Method {
	case "ANY":
		fileExtension = "*"
	case "IMAGE":
		fileExtension = "*"
	case "AUDIO":
		fileExtension = "*"
	case "VIDEO":
		fileExtension = "*"
	default:
		if strings.HasPrefix("method.Method", "__CUSTOM_") {
			resolveType := strings.Split(method.Method, "__CUSTOM_")
			fileExtension = resolveType[1]
			fmt.Println("handleFilePicker fileExtension:" + fileExtension)
		} else {
			fileExtension = "*"
		}
	}

	dialogProvider := dialogProvider{}
	fileDescriptor, err := p.filePicker(dialogProvider, false, fileExtension, multipleSelection)
	if err != nil {
		fmt.Println("user cancel")
		return nil, nil
	}

	// return the randomized Method Name
	return fileDescriptor, nil
}

func (p *FilePickerPlugin) filePicker(dialog dialog, isDirectory bool, fileExtension string, multipleSelection bool) (reply interface{}, err error) {
	fmt.Println("file Picker")

	switch multipleSelection {
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
