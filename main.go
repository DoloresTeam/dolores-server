package main

import "github.com/gin-gonic/gin"

// Department ...
type Department struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"pid"`
}

// Person ...
type Person struct {
	ID     string   `json:"id"`
	CN     string   `json:"name"`
	SN     string   `json:"realName"`
	Title  string   `json:"title"`
	UnitID []string `json:"departmentIDs"`
}

func main() {

	d1 := &Department{`1`, `D1`, ``}
	d2 := &Department{`2`, `D2`, ``}
	d3 := &Department{`3`, `D3`, `1`}
	d4 := &Department{`4`, `D4`, `1`}

	p1 := &Person{`1`, `m1`, `M1`, `developer`, []string{`1`}}
	p2 := &Person{`2`, `m2`, `M2`, `developer`, []string{`2`}}
	p3 := &Person{`3`, `m3`, `M3`, `developer`, []string{`1`, `3`}}

	router := gin.Default()

	router.GET(`/organization`, func(c *gin.Context) {
		c.JSON(200, map[string]interface{}{
			`departments`: []*Department{d1, d3, d4, d2},
			`members`:     []*Person{p1, p2, p3},
			`version`:     1,
		})
	})

	router.Run(`:3280`)
}
