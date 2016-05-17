package bstats

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/boltdb/bolt"
)

const (
	flagIncomplete = iota
	flagComplete
)

type entry struct {
	StartTime int64
	EndTime   int64
	Elapsed   time.Duration
	Status    int32
	Flags     int32
}

func (e *entry) timesAtMidnight() (time.Time, time.Time) {
	return toDay(fromNano(e.StartTime)), toDay(fromNano(e.EndTime))
}

var (
	entriesBucket = []byte("entries")

	timeNano = func() int64 {
		return time.Now().UnixNano()
	}

	errNoLastEntry = errors.New("no last entry")
)

func decodeEntry(data []byte, e *entry) error {
	r := bytes.NewReader(data)
	if err := binary.Read(r, binary.BigEndian, e); err != nil {
		return fmt.Errorf("unable to decode entry. err=%v", err)
	}
	return nil
}

func encodeEntry(e *entry) ([]byte, error) {
	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, e); err != nil {
		return nil, fmt.Errorf("unable to encode entry. err=%v", err)
	}
	return buf.Bytes(), nil
}

func insertEntry(tx *bolt.Tx, e *entry) error {
	b := tx.Bucket(entriesBucket)

	data, err := encodeEntry(e)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%d", e.StartTime)

	b.Put([]byte(key), data)

	return nil
}

func completeLastEntry(tx *bolt.Tx, status int32) error {
	b := tx.Bucket(entriesBucket)
	c := b.Cursor()

	key, val := c.Last()
	if val == nil {
		return errNoLastEntry
	}

	var e entry
	if err := decodeEntry(val, &e); err != nil {
		return err
	}

	e.EndTime = timeNano()
	e.Elapsed = time.Duration(e.EndTime - e.StartTime)
	e.Status = status
	e.Flags |= flagComplete

	data, err := encodeEntry(&e)
	if err != nil {
		return err
	}

	b.Put(key, data)

	return nil
}

func initializeDB(filename string) (*bolt.DB, error) {
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		return nil, err
	}

	return db, db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(entriesBucket)
		return err
	})
}

func Begin(filename string) error {
	db, err := initializeDB(filename)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		return insertEntry(tx, &entry{StartTime: timeNano()})
	})
}

func End(filename string, status int) error {
	db, err := initializeDB(filename)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		return completeLastEntry(tx, int32(status))
	})
}

type Timings struct {
	Slowest        time.Duration
	Fastest        time.Duration
	TotalBuildTime time.Duration

	count int
}

func (t Timings) String() string {
	return fmt.Sprintf("{slowest: %s, fastest: %s, average: %s, total: %s}", t.Slowest, t.Fastest, t.Average(), t.TotalBuildTime)
}

func (t *Timings) Update(elapsed time.Duration) {
	t.count++

	if t.Fastest == 0 || elapsed < t.Fastest {
		t.Fastest = elapsed
	}
	if elapsed > t.Slowest {
		t.Slowest = elapsed
	}

	t.TotalBuildTime += elapsed
}

func (t *Timings) Average() time.Duration {
	if t.count > 0 {
		return time.Duration(int64(t.TotalBuildTime) / int64(t.count))
	}
	return 0
}

type StatsGroup struct {
	Day     time.Time
	Build   int
	Timings Timings
}

func (s StatsGroup) IsValid() bool { return !s.Day.IsZero() }

func (s StatsGroup) String() string {
	return fmt.Sprintf("{day: %s, build: %d, timings: %s}", s.Day, s.Build, s.Timings)
}

func toDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

func isDayAfter(t1, t2 time.Time) bool {
	return toDay(t1).After(toDay(t2))
}

func fromNano(n int64) time.Time {
	return time.Unix(n/1e9, n%1e9)
}

const (
	SuccessfulBuilds  = 1
	FailedBuilds      = 2
	SlowestBuildGraph = 4
	BuildsGraph       = 8
)

