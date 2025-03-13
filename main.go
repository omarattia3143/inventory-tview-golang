package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"

	"github.com/rivo/tview"
)

type Item struct {
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}

var (
	inventory     []Item
	inventoryFile = "inventory.json"
)

func loadInventory() {
	data, err := os.ReadFile(inventoryFile)
	if err != nil {
		log.Fatal("error reading inventory file - ", err)
	}

	json.Unmarshal(data, &inventory)

}

func saveInventry() {
	data, err := json.MarshalIndent(inventory, "", " ")
	if err != nil {
		log.Fatal("failed to save the invetory file - ", err)

	}

	os.WriteFile(inventoryFile, data, 0644)
}

func deleteIndex(i int) {
	if i < 0 || i >= len(inventory) {
		fmt.Println("invalid index")
		return
	}

	inventory = slices.Delete(inventory, i, i+1)
	saveInventry()
}

func main() {
	app := tview.NewApplication()

	loadInventory()

	inventoryList := tview.NewTextView().
		SetDynamicColors(true). // Enable dynamic coloring of text
		SetRegions(true).       // Allows regions for interaction (not used here)
		SetWordWrap(true)       // Enables word wrapping to fit the TextView size

	inventoryList.SetBorder(true).SetTitle("Inventory Items") // Set border and title

	refreshInventory := func() {
		inventoryList.Clear()
		if len(inventory) == 0 {
			fmt.Fprintln(inventoryList, "no items in the inventory")
		} else {
			for i, item := range inventory {
				fmt.Fprintf(inventoryList, "[%d] %s (Stock:%d) \n", i+1, item.Name, item.Stock)
			}
		}
	}

	itemNameInput := tview.NewInputField().SetLabel("Item Name: ")
	itemStockInput := tview.NewInputField().SetLabel("Stock: ")
	itemIDInput := tview.NewInputField().SetLabel("Item ID: ")

	form := tview.NewForm().
		AddFormItem(itemNameInput).
		AddFormItem(itemStockInput).
		AddFormItem(itemIDInput).
		AddButton("Add Item", func() {
			name := itemNameInput.GetText()
			stock := itemStockInput.GetText()

			if name != "" && stock != "" {
				quantity, err := strconv.Atoi(stock)
				if err != nil {
					fmt.Fprintln(inventoryList, "Invalid stock value")
					return
				}
				inventory = append(inventory, Item{Name: name, Stock: quantity})
				saveInventry()
				refreshInventory()
				itemNameInput.SetText("")
				itemStockInput.SetText("")

			}
		}).
		AddButton("Delete Item", func() {
			idStr := itemIDInput.GetText()

			if idStr == "" {
				fmt.Fprintln(inventoryList, "Please enter item ID to delete.")
				return
			}

			id, err := strconv.Atoi(idStr)

			if err != nil || id < 1 || id > len(inventory) {
				fmt.Fprintln(inventoryList, "Invalid item ID")
				return
			}

			deleteIndex(id - 1)
			fmt.Fprintf(inventoryList, "Item [%d] deleted\n", id)
			refreshInventory()
			itemIDInput.SetText("")
		}).
		AddButton("Exit", func() {
			app.Stop()
		})

	form.SetBorder(true).SetTitle("Manage Inventory").SetTitleAlign(tview.AlignLeft)

	flex := tview.NewFlex().
		AddItem(inventoryList, 0, 1, false).
		AddItem(form, 0, 1, true)

	refreshInventory()

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
