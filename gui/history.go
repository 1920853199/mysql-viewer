package gui

type History struct {
	list []string
	index int
}

func (h *History) Add(str string)  {
	h.list = append(h.list, str)
	h.index = len(h.list) - 1
}
 
func (h *History) Prev() string {
	if len(h.list) == 0 {
		return ""
	}

	//h.index = (h.index - 1 + len(h.list)) % len(h.list)
	cmd := h.list[h.index]
	if h.index > 0 {
		h.index = h.index - 1
	}
	return cmd
}

func (h *History) Next() string {
	if len(h.list) == 0 {
		return ""
	}
	if h.index >= len(h.list) - 1 {
		return ""
	}
	//h.index = (h.index + 1) % len(h.list)
	h.index = h.index + 1
	cmd := h.list[h.index]
	return cmd
}