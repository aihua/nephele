package command

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	cat "github.com/ctripcorp/cat.go"
	"github.com/ctripcorp/nephele/imgsvr/img4g"
	"github.com/ctripcorp/nephele/util"
	"golang.org/x/net/context"
	"regexp"
)

var (
	generalArgPattern           = util.RegexpExt{regexp.MustCompile("^/images/(.*?)_(R|C|Z|W)_([0-9]+)_([0-9]+)(_R([0-9]+))?(_C([a-zA-Z]+))?(_Q(?P<n0>[0-9]+))?(_M((?P<wn>[a-zA-Z0-9]+)(_(?P<wl>[1-9]))?))?(_(?P<dwm>D))?.(?P<ext>jpg|jpeg|gif|png|Jpg)$")}
	sourceDigitalMarkArgPattern = util.RegexpExt{regexp.MustCompile("^/images/(.*?)(_(?P<dwm>D)).(?P<ext>jpg|jpeg|Jpg)$")}
)

type CommandArgument struct {
	StoragePath         string
	ImageExtension      string
	digitalMarkOnSource bool
}

type Command interface {
	Exec(context.Context, *img4g.Image) error
}

type CompositeCommand interface {
	Command
}

type GeneralComposite struct {
	commands []Command
}

func (g *GeneralComposite) Exec(c context.Context, img *img4g.Image) error {
	for _, cmd := range g.commands {
		err := cmd.Exec(c, img)
		if err != nil {
			return err
		}
	}

	return nil
}

func Parse(ctx context.Context) (*CommandArgument, error) {
	var (
		err        error
		commandArg *CommandArgument
		catVar     = ctx.Value("cat").(cat.Cat)
		imagePath  = ctx.Value("imagepath").(string)
	)

	params, ok := generalArgPattern.FindStringSubmatchMap(imagePath)
	if ok {
		commandArg, err = makeGeneralArg(imagePath, params)
	} else {
		params, ok = sourceDigitalMarkArgPattern.FindStringSubmatchMap(imagePath)
		if ok {
			commandArg.digitalMarkOnSource = true
			commandArg, err = makeDigitalMarkOnSourceArg(params)
		} else {
			err = errors.New("Image.Command.ParseError")
			log.WithFields(log.Fields{"uri": imagePath}).Warn(err.Error())
			util.LogErrorEvent(catVar, "URI.ParseError", "")
		}
	}

	return commandArg, err
}

func Get(arg CommandArgument) (CompositeCommand, error) {
	return nil, nil
}

func makeGeneralArg(imagePath string, params map[string]string) (*CommandArgument, error) {
	var ok bool
	arg := &CommandArgument{}
	arg.StoragePath, ok = params[":1"]
	if !ok {
		return arg, errors.New(fmt.Sprintf("invalid storage path on imagepath=%s", imagePath))
	}
	arg.ImageExtension, ok = params["ext"]
	if !ok {
		return arg, errors.New(fmt.Sprintf("invalid image extension on imagepath=%s", imagePath))
	}
	return arg, nil
}

func makeDigitalMarkOnSourceArg(params map[string]string) (*CommandArgument, error) {
	return nil, nil
}
