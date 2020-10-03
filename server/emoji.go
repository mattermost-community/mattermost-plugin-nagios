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
	smallRedTriangleDownEmoji = ":small_red_triangle_down:"
	whiteCheckMarkEmoji       = ":white_check_mark:"
	hourglassFlowingSandEmoji = ":hourglass_flowing_sand:"
	questionEmoji             = ":question:"
	mailboxWithNoMailEmoji    = ":mailbox_with_no_mail:"
	upEmoji                   = ":up:"
	warningEmoji              = ":warning:"
)

func emoji(state string) string {
	switch state {
	case criticalState:
		return bangBangEmoji
	case downState:
		return smallRedTriangleDownEmoji
	case okState:
		return whiteCheckMarkEmoji
	case pendingState:
		return hourglassFlowingSandEmoji
	case unknownState:
		return questionEmoji
	case unreachableState:
		return mailboxWithNoMailEmoji
	case upState:
		return upEmoji
	case warningState:
		return warningEmoji
	default:
		return questionEmoji
	}
}
