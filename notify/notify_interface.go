package notify

type NotifyType struct {
	NotifyName       string
	NotifyCode       string
	NotifyAppletCode string
}

type Notify interface {
	Notify(notifyType NotifyType, msg string)

	NotifyError(notifyType NotifyType, msg string)
}
