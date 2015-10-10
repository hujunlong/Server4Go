package i18n

import (
	"fmt"
	"testing"
)

func TestI18n(t *testing.T) {
	err := GetInit("locale_test.ini")
	if err != nil {
		t.Errorf("error:%s", err)
	}
	fmt.Println(Data["hi"])
	fmt.Println("bye")
	if Data["hi"] != "aaa" {
		t.Errorf("cout = %s %d", Data["hi"], len(Data["hi"]))
	}
	if Data["bye"] != "涨" {
		t.Error("data error")
	}

}
