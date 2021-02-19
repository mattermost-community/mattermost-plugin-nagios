package main

const (
	criticalState    = "critical"
	downState        = "down"
	okState          = "ok"
	pendingState     = "pending"
	unknownState     = "unknown"
	unreachableState = "unreachable"
	upState          = "up"
	warningState     = "warning"
)

const (
	bangBangEmoji             = ":bangbang:"
	barChartEmoji             = ":bar_chart:"
	smallRedTriangleDownEmoji = ":small_red_triangle_down:"
	whiteCheckMarkEmoji       = ":white_check_mark:"
	hourglassFlowingSandEmoji = ":hourglass_flowing_sand:"
	questionEmoji             = ":question:"
	mailboxWithNoMailEmoji    = ":mailbox_with_no_mail:"
	upEmoji                   = ":up:"
	warningEmoji              = ":warning:"
)

var emojiLookup = map[string]string{
	criticalState:    bangBangEmoji,
	downState:        smallRedTriangleDownEmoji,
	okState:          whiteCheckMarkEmoji,
	pendingState:     hourglassFlowingSandEmoji,
	unknownState:     questionEmoji,
	unreachableState: mailboxWithNoMailEmoji,
	upState:          upEmoji,
	warningState:     warningEmoji,
}

func emoji(state string) string {
	e, ok := emojiLookup[state]
	if !ok {
		return questionEmoji
	}

	return e
}