func Stats(filename string, w io.Writer, flags int) error {
	db, err := initializeDB(filename)
	if err != nil {
		return err
	}
	defer db.Close()

	var (
		totalComplete   int
		totalIncomplete int
		totalBuildTime  time.Duration

		successful struct {
			sg    StatsGroup
			graph graph
		}
		failed struct {
			sg    StatsGroup
			graph graph
		}
	)

	fn := func(tx *bolt.Tx) error {
		b := tx.Bucket(entriesBucket)
		c := b.Cursor()

		var firstEntry entry
		_, firstV := c.First()
		if err := decodeEntry(firstV, &firstEntry); err != nil {
			return err
		}

		var lastEntry entry
		_, lastV := c.Last()
		if err := decodeEntry(lastV, &lastEntry); err != nil {
			return err
		}

		firstStartTime, _ := firstEntry.timesAtMidnight()
		_, lastEndTime := lastEntry.timesAtMidnight()

		// daySpan := int(time.Duration(lastEntry.EndTime-firstEntry.StartTime) / (24 * time.Hour))
		daySpan := int(lastEndTime.Sub(firstStartTime).Hours() / 24)

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var e entry
			if err := decodeEntry(v, &e); err != nil {
				return err
			}

			elapsed := time.Duration(e.Elapsed)
			totalBuildTime += elapsed

			if e.Flags&flagComplete == flagComplete {
				totalComplete++

				start, end := e.timesAtMidnight()

				elapsedFromStart := end.Sub(firstStartTime)
				day := int(elapsedFromStart.Hours() / 24)
				graphIdx := remap(day, 0, daySpan, 0, maxWidth-1)

				if e.Status == 0 {
					successful.graph[graphIdx].Day = toDay(start)
					successful.graph[graphIdx].Build++
					successful.graph[graphIdx].Timings.Update(elapsed)
				} else {
					failed.graph[graphIdx].Day = toDay(start)
					failed.graph[graphIdx].Build++
					failed.graph[graphIdx].Timings.Update(elapsed)
				}
			} else {
				totalIncomplete++
			}

			if e.Status == 0 {
				successful.sg.Build++
				successful.sg.Timings.Update(elapsed)
			} else {
				failed.sg.Build++
				failed.sg.Timings.Update(elapsed)
			}
		}

		fmt.Fprintln(w, "Overview:")
		fmt.Fprintf(w, "  complete builds: %d\n", totalComplete)
		fmt.Fprintf(w, "  incomplete builds: %d\n", totalIncomplete)
		fmt.Fprintf(w, "  total build time: %s\n\n", totalBuildTime)

		if flags&SuccessfulBuilds == SuccessfulBuilds {
			fmt.Fprintln(w, "Successful builds")
			printSummary(successful.sg, w)
			fmt.Fprintln(w)

			if flags&SlowestBuildGraph == SlowestBuildGraph {
				printSlowestBuildGraph(successful.graph, w)
				fmt.Fprintln(w)
			}
			if flags&BuildsGraph == BuildsGraph {
				printBuildGraph(successful.graph, w)
				fmt.Fprintln(w)
			}
		}

		if flags&FailedBuilds == FailedBuilds {
			fmt.Fprintln(w, "Failed builds")
			printSummary(failed.sg, w)
			fmt.Fprintln(w)

			if flags&SlowestBuildGraph == SlowestBuildGraph {
				printSlowestBuildGraph(failed.graph, w)
				fmt.Fprintln(w)
			}
			if flags&BuildsGraph == BuildsGraph {
				printBuildGraph(failed.graph, w)
				fmt.Fprintln(w)
			}
		}

		return nil
	}

	return db.View(fn)
}

func printSummary(sg StatsGroup, w io.Writer) {
	fmt.Fprintf(w, "timings:\n")
	fmt.Fprintf(w, "  slowest: %s\n", sg.Timings.Slowest)
	fmt.Fprintf(w, "  fastest: %s\n", sg.Timings.Fastest)
	fmt.Fprintf(w, "  average: %s\n", sg.Timings.Average())
	fmt.Fprintf(w, "  total build time: %s\n", sg.Timings.TotalBuildTime)
	fmt.Fprintf(w, "build: %d\n", sg.Build)
}

type sgValFunc func(sg *StatsGroup) int

func printSlowestBuildGraph(g graph, w io.Writer) {
	var allSlowest time.Duration
	for _, sg := range g {
		if allSlowest == 0 || allSlowest < sg.Timings.Slowest {
			allSlowest = sg.Timings.Slowest
		}
	}

	fmt.Fprintf(w, "Slowest build graph over days\n\n")
	printGraph(g, int(allSlowest), func(sg *StatsGroup) int {
		return int(sg.Timings.Slowest)
	}, w)
}

func printBuildGraph(g graph, w io.Writer) {
	var maxBuilds int
	for _, sg := range g {
		if sg.Build > maxBuilds {
			maxBuilds = sg.Build
		}
	}

	fmt.Fprintf(w, "Builds graph over days\n\n")
	printGraph(g, maxBuilds, func(sg *StatsGroup) int {
		return int(sg.Build)
	}, w)
}

func printGraph(g graph, maxY int, fn sgValFunc, w io.Writer) {
	var maxX int
	for y := maxHeight; y > 0; y-- {
		fmt.Fprint(w, "| ")

		for x := 0; x < maxWidth; x++ {
			sg := g[x]
			if !sg.IsValid() {
				continue
			}
			if x > maxX {
				maxX = x
			}

			val := fn(&sg)
			mapped := remap(val, 0, maxY, 1, maxHeight)

			if mapped == y {
				fmt.Fprint(w, "*")
			} else {
				fmt.Fprint(w, " ")
			}
		}
		fmt.Fprintln(w)
	}
	fmt.Fprint(w, "+-")
	for x := 0; x < maxX+1; x++ {
		fmt.Fprint(w, "-")
	}
	fmt.Fprintln(w)
}

const (
	maxHeight = 10
	maxWidth  = 80
)

type graph [maxWidth]StatsGroup

func (g graph) LastValid() StatsGroup {
	var res StatsGroup
	for _, sg := range g {
		if sg.IsValid() {
			res = sg
		}
	}
	return res
}

func remap(oldV, oldMin, oldMax, newMin, newMax int) int {
	oldRange := float64(oldMax) - float64(oldMin)
	if oldRange == 0.0 {
		return newMin
	}
	newRange := float64(newMax) - float64(newMin)

	d := float64(oldV - oldMin)
	v := d * newRange
	v2 := v / oldRange

	return int(v2 + float64(newMin))
}
