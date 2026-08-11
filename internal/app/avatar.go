package app

// avatarView is the server-driven view model for the Avatar primitive. The
// avatar is a circular surface (surface-container + fg) that shows initials or
// an image at sm/md/lg size, backed by the scoped --ui-avatar-* tokens.
//
// Accessibility contract: an avatar that is decorative — paired with visible or
// accessible name text — must set Decorative=true so the whole element renders
// aria-hidden and the surrounding text supplies the meaning (never color-only,
// never initials-only). An image avatar that carries meaning on its own keeps
// Decorative=false and sets ImageAlt so the alt text is announced. Initials
// avatars are decorative by nature; recipes pair them with a name.
type avatarView struct {
	Initials   string // initials shown when ImageSrc is empty
	ImageSrc   string // optional image URL; when set it wins over Initials
	ImageAlt   string // alt text for the image when it carries meaning
	Decorative bool   // true → the whole avatar is aria-hidden
	Size       string // closed: "" (md) | sm | md | lg
}

// avatarSize sanitizes the size modifier against the closed sm/md/lg set.
func avatarSize(size string) string {
	switch size {
	case "sm", "md", "lg":
		return size
	}
	return "md"
}
