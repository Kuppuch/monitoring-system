package middleware

import "gorm.io/gorm"

type Tracker struct {
	gorm.Model
	Name        string
	Description string
	Issues      []Issue `gorm:"foreignKey:TrackerID" json:"issues"`
}

func GetTrackerList() []Tracker {
	var trackers []Tracker
	DB.Find(&trackers)
	return trackers
}

func (t *Tracker) GetTracker(id int) error {
	tx := DB.Where("id = ?", id).Find(&t)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (t *Tracker) InsertTracker() {

}
