-- Create new team
INSERT INTO "team" ("id", "name", "is_deleted", "created_at", "updated_at") VALUES ('933efe12-2219-42df-bd51-a2e84888432d', 'team1', DEFAULT, DEFAULT, DEFAULT);;

-- Create new users
INSERT INTO "user" ("id", "email", "name", "is_deleted", "created_at", "updated_at") VALUES ('8e159833-5078-4b0a-80a0-363d82bafd60', 'email1@gmail.com', 'name1', DEFAULT, DEFAULT, DEFAULT);
INSERT INTO "user" ("id", "email", "name", "is_deleted", "created_at", "updated_at") VALUES ('51241303-ebe0-4c2b-98be-44d93439f6d9', 'email2@gmail.com', 'name2', DEFAULT, DEFAULT, DEFAULT);
INSERT INTO "user" ("id", "email", "name", "is_deleted", "created_at", "updated_at") VALUES ('e0294804-ac72-4f2f-be54-2fb18a135709', 'email3@gmail.com', 'name3', DEFAULT, DEFAULT, DEFAULT);
INSERT INTO "user" ("id", "email", "name", "is_deleted", "created_at", "updated_at") VALUES ('0e49e11c-660c-43c5-954e-ef9e89b45833', 'email4@gmail.com', 'name4', DEFAULT, DEFAULT, DEFAULT);
INSERT INTO "user" ("id", "email", "name", "is_deleted", "created_at", "updated_at") VALUES ('bafdd20b-a2a5-41ca-b4a2-26fccbd029dd', 'email5@gmail.com', 'name5', DEFAULT, DEFAULT, DEFAULT);

-- Assign 3 users as team members
INSERT INTO "team_member" ("id", "team_id", "user_id", "is_deleted", "created_at", "updated_at") VALUES ('852da87d-2d16-4173-885a-e84476f2d0ba', '933efe12-2219-42df-bd51-a2e84888432d', '8e159833-5078-4b0a-80a0-363d82bafd60', DEFAULT, DEFAULT, DEFAULT);
INSERT INTO "team_member" ("id", "team_id", "user_id", "is_deleted", "created_at", "updated_at") VALUES ('8f0f3128-a777-480c-a763-53135548b573', '933efe12-2219-42df-bd51-a2e84888432d', '51241303-ebe0-4c2b-98be-44d93439f6d9', DEFAULT, DEFAULT, DEFAULT);
INSERT INTO "team_member" ("id", "team_id", "user_id", "is_deleted", "created_at", "updated_at") VALUES ('0ea664cf-5d20-4360-8688-1c3ab4ec065d', '933efe12-2219-42df-bd51-a2e84888432d', 'e0294804-ac72-4f2f-be54-2fb18a135709', DEFAULT, DEFAULT, DEFAULT);

-- Create a new wallet for team
INSERT INTO "wallet" ("id", "balance", "daily_limit", "monthly_limit", "team_id", "user_id", "is_deleted", "created_at", "updated_at") VALUES ('d4a6607a-1af7-4571-bdff-2672be72ba0e', DEFAULT, 500000, 500000, '933efe12-2219-42df-bd51-a2e84888432d', DEFAULT, DEFAULT, DEFAULT, DEFAULT);

-- Create a new card for team's wallet
INSERT INTO "card" ("id", "card_no", "expiry_month", "expiry_year", "cvv", "daily_limit", "monthly_limit", "wallet_id", "is_deleted", "created_at", "updated_at") VALUES ('cb21fe95-37e7-4e67-aac3-1b633fe1036d', '5200828282828210', '10', '2021', '123', 500000, 500000, 'd4a6607a-1af7-4571-bdff-2672be72ba0e', DEFAULT, DEFAULT, DEFAULT);

-- Create new wallets for 2 users
INSERT INTO "wallet" ("id", "balance", "daily_limit", "monthly_limit", "team_id", "user_id", "is_deleted", "created_at", "updated_at") VALUES ('370a9739-b90b-4264-81a2-f8d0d3236011', DEFAULT, DEFAULT, DEFAULT, DEFAULT, '0e49e11c-660c-43c5-954e-ef9e89b45833', DEFAULT, DEFAULT, DEFAULT);
INSERT INTO "wallet" ("id", "balance", "daily_limit", "monthly_limit", "team_id", "user_id", "is_deleted", "created_at", "updated_at") VALUES ('0ac7623a-f857-4e7a-856e-af42e8873651', DEFAULT, DEFAULT, DEFAULT, DEFAULT, 'bafdd20b-a2a5-41ca-b4a2-26fccbd029dd', DEFAULT, DEFAULT, DEFAULT);

-- Create new wallets for users' wallets
INSERT INTO "card" ("id", "card_no", "expiry_month", "expiry_year", "cvv", "daily_limit", "monthly_limit", "wallet_id", "is_deleted", "created_at", "updated_at") VALUES ('f7bee42b-c2b6-4d9b-b8bb-c9dff69240a4', '4242424242424242', '11', '2022', '234', 500000, 500000, '370a9739-b90b-4264-81a2-f8d0d3236011', DEFAULT, DEFAULT, DEFAULT);
INSERT INTO "card" ("id", "card_no", "expiry_month", "expiry_year", "cvv", "daily_limit", "monthly_limit", "wallet_id", "is_deleted", "created_at", "updated_at") VALUES ('2234abcc-d2f5-493a-b753-3b40c1f40ec4', '4000056655665556', '12', '2023', '345', 500000, 500000, '0ac7623a-f857-4e7a-856e-af42e8873651', DEFAULT, DEFAULT, DEFAULT);