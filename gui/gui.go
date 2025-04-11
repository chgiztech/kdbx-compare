package gui

import (
	"fmt"
	"kdbx-compare/compare"
	"kdbx-compare/database"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func CreateGUI() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Compare Two KeePass Databases")

	file1Entry := widget.NewEntry()
	file1Entry.SetPlaceHolder("File Path for Database 1")

	pass1Entry := widget.NewPasswordEntry()
	pass1Entry.SetPlaceHolder("Password for Database 1")

	file2Entry := widget.NewEntry()
	file2Entry.SetPlaceHolder("File Path for Database 2")

	pass2Entry := widget.NewPasswordEntry()
	pass2Entry.SetPlaceHolder("Password for Database 2")

	output := widget.NewMultiLineEntry()
	output.SetPlaceHolder("Differences will be displayed here")

	compareButton := widget.NewButton("Compare", func() {
		file1 := file1Entry.Text
		pass1 := pass1Entry.Text
		file2 := file2Entry.Text
		pass2 := pass2Entry.Text

		db1, err := database.LoadDatabase(file1, pass1)

		if err != nil {
			output.SetText(fmt.Sprintf("Error loading DB1: %v", err))
			return
		}

		db2, err := database.LoadDatabase(file2, pass2)
		if err != nil {
			output.SetText(fmt.Sprintf("Error loading DB2: %v", err))
			return
		}

		// Compare and display results
		result := compare.CompareDatabases(db1, db2)
		output.SetText(result)

		saveFileDialog := dialog.NewFileSave(func(file fyne.URIWriteCloser, err error) {
			if err != nil || file == nil {
				output.SetText("Error: unable to save the file.")
				return
			}

			// Записываем результат в файл
			err = os.WriteFile(file.URI().Path(), []byte(result), 0644)
			if err != nil {
				output.SetText(fmt.Sprintf("Error saving report: %v", err))
			} else {
				output.SetText("Report saved successfully.")
			}
		}, myWindow)

		saveFileDialog.Show()
	})

	myWindow.SetContent(container.NewVBox(
		file1Entry,
		pass1Entry,
		file2Entry,
		pass2Entry,
		compareButton,
		output,
	))

	myWindow.ShowAndRun()
}
