package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"io"
)

func HandleShowDirectorsPage(c buffalo.Context) error {
	c.Set("directors", LoadAllDirectors())
	return c.Render(200, r.HTML("directors.html"))
}

func HandleAddDirector(c buffalo.Context) error {
	request := ManageDirectorRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	AddDirectorToDB(request)
	return c.Redirect(302, "/directors")
}

func HandleEditDirector(c buffalo.Context) error {
	request := ManageDirectorRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	EditDirectorInDB(request)
	return c.Redirect(302, "/directors")
}

func HandleDirectorDelete(c buffalo.Context) error {
	request := OnlyID{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	DeleteDirector(request.ID)
	return c.Redirect(302, "/directors")
}

func HandleExportDirectorsToWord(c buffalo.Context) error {
	c.Response().Header().Set("Content-Disposition", "attachment; filename=directors.docx")
	return c.Render(200, r.Func("application/vnd.openxmlformats-officedocument.wordprocessingml.document", func(w io.Writer, d render.Data) error {
		data, err := ExportDirectorsToWord(DirectorsExportPayload{
			Directors: LoadAllDirectors(),
		})
		if err != nil {
			return err
		}
		_, writeError := w.Write(data)
		return writeError
	}))
}
