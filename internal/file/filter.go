package file

import "regexp"

// NewFilter returns a pointer to a new,
// initialised filter ready for use
func NewFilter() *Filter {

	f := &Filter{}
	f.init()
	return f
}

func (f *Filter) init() {
	ap := make(map[string]regexp.Regexp)
	dp := make(map[string]regexp.Regexp)
	f.AcceptPatterns = &ap
	f.DenyPatterns = &dp
}

// Reset replaces both AcceptPatterns and DenyPatterns
// with empty initialised maps, ready for use
func (f *Filter) Reset() {
	f.init()
}

// Pass returns whether or not a line should
// be passed by this filter
// which is true if the Filter
//

// Pass returns a bool indicating whether
// a line passes (true) or is blocked (false)
// by the filter
func (f *Filter) Pass(line string) bool {

	if f.AllPass() {
		return true
	}

	if f.Deny(line) {
		return false
	}

	if f.Accept(line) {
		return true
	}

	return false

}

// AllPass returns true if both AcceptPatterns and DenyPatterns
// are empty, i.e. all messages should pass.
// we do this for convenience and efficiency, rather than
// having an explict 'all pass' filter added to the AcceptList
// because we'd have to remove it the first time we add a filter
// and the second time we add a filter we'd have to check whether
// the first filter was the allpass one, and we might not know
// whether that was from initialisation or explicitly added by
// a user ....
func (f *Filter) AllPass() bool {
	return len(*f.AcceptPatterns) == 0 && len(*f.DenyPatterns) == 0
}

// match checks whether a string matches any patterns in the list of patterns
func match(line string, patterns *map[string]regexp.Regexp) bool {
	for _, p := range *patterns {
		if p.MatchString(line) {
			return true
		}
	}
	return false
}

// Deny returns true if this line is blocked by the filter
func (f *Filter) Deny(line string) bool {
	return match(line, f.DenyPatterns)
}

// Accept returns true if this line is passed by the filter
func (f *Filter) Accept(line string) bool {
	return match(line, f.AcceptPatterns)
}

// AddAcceptPattern adds a pattern to the AcceptPatterns
// that will be used to check if a message is accepted (passed)
func (f *Filter) AddAcceptPattern(p *regexp.Regexp) {
	(*f.AcceptPatterns)[p.String()] = *p
}

// AddDenyPattern adds a pattern to the DenyPatterns
// that will be used to check if a message is denied (blocked)
func (f *Filter) AddDenyPattern(p *regexp.Regexp) {
	(*f.DenyPatterns)[p.String()] = *p
}

// DeleteAcceptPattern will remove a given pattern from the
// list of patterns used to check for acceptance of a line
func (f *Filter) DeleteAcceptPattern(p *regexp.Regexp) {
	delete(*f.AcceptPatterns, p.String())
}

// DeleteDenyPattern will remove a given pattern from the
// list of patterns used to check for denial of a line
func (f *Filter) DeleteDenyPattern(p *regexp.Regexp) {
	delete(*f.DenyPatterns, p.String())
}
