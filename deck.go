package cardrank

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// Shuffler is an interface for a deck shuffler. Compatible with
// math/rand.Rand's Shuffle method.
type Shuffler interface {
	Shuffle(n int, swap func(int, int))
}

// DeckType is a deck type.
type DeckType uint8

// Deck types.
const (
	// DeckFrench is a standard deck of 52 playing cards.
	DeckFrench = DeckType(Two)
	// DeckShort is a deck of 36 playing cards of rank 6+ (see [Short]).
	DeckShort = DeckType(Six)
	// DeckManila is a deck of 32 playing cards of rank 7+ (see [Manila]).
	DeckManila = DeckType(Seven)
	// DeckSpanish is a deck of 28 playing cards of rank 8+ (see [Spanish]).
	DeckSpanish = DeckType(Eight)
	// DeckRoyal is a deck of 20 playing cards of rank 10+ (see [Royal]).
	DeckRoyal = DeckType(Ten)
	// DeckKuhn is a deck of 3 playing cards, a [King], [Queen], and a [Jack]
	// (see [Kuhn]).
	DeckKuhn = DeckType(^uint8(0) - 1)
	// DeckLeduc is a deck of 6 playing cards, a [King], [Queen], and a [Jack]
	// of the [Spade] and [Heart] suits (see [Leduc]).
	DeckLeduc = DeckType(^uint8(0) - 2)
)

// Name returns the deck name.
func (typ DeckType) Name() string {
	switch typ {
	case DeckFrench:
		return "French"
	case DeckShort:
		return "Short"
	case DeckManila:
		return "Manila"
	case DeckSpanish:
		return "Spanish"
	case DeckRoyal:
		return "Royal"
	case DeckKuhn:
		return "Kuhn"
	case DeckLeduc:
		return "Leduc"
	}
	return ""
}

// Desc returns the deck description.
func (typ DeckType) Desc(short bool) string {
	switch french := typ == DeckFrench; {
	case french && short:
		return ""
	case french, typ == DeckKuhn, typ == DeckLeduc:
		return typ.Name()
	}
	return typ.Name() + " (" + strconv.Itoa(int(typ+2)) + "+)"
}

// Ordinal returns the deck ordinal.
func (typ DeckType) Ordinal() int {
	return int(typ + 2)
}

// Format satisfies the [fmt.Formatter] interface.
func (typ DeckType) Format(f fmt.State, verb rune) {
	var buf []byte
	switch verb {
	case 'd':
		buf = []byte(strconv.Itoa(int(typ)))
	case 'n':
		buf = []byte(typ.Name())
	case 'o':
		buf = []byte(strconv.Itoa(typ.Ordinal()))
	case 's', 'S':
		buf = []byte(typ.Desc(verb != 's'))
	case 'v':
		buf = []byte("DeckType(" + Rank(typ).Name() + ")")
	default:
		buf = []byte(fmt.Sprintf("%%!%c(ERROR=unknown verb, deck: %d)", verb, int(typ)))
	}
	_, _ = f.Write(buf)
}

// Unshuffled returns a set of the deck's unshuffled cards.
func (typ DeckType) Unshuffled() []Card {
	switch typ {
	case DeckFrench, DeckShort, DeckManila, DeckSpanish, DeckRoyal:
		v := make([]Card, 4*(Ace-Rank(typ)+1))
		var i int
		for _, s := range []Suit{Spade, Heart, Diamond, Club} {
			for r := Rank(typ); r <= Ace; r++ {
				v[i] = New(r, s)
				i++
			}
		}
		return v
	case DeckKuhn:
		return []Card{
			New(King, Spade), New(Queen, Spade), New(Jack, Spade),
		}
	case DeckLeduc:
		return []Card{
			New(King, Spade), New(Queen, Spade), New(Jack, Spade),
			New(King, Heart), New(Queen, Heart), New(Jack, Heart),
		}
	}
	return nil
}

// deck cards.
var (
	deckFrench  []Card
	deckShort   []Card
	deckManila  []Card
	deckSpanish []Card
	deckRoyal   []Card
	deckKuhn    []Card
	deckLeduc   []Card
)

