

/*
Query Below returns all equitable items quality 2 and above
There are few additional columns included, but only column ITEMID will be used further
*/

SELECT   ITEMID
		,NAME
		,QUALITY
		,ICON
		,TOOLTIP
		,REPLACE(REPLACE(REPLACE(REPLACE(SUBSTRING(TOOLTIP,CHARINDEX('<br>Item Level <!--ilvl-->',TOOLTIP)+26,3),'<',''),'/',''),'-->',''),'le=','') ILVL
		
FROM WOWHEAD_ITEMLIST 
CROSS APPLY OPENJSON (JSONCOL)
WITH
(	
	 NAME    VARCHAR(1000) '$.name'	
	,QUALITY INT		   '$.quality'
	,ICON    varchar(1000) '$.icon'
	,TOOLTIP varchar(max)  '$.tooltip'
)
WHERE 
(
     TOOLTIP LIKE ('%<td>Head</td>%')
  or TOOLTIP like ('%<td>Neck</td>%')
  or TOOLTIP like ('%<td>Shoulder</td>%')
  or TOOLTIP like ('%<td>Back</td>%')
  or TOOLTIP like ('%<td>Chest</td>%')
  or TOOLTIP like ('%<td>Wrist</td>%')
  or TOOLTIP like ('%<td>Hands</td>%')
  or TOOLTIP like ('%<td>Waist</td>%')
  or TOOLTIP like ('%<td>Legs</td>%')
  or TOOLTIP like ('%<td>Feet</td>%')
  or TOOLTIP like ('%<td>Finger</td>%')
  or TOOLTIP like ('%<td>Trinket</td>%')
  or TOOLTIP like ('%<td>Ranged</td>%')
  or TOOLTIP like ('%<td>Thrown</td>%')
  or TOOLTIP like ('%<td>Relic</td>%')
  or TOOLTIP like ('%<td>Main Hand</td>%')
  or TOOLTIP like ('%<td>Two-Hand</td>%')
  or TOOLTIP like ('%<td>One-Hand</td>%')
  or TOOLTIP like ('%<td>Off Hand</td>%')
  or TOOLTIP like ('%<td>Held In Off-hand</td>%')
)
AND QUALITY NOT IN ('0','1')
AND NAME NOT LIKE ('%PH]%')
AND NAME NOT LIKE ('%TEST%')
AND NAME NOT LIKE ('%Bracer 3%')
AND NAME NOT LIKE ('%Bracer 2%')
AND NAME NOT LIKE ('%Bracer 1%')
AND NAME NOT LIKE ('%Boots 3%')
AND NAME NOT LIKE ('%Boots 2%')
AND NAME NOT LIKE ('%Boots 1%')
AND NAME NOT LIKE ('zzOLD%')
AND NAME NOT LIKE ('130 Epic%')
AND NAME NOT LIKE ('%Indalamar%')
AND NAME NOT LIKE ('%QR XXXX%')
AND NAME NOT LIKE ('%Deprecated: Keanna%')
AND NAME NOT LIKE ('%90 Epic%')
AND NAME NOT LIKE ('%66 Epic%')
AND NAME NOT LIKE ('%63 Blue%')
AND NAME NOT LIKE ('%90 Green%')
AND NAME NOT LIKE ('%63 Green%')
AND NAME NOT LIKE ('%Pattern:%')
AND NAME NOT LIKE ('%Design:%')
AND NAME NOT LIKE ('%Schematic:%')
AND NAME NOT LIKE ('%Plans:%')
AND NAME NOT LIKE ('%Recipe:%')
--AND ITEMID = ('24262')
AND REPLACE(REPLACE(REPLACE(REPLACE(SUBSTRING(TOOLTIP,CHARINDEX('<br>Item Level <!--ilvl-->',TOOLTIP)+26,3),'<',''),'/',''),'-->',''),'le=','') >= 70



SELECT   ITEMID
		,NAME
		,QUALITY
		,ICON
		,TOOLTIP
		,REPLACE(REPLACE(REPLACE(REPLACE(SUBSTRING(TOOLTIP,CHARINDEX('<br>Item Level <!--ilvl-->',TOOLTIP)+26,3),'<',''),'/',''),'-->',''),'le=','') ILVL
		
FROM WOWHEAD_ITEMLIST 
CROSS APPLY OPENJSON (JSONCOL)
WITH
(	
	 NAME    VARCHAR(1000) '$.name'	
	,QUALITY INT		   '$.quality'
	,ICON    varchar(1000) '$.icon'
	,TOOLTIP varchar(max)  '$.tooltip'
)
WHERE NAME LIKE ('%SPELLSTRIKE%')

/*
Below Query returns list of Gems
Only column ITEMID will be used further
*/



SELECT   ITEMID
		,NAME
		,QUALITY
		,ICON
		,TOOLTIP
		,REPLACE(REPLACE(REPLACE(REPLACE(SUBSTRING(TOOLTIP,CHARINDEX('<br>Item Level <!--ilvl-->',TOOLTIP)+26,3),'<',''),'/',''),'-->',''),'le=','') ILVL
FROM WOWHEAD_ITEMLIST 
CROSS APPLY OPENJSON (JSONCOL)
WITH
(	
	 NAME    VARCHAR(1000) '$.name'	
	,QUALITY INT		   '$.quality'
	,ICON    varchar(1000) '$.icon'
	,TOOLTIP varchar(max)  '$.tooltip'
)
WHERE TOOLTIP LIKE ('%Matches%')
  AND NAME NOT LIKE ('zzOLD%')
  AND NAME NOT LIKE ('%TEST%')
AND NAME NOT LIKE ('%Pattern:%')
AND NAME NOT LIKE ('%Design:%')
AND NAME NOT LIKE ('%Schematic:%')
AND NAME NOT LIKE ('%Plans:%')
AND NAME NOT LIKE ('%Recipe:%')

union all
  
SELECT   ITEMID
		,NAME
		,QUALITY
		,ICON
		,TOOLTIP
		,REPLACE(REPLACE(REPLACE(REPLACE(SUBSTRING(TOOLTIP,CHARINDEX('<br>Item Level <!--ilvl-->',TOOLTIP)+26,3),'<',''),'/',''),'-->',''),'le=','') ILVL
FROM WOWHEAD_ITEMLIST 
CROSS APPLY OPENJSON (JSONCOL)
WITH
(	
	 NAME    VARCHAR(1000) '$.name'	
	,QUALITY INT		   '$.quality'
	,ICON    varchar(1000) '$.icon'
	,TOOLTIP varchar(max)  '$.tooltip'
)
WHERE 
TOOLTIP LIKE ('%gemConditions%')
  AND NAME NOT LIKE ('zz%')
  AND NAME NOT LIKE ('%TEST%')
  AND NAME NOT LIKE ('%Pattern:%')
  AND NAME NOT LIKE ('%Design:%')
  AND NAME NOT LIKE ('%Schematic:%')
  AND NAME NOT LIKE ('%Plans:%')
  AND NAME NOT LIKE ('%Recipe:%')
