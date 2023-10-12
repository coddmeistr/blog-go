package types

type Article struct {
	ID        uint   `gorm:"primaryKey"`
	Header    string `json:"header"`
	Topic     string `json:"topic"`
	ShortText string `json:"short"`
	LongText  string `json:"long"`

	AuthorID uint      `gorm:"not null" json:"authorid"`
	Comment  []Comment `gorm:"constraint:OnDelete:CASCADE;foreignKey:ArticleID" json:"-"`
	Like     []Like    `gorm:"constraint:OnDelete:CASCADE;foreignKey:ArticleID" json:"-"`
}

type Comment struct {
	ID      uint   `gorm:"primaryKey"`
	RawText string `json:"text"`

	AuthorID       uint     `gorm:"not null" json:"authorid"`
	ArticleID      uint     `gorm:"not null" json:"articleid"`
	ReplyCommentID *uint    `gorm:"default:null" json:"replyid"`
	ReplyComment   *Comment `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	Like           []Like   `gorm:"constraint:OnDelete:CASCADE;foreignKey:CommentID" json:"-"`
}

type Like struct {
	ID uint `gorm:"primaryKey"`

	UserID    uint `gorm:"not null"`
	CommentID uint `gorm:"default:null"`
	ArticleID uint `gorm:"default:null"`
}