func init() {
	deckFrench = DeckFrench.Unshuffled()
	deckShort = DeckShort.Unshuffled()
	deckManila = DeckManila.Unshuffled()
	deckSpanish = DeckSpanish.Unshuffled()
	deckRoyal = DeckRoyal.Unshuffled()
	deckKuhn = DeckKuhn.Unshuffled()
	deckLeduc = DeckLeduc.Unshuffled()
}

// v returns the cards for the type.
func (typ DeckType) v() []Card {
	switch typ {
	case DeckFrench:
		return deckFrench
	case DeckShort:
		return deckShort
	case DeckManila:
		return deckManila
	case DeckSpanish:
		return deckSpanish
	case DeckRoyal:
		return deckRoyal
	case DeckKuhn:
		return deckKuhn
	case DeckLeduc:
		return deckLeduc
	}
	return nil
}

// Shoe creates a card shoe composed of count number of decks of unshuffled
// cards.
func (typ DeckType) Shoe(count int) *Deck {
	v := typ.v()
	n := len(v)
	d := &Deck{
		V: make([]Card, n*count),
		L: count * n,
	}
	for i := range count {
		copy(d.V[i*n:], v)
	}
	return d
}

// New returns a new deck.
func (typ DeckType) New() *Deck {
	return typ.Shoe(1)
}

// Shuffle returns a new deck, shuffled by the shuffler.
func (typ DeckType) Shuffle(shuffler Shuffler, shuffles int) *Deck {
	d := typ.Shoe(1)
	d.Shuffle(shuffler, shuffles)
	return d
}

// Exclude returns a set of unshuffled cards excluding any supplied cards.
func (typ DeckType) Exclude(ex ...[]Card) []Card {
	return Exclude(typ.v(), ex...)
}

// Deck is a set of playing cards.
type Deck struct {
	I int    `json:"index"`
	L int    `json:"length"`
	V []Card `json:"cards"`
}

// DeckOf creates a deck for the provided cards.
func DeckOf(cards ...Card) *Deck {
	return &Deck{
		V: cards,
		L: len(cards),
	}
}

// NewDeck creates a French deck of 52 unshuffled cards.
func NewDeck() *Deck {
	return DeckFrench.New()
}

// NewShoe creates a card shoe with multiple sets of 52 unshuffled cards.
func NewShoe(count int) *Deck {
	return DeckFrench.Shoe(count)
}

// Limit limits the cards for the deck, for use with card shoes composed of
// more than one deck of cards.
func (d *Deck) Limit(limit int) {
	d.L = limit
}

// Empty returns true when there are no cards remaining in the deck.
func (d *Deck) Empty() bool {
	return d.L <= d.I
}

// Remaining returns the number of remaining cards in the deck.
func (d *Deck) Remaining() int {
	if n := d.L - d.I; 0 <= n {
		return n
	}
	return 0
}

// All returns a copy of all cards in the deck, without advancing.
func (d *Deck) All() []Card {
	v := make([]Card, d.L)
	copy(v, d.V)
	return v
}

// Reset resets the deck.
func (d *Deck) Reset() {
	d.I = 0
}

// Draw draws count cards from the top (front) of the deck.
func (d *Deck) Draw(count int) []Card {
	if count < 0 {
		return nil
	}
	var cards []Card
	for l := min(d.I+count, d.L); d.I < l; d.I++ {
		cards = append(cards, d.V[d.I])
	}
	return cards
}

// Shuffle shuffles the deck's cards using the shuffler.
func (d *Deck) Shuffle(shuffler Shuffler, shuffles int) {
	for range shuffles {
		shuffler.Shuffle(len(d.V), func(i, j int) {
			d.V[i], d.V[j] = d.V[j], d.V[i]
		})
	}
}

// Dealer maintains deal state for a type, streets, deck, positions, runs,
// results, and wins. Use as a street and run iterator for a [Type]. See usage
// details in the [package example].
//
// [package example]: https://pkg.go.dev/github.com/cardrank/cardrank#example-package
type Dealer struct {
	TypeDesc `json:"type"`
	Deck     *Deck        `json:"deck"`
	Count    int          `json:"count"`
	Active   map[int]bool `json:"active"`
	Runs     []*Run       `json:"runs"`
	Results  []*Result    `json:"results"`

	RunCount int `json:"runCount"`
	ST       int `json:"st"`
	S        int `json:"s"`
	R        int `json:"r"`
	E        int `json:"e"`
}

