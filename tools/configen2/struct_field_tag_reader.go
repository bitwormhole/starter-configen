package configen2

import "strings"

type FieldTagReader struct {
	text string
	size int
	ptr  int
}

func (inst *FieldTagReader) init(str string) {
	const token = "`"
	str = strings.TrimSpace(str)
	if strings.HasPrefix(str, token) && strings.HasSuffix(str, token) {
		from := 1
		to := len(str) - 1
		if from < to {
			str = str[from:to]
		}
	}
	inst.size = len(str)
	inst.text = str
	inst.ptr = 0
}

// @return( key, value, count) , 如果到达结尾，count=-1
func (inst *FieldTagReader) read() (string, string, int) {

	const m1 = ':'
	const m2 = '"'
	const m3 = '\\'

	ptr := inst.ptr
	end := inst.size
	text := inst.text

	m1count := 0
	m2count := 0
	i0 := ptr
	i1 := 0 // to m1
	i2 := 0 // to m2.1
	i3 := 0 // to m2.2
	skip := false

	for ; ptr < end; ptr++ {
		if skip {
			skip = false
			continue
		}
		ch := rune(text[ptr])
		if ch == m1 {
			if m1count == 0 {
				m1count++
				i1 = ptr
			}
		} else if ch == m2 {
			if m2count == 0 {
				m2count++
				i2 = ptr
			} else if m2count == 1 {
				m2count++
				i3 = ptr
				break
			}
		} else if ch == m3 {
			skip = true
		}
	}

	if i0 < i1 && i1 < i2 && i2 < i3 {
		key := strings.TrimSpace(text[i0:i1])
		val := strings.TrimSpace(text[i2+1 : i3])
		inst.ptr = i3 + 1
		return key, val, 1
	}

	return "", "", -1
}
