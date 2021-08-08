package db

import (
    "database/sql/driver"
    "strings"
    "time"
)

type CompetitionType string

const (
    CompetitionOD    CompetitionType = "od"
    CompetitionOS    CompetitionType = "os"
    CompetitionAD    CompetitionType = "ad"
    CompetitionAS    CompetitionType = "as"
    CompetitionND    CompetitionType = "nd"
    CompetitionNS    CompetitionType = "ns"
    CompetitionBD    CompetitionType = "bd"
    CompetitionBS    CompetitionType = "bs"
    CompetitionSPD   CompetitionType = "spd"
    CompetitionSPS   CompetitionType = "sps"
    CompetitionWD    CompetitionType = "wd"
    CompetitionWS    CompetitionType = "ws"
    CompetitionMD    CompetitionType = "md"
    CompetitionProAm CompetitionType = "pro-am"
    CompetitionCOD   CompetitionType = "cod"
    CompetitionDYP   CompetitionType = "dyp"
)

func (p *CompetitionType) Scan(value interface{}) error {
    *p = CompetitionType(value.([]byte))
    return nil
}

func (p CompetitionType) Value() (driver.Value, error) {
    return string(p), nil
}

func (p CompetitionType) String() string {
    return strings.ToUpper(string(p))
}

type Tournament struct {
    Id   uint   `gorm:"primary_key"`
    Name string `gorm:"type:varchar(255);not null"`

    Competitions []*Competition `gorm:"ForeignKey:TournamentId"`
}

type Competition struct {
    Id           uint            `gorm:"primary_key"`
    TournamentId uint            `gorm:"not null"`
    Type         CompetitionType `sql:"type:competition_type"`
    Name         *string         `gorm:"type:varchar(255)"`
    Date         time.Time       `gorm:"not null"`
    Order        uint            `gorm:"not null"`
    Importance   float64         `gorm:"not null"`

    Tournament *Tournament
    Matches    []*Match `gorm:"ForeignKey:CompetitionId"`
    Teams      []*Team  `gorm:"ForeignKey:CompetitionId"`
}

type Match struct {
    Id            uint `gorm:"primary_key"`
    CompetitionId uint `gorm:"not null"`
    Team1Id       uint `gorm:"not null"`
    Team2Id       uint `gorm:"not null"`
    Order         uint `gorm:"not null"`
    Forfeit       bool `gorm:"not null"`
    Team1MaxSets  uint `gorm:`
    Team2MaxSets  uint `gorm:`

    Team1 *Team
    Team2 *Team

    Sets        []Set `gorm:"ForeignKey:MatchId"`
    Competition *Competition
}

type Team struct {
    Id            uint `gorm:"primary_key"`
    CompetitionId uint `gorm:"not null"`
    Player1Id     uint `gorm:"not null"`
    Player2Id     *uint
    Order         uint `gorm:"not null"`
    Position      uint `gorm:"not null"`

    Player1     *Player
    Player2     *Player
    Competition *Competition
}

type Set struct {
    Id         uint `gorm:"primary_key"`
    MatchId    uint `gorm:"not null"`
    Team1Score uint `gorm:"not null"`
    Team2Score uint `gorm:"not null"`
    Order      uint `gorm:"not null"`
}

type Player struct {
    Id uint `gorm:"primary_key"`

    FirstName string `gorm:"type:varchar(255);not null"`
    LastName  string `gorm:"type:varchar(255);not null"`

    ItsfFirstName *string `gorm:"type:varchar(255)"`
    ItsfLastName  *string `gorm:"type:varchar(255)"`
    ItsfLicense   *uint
    ItsfRating    *uint

    EvksInitialRating       uint `gorm:"not null"`
    EvksInitialMatchesCount uint `gorm:"not null"`
    EvksInitialMatchesWin   uint `gorm:"not null"`

    Foreigner bool `gorm:"not null"`
}
