-- Author      : generalwrex (Natop on Myzrael TBC)
-- Create Date : 1/28/2022 9:30:08 AM
--
-- Update Date : 2023-04-16 Riotdog-GehennasEU: v2.5 - exporting bag items for bulk sim, fixes use of legacy APIs in libs and corrects link order (LibStub must come first).
--

WowSimsExporter = LibStub("AceAddon-3.0"):NewAddon("WowSimsExporter", "AceConsole-3.0", "AceEvent-3.0")


WowSimsExporter.Character = ""
WowSimsExporter.Link = "https://wowsims.github.io/wotlk/"

local AceGUI = LibStub("AceGUI-3.0")
local LibParse = LibStub("LibParse")

local version = "2.5"

-- Ulduar WotLK classic patch moved to the retail API.
local GetContainerNumSlots = C_Container and C_Container.GetContainerNumSlots or _G.GetContainerNumSlots
local GetContainerItemLink = C_Container and C_Container.GetContainerItemLink or _G.GetContainerItemLink

local defaults = {
	profile = {
		--updateGearChange = true,
	},
}

local options = { 
	name = "WowSimsExporter",
	handler = WowSimsExporter,
	type = "group",
	args = {
		--updateGearChange = {
			--type = "toggle",
			--name = "Update on Gear Change",
			--desc = "Update your data when you change gear pieces.",
			--get = "isGearChangeSet",
			--set = "setGearChange"
		--},
		openExporterButton = {
			type = "execute",
			name = "Open Exporter Window",
			desc = "Opens the exporter window",
			func = function() WowSimsExporter:CreateWindow() end
		},
	},
}


function WowSimsExporter:CreateCharacterStructure(unit)
    local name, realm = UnitFullName(unit)
    local locClass, engClass, locRace, engRace, gender, name, server = GetPlayerInfoByGUID(UnitGUID(unit))
    local level = UnitLevel(unit)

    self.Character = {
        name = name,
        realm = realm,
        race = engRace:gsub( "Scourge", "Undead"), -- hack? lol
        class = engClass:lower(),
		level = tonumber(level),
        talents = "",
		glyphs = { major = { }, minor = { } }, --wotlk
        professions = { }, --{ name = "", level = "" }, --wotlk
		spec  =  self:CheckCharacterSpec(engClass:lower()),
        gear = { items = { } } 
	}

    return self.Character
end

function WowSimsExporter:CreateGlyphEntry()
	local minor = {}
	local major = {}
	for t = 1, 6 do
		local enabled, glyphType, glyphSpellID = GetGlyphSocketInfo(t)
		if enabled and glyphSpellID then
			local localizedName = GetSpellInfo(glyphSpellID)
			if localizedName then
				local t = glyphType == 1 and major or minor
				table.insert(t, {["name"] = localizedName, ["spellID"] = glyphSpellID})
			end
		end
		self.Character.glyphs.major = major
		self.Character.glyphs.minor = minor
	end
end

function WowSimsExporter:CreateProfessionEntry()
	local names = WowSimsExporter.professionNames
	local names_inv = tInvert(names)
	local professions = {}
 
	for i = 1, GetNumSkillLines() do
		local name, _, _, skillLevel = GetSkillLineInfo(i)
		if names_inv[name] then
			table.insert(professions, { name = name, level = skillLevel })		
		end
	end
	self.Character.professions = professions
end

function WowSimsExporter:CreateTalentEntry()
    local talents = {}

    local numTabs = GetNumTalentTabs()
    for t = 1, numTabs do
        local numTalents = GetNumTalents(t)
        for i = 1, numTalents do
            local nameTalent, icon, tier, column, currRank, maxRank = self:GetOrderedTalentInfo(t, i)

            table.insert(talents, tostring(currRank))
        end
        if (t < 3) then
            table.insert(talents, "-")
        end
    end

    return table.concat(talents)
end

function WowSimsExporter:CheckCharacterSpec(class)

	local specs = self.specializations

	T1 = GetNumTalents(1)
	T2 = GetNumTalents(2)
	T3 = GetNumTalents(3)

	local spec = ""

	for i, character in ipairs(specs) do	
		if character then				
			if (character.class == class) then																			
				if character.comparator(T1,T2,T3) then
					spec = character.spec											
					break
				end		
			end																
		end
	end	
	return spec
end

function WowSimsExporter:OpenWindow(input)
    if not input or input:trim() == "" then
		self:CreateWindow()
	elseif (input == "export") then        
        self:CreateWindow(true)
    elseif (input=="options") then           
		InterfaceOptionsFrame_OpenToCategory(self.optionsFrame)
		InterfaceOptionsFrame_OpenToCategory(self.optionsFrame)
    end
end

function WowSimsExporter:createItemFromItemLink(itemLink)
	local Id, Enchant, Gem1, Gem2, Gem3, Gem4 = self:ParseItemLink(itemLink)
	item = {}
	item.id = tonumber(Id)
	item.enchant = tonumber(Enchant)
	item.gems = {tonumber(Gem1), tonumber(Gem2), tonumber(Gem3), tonumber(Gem4)}
	return item
