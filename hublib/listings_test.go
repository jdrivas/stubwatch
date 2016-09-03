package hublib

import(
  "fmt"
  "testing"
  "sort"
  "github.com/stretchr/testify/assert"
)


func TestAlphaNumLess(t *testing.T) {

  trueCases := [][2]string {
    {"0", "1",},
    {"2", "13",},
    {"200","1000"},
    {"a", "b",},
    {"General Admission", "100",},
  }
  falseCases := [][2]string {
    {"500", "300"},
  }
  for i, c := range trueCases {
    s := fmt.Sprintf("Case %d. %s < %s", i+1, c[0], c[1])
    assert.True(t, alphaNumLess(c[0], c[1]), s)
  }
  for i, c := range falseCases {
    s := fmt.Sprintf("Case %d. %s < %s", i+1, c[0], c[1])
    assert.False(t, alphaNumLess(c[0], c[1]), s)
  }
}

func lStr(l Listing) (string) {
  return fmt.Sprintf("%s %s %s:%s", l.ZoneName, l.SectionName, l.Row, l.SeatNumbers)
}

func TestSortListings(t *testing.T) {
  l1 := Listing{
    ZoneName: "Dugout Club",
    SectionName: "Dugout Club 112",
    Row: "CCC",
    SeatNumbers: "1,2",
  }
  l2 := Listing{
    ZoneName: "Dugout Club",
    SectionName: "Dugout Club 112",
    Row: "BBB",
    SeatNumbers: "5,6,7,8",
  }
  l3 := Listing{
    ZoneName: "Premium Field Club",
    SectionName: "Premium Field Club 121",
    Row: "A",
    SeatNumbers: "9,10,12",
  }
  l4 := Listing{
    ZoneName: "Premium Field Club",
    SectionName: "Premium Field Club 121",
    Row: "B",
    SeatNumbers:"1,2",
  }
  l6 := Listing{
    ZoneName: "Premium Lower Box",
    SectionName: "Premium Lower Box 105",
    Row: "4",
    SeatNumbers:"1,2",
  }
  l8 := Listing{
    ZoneName: "Premium Lower Box",
    SectionName: "Premium Lower Box 105",
    Row: "4",
    SeatNumbers:"3,4",
  }  
  l5 := Listing{
    ZoneName: "Premium Lower Box",
    SectionName: "Premium Lower Box 106",
    Row: "13",
    SeatNumbers:"1,2",
  }

  l10 := Listing{
    ZoneName: "Premium Lower Box",
    SectionName: "Premium Lower Box 106",
    Row: "2",
    SeatNumbers: "General Admission",
  }
  l7 := Listing{
    ZoneName: "Premium Lower Box",
    SectionName: "Premium Lower Box 106",
    Row: "2",
    SeatNumbers:"2,3",
  }
  l9 := Listing{
    ZoneName: "Premium Lower Box",
    SectionName: "Premium Lower Box 106",
    Row: "2",
    SeatNumbers:"11,12",
  }

  ls := []Listing{l5, l3, l6, l1, l4, l2, l10, l7,l8, l9,}

  // Test Swap
  assert.Equal(t, ls[0].Row, "13")
  assert.Equal(t, ls[1].Row, "A")
  ByZoneSectionRowSeat(ls).Swap(0,1)
  assert.Equal(t, ls[1].Row, "13")
  assert.Equal(t, ls[0].Row, "A")

  // Test sort.
  ex := []Listing{l2, l1, l3, l4, l6, l8, l10, l7, l9, l5,}
  sort.Sort(ByZoneSectionRowSeat(ls))
  for i, _ := range ls {
    // fmt.Printf("%d: %s, %s, %s\n", i+1, l.ZoneName, l.SectionName, l.Row)
    assert.Equal(t, ex[i].ZoneName, ls[i].ZoneName, "(%d.) Zone names should be equal. exp: %s, act: %s", i, lStr(ex[i]), lStr(ls[i]))
    assert.Equal(t, ex[i].SectionName, ls[i].SectionName, "(%d). Section names should be equal. exp: %s, act: %s", i, lStr(ex[i]), lStr(ls[i]))
    assert.Equal(t, ex[i].Row, ls[i].Row, "(%d). Rows should be equal. exp: %s, act: %s", i, lStr(ex[i]), lStr(ls[i]))
    assert.Equal(t, ex[i].SeatNumbers, ls[i].SeatNumbers, "(%d). Seats should be equal. exp: %s, act: %s", i, lStr(ex[i]), lStr(ls[i]))
  }

}