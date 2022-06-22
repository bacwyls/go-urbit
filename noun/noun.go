package noun

import (
	"fmt"
	"math/big"
	"math/bits"
	"strconv"
	"unicode"
)

type MatTupl [2]*big.Int

type nounMap map[string]int64
type cueNounMap map[int64]Noun

type InvalidAtomError struct {
	Message string
}

func (e *InvalidAtomError) Error() string {
	return e.Message
}

// Noun is data
type Noun interface {
	isNoun() bool
	String() string
}

type Atom struct {
	Value *big.Int
}

type Cell struct {
	Head Noun
	Tail Noun
}

func Head(n Noun) Noun {
	switch t := n.(type) {
	case Cell:
		{
			return t.Head
		}
	default:
		return MakeNoun(0)
	}
}

func Tail(n Noun) Noun {
	switch t := n.(type) {
	case Cell:
		{
			return t.Tail
		}
	default:
		return MakeNoun(0)
	}
}

func B(i int64) *big.Int {
	return big.NewInt(i)
}

func (a Atom) isNoun() bool {
	return true
}
func (a Atom) String() string {
	return a.Value.Text(10)
}
func (a Cell) isNoun() bool {
	return true
}
func (a Cell) String() string {
	return "[" + a.innerString() + "]"
}
func (a Cell) innerString() string {
	switch t := a.Tail.(type) {
	case Cell:
		return a.Head.String() + " " + t.innerString()
	default:
		return a.Head.String() + " " + a.Tail.String()
	}
}

// AssertAtom returns error if not a valid Atom
func AssertAtom(n Noun) (Atom, error) {
	switch t := n.(type) {
	case Atom:
		{
			return t, nil
		}
	default:
		{
			return Atom{}, &InvalidAtomError{
				Message: fmt.Sprintf("Expected Atom. Received: %s", t),
			}
		}
	}
}

// Cut takes value from start including run
func Cut(start, run int64, b *big.Int) *big.Int {
	b1 := B(0).Rsh(b, uint(start))
	c1 := B(0).Mod(b1, B(0).Lsh(B(1), uint(run)))
	return c1
}

// Mat is Jam on atoms
func Mat(arg *big.Int) MatTupl {
	if arg.Cmp(B(0)) == 0 {
		return MatTupl{B(1), B(1)}
	}
	b := int64(arg.BitLen())
	fmt.Println("MAT arg: ", arg)
	fmt.Println("MAT arg bitlen: ", b)
	c := int64(len(fmt.Sprintf("%b", b)))
	tup1 := B(b + c + c)

	d1 := 1 << c // 2 ** c
	var d2 int64 = b % (1 << (c - 1))
	d3 := B(0).Lsh(arg, uint(c-1))
	d4 := B(0).Xor(d3, B(d2))
	d5 := B(0).Lsh(d4, uint(len(fmt.Sprintf("%b", d1))))
	tup2 := B(0).Add(d5, B(int64(d1)))

	return MatTupl{tup1, tup2}
}

func Rub(index int64, b *big.Int) (int64, Atom) {
	var c int64 = 0

	for ; b.Bit(int(index+c)) == 0; c++ {
	}

	if c == 0 {
		return 1, Atom{Value: B(0)}
	}

	d := index + c + 1
	d1 := Cut(d, c-1, b)
	e := B(0).Add(d1, B(0).Lsh(B(1), uint(c-1)))
	return c + c + e.Int64(), Atom{Value: Cut(d+c-1, e.Int64(), b)}
}

// CatLen is Cat but with a provided length
func CatLen(a, b *big.Int, length uint) *big.Int {
	b2 := B(0).Lsh(b, length)
	a2 := B(0).Xor(a, b2)
	return a2
}

// Cat concats two big ints
func Cat(a, b *big.Int) *big.Int {
	l := uint(a.BitLen())
	b2 := B(0).Lsh(b, l)
	a2 := B(0).Xor(a, b2)
	return a2
}

// MakeNoun takes an input and turns it into a Noun
func MakeNoun(arg interface{}) Noun {
	switch t := arg.(type) {
	case int:
		{
			return Atom{Value: B(int64(t))}
		}
	case int64:
		{
			return Atom{Value: B(t)}
		}
	case *big.Int:
		{
			return Atom{Value: t}
		}
	case Noun:
		return t
	case []string:
		// assume it is a `path`
		l := len(t)
		if l == 0 {
			return Atom{Value: B(0)}
		}
		return Cell{
			Head: MakeNoun(t[0]),
			Tail: MakeNoun(t[1:]),
		}
	case []interface{}:
		{
			l := len(t)
			if l == 0 {
				return Atom{Value: B(0)}
			}
			if l == 1 {
				return MakeNoun(t[0])
			}
			c := Cell{
				Head: MakeNoun(t[l-2]),
				Tail: MakeNoun(t[l-1]),
			}

			for k := range t[:l-2] {
				c = Cell{
					Head: MakeNoun(t[l-k-3]),
					Tail: c,
				}
			}
			return c
		}
	case string:
		{
			return StringToCord(t)
		}
	default:
		return Atom{Value: B(0)}
	}
}

