
CREATE TABLE WOWHEAD_ITEMLIST
(
   JSONCOL VARCHAR(MAX)
  ,ITEMID INT
  ,DATE_LOAD DATETIME DEFAULT GETDATE()
)




DECLARE @UID INT
SET @UID = 187815


WHILE @UID > 180088
BEGIN


    DECLARE @url varchar(max)
	
	SET @URL = ''

    DECLARE @authHeader NVARCHAR(64);
	DECLARE @contentType NVARCHAR(64);
	DECLARE @ret INT;
	DECLARE @status NVARCHAR(32);
	DECLARE @token INT;
	DECLARE @JSONTABLE TABLE (JSONCOL text) 

	SET @authHeader = 'No Auth';
	SET @contentType = 'application/json';

	SET @url = CONCAT('https://tbc.wowhead.com/tooltip/item/', @UID,'&json')

	-- Open the connection.
	EXEC @ret = sp_OACreate 'MSXML2.ServerXMLHTTP', @token OUT;
	IF @ret <> 0 RAISERROR('Unable to open HTTP connection.', 10, 1);

	-- Send the request.
	EXEC @ret = sp_OAMethod @token, 'open', NULL, 'GET', @url, 'false';
	EXEC @ret = sp_OAMethod @token, 'setRequestHeader', NULL, 'Authentication', @authHeader;
	EXEC @ret = sp_OAMethod @token, 'setRequestHeader', NULL, 'Content-type', @contentType;
	EXEC @ret = sp_OAMethod @token, 'send'

	-- Handle the response.
	EXEC @ret = sp_OAGetProperty @token, 'status', @status OUT;


	INSERT INTO @JSONTABLE (JSONCOL)
	EXEC @ret = sp_OAMethod @token, 'responseText'

	-- Show the response.
	--PRINT 'Status: ' + @status + ' (' + @statusText + ')';

	-- Close the connection.
	EXEC @ret = sp_OADestroy @token;
	IF @ret <> 0 RAISERROR('Unable to close HTTP connection.', 10, 1);


	INSERT INTO WOWHEAD_ITEMLIST (JSONCOL, ITEMID, URL_STRING)
	SELECT JSONCOL, @UID, @url
	FROM @JSONTABLE

	DELETE FROM @JSONTABLE

	SET @UID = @UID - 1

END










DECLARE @UID INT
SET @UID = 39656


WHILE @UID > 0
BEGIN


    DECLARE @url varchar(max)
	
	SET @URL = ''

    DECLARE @authHeader NVARCHAR(64);
	DECLARE @contentType NVARCHAR(64);
	DECLARE @ret INT;
	DECLARE @status NVARCHAR(32);
	DECLARE @token INT;
	DECLARE @JSONTABLE TABLE (JSONCOL text) 

	SET @authHeader = 'No Auth';
	SET @contentType = 'application/json';

	SET @url = CONCAT('https://tbc.wowhead.com/tooltip/item/', @UID,'&json')

	-- Open the connection.
	EXEC @ret = sp_OACreate 'MSXML2.ServerXMLHTTP', @token OUT;
	IF @ret <> 0 RAISERROR('Unable to open HTTP connection.', 10, 1);

	-- Send the request.
	EXEC @ret = sp_OAMethod @token, 'open', NULL, 'GET', @url, 'false';
	EXEC @ret = sp_OAMethod @token, 'setRequestHeader', NULL, 'Authentication', @authHeader;
	EXEC @ret = sp_OAMethod @token, 'setRequestHeader', NULL, 'Content-type', @contentType;
	EXEC @ret = sp_OAMethod @token, 'send'

	-- Handle the response.
	EXEC @ret = sp_OAGetProperty @token, 'status', @status OUT;


	INSERT INTO @JSONTABLE (JSONCOL)
	EXEC @ret = sp_OAMethod @token, 'responseText'

	-- Show the response.
	--PRINT 'Status: ' + @status + ' (' + @statusText + ')';

	-- Close the connection.
	EXEC @ret = sp_OADestroy @token;
	IF @ret <> 0 RAISERROR('Unable to close HTTP connection.', 10, 1);


	INSERT INTO WOWHEAD_ITEMLIST (JSONCOL, ITEMID, URL_STRING)
	SELECT JSONCOL, @UID, @url
	FROM @JSONTABLE

	DELETE FROM @JSONTABLE

	SET @UID = @UID - 1

END
