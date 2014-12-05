package deliver

import (
	util "github.com/xiangshouding/go-deliver/util"
	"log"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func NewDeliver(from, to string) *Deliver {
	return &Deliver{
		from:   from,
		to:     to,
		maps:   make([]Roadmap, 0),
		silent: true,
	}
}

type Roadmap struct {
	reg     string
	_reg    *regexp.Regexp
	release string
}

func (r *Roadmap) Parse() {
	if r._reg == nil {
		r._reg = regexp.MustCompile(r.reg)
	}
}

func (r *Roadmap) Fill(s string) string {
	r.Parse()
	m := r._reg.FindStringSubmatch(s)
	if len(m) == 0 {
		return ""
	}
	str := r.release
	for i, v := range m {
		str = strings.Replace(str, "$"+strconv.Itoa(i), v, -1)
	}
	str = strings.Replace(str, "$&", m[0], -1)
	return str
}

type Deliver struct {
	maps   []Roadmap
	from   string
	to     string
	silent bool
}

func (d *Deliver) ShowLog() {
	d.silent = false
}

func (d *Deliver) Push(r map[string]string) {
	var reg, release string
	var ok bool
	reg, ok = r["reg"]
	if !ok {
		panic("A roadmap rule must given key `reg`")
	}
	release, ok = r["release"]
	if !ok {
		panic("A roadmap rule must given key `release`")
	}
	roadmap := Roadmap{
		reg:     reg,
		release: release,
	}
	roadmap.Parse()
	d.maps = append(d.maps, roadmap)
}

func (d *Deliver) release(srcSubpath string, r Roadmap) bool {
	dstSubpath := r.Fill(srcSubpath)
	if dstSubpath == "" {
		return false
	}
	from := path.Join(d.from, srcSubpath)
	to := path.Join(d.to, dstSubpath)
	util.CopyFile(from, to, false)
	if !d.silent {
		log.Printf("COPY: %s TO %s\n", from, to)
	}
	return true
}

func (d *Deliver) Release(r ...map[string]string) {
	for _, v := range r {
		d.Push(v)
	}
	files := util.Find(d.from)
	for _, file := range files {
		subpath := strings.Replace(file, d.from, "", -1)
		for _, rule := range d.maps {
			if d.release(subpath, rule) {
				break
			}
		}
	}
}
