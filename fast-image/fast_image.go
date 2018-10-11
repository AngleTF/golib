package fastImage

import (
	"image/jpeg"
	"os"
	"fmt"
	"image"
	"image/png"
	"strings"
)

const(
	JPEG = "jpeg"
	PNG  = "png"
	JPG  = "jpg"
)



type Setting struct{
	//图片质量 0 ~ 100
	Quality	int
	//图片地址
	PicUrl  string
	//图片宽度
	PicWidth int
	//图片高度
	PicHeight int
	//io.Read
	Io *os.File
	//图片类型
	PicType string

	//导出的图片格式
	PicToType string
}

type PicError struct{
	Info string
}

func (ctx PicError) Error() string{
	return fmt.Sprintf("PicError: %s", ctx.Info)
}


func NewPic(url string, width, height int , picType, picToType string, quality int) (*Setting, error){

	f, err := os.Open(url)
	if err != nil{
		return nil, PicError{"open file error"}
	}

	picType = strings.ToLower(picType)

	if picType != JPEG && picType != PNG && picType != JPG{
		return nil, PicError{"picType range JPEG, JPG, PNG const"}
	}

	if quality < 0 || quality > 100{
		return nil, PicError{"quality range 0 ~ 100"}
	}

	return &Setting{
		Quality: quality,
		PicUrl:url,
		PicWidth: width,
		PicHeight: height,
		Io:f,
		PicType:picType,
		PicToType:picToType,
	}, nil
}

func (ctx *Setting) Compress(toUrl string) error{
	var ig image.Image
	var err error
	var file *os.File
	//var cnf image.Config
	//var currentHeight, currentWidth int
	switch ctx.PicType{
	case PNG:
		ig, err = png.Decode(ctx.Io)
		if err != nil{
			return err
		}
	case JPEG:
		fallthrough
	case JPG:
		ig, err = jpeg.Decode(ctx.Io)
		if err != nil{
			return err
		}

		ctx.Io.Close()
		file , err = os.Create(toUrl)
		defer file.Close()
		if err != nil{
			return err
		}
	}

	switch ctx.PicToType{
	case PNG:
		err = png.Encode(file, ig)
		if err != nil{
			return err
		}
	case JPEG:
		fallthrough
	case JPG:
		err = jpeg.Encode(file, ig, &jpeg.Options{ctx.Quality})
		if err != nil{
			return err
		}
	}
	return nil
}