package intraproces

import (
	"github.com/gpioblink/go-auto-clean-fme-editor/editor/application"
	"github.com/gpioblink/go-auto-clean-fme-editor/editor/domain/lyric"
	"github.com/pkg/errors"
)

type LyricInterface struct {
	service application.LyricService
}

type AddLyricLyric struct {
	Point  AddLyricPoint
	Colors AddLyricColorPicker
	Lyric  AddLyricLyricString
}

type AddLyricPoint struct {
	X int
	Y int
}

type AddLyricColorPicker struct {
	BeforeCharColor    AddLyricColorPickerColor
	AfterCharColor     AddLyricColorPickerColor
	BeforeOutlineColor AddLyricColorPickerColor
	AfterOutlineColor  AddLyricColorPickerColor
}

type AddLyricColorPickerColor struct {
	Red   int
	Green int
	Blue  int
}

type AddLyricLyricString []AddLyricLyricChar

type AddLyricLyricChar struct {
	Furigana  string `json:"furigana"`
	Length    int    `json:"length"`
	LyricChar string `json:"char"`
}

func NewLyricsInterface(service application.LyricService) LyricInterface {
	return LyricInterface{service}
}

func (p LyricInterface) AddLyric(l AddLyricLyric) error {
	var lyricString lyric.LyricString
	for _, lst := range l.Lyric {
		lyricChar, err := lyric.NewLyricChar(lst.Furigana, lst.Length, lst.LyricChar)
		if err != nil {
			return errors.Wrap(err, "cannot parse lyricChar")
		}
		lyricString = append(lyricString, *lyricChar)
	}

	newPoint, err := lyric.NewPoint(l.Point.X, l.Point.Y)
	if err != nil {
		return errors.Wrap(err, "cannot parse point")
	}

	newColorBC, err := lyric.NewColor(l.Colors.BeforeCharColor.Red, l.Colors.BeforeCharColor.Green, l.Colors.BeforeCharColor.Blue)
	if err != nil {
		return errors.Wrap(err, "cannot parse color")
	}

	newColorAC, err := lyric.NewColor(l.Colors.AfterCharColor.Red, l.Colors.AfterCharColor.Green, l.Colors.AfterCharColor.Blue)
	if err != nil {
		return errors.Wrap(err, "cannot parse color")
	}

	newColorBO, err := lyric.NewColor(l.Colors.BeforeOutlineColor.Red, l.Colors.BeforeOutlineColor.Green, l.Colors.BeforeOutlineColor.Blue)
	if err != nil {
		return errors.Wrap(err, "cannot parse color")
	}

	newColorAO, err := lyric.NewColor(l.Colors.AfterOutlineColor.Red, l.Colors.AfterOutlineColor.Green, l.Colors.AfterOutlineColor.Blue)
	if err != nil {
		return errors.Wrap(err, "cannot parse color")
	}

	newColorPicker, err := lyric.NewColorPicker(*newColorBC, *newColorAC, *newColorBO, *newColorAO)
	if err != nil {
		return errors.Wrap(err, "cannot merge colors")
	}

	newLyric, err := lyric.NewLyric(*newPoint, *newColorPicker, lyricString)
	if err != nil {
		return errors.Wrap(err, "cannot make lyric")
	}

	return p.service.AddLyric(*newLyric)
}

// TODO: fme周りのintraprocessで使うためにprivateかpublicでlistLyricsを実装
