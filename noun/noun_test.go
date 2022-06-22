package noun

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestMat(t *testing.T) {
	r1 := Mat(B(100))
	c1 := [2]*big.Int{B(13), B(6456)}
	if !reflect.DeepEqual(r1[0], c1[0]) && !reflect.DeepEqual(r1[1], c1[1]) {
		t.Errorf("expected %s got %s", c1, r1)
	}

	r2 := Mat(B(100000000))
	c2 := [2]*big.Int{B(37), B(102400000736)}
	if !reflect.DeepEqual(r2[0], c2[0]) && !reflect.DeepEqual(r2[1], c2[1]) {
		t.Errorf("expected %s got %s", c2, r2)
	}

	b1 := B(0)
	b1.SetString("100000000000000000001", 10)
	b2 := B(0)
	b2.SetString("1638400000000000000017280", 10)
	r3 := Mat(b1)
	c3 := [2]*big.Int{B(18), b2}
	if !reflect.DeepEqual(r3[0], c3[0]) && !reflect.DeepEqual(r3[1], c3[1]) {
		t.Errorf("expected %s got %s", c3, r3)
	}

	r4 := Mat(B(0))
	c4 := [2]*big.Int{B(1), B(1)}
	if !reflect.DeepEqual(r4[0], c4[0]) && !reflect.DeepEqual(r4[1], c4[1]) {
		t.Errorf("expected %s got %s", c4, r4)
	}

	r5 := Mat(B(1))
	c5 := [2]*big.Int{B(3), B(6)}
	if !reflect.DeepEqual(r5[0], c5[0]) && !reflect.DeepEqual(r5[1], c5[1]) {
		t.Errorf("expected %s got %s", c5, r5)
	}
}

func TestMakeNoun(t *testing.T) {
	r1 := MakeNoun(100)
	c1 := "100"
	if r1.String() != c1 {
		t.Errorf("expected %s got %s", c1, r1)
	}
	r2 := MakeNoun([]interface{}{100, []interface{}{12, 16}, B(255), 24})
	c2 := "[100 [12 16] 255 24]"
	if r2.String() != c2 {
		t.Errorf("expected %s got %s", c2, r2)
	}
}

func TestFas(t *testing.T) {
	a := []interface{}{12, 16, 19, 23}
	b := []interface{}{13, 17, 20, 24}
	n := MakeNoun([]interface{}{a, b})
	s4 := Fas(4, n)

	if s4.String() != "12" {
		t.Errorf("expected 12 got %s", s4)
	}

	s5 := Fas(5, n)
	s5_truth := "[16 19 23]"
	if s5.String() != s5_truth {
		t.Errorf("expected %s got %s", s5_truth, s5)
	}

	s10 := Fas(10, n)
	if s10.String() != "16" {
		t.Errorf("expected 16 got %s", s10)
	}

	s6 := Fas(6, n)
	if s6.String() != "13" {
		t.Errorf("expected 13 got %s", s6)
	}

	s14 := Fas(14, n)
	if s14.String() != "17" {
		t.Errorf("expected 17 got %s", s14)
	}

}

func TestHax(t *testing.T) {
	a := []interface{}{12, 16, 19, 23}
	b := []interface{}{13, 17, 20, 24}
	n := MakeNoun([]interface{}{a, b})

	r := 33

	h5 := Hax(5, MakeNoun(r), n)
	amod := []interface{}{12, r}
	h5_truth := MakeNoun([]interface{}{amod, b})
	if h5.String() != h5_truth.String() {
		t.Errorf("expected %s got %s", h5_truth, h5)
	}

	h7 := Hax(7, MakeNoun(r), n)
	bmod := []interface{}{13, r}
	h7_truth := MakeNoun([]interface{}{a, bmod})
	if h7.String() != h7_truth.String() {
		t.Errorf("expected %s got %s", h7_truth, h7)
	}

	/*
		//	debugging hax TODO remove
			n = MakeNoun([]interface{}{55, []interface{}{66, 77}, 88})
			//	alt := MakeNoun([]interface{}{55, []interface{}{66, 33}, 88})
			n = MakeNoun([]interface{}{n, n})
			h13 := Hax(20, MakeNoun(r), n)
			if h13.String() != "" {
				t.Errorf("expected %s got %s", h7_truth, h13)
			}
	*/

}

