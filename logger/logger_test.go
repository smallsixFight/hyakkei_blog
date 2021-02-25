package logger

import (
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	Debug("deug")
	Debugf("detbug %s", "a")
	var a int
	for a < 10 {
		Info("info")
		Infof("info: %d, %s", 1, "4444")
		Print("Print")
		Printf("Print: %d, %s", 3, "43434")
		Println("Print: !d, !s", 3, "43434")
		Warn("Warn")
		Warnf("Warnf: %d, %s", 1, "4444")
		Error("error")
		Errorf("errorf: %d, %s", 1, "4444")
		time.Sleep(time.Millisecond)
		a++
	}

	//Fatal("Fatal")
	//Fatalln("Fatalln")
	//Fatalf("Fatalf: %d, %s", 1, "4444")

	//Panic("Panic")
	//Panicln("Panicln")
	//Panicf("Panicf: %d, %s", 1, "4444")
}

func TestNewLogger(t *testing.T) {
	Init(LowestLevel(WarnLevel), IsPrint(false), FileNamePre("blog"))
	Debug("deug")
	Debugf("detbug %s", "a")
	Info("info")
	Infof("info: %d, %s", 1, "4444")
	Print("Print")
	Printf("Print: %d, %s", 3, "43434")
	Error("error")
}