// NewDealer creates a new dealer for a provided deck and pocket count.
func NewDealer(desc TypeDesc, deck *Deck, count int) *Dealer {
	d := &Dealer{
		TypeDesc: desc,
		Deck:     deck,
		Count:    count,
	}
	d.init()
	return d
}

// NewShuffledDealer creates a new deck and dealer, shuffling the deck multiple
// times and returning the dealer with the created deck and pocket count.
func NewShuffledDealer(desc TypeDesc, shuffler Shuffler, shuffles, count int) *Dealer {
	return NewDealer(desc, desc.Deck.Shuffle(shuffler, shuffles), count)
}

// init inits the street position and active positions.
func (d *Dealer) init() {
	d.Active = make(map[int]bool)
	d.Runs = []*Run{NewRun(d.Count)}
	d.Results = nil
	d.RunCount = 1
	d.ST = -1
	d.S = -1
	d.R = -1
	d.E = -1
	for i := range d.Count {
		d.Active[i] = true
	}
}

// Format satisfies the [fmt.Formatter] interface.
func (d *Dealer) Format(f fmt.State, verb rune) {
	var buf []byte
	switch verb {
	case 'n': // name
		buf = []byte(d.Streets[d.S].Name)
	case 's':
		buf = []byte(d.Streets[d.S].Desc())
	default:
		buf = []byte(fmt.Sprintf("%%!%c(ERROR=unknown verb, dealer)", verb))
	}
	_, _ = f.Write(buf)
}

// Inactive returns the inactive positions.
func (d *Dealer) Inactive() []int {
	var v []int
	for i := range d.Count {
		if !d.Active[i] {
			v = append(v, i)
		}
	}
	return v
}

// Deactivate deactivates positions, which will not be dealt further cards and
// will not be included during eval.
func (d *Dealer) Deactivate(positions ...int) bool {
	if d.R != -1 && d.R != 0 {
		return false
	}
	for _, position := range positions {
		delete(d.Active, position)
	}
	return true
}

// Id returns the current street id.
func (d *Dealer) Id() byte {
	if 0 <= d.S && d.S < len(d.Streets) {
		return d.Streets[d.S].Id
	}
	return 0
}

// Name returns the current street name.
func (d *Dealer) Name() string {
	if 0 <= d.S && d.S < len(d.Streets) {
		return d.Streets[d.S].Name
	}
	return ""
}

// NextId returns the next street id.
func (d *Dealer) NextId() byte {
	if -1 <= d.S && d.S < len(d.Streets)-1 {
		return d.Streets[d.S+1].Id
	}
	return 0
}

// HasNext returns true when there is one or more remaining streets.
func (d *Dealer) HasNext() bool {
	n := len(d.Streets)
	return n != 0 && d.S < n-1
}

// HasPocket returns true when one or more pocket cards are dealt for the
// current street.
func (d *Dealer) HasPocket() bool {
	return 0 <= d.S && d.S < len(d.Streets) && 0 < d.Streets[d.S].Pocket
}

// HasBoard returns true when one or more board cards are dealt for the
// current street.
func (d *Dealer) HasBoard() bool {
	return 0 <= d.S && d.S < len(d.Streets) && 0 < d.Streets[d.S].Board
}

// HasActive returns true when there is more than 1 active positions.
func (d *Dealer) HasActive() bool {
	return 0 <= d.S && (d.Type.Max() == 1 || 1 < len(d.Active))
}

// HasCalc returns true when odds are available for calculation.
func (d *Dealer) HasCalc() bool {
	if d.Count != 0 && 0 <= d.R && d.R < d.RunCount && d.Type.Cactus() {
		p, b := d.Type.Pocket(), d.Type.Board()
		if p != 2 && d.S == 0 {
			return false
		}
		return b != 0 && len(d.Runs[d.R].Pockets[0]) >= p
	}
	return false
}

// Pocket returns the number of pocket cards to be dealt on the current street.
func (d *Dealer) Pocket() int {
	if 0 <= d.S && d.S < len(d.Streets) {
		return d.Streets[d.S].Pocket
	}
	return 0
}