func TestNounFromString(t *testing.T) {
	a := []interface{}{12, 16, 19, 23}
	b := []interface{}{13, 17, 20, 24}
	n := MakeNoun([]interface{}{a, b})
	ns := NounFromString(n.String())

	if ns.String() != n.String() {
		t.Errorf("expected %s got %s", n.String(), ns.String())
	}

	// a sample landscape dm
	// coverage for bigints and deeper nesting
	dm := "[170141184505687446085187436250516160512 2128681655491098862689 [101844224 8678263233316744548] [[11093637012246365440 170141184505687446085091054714890804527 0] [0 101844224 [11093637012246365440 170141184505687446085091054714890804527 0] 170141184505687446085091054714890804527 [[1954047348 5018395299700377944449393971777] 0] [0 49678085277749340053015205300505906542] [1430383508112382614187095719319023241763957391717204612572259633867244619751332669867972417882692126456701976374054689233573579725316335638796765440959819871434272292658192524339426869369119827168292865 101844224 1] 0 0] 521644567909 0] 0 0]"
	dmn := NounFromString(dm)
	if dmn.String() != dm {
		t.Errorf("expected %s got %s", dm, dmn.String())
	}

}

func TestJam(t *testing.T) {
	n1 := MakeNoun([]interface{}{12, 16})
	r1 := Jam(n1)
	if r1.Int64() != 17176641 {
		t.Errorf("expected %s got %s", n1, r1)
	}

	a := []interface{}{12, 16, 19, 23}
	b := []interface{}{12, 16, 19, 23}
	n2 := MakeNoun([]interface{}{a, b})
	r2 := Jam(n2)
	if r2.Int64() != 5322556398681252101 {
		t.Errorf("expected %s got %s", n2, r2)
	}

	n3 := MakeNoun([]interface{}{[]string{"ge", "hood"}, 0, "m", "helm-hi", MakeNoun("ping")})
	c1 := "83103842581186151537609419784725107274636599623840339663322629"
	r3 := Jam(n3)
	if r3.Text(10) != c1 {
		t.Errorf("expected %s got %s", c1, r3)
	}
}

func TestCue(t *testing.T) {
	n1 := B(17176641)
	r1 := Cue(n1)
	c1 := MakeNoun([]interface{}{12, 16})

	if r1.String() != c1.String() {
		t.Errorf("expected %s got %s", c1, n1)
	}

	a := []interface{}{12, 16, 19, 23}
	b := []interface{}{12, 16, 19, 23}
	n2 := MakeNoun([]interface{}{a, b})
	r2 := Cue(B(5322556398681252101))
	if r2.String() != n2.String() {
		t.Errorf("expected %s got %s", n2, r2)
	}
}

func TestStringToCord(t *testing.T) {
	n1 := "ping"
	c1 := "676e6970"
	r1 := StringToCord(n1)
	if r1.Value.Text(16) != c1 {
		t.Errorf("expected %s got %s", c1, r1)
	}
}

func TestAssertAtom(t *testing.T) {
	_, err := AssertAtom(MakeNoun("12"))
	if err != nil {
		t.Errorf("expected %v got %e", nil, err)
	}
	_, err = AssertAtom(Atom{})
	if err != nil {
		t.Errorf("expected %v got %e", nil, err)
	}
	_, err = AssertAtom(Cell{})
	if err == nil {
		t.Errorf("expected error got nil")
	}
}

func ExampleMakeNoun() {
	stringNoun := MakeNoun("string value")

	fmt.Println(stringNoun)
	// Output: 31399942126277005645796504691
}
