-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS MESSAGE_I1 ON MESSAGE (CHANNEL_ID, CHUNK_ID, IS_PARENT);
CREATE INDEX IF NOT EXISTS MESSAGE_I2 ON MESSAGE (CHANNEL_ID, PARENT_ID);
CREATE INDEX IF NOT EXISTS MESSAGE_I3 ON MESSAGE (IS_PARENT, JSON_EXTRACT(DATA, '$.latest_reply'));

-- V_CHANNEL_THREADS CONTAINS A THREAD COUNT PER EACH CHANNEL.
CREATE VIEW IF NOT EXISTS V_CHANNEL_THREADS AS
WITH MSG AS (SELECT ID, CHANNEL_ID, IS_PARENT, CHUNK_ID
             FROM MESSAGE
             WHERE IS_PARENT = TRUE
             EXCEPT
             SELECT ID, CHANNEL_ID, IS_PARENT, CHUNK_ID
             FROM MESSAGE
             WHERE JSON_EXTRACT(DATA, '$.latest_reply') = '0000000000.000000' -- EMPTY THREADS
               AND IS_PARENT = TRUE)
SELECT C.SESSION_ID, M.CHANNEL_ID, SUM(M.IS_PARENT) THREADS
FROM MSG M
         JOIN CHUNK C ON C.ID = M.CHUNK_ID
WHERE IS_PARENT = TRUE
GROUP BY M.CHANNEL_ID;

-- V_THREAD_COUNT HAS THE COUNT OF ACTUALLY DOWNLOADED THREADS.
CREATE VIEW IF NOT EXISTS V_THREAD_COUNT AS
SELECT C.SESSION_ID, M.CHANNEL_ID, COUNT(DISTINCT M.PARENT_ID) PARENT_COUNT
FROM CHUNK C
         JOIN MESSAGE M ON C.ID = M.CHUNK_ID
WHERE C.TYPE_ID = 1
GROUP BY C.SESSION_ID, M.CHANNEL_ID;

-- V_UNFINISHED_THREADS IS THE COUNT OF UNFINISHED THREADS FOR EACH CHANNEL.
-- SUCCESSFUL EXPORT MUST HAVE REF_COUNT = 0 FOR ALL CHANNELS.
CREATE VIEW IF NOT EXISTS V_UNFINISHED_THREADS AS
SELECT CT.SESSION_ID, CT.CHANNEL_ID, CT.THREADS - TC.PARENT_COUNT REF_COUNT
FROM V_CHANNEL_THREADS AS CT
         JOIN V_THREAD_COUNT TC
              ON (CT.CHANNEL_ID = TC.CHANNEL_ID
                  AND CT.SESSION_ID = TC.SESSION_ID);

-- ORPHAN THREADS FINDS ALL THE THREADS WHICH HAVE A PARENT, BUT NOT A CHILD
-- (THREAD_MESSAGES HAS NOT (YET) FETCHED THE THREAD.  IN THE NORMAL CONDITIONS
-- ON A COMPLETED DUMP, ALL COUNTS FOR A SESSION SHOULD BE 0.
CREATE VIEW IF NOT EXISTS V_ORPHAN_THREADS AS
WITH UNFINISHED_THREADS AS (SELECT CHANNEL_ID
                            FROM V_UNFINISHED_THREADS
                            WHERE REF_COUNT > 0),
     DIFF AS (SELECT M.CHANNEL_ID, M.THREAD_TS
              FROM MESSAGE M
                       JOIN UNFINISHED_THREADS P ON P.CHANNEL_ID = M.CHANNEL_ID
              WHERE IS_PARENT = TRUE
              EXCEPT
              SELECT DISTINCT M.CHANNEL_ID, M.THREAD_TS
              FROM MESSAGE M
                       JOIN UNFINISHED_THREADS P ON P.CHANNEL_ID = M.CHANNEL_ID
                       JOIN CHUNK C ON C.ID = M.CHUNK_ID
              WHERE C.TYPE_ID = 1
              GROUP BY M.CHANNEL_ID, M.THREAD_TS)
SELECT *
FROM DIFF;

-- V_EMPTY_THREADS SHOWS ALL EMPTY THREADS (THREADS THAT HAVE PARENT MESSAGE,
-- BUT NO MESSAGES IN THEM).
CREATE VIEW V_EMPTY_THREADS AS
SELECT SESSION_ID,
       CHUNK_ID,
       CHANNEL_ID,
       THREAD_TS,
       CONCAT('archives/', CHANNEL_ID, '/p', cast(PARENT_ID AS TEXT)) as PATH
FROM MESSAGE M
         JOIN CHUNK C ON M.CHUNK_ID = C.ID
WHERE IS_PARENT = TRUE
  AND JSON_EXTRACT(DATA, '$.latest_reply') = '0000000000.000000';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS V_EMPTY_THREADS;
DROP VIEW IF EXISTS V_ORPHAN_THREADS;
DROP VIEW IF EXISTS V_UNFINISHED_THREADS;
DROP VIEW IF EXISTS V_THREAD_COUNT;
DROP VIEW IF EXISTS V_CHANNEL_THREADS;
DROP INDEX IF EXISTS MESSAGE_I3;
DROP INDEX IF EXISTS MESSAGE_I2;
DROP INDEX IF EXISTS MESSAGE_I1;
-- +goose StatementEnd