// PocketUp returns the number of pocket cards to be turned up on the current
// street.
func (d *Dealer) PocketUp() int {
	if 0 <= d.S && d.S < len(d.Streets) {
		return d.Streets[d.S].PocketUp
	}
	return 0
}

// PocketDiscard returns the number of cards to be discarded prior to dealing
// pockets on the current street.
func (d *Dealer) PocketDiscard() int {
	if 0 <= d.S && d.S < len(d.Streets) {
		return d.Streets[d.S].PocketDiscard
	}
	return 0
}

// PocketDraw returns the number of pocket cards that can be drawn on the
// current street.
func (d *Dealer) PocketDraw() int {
	if 0 <= d.S && d.S < len(d.Streets) {
		return d.Streets[d.S].PocketDraw
	}
	return 0
}

// Board returns the number of board cards to be dealt on the current street.
func (d *Dealer) Board() int {
	if 0 <= d.S && d.S < len(d.Streets) {
		return d.Streets[d.S].Board
	}
	return 0
}

// BoardDiscard returns the number of board cards to be discarded prior to
// dealing a board on the current street.
func (d *Dealer) BoardDiscard() int {
	if 0 <= d.S && d.S < len(d.Streets) {
		return d.Streets[d.S].BoardDiscard
	}
	return 0
}

// Street returns the current street.
func (d *Dealer) Street() int {
	return d.S
}

// Discarded returns the cards discarded on the current street and run.
func (d *Dealer) Discarded() []Card {
	if 0 <= d.S && d.S <= len(d.Streets) && 0 <= d.R && d.R < d.RunCount {
		return d.Runs[d.R].Discard
	}
	return nil
}

// Run returns the current run.
func (d *Dealer) Run() (int, *Run) {
	if 0 <= d.R && d.R < d.RunCount {
		return d.R, d.Runs[d.R]
	}
	return -1, nil
}

// Calc calculates the run odds, including whether or not to include folded
// positions.
func (d *Dealer) Calc(ctx context.Context, folded bool, opts ...CalcOption) (*Odds, *Odds, bool) {
	if 0 <= d.R && d.R < d.RunCount {
		return NewOddsCalc(
			d.Type,
			append(
				opts,
				WithRuns(d.Runs[:d.R+1]),
				WithActive(d.Active, folded),
			)...,
		).Calc(ctx)
	}
	return nil, nil, false
}

// Result returns the current result.
func (d *Dealer) Result() (int, *Result) {
	if 0 <= d.E && d.E < d.RunCount {
		return d.E, d.Results[d.E]
	}
	return -1, nil
}

// Reset resets the dealer and deck.
func (d *Dealer) Reset() {
	d.Deck.Reset()
	d.init()
}

// ChangeRuns changes the number of runs, returning true if successful.
func (d *Dealer) ChangeRuns(runs int) bool {
	switch {
	// check state
	case d.R != 0,
		d.RunCount != 1,
		len(d.Runs) != 1,
		len(d.Streets) <= d.S,
		!d.HasActive():
		return false
	}
	d.Runs = append(d.Runs, make([]*Run, runs-1)...)
	for run := 1; run < runs; run++ {
		d.Runs[run] = d.Runs[0].Dupe()
	}
	d.ST, d.RunCount = d.S, runs
	return true
}

// Next iterates the current street and run, discarding cards prior to dealing
// additional pocket and board cards for each street and run. Returns true when
// there are at least 2 active positions for a [Type] having Max greater than 1
// and when there are additional streets or runs.
func (d *Dealer) Next() bool {
	switch {
	case d.S == -1 && d.R == -1:
		d.S, d.R = 0, 0
	default:
		d.S++
	}
	switch n := len(d.Streets); {
	case n <= d.S && d.R == d.RunCount-1, !d.HasActive():
		return false
	case len(d.Streets) <= d.S && d.R < d.RunCount:
		d.S, d.R = d.ST+1, d.R+1
	}
	d.Deal(d.S, d.Runs[d.R])
	return d.S < len(d.Streets) || d.R < d.RunCount-1
}

