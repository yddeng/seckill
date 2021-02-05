package sdk

import (
	"fmt"
	"testing"
)

func TestGetCookie(t *testing.T) {
	LoadCookie("../my.cookies")
	fmt.Println(GetCookie("ipLoc-djd"))
	fmt.Println(GetCookie("wlfstk_smdl"))
}