end

	-- TODO(Riotdog-GehennasEU): Is this sufficient? This seems to be what simc uses:
	-- https://github.com/simulationcraft/simc-addon/blob/master/core.lua
	-- Except we don't need the artifact check for wotlk classic.
function considerItemReplacement(itemLink)
	if not IsEquippableItem(itemLink) then
		return false
	end

	local _, _, itemRarity, itemLevel = GetItemInfo(itemLink)

	-- Ignore TBC items like Rocket Boots Xtreme (Lite). The ilvl limit is intentionally set low
	-- to limit accidental filtering.
	if itemLevel <= 112 then
		return false
	end

	-- https://wowwiki-archive.fandom.com/wiki/API_TYPE_Quality
	-- 3 = Rare, 4 = Epic, 5 = Legendary
	return itemRarity == 3 or itemRarity == 4 or itemRarity == 5
end

function WowSimsExporter:GetGearEnchantGems(withBags)
	self.Character.gear = {}
	self.Character.glyphs = {}
	self.Character.bagItems = {}
	self.Character.items = nil

	if not withBags then
		local equippedGear = {}
		local slotNames = WowSimsExporter.slotNames
		for slotNum = 1, #slotNames do
				local slotId = GetInventorySlotInfo(slotNames[slotNum])
				local itemLink = GetInventoryItemLink("player", slotId)
				if itemLink then
					equippedGear[slotNum] = self:createItemFromItemLink(itemLink)
				end
		end

		self.Character.spec = self:CheckCharacterSpec(self.Character.class)
		self.Character.talents = self:CreateTalentEntry()
		self:CreateGlyphEntry() -- wotlk
		self:CreateProfessionEntry() -- wotlk
		self.Character.gear.items = equippedGear
		return self.Character
	end

	local bagGear = {}
	for bag = 0, 4 do
		for slot = 1, GetContainerNumSlots(bag) do
			local itemLink = GetContainerItemLink(bag, slot)
			if itemLink and considerItemReplacement(itemLink) then
				table.insert(bagGear, self:createItemFromItemLink(itemLink))
			end
		end
	end
	DEFAULT_CHAT_FRAME:AddMessage(("[|cffFFFF00WowSimsExporter|r] Exported %d items from bags."):format(#bagGear))
  return {["items"] = bagGear}
end


function WowSimsExporter:ParseItemLink(itemLink)
    local _, _, Color, Ltype, Id, Enchant, Gem1, Gem2, Gem3, Gem4, Suffix, Unique, LinkLvl, Name =
        string.find(
        itemLink,
        "|?c?f?f?(%x*)|?H?([^:]*):?(%d+):?(%d*):?(%d*):?(%d*):?(%d*):?(%d*):?(%-?%d*):?(%-?%d*):?(%d*):?(%d*):?(%-?%d*)|?h?%[?([^%[%]]*)%]?|?h?|?r?"
    )
    return Id, Enchant, Gem1, Gem2, Gem3, Gem4
end

function WowSimsExporter:OnInitialize()

	self.db = LibStub("AceDB-3.0"):New("WSEDB", defaults, true)

	LibStub("AceConfig-3.0"):RegisterOptionsTable("WowSimsExporter", options)
	self.optionsFrame = LibStub("AceConfigDialog-3.0"):AddToBlizOptions("WowSimsExporter", "WowSimsExporter")

	local profiles = LibStub("AceDBOptions-3.0"):GetOptionsTable(self.db)
	LibStub("AceConfig-3.0"):RegisterOptionsTable("WowSimsExporter_Profiles", profiles)
	LibStub("AceConfigDialog-3.0"):AddToBlizOptions("WowSimsExporter_Profiles", "Profiles", "WowSimsExporter")

    self:RegisterChatCommand("wse", "OpenWindow")
    self:RegisterChatCommand("wowsimsexporter", "OpenWindow")
    self:RegisterChatCommand("wsexporter", "OpenWindow")

    self:Print("WowSimsExporter v" .. version .. " Initialized. use /wse For Window.")

end


-- UI
function WowSimsExporter:BuildLinks(frame, character)
	local specs = self.specializations
	local supportedsims =  self.supportedSims 
	local class = character.class
	local spec  = character.spec

	if table.contains(supportedsims, class) then

		for i, char in ipairs(specs) do
			if char and char.class == class and char.spec == spec then

				local link = WowSimsExporter.prelink..(char.url)..WowSimsExporter.postlink

				local l = AceGUI:Create("InteractiveLabel")
				l:SetText("Click to copy: "..link.."\r\n")
				l:SetFullWidth(true)
				l:SetCallback("OnClick", function()		
					WowSimsExporter:CreateCopyDialog(link)
				end)
				frame:AddChild(l)
			end
		end
	end
end


function WowSimsExporter:CreateCopyDialog(text)

	local frame = AceGUI:Create("Frame")
	frame:SetTitle("WSE Copy Dialog")
    frame:SetStatusText("Use CTRL+C to copy link")
    frame:SetLayout("Flow")
	frame:SetWidth(400)
	frame:SetHeight(100)
	frame:SetCallback(
        "OnClose",
        function(widget)
            AceGUI:Release(widget)
        end
    )

	local editbox = AceGUI:Create("EditBox")
    editbox:SetText(text)
    editbox:SetFullWidth(true)
    editbox:DisableButton(true)

	editbox:SetFocus()
	editbox:HighlightText()
	
	frame:AddChild(editbox)

end

function WowSimsExporter:CreateWindow(generate)
	local char = self:CreateCharacterStructure("player")
	
    local frame = AceGUI:Create("Frame")
    frame:SetCallback(
        "OnClose",
        function(widget)
            AceGUI:Release(widget)
        end
    )
    frame:SetTitle("WowSimsExporter V" .. version .. "")
    frame:SetStatusText("Click 'Generate Data' to generate exportable data")
    frame:SetLayout("Flow")


    local jsonbox = AceGUI:Create("MultiLineEditBox")
    jsonbox:SetLabel("Copy and paste into the websites importer!")
    jsonbox:SetFullWidth(true)
    jsonbox:SetFullHeight(true)
    jsonbox:DisableButton(true)
   
	local function l_Generate(withBags)
		jsonbox:SetText('')
		if not withBags then
			WowSimsExporter.Character = WowSimsExporter:GetGearEnchantGems(withBags)
			jsonbox:SetText(LibParse:JSONEncode(WowSimsExporter.Character))
		else
			local bagData = WowSimsExporter:GetGearEnchantGems(withBags)
			jsonbox:SetText(LibParse:JSONEncode(bagData))
		end
		jsonbox:HighlightText()
		jsonbox:SetFocus()

		frame:SetStatusText("Data Generated!")
	end

	if generate then l_Generate() end

  local button = AceGUI:Create("Button")
    button:SetText("Generate Data (Equipped Only)")
    button:SetWidth(300)
	button:SetCallback("OnClick", function()		
		l_Generate(false)
	end)

	local extraButton = AceGUI:Create("Button")
	extraButton:SetText("Batch: Export Bag Items")
	extraButton:SetWidth(300)
	extraButton:SetCallback("OnClick", function()		
		l_Generate(true)
	end)
	
	local icon = AceGUI:Create("Icon")
	icon:SetImage("Interface\\AddOns\\wowsimsexporter\\Skins\\wowsims.tga") 
	icon:SetImageSize(32, 32)
	icon:SetFullWidth(true)

  local label = AceGUI:Create("Label")
	label:SetFullWidth(true)
  label:SetText([[

To upload your character to the simuator, click on the url below that leads to the simuator website.

You will find an Import button on the top right of the simulator named "Import". Click that and select the "Addon" tab, paste the data
into the provided box and click "Import"

]])

	if not table.contains(self.supportedSims, char.class) then

		frame:AddChild(icon)

		local l1 = AceGUI:Create("Heading")
		l1:SetText("")
		l1:SetColor(255,0,0)
		l1:SetFullWidth(true)
		frame:AddChild(l1)


		local l = AceGUI:Create("Label")
		l:SetText("Your characters class is currently unsupported. The supported classes are currently;\n"..table.concat(self.supportedSims,"\n"))
		l:SetColor(255,0,0)
		l:SetFullWidth(true)
		frame:AddChild(l)
	else

		frame:AddChild(icon)
		frame:AddChild(label)
		WowSimsExporter:BuildLinks(frame, char)
		frame:AddChild(button)
		frame:AddChild(extraButton)
		frame:AddChild(jsonbox)

	end

end
-- Borrowed from rating buster!!
-- As of Classic Patch 3.4.0, GetTalentInfo indices no longer correlate
-- to their positions in the tree. Building a talent cache ordered by
-- tier then column allows us to replicate the previous behavior.
local orderedTalentCache = {}
do
	local f = CreateFrame("Frame")
	f:RegisterEvent("SPELLS_CHANGED")
	f:SetScript("OnEvent", function()
		local temp = {}
		for tab = 1, GetNumTalentTabs() do
			temp[tab] = {}
			local products = {}
			for i = 1,GetNumTalents(tab) do
				local name, _, tier, column = GetTalentInfo(tab,i)
				local product = (tier - 1) * 4 + column
				temp[tab][product] = i
				table.insert(products, product)
			end

			table.sort(products)

			orderedTalentCache[tab] = {}
			local j = 1
			for _, product in ipairs(products) do
				orderedTalentCache[tab][j] = temp[tab][product]
				j = j + 1
			end
		end
		f:UnregisterEvent("SPELLS_CHANGED")
	end)
end

function WowSimsExporter:GetOrderedTalentInfo(tab, num)
	return GetTalentInfo(tab, orderedTalentCache[tab][num])
end



function WowSimsExporter:OnEnable()
end

function WowSimsExporter:OnDisable()
end

function WowSimsExporter:isGearChangeSet(info)
	return self.db.profile.updateGearChange
end

function WowSimsExporter:setGearChange(info, value)
	self.db.profile.updateGearChange = value
end