// NextResult iterates the next result.
func (d *Dealer) NextResult() bool {
	if d.Results == nil {
		switch n := len(d.Active); {
		case d.Results != nil:
		case n == 1 && d.RunCount == 1 && d.Max != 1:
			// only one active position
			var i int
			for ; i < d.Count && !d.Active[i]; i++ {
			}
			res := &Result{
				Evals:   []*Eval{EvalOf(d.Type)},
				HiOrder: []int{i},
				HiPivot: 1,
			}
			if d.Low || d.Double {
				res.LoOrder, res.LoPivot = res.HiOrder, res.HiPivot
			}
			d.Results = []*Result{res}
		case n > 1 || d.Max == 1:
			d.Results = make([]*Result, d.RunCount)
			for i := range d.RunCount {
				d.Results[i] = NewResult(d.Type, d.Runs[i], d.Active, false)
			}
		}
	}
	if d.RunCount <= d.E {
		return false
	}
	d.E++
	return d.E < d.RunCount
}

// Deal deals pocket and board cards for the street and run, discarding cards
// accordingly.
func (d *Dealer) Deal(street int, run *Run) {
	desc := d.Streets[street]
	// pockets
	if p := desc.Pocket; 0 < p {
		if n := desc.PocketDiscard; 0 < n {
			run.Discard = append(run.Discard, d.Deck.Draw(n)...)
		}
		for range p {
			for i := range d.Count {
				run.Pockets[i] = append(run.Pockets[i], d.Deck.Draw(1)...)
			}
		}
	}
	// board
	if b := desc.Board; 0 < b {
		// hi
		disc := desc.BoardDiscard
		if 0 < disc {
			run.Discard = append(run.Discard, d.Deck.Draw(disc)...)
		}
		run.Hi = append(run.Hi, d.Deck.Draw(b)...)
		// lo
		if d.Double {
			if 0 < disc {
				run.Discard = append(run.Discard, d.Deck.Draw(disc)...)
			}
			run.Lo = append(run.Lo, d.Deck.Draw(b)...)
		}
	}
}

// Run holds pockets, and a Hi/Lo board for a deal.
type Run struct {
	Discard []Card
	Pockets [][]Card
	Hi      []Card
	Lo      []Card
}

// NewRun creates a new run for the pocket count.
func NewRun(count int) *Run {
	return &Run{
		Pockets: make([][]Card, count),
	}
}

// Dupe creates a duplicate of run, with a copy of the pockets and Hi and Lo
// board.
func (run *Run) Dupe() *Run {
	r := new(Run)
	if run.Pockets != nil {
		r.Pockets = make([][]Card, len(run.Pockets))
		for i := range len(run.Pockets) {
			r.Pockets[i] = make([]Card, len(run.Pockets[i]))
			copy(r.Pockets[i], run.Pockets[i])
		}
	}
	if run.Hi != nil {
		r.Hi = make([]Card, len(run.Hi))
		copy(r.Hi, run.Hi)
	}
	if run.Lo != nil {
		r.Lo = make([]Card, len(run.Lo))
		copy(r.Lo, run.Lo)
	}
	return r
}

// Eval returns the evals for the run.
func (run *Run) Eval(typ Type, active map[int]bool, calc bool) []*Eval {
	n := len(run.Pockets)
	evs := make([]*Eval, n)
	var f EvalFunc
	if calc {
		f = calcs[typ]
	} else {
		f = evals[typ]
	}
	for i, double := 0, typ.Double(); i < n; i++ {
		if active == nil || active[i] {
			evs[i] = EvalOf(typ)
			f(evs[i], run.Pockets[i], run.Hi)
			if double {
				ev := EvalOf(typ)
				f(ev, run.Pockets[i], run.Lo)
				evs[i].LoRank, evs[i].LoBest, evs[i].LoUnused = ev.HiRank, ev.HiBest, ev.HiUnused
			}
		}
	}
	return evs
}

// CalcStart returns the run's starting odds.
func (run *Run) CalcStart(low bool) (*Odds, *Odds) {
	count := len(run.Pockets)
	hi := NewOdds(count, nil)
	hi.Total = startingTotal
	var lo *Odds
	if low {
		lo = NewOdds(count, nil)
		lo.Total = startingTotal
	}
	for i, pocket := range run.Pockets {
		expv := StartingExpValue(pocket)
		if expv == nil {
			return nil, nil
		}
		hi.Counts[i] = int(expv.Wins + expv.Losses)
		if low {
			lo.Counts[i] = int(expv.Wins + expv.Losses)
		}
	}
	return hi, lo
}

