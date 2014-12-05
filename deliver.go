package deliver

import (
	"log"
	"strings"
	"regexp"
	"strconv"
	"io/ioutil"
	"path"
	util "./util"
)

func NewDeliver(from, to string) *Deliver {
	return &Deliver {
		from: from,
		to: to,
		maps: make([]Roadmap,0),
	}
}

type Roadmap struct {
	reg string
	_reg *regexp.Regexp
	release string
}

func (r *Roadmap) Parse() {
	if r._reg == nil {
		r._reg = regexp.MustCompile(r.reg)
	}
}

func (r *Roadmap) Fill (s string) string {
	r.Parse();
	m := r._reg.FindStringSubmatch(s)
	if len(m) == 0 {
		return ""
	}
	str := r.release
	for i,v := range m {
		str = strings.Replace(str, "$"+strconv.Itoa(i), v, -1);
	}
	str = strings.Replace(str, "$&", m[0], -1);
	return str
}


type Deliver struct {
	maps []Roadmap
	from string
	to string
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
		reg: reg,
		release: release,
	}
	roadmap.Parse();
	d.maps = append(d.maps, roadmap)
}

func (d *Deliver) release(srcSubpath string, r Roadmap) bool {
	dstSubpath := r.Fill(srcSubpath)
	if dstSubpath == "" {
		return false
	}
	util.CopyFile(path.Join(d.from, srcSubpath), path.Join(d.to, dstSubpath), false);
	return true
}

func (d *Deliver) Release(r ...map[string]string) {
	for _, v := range r {
		log.Println(v)
		d.Push(v)
	}
	files, err := ioutil.ReadDir(d.from)
	if err != nil {
		panic(err)
	}
	var filepath_ string
	for _, file := range files {
		filepath_ = file.Name() //path.Join(d.from, file.Name())
		log.Println(filepath_)
		for _, rule := range d.maps {
			if d.release(filepath_, rule) {
				break
			}
		}
	}
}
