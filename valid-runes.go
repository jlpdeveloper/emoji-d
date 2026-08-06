package emojid

// validRunes is a curated list of emojis that are allowed for use in IDs.
// Source: https://unicode.org/emoji/charts/full-emoji-list.html
var validRunes = map[rune]struct{}{
	'😀': {},
	'🚀': {},
	'🔥': {},
	'🍕': {},
	'😅': {},
	'😐': {},
	'🙄': {},
	'🤮': {},
	'🤠': {},
	'🌮': {},
	//10
	'😎': {},
	'💩': {},
	'👻': {},
	'😺': {},
	'🍝': {},
	'💯': {},
	'🖕': {},
	'👍': {},
	'💪': {},
	'👀': {},
	//20
	'🦆': {},
	'🦉': {},
	'🐔': {},
	'🐸': {},
	'🐢': {},
	'🐌': {},
	'🌭': {},
	'🌲': {},
	'🍄': {},
	'🍌': {},
	//30
	'🍆': {},
	'🧀': {},
}
