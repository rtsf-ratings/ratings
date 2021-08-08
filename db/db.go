package db

import (
    "gorm.io/gorm"
)

type Ratings struct {
    Conn *gorm.DB

    Competitions []*Competition
    Players      []*Player
}

func RatingsConnect(dialector gorm.Dialector, config *gorm.Config) (*gorm.DB, error) {
    conn, err := gorm.Open(dialector, config)

    if err != nil {
        return nil, err
    }

    conn.AutoMigrate(&Player{})
    conn.AutoMigrate(&Tournament{}, &Competition{}, &Match{}, &Team{}, &Set{})

    return conn, nil
}

func RatingsOpen(dialector gorm.Dialector, config *gorm.Config) (*Ratings, error) {
    conn, err := RatingsConnect(dialector, config)
    if err != nil {
        return nil, err
    }

    ratingsDb := &Ratings{
        Conn: conn,
    }

    err = ratingsDb.load()
    if err != nil {
        return nil, err
    }

    return ratingsDb, nil
}

func (ratings *Ratings) load() error {
    r := ratings.Conn.Order("competitions.date ASC").Find(&ratings.Competitions)
    if r.Error != nil {
        return r.Error
    }

    var tounaments []*Tournament
    r = ratings.Conn.Find(&tounaments)
    if r.Error != nil {
        return r.Error
    }

    var matches []*Match
    r = ratings.Conn.Preload("Sets").Order("matches.order ASC").Find(&matches)
    if r.Error != nil {
        return r.Error
    }

    var teams []*Team
    r = ratings.Conn.Order("teams.order ASC").Find(&teams)
    if r.Error != nil {
        return r.Error
    }

    var players []*Player
    r = ratings.Conn.Find(&players)
    if r.Error != nil {
        return r.Error
    }

    ratings.Players = players

    for _, tournament := range tounaments {
        for _, competition := range ratings.Competitions {
            if competition.TournamentId == tournament.Id {
                tournament.Competitions = append(tournament.Competitions, competition)
                competition.Tournament = tournament
            }
        }
    }

    for _, competition := range ratings.Competitions {
        for _, match := range matches {
            if match.CompetitionId == competition.Id {
                competition.Matches = append(competition.Matches, match)
                match.Competition = competition
            }
        }

        for _, team := range teams {
            if team.CompetitionId == competition.Id {
                competition.Teams = append(competition.Teams, team)
                team.Competition = competition
            }
        }
    }

    for _, match := range matches {
        for _, team := range teams {
            if match.Team1Id == team.Id {
                match.Team1 = team
            }

            if match.Team2Id == team.Id {
                match.Team2 = team
            }
        }
    }

    for _, team := range teams {
        for _, player := range players {
            if team.Player1Id == player.Id {
                team.Player1 = player
            }

            if team.Player2Id != nil && *team.Player2Id == player.Id {
                team.Player2 = player
            }
        }
    }

    return r.Error
}