func jamIn(nmap nounMap, n Noun, index int64) (int64, *big.Int) {
	if pIndex, ok := nmap[n.String()]; ok {
		switch t := n.(type) {
		case Atom:
			{
				if t.Value.BitLen() < bits.Len64(uint64(pIndex)) {
					d := Mat(t.Value)
					return 1 + d[0].Int64(), B(0).Lsh(d[1], 1)
				} else {
					fmt.Println("atom ignored: ", n)
				}
			}
		}
		fmt.Println("seen this noun b4: ", n)

		d1 := Mat(B(pIndex))
		d2 := B(0).Lsh(d1[1], 2)
		d3 := B(0).Xor(d2, B(3))
		return 2 + d1[0].Int64(), d3
	}

	nmap[n.String()] = index

	switch t := n.(type) {
	case Atom:
		{
			d := Mat(t.Value)
			return 1 + d[0].Int64(), d[1].Lsh(d[1], 1)
		}
	case Cell:
		{
			index += 2
			hidx, d1 := jamIn(nmap, t.Head, index)
			index += hidx
			tidx, d2 := jamIn(nmap, t.Tail, index)
			index += tidx
			d3 := Cat(d1, d2)
			d4 := B(0).Lsh(d3, 2)
			d5 := B(0).Xor(d4, B(1))
			return index, d5
		}
	}
	return index, B(0)
}

// Jam jams noun into a new NounMap
func Jam(n Noun) *big.Int {
	var nmap nounMap = make(nounMap)
	var index int64 = 0

	_, q1 := jamIn(nmap, n, index)
	return q1
}

func cueIn(nmap cueNounMap, b *big.Int, index int64) (int64, Noun) {
	a := b.Bit(int(index))
	// a == 0 > a is an atom
	index1 := index + 1
	if a == 0 {
		i, a1 := Rub(index1, b)
		nmap[index] = a1
		return i + 1, a1
	}

	index2 := index + 2
	a2 := b.Bit(int(index1))
	// when it is a Cell
	if a2 == 0 {
		i1, n1 := cueIn(nmap, b, index2)
		i2, n2 := cueIn(nmap, b, index2+i1)
		cell := Cell{
			Head: n1,
			Tail: n2,
		}
		nmap[index] = cell
		return i1 + i2 + 2, cell
	}

	// when it is a pointer, not atom or cell
	i3, a3 := Rub(index2, b)
	n3 := nmap[a3.Value.Int64()]

	return i3 + 2, n3
}

// Cue is the opposite of Jam
func Cue(b *big.Int) Noun {
	if b.Cmp(B(0)) == 0 {
		return MakeNoun(0)
	}
	var nmap cueNounMap = make(cueNounMap)
	var index int64 = 0

	_, q1 := cueIn(nmap, b, index)
	return q1
}

// StringToCord returns Atom of type cord
func StringToCord(str string) Atom {
	a := LittleToBig([]byte(str))
	return Atom{
		Value: a,
	}
}

// ByteLen returns the length of the big int in bytes
func ByteLen(b *big.Int) int64 {
	return int64(b.BitLen()-1)/8 + 1
}

// Snag gets the Head at the position pos
func Snag(n Noun, pos int) Noun {
	return Head(Slag(n, pos))
}

// Slag gets the Tail at the position pos
func Slag(n Noun, pos int) Noun {
	cur := n
	i := 0

	for i < pos {
		cur = Tail(cur)
		i++
	}
	return cur
}

// nock operator /
// get the noun at position a of n
func Fas(a uint64, n Noun) Noun {
	if a == 1 {
		return n
	}
	if a == 2 {
		return Head(n)
	}
	if a == 3 {
		return Tail(n)
	}
	if a%2 == 0 {
		return Fas(2, Fas(a/2, n))
	}
	if a%2 == 1 {
		return Fas(3, Fas((a-1)/2, n))
	}
	// This won't happen
	return MakeNoun(0)
}

// nock operator #
// replace the noun at position a of c with b
func Hax(a uint64, b Noun, c Noun) Noun {
	if a == 1 {
		return b
	}
	if a%2 == 0 {
		return Hax(a/2, MakeNoun([]interface{}{b, Fas(a+1, c)}), c)
	}
	if a%2 == 1 {
		return Hax((a-1)/2, MakeNoun([]interface{}{Fas(a-1, c), b}), c)
	}
	// This won't happen
	return MakeNoun(0)
}

func IsStrAtom(str string) bool {
	for _, c := range str {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

// is this string noun an explicit cell?
// explicit: begins and ends with [ and ] respectively
//           and has valid inner brackets
func IsStrXCell(str string) bool {
	if str[0] != '[' {
		return false
	}

	brax := 0
	for i, c := range str {
		if c == '[' {
			brax += 1
		}
		if c == ']' {
			brax -= 1
		}

		if brax == 0 {
			return i == len(str)-1
		}
	}
	return false
}

// return split index
func SplitStrCell(str string) int {
	brax := 0
	for i, c := range str {
		if c == '[' {
			brax += 1
		}
		if c == ']' {
			brax -= 1
		}
		if c == ' ' && brax == 0 {
			return i
		}
	}

	// TODO error here
	return -99
}

func NounFromString(str string) Noun {
	if IsStrAtom(str) {
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			bi := B(0)
			bi, ok := bi.SetString(str, 10)
			if !ok {
				fmt.Println("SetString: error")
				panic(err)
			}
			return MakeNoun(bi)
		}
		return MakeNoun(i)
	}

	if IsStrXCell(str) {
		str = str[1 : len(str)-1]
	}

	split := SplitStrCell(str)

	headRes := NounFromString(str[:split])
	tailRes := NounFromString(str[split+1:])
	res := MakeNoun([]interface{}{headRes, tailRes})
	return res
}
