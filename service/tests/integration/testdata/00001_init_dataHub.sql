-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE Account (
    AccountID int  NOT NULL,
    AccountType nvarchar(50),
    CONSTRAINT PK_Account PRIMARY KEY (AccountID)
);

INSERT INTO Account (AccountID ,  AccountType)
VALUES (1, 'GOLDENAYN');

Create TABLE Portfolio (
    PortfolioId int  NOT NULL,
    AccountID int  NOT NULL,
    PortfolioRef nvarchar(50)  NOT NULL,
    CONSTRAINT PK_Portfolio PRIMARY KEY (PortfolioId),
    CONSTRAINT FK_Portfolio_Account FOREIGN KEY (AccountID) REFERENCES Account(AccountID)
);

INSERT INTO Portfolio (PortfolioId ,AccountID, PortfolioRef)
VALUES (1,1, 'ORE1246.001');

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE accounts;

