package gui

import (
	"fmt"
	"github.com/1920853199/mysql-viewer/mysql"
	"github.com/jroimartin/gocui"
	"log"
	"strings"
)

type Window struct {
	g		*gocui.Gui
	sql 	mysql2.Db
}
var (
	viewArr = []string{"v1", "v2", "v3", "v4"}
	active  = 0
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	out, err := g.View("v2")
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "Going from view "+v.Name()+" to "+name)

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	if nextIndex == 0 || nextIndex == 3 {
		g.Cursor = true
	} else {
		g.Cursor = false
	}

	active = nextIndex
	return nil
}


func (w *Window)layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("v1", 0, 0, maxX/5*3-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Mysql"
		v.Editable = true
		v.Wrap = true
		v.Autoscroll = true
		w.writeArrow(v)

		if _, err = setCurrentViewOnTop(g, "v1"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("v2", maxX/5*3, 0, maxX-1, maxY/3-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Explain"
		v.Wrap = true
		v.Autoscroll = true
	}
	if v, err := g.SetView("v3", 0, maxY/2, maxX/5*3-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Result"
		v.Autoscroll = true
	}
	if v, err := g.SetView("v4", maxX/5*3, maxY/3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Show Profile"
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (w *Window)execCmd(g *gocui.Gui, v *gocui.View) error {
	var line string
	var err error
	defer func() {
		v.EditNewLine()
		w.writeArrow(v)
	}()

	_, cy := v.Cursor()
	if line, err = v.Line(cy); err != nil {
		line = ""
		return nil
	}
	v.SetCursor(len(line), cy)

	inputs := ""
	splits := strings.Split(line, "mysql>")
	if len(splits) > 1 {
		inputs = splits[1]
	}
	if inputs == "" {
		return nil
	}

	v2, _ := g.View("v2")
	v2.Clear()
	v2.SetCursor(0, 0)
	v3, _ := g.View("v3")
	v3.Clear()
	v3.SetCursor(0, 0)
	v4, _ := g.View("v4")
	v4.Clear()
	v4.SetCursor(0, 0)

	columns,result, err := w.sql.Explain(inputs)
	if err != nil {
		fmt.Fprint(v3, err.Error())
		//w.writeText(v3, err.Error())
		return nil
	}
	var id string
	id = result[0]["id"]
	w.writeExplain(v2,result,columns)


	columns,result, err = w.sql.Result(inputs)
	if err != nil {
		fmt.Fprint(v3, err.Error())
		//w.writeText(v2, err.Error())
		return nil
	}
	w.writeResult(v3,result,columns)


	_,result, err = w.sql.Profile(id)
	if err != nil {
		fmt.Fprint(v3, err.Error())
		//w.writeText(v2, err.Error())
		return nil
	}
	w.writeProfile(v4,result)

	//fmt.Printf("%s",result)
	return nil
}

func (w *Window) writeExplain(v *gocui.View,result []map[string]string,columns []string) {
	//select * from article
	//v.EditNewLine()
	str := fmt.Sprintf("+---------------+---------------------------+")

	w.writeText(v, str)
	v.EditNewLine()
	for _,item := range result {
		for _,value := range columns {
			str = fmt.Sprintf("|")
			w.writeText(v, str)
			str = fmt.Sprintf("%-15s",value)
			w.writeText(v, str)
			str = fmt.Sprintf("|")
			w.writeText(v, str)
			str = fmt.Sprintf("%-27s",item[value])
			w.writeText(v, str)
			str = fmt.Sprintf("|")
			w.writeText(v, str)
			v.EditNewLine()
		}
	}
	str = fmt.Sprintf("+---------------+---------------------------+")
	w.writeText(v, str)
	v.EditNewLine()
}

func (w *Window) writeResult(v *gocui.View,result []map[string]string,columns []string) {
	//select * from article
	//v.EditNewLine()
	//str := fmt.Sprintf("+---------------+---------------------------+")

	//w.writeText(v, str)
	//v.EditNewLine()
	var str string
	for k,item := range result {
		for _,value := range columns {
			if k == 0 {
				str = fmt.Sprintf("|")
				w.writeText(v, str)
				str = fmt.Sprintf("%-10s",value)
				w.writeText(v, str)
			}else {
				str = fmt.Sprintf("|")
				w.writeText(v, str)
				str = fmt.Sprintf("%-10s", item[value])
				w.writeText(v, str)
			}

		}
		str = fmt.Sprintf("|")
		w.writeText(v, str)
		v.EditNewLine()
	}
	//str = fmt.Sprintf("+---------------+---------------------------+")
	//w.writeText(v, str)
	v.EditNewLine()
}

func (w *Window) writeProfile(v *gocui.View,result []map[string]string) {
	//select * from article
	//v.EditNewLine()
	str := fmt.Sprintf("+---------------+---------------------------+")

	w.writeText(v, str)
	v.EditNewLine()
	for _,item := range result {
		for key,value := range item {
			str = fmt.Sprintf("|")
			w.writeText(v, str)
			str = fmt.Sprintf("%-15s",key)
			w.writeText(v, str)
			str = fmt.Sprintf("|")
			w.writeText(v, str)
			str = fmt.Sprintf("%-27s",value)
			w.writeText(v, str)
			str = fmt.Sprintf("|")
			w.writeText(v, str)
			v.EditNewLine()
		}
	}
	str = fmt.Sprintf("+---------------+---------------------------+")
	w.writeText(v, str)
	v.EditNewLine()
}


func (w *Window) writeArrow(v *gocui.View) {
	line := fmt.Sprintf("mysql>")
	w.writeText(v, line)
}

func (w *Window) writeText(v *gocui.View, str string) {
	runes := []rune(str)
	for _, c := range runes {
		v.EditWrite(c)
	}
}



func Run() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	w := Window{g:g}
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(w.layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("v1", gocui.KeyEnter, gocui.ModNone,w.execCmd); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
