package models

type Post struct {
	ID            int    `json:"id"`
	VerseID       string `json:"verse"`
	Published     bool   `json:"published"`
	State         string `json:"state"`
	PublishmentID string `json:"publishmentid"`
	PostURL       string `json:"posturl"`
	ReelURL       string `json:"reelurl"`
}

// CREATE Table Post (
//  ID int,
//  VerseID varchar(255),
//  Published BOOLEAN,
//  State varchar(255),
//  PublishmentID varchar(255),
//  PostURL varchar(255),
//  ReelURL varchar(255),
//  PRIMARY KEY (ID)
// );
	
