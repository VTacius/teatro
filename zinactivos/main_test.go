package main
import "testing"

func TestObtieneElemento(t *testing.T)  {
	resultado := obtieneElemento("200802", 0, 4)
	if resultado != 2008 {
		t.Errorf("fechador(2008) no es 2008, es %d", resultado)
	}
}

func TestFechador(t *testing.T){
	resultado := fechador("2020 06 19 13 47 02.117Z")
	if resultado.String() != "2020-06-19 13:47:02 +0000 UTC"{
		t.Errorf("fechador(20200619134702.117Z) es 2020-06-19 13:47:02 +0000 UTC, no %s", resultado.String())
	}
}