	"gorm.io/gorm"
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	b.Id = uuid.New()
	return nil
}
