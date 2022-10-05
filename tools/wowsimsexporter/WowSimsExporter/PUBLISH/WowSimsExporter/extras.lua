-- Author      : generalwrex (Natop on Myzrael)
-- Create Date : 2/6/2022 10:35:32 AM

if not WowSimsExporter then WowSimsExporter = {} end

WowSimsExporter.supportedSims = {
"hunter",
"mage",
"shaman",
"priest",
"rogue",
"druid",
"warrior",
"warlock",
"paladin"
}

WowSimsExporter.slotNames = {
    "HeadSlot",
    "NeckSlot",
    "ShoulderSlot",
    "BackSlot",
    "ChestSlot",
    "WristSlot",
    "HandsSlot",
    "WaistSlot",
    "LegsSlot",
    "FeetSlot",
    "Finger0Slot",
    "Finger1Slot",
    "Trinket0Slot",
    "Trinket1Slot",
    "MainHandSlot",
    "SecondaryHandSlot",
	"RangedSlot",
    "AmmoSlot",
}


WowSimsExporter.prelink = "https://wowsims.github.io/tbc/"
WowSimsExporter.postlink = ""
WowSimsExporter.specializations = {

	-- shaman
	{comparator = function(A,B,C) return A > B and A > C end, spec="elemental", class="shaman", url="elemental_shaman"},
	{comparator = function(A,B,C) return B > A and B > C end, spec="enhancement", class="shaman",url="enhancement_shaman"},
	-- hunter
	{comparator = function(A,B,C) return A > B and A > C end, spec="beast_mastery", class="hunter",url="hunter"},
	{comparator = function(A,B,C) return B > A and B > C end, spec="marksman", class="hunter",url="hunter"},
	{comparator = function(A,B,C) return C > A and C > B end, spec="survival", class="hunter",url="hunter"},
	-- druid
	{comparator = function(A,B,C) return A > B and A > C end, spec="balance", class="druid",url="balance_druid"},
	--{comparator = function(A,B,C) return B > A and B > C end, spec="feral", class="druid",url="feral_druid"},
	-- warlock
	{comparator = function(A,B,C) return true end, spec="warlock", class="warlock",url="warlock"},
	--{comparator = function(A,B,C) return B > A and B > C end, spec="demonology", class="warlock", url="affliction_warlock"},	
	--{comparator = function(A,B,C) return B > A and B > C end, spec="demonology", class="warlock", url="demonology_warlock"},	
	--{comparator = function(A,B,C) return C > A and C > B end, spec="destruction",class="warlock", url="destruction_warlock"},
	-- rogue
	{comparator = function(A,B,C) return A > B and A > C end, spec="assassination", class="rogue", url="rogue"},
	{comparator = function(A,B,C) return B > A and B > C end, spec="combat", class="rogue", url="rogue"},
	{comparator = function(A,B,C) return C > A and C > B end, spec="subtlety", class="rogue", url="rogue"},
	-- mage
	{comparator = function(A,B,C) return A > B and A > C end, spec="arcane", class="mage", url="mage"},
	{comparator = function(A,B,C) return B > A and B > C end, spec="fire", class="mage", url="mage"},
	{comparator = function(A,B,C) return C > A and C > B end, spec="frost", class="mage", url="mage"},
	-- warrior
	{comparator = function(A,B,C) return A > B and A > C end, spec="arms", class="warrior", url="arms_warrior"},
	{comparator = function(A,B,C) return B > A and B > C end, spec="fury", class="warrior", url="fury_warrior"},
	-- paladin
	{comparator = function(A,B,C) return true            end, spec="retribution", class="paladin", url="retribution_paladin"},
	-- priest
	{comparator = function(A,B,C) return C > A and C > B end, spec="shadow", class="priest", url="shadow_priest"},
	{comparator = function(A,B,C) return B > A and B > C end, spec="holy", class="priest", url="smite_priest"},
}

-- table extension contains
function table.contains(table, element)
  for _, value in pairs(table) do
    if value == element then
      return true
    end
  end
  return false
end