// Result contains dealer eval results.
type Result struct {
	Evals   []*Eval
	HiOrder []int
	HiPivot int
	LoOrder []int
	LoPivot int
}

// NewResult creates a result for the run, storing the calculated or evaluated
// result.
func NewResult(typ Type, run *Run, active map[int]bool, calc bool) *Result {
	evs := run.Eval(typ, active, calc)
	hiOrder, hiPivot := Order(evs, false)
	var loOrder []int
	var loPivot int
	if typ.Low() || typ.Double() {
		loOrder, loPivot = Order(evs, true)
	}
	return &Result{
		Evals:   evs,
		HiOrder: hiOrder,
		HiPivot: hiPivot,
		LoOrder: loOrder,
		LoPivot: loPivot,
	}
}

// Win returns the Hi and Lo win.
func (res *Result) Win(names ...string) (*Win, *Win) {
	low := res.Evals[res.HiOrder[0]].Type.Low()
	var lo *Win
	if res.LoOrder != nil && res.LoPivot != 0 {
		lo = NewWin(res.Evals, res.LoOrder, res.LoPivot, true, false, names)
	}
	hi := NewWin(res.Evals, res.HiOrder, res.HiPivot, false, low && lo == nil, names)
	return hi, lo
}

// Win formats win information.
type Win struct {
	Evals []*Eval
	Order []int
	Pivot int
	Low   bool
	Scoop bool
	Names []string
}

// NewWin creates a new win.
func NewWin(evs []*Eval, order []int, pivot int, low, scoop bool, names []string) *Win {
	return &Win{
		Evals: evs,
		Order: order,
		Pivot: pivot,
		Low:   low,
		Scoop: scoop,
		Names: names,
	}
}

// Desc returns the eval descriptions.
func (win *Win) Desc() []*EvalDesc {
	var v []*EvalDesc
	for i := range win.Pivot {
		if d := win.Evals[win.Order[i]].Desc(win.Low); d != nil && d.Rank != 0 && d.Rank != Invalid {
			v = append(v, d)
		}
	}
	return v
}

// Invalid returns true when there are no valid winners.
func (win *Win) Invalid() bool {
	switch {
	case win == nil, win.Pivot == 0,
		len(win.Evals) == 0, len(win.Order) == 0:
		return false
	}
	d := win.Evals[win.Order[0]].Desc(win.Low)
	return d == nil || d.Rank == 0 || d.Rank == Invalid
}

// Format satisfies the [fmt.Formatter] interface.
func (win *Win) Format(f fmt.State, verb rune) {
	switch verb {
	case 'd':
		var v []string
		for i := range win.Pivot {
			v = append(v, strconv.Itoa(win.Order[i]))
		}
		fmt.Fprint(f, strings.Join(v, ", ")+" "+win.Verb())
	case 's':
		win.Evals[win.Order[0]].Desc(win.Low).Format(f, 's')
	case 'S':
		if !win.Invalid() {
			var v []string
			for i := range win.Pivot {
				pos := win.Order[i]
				if pos < len(win.Names) {
					v = append(v, win.Names[win.Order[i]])
				} else {
					v = append(v, strconv.Itoa(win.Order[i]))
				}
			}
			fmt.Fprintf(f, "%s %s with %s", strings.Join(v, ", "), win.Verb(), win)
		} else {
			fmt.Fprint(f, "None")
		}
	case 'V':
		fmt.Fprint(f, win.Verb())
	case 'v':
		var v []string
		for i := range win.Pivot {
			desc := win.Evals[win.Order[i]].Desc(win.Low)
			v = append(v, fmt.Sprintf("%v", desc.Best))
		}
		fmt.Fprint(f, strings.Join(v, ", "))
	default:
		fmt.Fprintf(f, "%%!%c(ERROR=unknown verb, win)", verb)
	}
}

// Verb returns the win verb.
func (win *Win) Verb() string {
	switch {
	case win.Scoop:
		return "scoops"
	case win.Pivot > 2:
		return "push"
	case win.Pivot == 2:
		return "split"
	case win.Pivot == 0:
		return "none"
	}
	return "wins"
}
