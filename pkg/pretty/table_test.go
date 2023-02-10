package pretty

import (
	"fmt"
	"testing"
)

func TestTableRender(t *testing.T) {
	table := NewTable(5)
	table.SetColumnMinSize(0, 8)
	table.SetColumnMinSize(1, 7)
	table.SetColumnMinSize(2, 7)
	table.SetColumnMinSize(3, 10)
	table.SetColumnMinSize(4, 8)
	table.SetColumnAlign(0, TextAlignRight)
	table.SetHeader([]string{"symbol", "low", "high", "date", "time utc"})
	table.AddRow([]string{"Item1", "0.2345", "0.2453", "2022-11-28", "14:00:00"})
	table.AddRow([]string{"Item2", "0.7988", "0.8311", "2022-11-28", "14:00:00"})
	s := table.Render()

	fmt.Printf("%s\n", s)
}
